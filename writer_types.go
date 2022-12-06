package dj

import "context"

// ChanWriter returns a writer that writes to the given channel.
func ChanWriter[T any](ch chan<- T) Writer[T] {
	return NewWriter[T](&chanWriter[T]{ctx: context.Background(), ch: ch})
}

// ChanWriterCtx returns a writer that writes to the given channel until the given context is canceled.
func ChanWriterCtx[T any](ctx context.Context, ch chan<- T) Writer[T] {
	return NewWriter[T](&chanWriter[T]{ctx: ctx, ch: ch})
}

type chanWriter[T any] struct {
	ctx context.Context
	ch  chan<- T
}

func (i *chanWriter[T]) Write(v T) error {
	select {
	case <-i.ctx.Done():
		return i.ctx.Err()

	case i.ch <- v:
		return nil
	}
}

// FuncWriter returns a writer that calls the given function for each value to be written.
func FuncWriter[T any](fn func(T) error) Writer[T] {
	return NewWriter[T](&funcWriter[T]{fn: fn})
}

type funcWriter[T any] struct {
	fn func(T) error
}

func (i *funcWriter[T]) Write(v T) error {
	return i.fn(v)
}
