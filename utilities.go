package progressbar

import (
	"fmt"
	"sync"
)

// ANSI escape codes - print with fmt.Print() (no newline)
const (
	// use fmt.Sprintf to fill in the number of cells to move
	MoveLeftFString = "\033[%dD"
	// erases through the end of the line
	EraseRightInLineCode = "\033[0K"
)

// formats the left move command
func GetMoveLeftCode(cells int) string {
	return fmt.Sprintf(MoveLeftFString, cells)
}

// AtomicGeneric is a simple atomic container around an interface{}
// it's used here to provide atomic strings
type AtomicGeneric[T any] struct {
	mu sync.RWMutex
	t  T
}

func (a *AtomicGeneric[T]) Load() T    { a.mu.RLock(); defer a.mu.RUnlock(); return a.t }
func (a *AtomicGeneric[T]) Store(t T)  { a.mu.Lock(); a.mu.Unlock(); a.t = t }
func (a *AtomicGeneric[T]) Swap(t T) T { a.mu.Lock(); a.mu.Unlock(); tmp := a.t; a.t = t; return tmp }
