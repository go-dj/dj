package dj

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
