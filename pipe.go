package dj

import (
	"sync"
	"sync/atomic"
)

type Pipe[T any] struct {
	// items contains the items in the pipe.
	items []T

	// inCh and outCh are the channels for writing and reading from the pipe.
	inCh, outCh chan T

	// cond signals when items are added to the pipe.
	cond *sync.Cond

	// done holds whether the pipe is closed.
	done atomic.Bool
}

// NewPipe returns a new buffered pipe.
func NewPipe[T any]() (chan<- T, <-chan T) {
	p := &Pipe[T]{
		items: make([]T, 0),
		inCh:  make(chan T),
		outCh: make(chan T),
		cond:  sync.NewCond(&sync.Mutex{}),
	}

	go func() {
		for p.in() {
			// ...
		}
	}()

	go func() {
		defer close(p.outCh)

		for p.out() {
			// ...
		}
	}()

	return p.inCh, p.outCh
}

// in writes the next value from the in channel to the items.
func (p *Pipe[T]) in() bool {
	defer p.cond.Broadcast()

	v, ok := <-p.inCh
	if !ok {
		p.done.Store(true)
	} else {
		p.cond.L.Lock()
		defer p.cond.L.Unlock()

		p.items = append(p.items, v)
	}

	return ok
}

// out reads the next value from the items and sends it to the out channel.
func (p *Pipe[T]) out() bool {
	v, ok := p.pop()
	if !ok {
		return false
	}

	p.outCh <- v

	return true
}

// pop removes the first item from the items and returns it.
func (p *Pipe[T]) pop() (T, bool) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	for len(p.items) == 0 {
		if p.done.Load() {
			return Zero[T](), false
		}

		p.cond.Wait()
	}

	var v T

	v, p.items = p.items[0], p.items[1:]

	return v, true
}
