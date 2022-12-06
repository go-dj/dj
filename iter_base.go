package dj

// Readable is a type that can be iterated over by repeatedly calling Next.
// The returned value is the next value in the iterator, or an error if one occurred while reading.
// The boolean return value indicates whether the iterator has been exhausted.
type Readable[T any] interface {
	Read() (Result[T], bool)
}

// iter is a base type for implementing Iter which wraps a readable.
type iter[T any] struct {
	Readable[T]
}

func (i iter[T]) Take(n int) ([]T, error) {
	out := make([]T, 0, n)

	for {
		v, ok := i.Read()
		if !ok {
			return out, v.Err()
		}

		if out = append(out, v.Val()); len(out) == n {
			return out, nil
		}
	}
}

func (i iter[T]) Collect() ([]T, error) {
	out := make([]T, 0)

	for {
		v, ok := i.Read()
		if !ok {
			return out, v.Err()
		}

		out = append(out, v.Val())
	}
}

func (i iter[T]) Send(ch chan<- T) (int, error) {
	var n int

	for {
		v, ok := i.Read()
		if !ok {
			return n, v.Err()
		}

		ch <- v.Val()

		n++
	}
}

func (i iter[T]) Recv() <-chan T {
	ch := make(chan T)

	go func() {
		defer close(ch)

		for v, ok := i.Read(); ok; v, ok = i.Read() {
			ch <- v.Val()
		}
	}()

	return ch
}

func (i iter[T]) WriteTo(w Writable[T]) (int, error) {
	var n int

	for {
		v, ok := i.Read()
		if !ok {
			return n, v.Err()
		}

		if err := w.Write(v.Val()); err != nil {
			return n, err
		}

		n++
	}
}
