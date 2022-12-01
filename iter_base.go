package xn

// Iterable is a type that can be iterated over by repeatedly calling Next.
type Iterable[T any] interface {
	Next() (T, bool)
}

// Writable is a type that can be written to by repeatedly calling Write.
type Writable[T any] interface {
	Write(T) bool
}

// iter is a base type for implementing Iter which wraps a readable.
type iter[T any] struct {
	Iterable[T]
}

func (i iter[T]) For(fn func(T)) {
	i.ForIdx(func(_ int, v T) {
		fn(v)
	})
}

func (i iter[T]) ForIdx(fn func(int, T)) {
	for idx := 0; ; idx++ {
		v, ok := i.Next()
		if !ok {
			return
		}

		fn(idx, v)
	}
}

func (i iter[T]) Collect() []T {
	out := make([]T, 0)

	for v, ok := i.Next(); ok; v, ok = i.Next() {
		out = append(out, v)
	}

	return out
}

func (i iter[T]) Chan() <-chan T {
	ch := make(chan T)

	go func() {
		defer close(ch)

		for v, ok := i.Next(); ok; v, ok = i.Next() {
			ch <- v
		}
	}()

	return ch
}

func (i iter[T]) WriteTo(w Writable[T]) (int, bool) {
	var n int

	for v, ok := i.Next(); ok; v, ok = i.Next() {
		if !w.Write(v) {
			return n, false
		}

		n++
	}

	return n, true
}
