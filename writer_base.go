package dj

// Writable is a type that can be written to by repeatedly calling Write.
type Writable[T any] interface {
	Write(T) error
}

// writer is a base type for implementing Writer which wraps a writable.
type writer[T any] struct {
	Writable[T]
}

func (w writer[T]) WriteFrom(r Readable[T]) (int, error) {
	var n int

	for v, ok := r.Read(); ok; v, ok = r.Read() {
		if err := v.Err(); err != nil {
			return n, err
		} else if err := w.Write(v.Val()); err != nil {
			return n, err
		}

		n++
	}

	return n, nil
}
