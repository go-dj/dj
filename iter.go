package xn

// Iter is a type that can be iterated over.
type Iter[T any] interface {
	iterable[T]

	// For calls the given function for each value in the iterator.
	For(func(T))

	// ForIdx calls the given function for each value in the iterator, along with the index of the value.
	ForIdx(func(int, T))

	// Read reads up to n values from the iterator.
	Read(int) []T

	// Collect returns a slice containing all the values in the iterator.
	Collect() []T

	// Chan returns a channel that will receive all the values in the iterator.
	Chan() <-chan T
}

// PeekIter is a peekable iterator.
type PeekIter[T any] interface {
	Iter[T]

	// Peek returns the next value in the iterator without advancing it.
	Peek() (T, bool)
}

// WithPeek wraps the given iterator with a peekable iterator.
func WithPeek[T any](iter Iter[T]) PeekIter[T] {
	return &peekIter[T]{Iter: iter}
}

type peekIter[T any] struct {
	Iter[T]

	next T
	has  bool
}

func (i *peekIter[T]) Peek() (T, bool) {
	if !i.has {
		v, ok := i.Next()
		if !ok {
			return zero[T](), false
		}

		i.next = v
		i.has = true
	}

	return i.next, true
}

func (i *peekIter[T]) Next() (T, bool) {
	if !i.has {
		return i.Iter.Next()
	}

	i.has = false

	return i.next, true
}

// MapIter applies the given function to each value returned by the given iterator.
func MapIter[T, U any](iter Iter[T], fn func(T) U) Iter[U] {
	return &iterBase[U]{iterable: &mapIter[T, U]{iter: iter, fn: fn}}
}

type mapIter[T, U any] struct {
	iter Iter[T]
	fn   func(T) U
}

func (i *mapIter[T, U]) Next() (U, bool) {
	v, ok := i.iter.Next()
	if !ok {
		return zero[U](), false
	}

	return i.fn(v), true
}

// ChunkIter returns an iterator over the given iterator, chunking the values into slices of the given size.
func ChunkIter[T any](iter Iter[T], size int) Iter[[]T] {
	return &iterBase[[]T]{iterable: &chunkIter[T]{iter: iter, size: size}}
}

type chunkIter[T any] struct {
	iter Iter[T]
	size int
}

func (i *chunkIter[T]) Next() ([]T, bool) {
	out := make([]T, 0, i.size)

	for len(out) < i.size {
		v, ok := i.iter.Next()
		if !ok {
			break
		}

		out = append(out, v)
	}

	return out, len(out) > 0
}

// FilterIter returns an iterator over the given iterator, filtering out values that do not match the given predicate.
func FilterIter[T any](iter Iter[T], fn func(T) bool) Iter[T] {
	return &iterBase[T]{iterable: &filterIter[T]{iter: iter, fn: fn}}
}

type filterIter[T any] struct {
	iter Iter[T]
	fn   func(T) bool
}

func (i *filterIter[T]) Next() (T, bool) {
	for {
		v, ok := i.iter.Next()
		if !ok {
			return zero[T](), false
		}

		if i.fn(v) {
			return v, true
		}
	}
}

// FlattenIter returns an iterator over the given iterator, flattening nested iterators.
// That is, it converts an iterator over iterators into an iterator over the values of those iterators.
func FlattenIter[T any](iter Iter[Iter[T]]) Iter[T] {
	return &iterBase[T]{iterable: &flattenIter[T]{iter: iter}}
}

type flattenIter[T any] struct {
	iter Iter[Iter[T]]
	curr Iter[T]
}

func (i *flattenIter[T]) Next() (T, bool) {
	for {
		if i.curr == nil {
			next, ok := i.iter.Next()
			if !ok {
				return zero[T](), false
			}

			i.curr = next
		}

		v, ok := i.curr.Next()
		if ok {
			return v, true
		}

		i.curr = nil
	}
}

// JoinIter concatenates the given iterators into a single iterator.
func JoinIter[T any](iters ...Iter[T]) Iter[T] {
	return FlattenIter(IterSlice(iters...))
}
