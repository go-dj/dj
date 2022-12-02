package dj

// Writer is a type that can be written to by repeatedly calling Write.
type Writer[T any] interface {
	Writable[T]

	// WriteFrom writes all the values from the given Iterable to the writer.
	WriteFrom(Iterable[T]) (int, bool)
}

// NewWriter returns a new Writer that writes to the given Writable.
func NewWriter[T any](w Writable[T]) Writer[T] {
	return &writer[T]{Writable: w}
}
