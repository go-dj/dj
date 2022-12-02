package dj

import (
	"sync"
	"sync/atomic"
)

// NewBufPipe returns a pair of connected channels, one for sending and one for receiving.
// Writes to the in channel block when the buffer is full.
// The out channel is closed when the in channel is closed.
func NewBufPipe[T any](size int) (chan<- T, <-chan T) {
	ch := make(chan T, size)
	return ch, ch
}

// NewPipe returns a pair of connected channels, one for sending and one for receiving.
// Sent values are buffered until they are received.
// The out channel is closed when the in channel is closed.
func NewPipe[T any]() (chan<- T, <-chan T) {
	p := &pipe[T]{
		inCh:  make(chan T),
		outCh: make(chan T),
		cond:  sync.NewCond(&sync.Mutex{}),
	}

	// in reads the next value from the in channel and adds it to the items.
	go func() {
		for p.in() {
			// ...
		}
	}()

	// Send an initial value on the in channel to ensure the in goroutine is running.
	p.inCh <- Zero[T]()

	// out reads the next value from the items and sends it to the out channel.
	go func() {
		defer close(p.outCh)

		for p.out() {
			// ...
		}
	}()

	// Receive an initial value from the out channel to ensure the out goroutine is running.
	<-p.outCh

	return p.inCh, p.outCh
}

// pipe is a pair of unbouded channels.
type pipe[T any] struct {
	// buf contains the buffered values in the pipe.
	buf []T

	// inCh and outCh are the channels for writing and reading from the pipe.
	inCh, outCh chan T

	// cond signals when items are added to the pipe.
	cond *sync.Cond

	// done holds whether the pipe is closed.
	done atomic.Bool
}

// in writes the next value from the in channel to the items.
func (p *pipe[T]) in() bool {
	defer p.cond.Broadcast()

	v, ok := <-p.inCh
	if !ok {
		p.done.Store(true)
	} else {
		p.cond.L.Lock()
		defer p.cond.L.Unlock()

		p.buf = append(p.buf, v)
	}

	return ok
}

// out reads the next value from the items and sends it to the out channel.
func (p *pipe[T]) out() bool {
	v, ok := p.pop()
	if !ok {
		return false
	}

	p.outCh <- v

	return true
}

// pop removes the first item from the items and returns it.
func (p *pipe[T]) pop() (T, bool) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	for len(p.buf) == 0 {
		if p.done.Load() {
			return Zero[T](), false
		}

		p.cond.Wait()
	}

	var v T

	v, p.buf = p.buf[0], p.buf[1:]

	return v, true
}
