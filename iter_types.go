package dj

import "context"

// SliceIter returns an iterator over the given slice.
func SliceIter[T any](slice ...T) Iter[T] {
	return NewIter[T](&sliceIter[T]{slice: slice})
}

type sliceIter[T any] struct {
	slice []T
	index int
}

func (i *sliceIter[T]) Read() (Result[T], bool) {
	if i.index >= len(i.slice) {
		return Ok(Zero[T]()), false
	}

	v := i.slice[i.index]

	i.index++

	return Ok(v), true
}

// ChanIter returns an iterator over the given channel.
func ChanIter[T any](ch <-chan T) Iter[T] {
	return NewIter[T](&chanIter[T]{ctx: context.Background(), ch: ch})
}

// ChanIterCtx returns an iterator over the given channel, which will be closed when the given context is canceled.
func ChanIterCtx[T any](ctx context.Context, ch <-chan T) Iter[T] {
	return NewIter[T](&chanIter[T]{ctx: ctx, ch: ch})
}

type chanIter[T any] struct {
	ctx context.Context
	ch  <-chan T
}

func (i *chanIter[T]) Read() (Result[T], bool) {
	select {
	case <-i.ctx.Done():
		return Err[T](i.ctx.Err()), false

	case v, ok := <-i.ch:
		return Ok(v), ok
	}
}

// FuncIter returns an iterator over the given function.
func FuncIter[T any](fn func() (Result[T], bool)) Iter[T] {
	return NewIter[T](&funcIter[T]{fn: fn})
}

type funcIter[T any] struct {
	fn func() (Result[T], bool)
}

func (i *funcIter[T]) Read() (Result[T], bool) {
	return i.fn()
}
