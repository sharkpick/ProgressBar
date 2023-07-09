package progressbar

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestBasicProgress(t *testing.T) {
	progressbar := NewProgressBar()
	defer progressbar.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	RenderProgress(ctx, progressbar)
	for i := 1; i <= 5; i++ {
		progress := BasicProgress(float64(i) / float64(5))
		progressbar.UpdateProgress(progress)
		time.Sleep(time.Second)
	}
	fmt.Print(progressbar.GetClearLineCode())
	fmt.Println("finished BasicProgress")
}

func TestProgressWithRemaining(t *testing.T) {
	progressbar := NewProgressBar()
	defer progressbar.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	RenderProgress(ctx, progressbar)
	for i := 1; i <= 5; i++ {
		progress := ProgressWithRemaining{
			Progress:  float64(i) / float64(5),
			Remaining: time.Second * time.Duration(5-i),
		}
		progressbar.UpdateProgress(progress)
		time.Sleep(time.Second)
	}
	fmt.Print(progressbar.GetClearLineCode())
	fmt.Println("finished ProgressWithRemaining")
}

func TestLongDuration(t *testing.T) {
	progressbar := NewProgressBar()
	defer progressbar.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	RenderProgress(ctx, progressbar)
	for i := 1; i < 10000; i++ {
		progress := BasicProgress(float64(i) / float64(10000))
		progressbar.UpdateProgress(progress)
		n := rand.Intn(500)
		sleepTime := time.Microsecond * time.Duration(n)
		time.Sleep(sleepTime)
	}
	fmt.Print(progressbar.GetClearLineCode())
	fmt.Println("finished TestLongDuration")
}

func TestShortDuration(t *testing.T) {
	progressbar := NewProgressBar()
	defer progressbar.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	RenderProgress(ctx, progressbar)
	for i := 0; i < 5000; i++ {
		progress := BasicProgress(float64(i) / float64(5000))
		progressbar.UpdateProgress(progress)
		time.Sleep(time.Nanosecond)
	}
	fmt.Print(progressbar.GetClearLineCode())
	fmt.Println("finished TestShortDuration")
}
