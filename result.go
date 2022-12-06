package dj

import "fmt"

// Result is a type that represents a value or an error.
type Result[T any] struct {
	val T
	err error
}

// Ok returns a new Result with the given value and no error.
func Ok[T any](val T) Result[T] {
	return Result[T]{val, nil}
}

// Err returns a new Result with the given error and no value.
func Err[T any](err error) Result[T] {
	return Result[T]{Zero[T](), err}
}

// NewResult returns a new Result with the given value and error.
func NewResult[T any](val T, err error) Result[T] {
	return Result[T]{val, err}
}

// Val returns the value of the Result.
// If the Result is an error, the value is undefined.
func (r Result[T]) Val() T {
	return r.val
}

// Err returns the error of the Result.
// If the Result is a value, the error is nil.
func (r Result[T]) Err() error {
	return r.err
}

// Unpack returns the value and error of the Result.
func (r Result[T]) Unpack() (T, error) {
	return r.val, r.err
}

// String returns a string representation of the Result.
func (r Result[T]) String() string {
	if r.err != nil {
		return fmt.Sprintf("Err(%v)", r.err)
	}

	return fmt.Sprintf("Ok(%v)", r.val)
}
