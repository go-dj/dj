package dj

import "sync"

// Future is a future value.
// Its value is not available until the future is resolved.
type Future[T any] struct {
	value T
	done  chan struct{}
}

// New creates a new Future.
func New[T any]() *Future[T] {
	return &Future[T]{done: make(chan struct{})}
}

// Set sets the value of the future.
// Set can only be called once; subsequent calls will panic.
func (f *Future[T]) Set(v T) {
	f.value = v
	close(f.done)
}

// Get returns the value of the future. It blocks until the value has been set.
func (f *Future[T]) Get() T {
	<-f.done
	return f.value
}

// Singular obtains a value exactly once.
// Any subsequent attempts to set the value will be ignored.
type Singular[T any] struct {
	once sync.Once
	val  T
}

// Set sets the value if it has not already been set.
func (s *Singular[T]) Set(val T) {
	s.once.Do(func() {
		s.val = val
	})
}

// Get returns the value if it has been set.
func (s *Singular[T]) Get() T {
	return s.val
}
