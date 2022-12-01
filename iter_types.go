package xn

import "context"

// Iter is a type that can be read sequentially.
type Iter[T any] interface {
	Iterable[T]

	// For calls the given function for each value in the iterator.
	For(func(T))

	// ForIdx calls the given function for each value in the iterator, along with the index of the value.
	ForIdx(func(int, T))

	// Collect returns a slice containing all the values in the iterator.
	Collect() []T

	// Chan returns a channel that will receive all the values in the iterator.
	Chan() <-chan T

	// WriteTo writes all the values in the iterator to the given writable.
	WriteTo(Writable[T]) (int, bool)
}

// NewIter returns a new Iter that reads from the given Iterable.
func NewIter[T any](r Iterable[T]) Iter[T] {
	return &iter[T]{Iterable: r}
}

// SliceIter returns an iterator over the given slice.
func SliceIter[T any](slice ...T) Iter[T] {
	return NewIter[T](&sliceIter[T]{slice: slice})
}

type sliceIter[T any] struct {
	slice []T
	index int
}

func (i *sliceIter[T]) Next() (T, bool) {
	if i.index >= len(i.slice) {
		return zero[T](), false
	}

	v := i.slice[i.index]

	i.index++

	return v, true
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

func (i *chanIter[T]) Next() (T, bool) {
	select {
	case <-i.ctx.Done():
		return zero[T](), false

	case v, ok := <-i.ch:
		return v, ok
	}
}

// FuncIter returns an iterator over the given function.
func FuncIter[T any](fn func() (T, bool)) Iter[T] {
	return NewIter[T](&funcIter[T]{fn: fn})
}

type funcIter[T any] struct {
	fn func() (T, bool)
}

func (i *funcIter[T]) Next() (T, bool) {
	return i.fn()
}
