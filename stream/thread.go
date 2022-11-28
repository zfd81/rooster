package stream

import "sync"

// Recover is used with defer to do cleanup on panics.
func Recover(cleanups ...func()) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	for _, cleanup := range cleanups {
		cleanup()
	}
}

// GoSafe runs the given fn using another goroutine, recovers if fn panics.
func GoSafe(fn func()) {
	go RunSafe(fn)
}

// RunSafe runs the given fn, recovers if fn panics.
func RunSafe(fn func()) {
	defer Recover()
	fn()
}

// RoutineGroup is used to group goroutines together and all wait all goroutines to be done.
type RoutineGroup struct {
	waitGroup sync.WaitGroup
}

// Run runs the given fn in RoutineGroup.
func (g *RoutineGroup) Run(fn func()) {
	g.waitGroup.Add(1)
	go func() {
		defer g.waitGroup.Done()
		fn()
	}()
}

// RunSafe runs the given fn in RoutineGroup, and avoid panics.
func (g *RoutineGroup) RunSafe(fn func()) {
	g.waitGroup.Add(1)
	GoSafe(func() {
		defer g.waitGroup.Done()
		fn()
	})
}

// Wait waits all running functions to be done.
func (g *RoutineGroup) Wait() {
	g.waitGroup.Wait()
}

// NewRoutineGroup returns a RoutineGroup.
func NewRoutineGroup() *RoutineGroup {
	return new(RoutineGroup)
}
