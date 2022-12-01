package xn

// iterable is a type that can be iterated over by repeatedly calling Next.
type iterable[T any] interface {
	Next() (T, bool)
}

// iterBase is a base type for implementing Iter.
// It provides implementations of For and Collect,
// and wraps an underlying iterable.
type iterBase[T any] struct {
	iterable[T]
}

func (base iterBase[T]) For(fn func(T)) {
	base.ForIdx(func(_ int, v T) {
		fn(v)
	})
}

func (base iterBase[T]) ForIdx(fn func(int, T)) {
	for idx := 0; ; idx++ {
		v, ok := base.Next()
		if !ok {
			return
		}

		fn(idx, v)
	}
}

func (base iterBase[T]) Read(n int) []T {
	out := make([]T, 0, n)

	for i := 0; i < n; i++ {
		v, ok := base.Next()
		if !ok {
			break
		}

		out = append(out, v)
	}

	return out
}

func (base iterBase[T]) Collect() []T {
	out := make([]T, 0)

	for v, ok := base.Next(); ok; v, ok = base.Next() {
		out = append(out, v)
	}

	return out
}

func (base iterBase[T]) Chan() <-chan T {
	ch := make(chan T)

	go func() {
		defer close(ch)

		for v, ok := base.Next(); ok; v, ok = base.Next() {
			ch <- v
		}
	}()

	return ch
}
