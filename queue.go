package xn

type Queue[T any] interface {
	// Enqueue adds the given value to the queue.
	Enqueue(v T)

	// Dequeue removes and returns the next value from the queue.
	Dequeue() (T, bool)

	// Peek returns the next value in the queue without removing it.
	Peek() (T, bool)

	// Len returns the number of values in the queue.
	Len() int
}

// NewQueue returns a new queue.
func NewQueue[T any]() Queue[T] {
	return &queue[T]{}
}

type queue[T any] struct {
	values []T
}

// Enqueue adds the given value to the queue.
func (q *queue[T]) Enqueue(v T) {
	panic("not implemented") // TODO: Implement
}

// Dequeue removes and returns the next value from the queue.
func (q *queue[T]) Dequeue() (T, bool) {
	panic("not implemented") // TODO: Implement
}

// Peek returns the next value in the queue without removing it.
func (q *queue[T]) Peek() (T, bool) {
	panic("not implemented") // TODO: Implement
}

// Len returns the number of values in the queue.
func (q *queue[T]) Len() int {
	panic("not implemented") // TODO: Implement
}
