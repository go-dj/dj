package dj

type RWBuffer[T any] struct {
	items []T
}

// NewRWBuffer returns a new RWBuffer.
func NewRWBuffer[T any](items []T) *RWBuffer[T] {
	return &RWBuffer[T]{items: items}
}

// Write writes the given value to the buffer.
func (b *RWBuffer[T]) Write(v T) error {
	b.items = append(b.items, v)

	return nil
}

// Next returns the next value from the buffer.
func (b *RWBuffer[T]) Read() (Result[T], bool) {
	if len(b.items) == 0 {
		return Ok(Zero[T]()), false
	}

	v := b.items[0]

	b.items = b.items[1:]

	return Ok(v), true
}
