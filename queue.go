package dj

import "sync"

type Queue[T any] struct {
	items []T
	mu    sync.Mutex
}

func NewQueue[T any](items ...T) *Queue[T] {
	return &Queue[T]{items: items}
}

func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.items)
}

func (q *Queue[T]) Push(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, v)
}

func (q *Queue[T]) PopFront() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return Zero[T](), false
	}

	v := q.items[0]

	q.items = q.items[1:]

	return v, true
}

func (q *Queue[T]) PopBack() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return Zero[T](), false
	}

	v := q.items[len(q.items)-1]

	q.items = q.items[:len(q.items)-1]

	return v, true
}
