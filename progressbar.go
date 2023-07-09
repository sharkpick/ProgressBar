package progressbar

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"strings"

	"github.com/sharkpick/stringcolor"
)

// check for compatibility - both this and the stringcolor package are Linux-only (at the moment at least)
func init() {
	if os := runtime.GOOS; !strings.HasPrefix(os, "linux") {
		panic("error importing progressbar package: unsupported OS " + os)
	}
}

type PreviousProgressReport struct {
	s string
	t Progress
}

// ProgressBar struct
type ProgressBar struct {
	// unbuffered channel of Progress reports
	c chan Progress
	// saves a copy of the previous log statement to allow movement
	// and clearing the console.
	previous AtomicGeneric[PreviousProgressReport]
}

func (p *ProgressBar) UpdateProgress(progress Progress) {
	p.c <- progress
}

// returns the statement to fmt.Print() (no newline) - moves the cursor to the left
// then clears through the end of the line.
func (p *ProgressBar) GetClearLineCode() string {
	return GetMoveLeftCode(len(p.previous.Load().s)) + EraseRightInLineCode
}

// closes the unbuffered channel
func (p *ProgressBar) Close() {
	close(p.c)
}

// necessary to construct the ProgressBar's channel
func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		c: make(chan Progress),
	}
}

// launch to monitor progress. when the context ends or the channel
// closes it will clear the previous output and return. contains its
// own goroutine to make it simple to call/use.
func (progress *ProgressBar) RenderProgress(ctx context.Context) {
	go func() {
		defer func() {
			fmt.Print(progress.GetClearLineCode())
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case got, ok := <-progress.c:
				if !ok {
					return
				}
				fmt.Print(progress.GetClearLineCode())
				line := progress.ConstructBar(got)
				progress.previous.Store(
					PreviousProgressReport{
						s: strings.Clone(line),
						t: got,
					},
				)
				fmt.Print(line)
			}
		}
	}()
}

func (p *ProgressBar) ConstructBar(progress Progress) string {
	results := strings.Builder{}
	results.Grow(125)
	currentProgressPercent := progress.ProgressPercent()
	currentProgress := int(math.Floor(currentProgressPercent)) / 2
	max := 50
	for i := 0; i < int(currentProgress); i++ {
		results.WriteByte('|')
	}
	for i := int(currentProgress); i < max; i++ {
		results.WriteByte(' ')
	}
	func() {
		buffer := stringcolor.ColorWrapString(stringcolor.Green, results.String())
		results.Reset()
		results.WriteByte('[')
		results.WriteString(buffer)
		results.WriteByte(']')
	}()
	results.WriteString(fmt.Sprintf(" %s", progress.String()))
	return results.String()
}
