package dj

type Result[T any] struct {
	val T
	err error
}

func Ok[T any](val T) Result[T] {
	return Result[T]{val, nil}
}

func Err[T any](err error) Result[T] {
	return Result[T]{Zero[T](), err}
}

func NewResult[T any](val T, err error) Result[T] {
	return Result[T]{val, err}
}

func (r Result[T]) Value() T {
	return r.val
}

func (r Result[T]) Error() error {
	return r.err
}
