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

func (q *Queue[T]) Items() []T {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.items
}

func (q *Queue[T]) PushFront(items ...T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(items, q.items...)
}

func (q *Queue[T]) PushBack(v ...T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, v...)
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

type Map[Key comparable, Value any] struct {
	items map[Key]Value
	mu    sync.Mutex
}

func NewMap[Key comparable, Value any]() *Map[Key, Value] {
	return &Map[Key, Value]{items: make(map[Key]Value)}
}

func (m *Map[Key, Value]) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.items)
}

func (m *Map[Key, Value]) Items() map[Key]Value {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.items
}

func (m *Map[Key, Value]) Set(key Key, value Value) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = value
}

func (m *Map[Key, Value]) Get(key Key) (Value, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.items[key]

	return v, ok
}

func (m *Map[Key, Value]) Delete(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
}
