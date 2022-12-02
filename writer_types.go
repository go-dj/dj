package xn

import "context"

// SliceWriter returns a writer that writes to the given slice.
func SliceWriter[T any](slice []T) Writer[T] {
	return NewWriter[T](&sliceWriter[T]{slice: slice})
}

type sliceWriter[T any] struct {
	slice []T
	index int
}

func (i *sliceWriter[T]) Write(v T) bool {
	if i.index >= len(i.slice) {
		return false
	}

	i.slice[i.index] = v

	i.index++

	return true
}

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

func (i *chanWriter[T]) Write(v T) bool {
	select {
	case <-i.ctx.Done():
		return false

	case i.ch <- v:
		return true
	}
}

// FuncWriter returns a writer that calls the given function for each value to be written.
func FuncWriter[T any](fn func(T) bool) Writer[T] {
	return NewWriter[T](&funcWriter[T]{fn: fn})
}

type funcWriter[T any] struct {
	fn func(T) bool
}

func (i *funcWriter[T]) Write(v T) bool {
	return i.fn(v)
}
