package dj

import (
	"context"
	"sync"
)

// Group is a group of goroutines.
type Group struct {
	// ctx is the context in which the goroutines are running.
	ctx    context.Context
	cancel context.CancelFunc

	// sem limits the number of goroutines that can run concurrently.
	sem Sem

	// wg is the wait group for the goroutines.
	wg *sync.WaitGroup

	// parent is the parent group, if any.
	parent *Group
}

// NewGroup returns a new group of goroutines.
// The given context is used to cancel all goroutines in the group.
// The number of goroutines permitted to run concurrently is given by n.
func NewGroup(ctx context.Context, sem Sem) *Group {
	ctx, cancel := context.WithCancel(ctx)

	return &Group{
		ctx:    ctx,
		cancel: cancel,
		sem:    sem,
		wg:     &sync.WaitGroup{},
	}
}

// Wait waits for all goroutines in the group to finish.
func (s *Group) Wait() {
	s.wg.Wait()
}

// Cancel cancels all goroutines in the group.
// This is destructive and cannot be undone; the group cannot be reused.
func (s *Group) Cancel() {
	s.cancel()
}

// Go runs the given function n times concurrently in a child group
// and returns the child group.
func (s *Group) Go(n int, fn func(context.Context, int)) *Group {
	child := s.Child()

	ForN(n, func(i int) {
		child.exec(func(ctx context.Context) {
			fn(ctx, i)
		})
	})

	return child
}

// Child returns a sub group inside the group.
func (s *Group) Child() *Group {
	ctx, cancel := context.WithCancel(s.ctx)

	return &Group{
		ctx:    ctx,
		cancel: cancel,
		sem:    s.sem,
		wg:     s.wg,
		parent: s,
	}
}

// Parent returns the parent group, if any.
func (s *Group) Parent() *Group {
	return s.parent
}

// exec executes the given function in the group.
func (s *Group) exec(fn func(context.Context)) {
	s.sem.Acquire()
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		defer s.sem.Release()

		select {
		case <-s.ctx.Done():
			// Context was canceled, don't run the function.

		default:
			fn(s.ctx)
		}
	}()
}

// Sem is a semaphore.
// Goroutines can acquire and release permits from the semaphore.
// When the semaphore has no permits available, goroutines will block until a permit is released.
type Sem interface {
	Acquire()
	Release()
}

// chanSem is a semaphore implemented using a channel.
type chanSem struct {
	ch chan struct{}
}

// NewSem returns a new semaphore with the given number of permits.
// The semaphore is implemented using a channel.
// If n is less than zero, the semaphore is unbounded.
func NewSem(n int) Sem {
	var ch chan struct{}

	if n >= 0 {
		ch = make(chan struct{}, n)
	}

	return &chanSem{ch: ch}
}

// Acquire acquires a permit from the semaphore.
func (s *chanSem) Acquire() {
	if s.ch != nil {
		s.ch <- struct{}{}
	}
}

// Release releases a permit to the semaphore.
func (s *chanSem) Release() {
	if s.ch != nil {
		<-s.ch
	}
}
