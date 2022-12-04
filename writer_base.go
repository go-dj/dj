package dj

// Writable is a type that can be written to by repeatedly calling Write.
type Writable[T any] interface {
	Write(T) bool
}

// writer is a base type for implementing Writer which wraps a writable.
type writer[T any] struct {
	Writable[T]
}

func (w writer[T]) WriteFrom(r Readable[T]) (int, bool) {
	var count int

	for v, ok := r.Read(); ok; v, ok = r.Read() {
		if !w.Write(v) {
			return count, false
		}

		count++
	}

	return count, true
}
