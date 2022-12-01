package xn

// IterSlice returns an iterator over the given slice.
func IterSlice[T any](slice ...T) Iter[T] {
	return &iterBase[T]{iterable: &iterSlice[T]{slice: slice}}
}

type iterSlice[T any] struct {
	slice []T
	index int
}

func (i *iterSlice[T]) Next() (T, bool) {
	if i.index >= len(i.slice) {
		return zero[T](), false
	}

	v := i.slice[i.index]

	i.index++

	return v, true
}

// IterChan returns an iterator over the given channel.
func IterChan[T any](ch <-chan T) Iter[T] {
	return &iterBase[T]{iterable: &iterChan[T]{ch: ch}}
}

type iterChan[T any] struct {
	ch <-chan T
}

func (i *iterChan[T]) Next() (T, bool) {
	v, ok := <-i.ch
	return v, ok
}

// IterFunc returns an iterator over the given function.
func IterFunc[T any](fn func() (T, bool)) Iter[T] {
	return &iterBase[T]{iterable: &iterFunc[T]{fn: fn}}
}

type iterFunc[T any] struct {
	fn func() (T, bool)
}

func (i *iterFunc[T]) Next() (T, bool) {
	return i.fn()
}
