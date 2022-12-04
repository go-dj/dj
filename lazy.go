package dj

import "sync"

// Lazy is a lazy value. It is initialized once via fn and then cached.
// It is safe for concurrent use.
func Lazy[T any](fn func() T) func() T {
	var (
		once sync.Once
		lazy T
	)

	return func() T {
		once.Do(func() {
			lazy = fn()
		})

		return lazy
	}
}
