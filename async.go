package dj

import (
	"context"
	"sync"
)

// Group is a group of goroutines.
type Group struct {
	ctx context.Context
	sem Sem
	wg  sync.WaitGroup
}

// NewGroup returns a new group of goroutines.
// The given context is used to cancel all goroutines in the group.
// The number of goroutines permitted to run concurrently is given by n.
func NewGroup(ctx context.Context, sem Sem) *Group {
	return &Group{
		ctx: ctx,
		sem: sem,
	}
}

func (s *Group) Wait() {
	s.wg.Wait()
}

func (s *Group) Go(fn func(context.Context)) {
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

func (s *Group) GoN(n int, fn func(context.Context, int)) {
	ForN(n, func(i int) {
		s.Go(func(ctx context.Context) {
			fn(ctx, i)
		})
	})
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
