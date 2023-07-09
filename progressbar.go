package progressbar

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

// check for compatibility
func init() {
	if os := runtime.GOOS; !strings.HasPrefix(os, "linux") {
		panic("error importing progressbar package: unsupported OS " + os)
	}
}

// ProgressBar struct
type ProgressBar struct {
	// unbuffered channel of Progress reports
	c chan Progress
	// saves a copy of the previous log statement to allow movement
	// and clearing the console.
	previousStatement AtomicGeneric[string]
}

func (p *ProgressBar) UpdateProgress(progress Progress) {
	p.c <- progress
}

// returns the statement to fmt.Print() (no newline) - moves the cursor to the left
// then clears through the end of the line.
func (p *ProgressBar) GetClearLineCode() string {
	return GetMoveLeftCode(len(p.previousStatement.Load())) + EraseRightInLineCode
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
func RenderProgress(ctx context.Context, progress *ProgressBar) {
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
				line := ConstructBar(got)
				progress.previousStatement.Store(line)
				fmt.Print(line)
			}
		}
	}()
}
