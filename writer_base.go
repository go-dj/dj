package xn

// Writable is a type that can be written to by repeatedly calling Write.
type Writable[T any] interface {
	Write(T) bool
}

// writer is a base type for implementing Writer which wraps a writable.
type writer[T any] struct {
	Writable[T]
}

func (w writer[T]) WriteFrom(r Iterable[T]) (int, bool) {
	var count int

	for v, ok := r.Next(); ok; v, ok = r.Next() {
		if !w.Write(v) {
			return count, false
		}

		count++
	}

	return count, true
}
