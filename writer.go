package dj

// Writer is a type that can be written to by repeatedly calling Write.
type Writer[T any] interface {
	Writable[T]

	// WriteFrom writes all the values from the given Readable to the writer.
	WriteFrom(Readable[T]) (int, error)
}

// NewWriter returns a new Writer that writes to the given Writable.
func NewWriter[T any](w Writable[T]) Writer[T] {
	return &writer[T]{Writable: w}
}
