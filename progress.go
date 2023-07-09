package progressbar

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// interface that defines the Progress struct
type Progress interface {
	// also a string interface - here is where you put any additional information
	// about the given progress report to write to the console.
	String() string
	// will return the progress as a percent representation of the current progress
	ProgressPercent() float64
}

// construct the progress bar from a given progress report
func ConstructBar(progress Progress) string {
	results := strings.Builder{}
	results.Grow(125)
	results.WriteByte('[')
	currentProgress := int(math.Floor(progress.ProgressPercent())) / 2
	max := 50
	for i := 0; i < int(currentProgress); i++ {
		results.WriteByte('|')
	}
	for i := int(currentProgress); i < max; i++ {
		results.WriteByte(' ')
	}
	results.WriteByte(']')
	results.WriteString(fmt.Sprintf(" %s", progress.String()))
	return results.String()
}

// float-only progress
type BasicProgress float64

// use when you have an ETA
type ProgressWithRemaining struct {
	Progress  float64
	Remaining time.Duration
}

func (b BasicProgress) String() string { return fmt.Sprintf("%0f%%", b.ProgressPercent()) }

func (b BasicProgress) ProgressPercent() float64 { return float64(b) * 100 }

func (p ProgressWithRemaining) String() string {
	return fmt.Sprintf("%0f%% (%v)", p.ProgressPercent(), p.Remaining)
}

func (p ProgressWithRemaining) ProgressPercent() float64 { return p.Progress * 100 }
