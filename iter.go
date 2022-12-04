package dj

// Iter is a type that can be read sequentially.
type Iter[T any] interface {
	Readable[T]

	// For calls the given function for each value in the iterator.
	For(func(T))

	// ForIdx calls the given function for each value in the iterator, along with the index of the value.
	ForIdx(func(int, T))

	// Take reads up to n values from the iterator and returns them.
	Take(int) []T

	// Collect returns a slice containing all the values in the iterator.
	Collect() []T

	// Send sends all the values in the iterator to the given channel.
	Send(ch chan<- T)

	// Recv returns a channel that will receive all the values in the iterator.
	Recv() <-chan T

	// WriteTo writes all the values in the iterator to the given writable.
	WriteTo(Writable[T]) (int, bool)
}

// NewIter returns a new Iter that reads from the given Readable.
func NewIter[T any](r Readable[T]) Iter[T] {
	return &iter[T]{Readable: r}
}

// Peeker is a peekable iterator.
type Peeker[T any] interface {
	Iter[T]

	// Peek returns the next value in the iterator without advancing it.
	Peek() (T, bool)
}

// WithPeek wraps the given iterator with a peekable iterator.
func WithPeek[T any](iter Iter[T]) Peeker[T] {
	return &peekIter[T]{Iter: iter}
}

type peekIter[T any] struct {
	Iter[T]

	next T
	has  bool
}

func (i *peekIter[T]) Peek() (T, bool) {
	if !i.has {
		v, ok := i.Read()
		if !ok {
			return Zero[T](), false
		}

		i.next = v
		i.has = true
	}

	return i.next, true
}

func (i *peekIter[T]) Read() (T, bool) {
	if !i.has {
		return i.Iter.Read()
	}

	i.has = false

	return i.next, true
}

// MapIter applies the given function to each value returned by the given iterator.
func MapIter[T, U any](iter Iter[T], fn func(T) U) Iter[U] {
	return NewIter[U](&mapIter[T, U]{iter: iter, fn: fn})
}

type mapIter[T, U any] struct {
	iter Iter[T]
	fn   func(T) U
}

func (i *mapIter[T, U]) Read() (U, bool) {
	v, ok := i.iter.Read()
	if !ok {
		return Zero[U](), false
	}

	return i.fn(v), true
}

// ChunkIter returns an iterator over the given iterator, chunking the values into slices of the given size.
func ChunkIter[T any](iter Iter[T], size int) Iter[[]T] {
	return NewIter[[]T](&chunkIter[T]{iter: iter, size: size})
}

type chunkIter[T any] struct {
	iter Iter[T]
	size int
}

func (i *chunkIter[T]) Read() ([]T, bool) {
	out := make([]T, 0, i.size)

	for len(out) < i.size {
		v, ok := i.iter.Read()
		if !ok {
			break
		}

		out = append(out, v)
	}

	return out, len(out) > 0
}

// FilterIter returns an iterator over the given iterator, filtering out values that do not match the given predicate.
func FilterIter[T any](iter Iter[T], fn func(T) bool) Iter[T] {
	return NewIter[T](&filterIter[T]{iter: iter, fn: fn})
}

type filterIter[T any] struct {
	iter Iter[T]
	fn   func(T) bool
}

func (i *filterIter[T]) Read() (T, bool) {
	for {
		v, ok := i.iter.Read()
		if !ok {
			return Zero[T](), false
		}

		if i.fn(v) {
			return v, true
		}
	}
}

// FlattenIter returns an iterator over the given iterator, flattening nested iterators.
// That is, it converts an iterator over iterators into an iterator over the values of those iterators.
func FlattenIter[T any](iter Iter[Iter[T]]) Iter[T] {
	return NewIter[T](&flattenIter[T]{iter: iter})
}

type flattenIter[T any] struct {
	iter Iter[Iter[T]]
	curr Iter[T]
}

func (i *flattenIter[T]) Read() (T, bool) {
	for {
		if i.curr == nil {
			next, ok := i.iter.Read()
			if !ok {
				return Zero[T](), false
			}

			i.curr = next
		}

		v, ok := i.curr.Read()
		if ok {
			return v, true
		}

		i.curr = nil
	}
}

// JoinIter concatenates the given iterators into a single iterator.
func JoinIter[T any](iters ...Iter[T]) Iter[T] {
	return FlattenIter(SliceIter(iters...))
}

// ZipIter returns an iterator over the given iterators, zipping the values together.
// That is, it converts an iterator over iterators into an iterator over tuples of the values of those iterators.
func ZipIter[T any](iters ...Iter[T]) Iter[[]T] {
	return NewIter[[]T](&zipIter[T]{iters: iters})
}

type zipIter[T any] struct {
	iters []Iter[T]
}

func (i *zipIter[T]) Read() ([]T, bool) {
	if len(i.iters) == 0 {
		return nil, false
	}

	out := make([]T, len(i.iters))

	for j, iter := range i.iters {
		v, ok := iter.Read()
		if !ok {
			return nil, false
		}

		out[j] = v
	}

	return out, true
}
