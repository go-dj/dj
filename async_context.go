package dj

import (
	"context"
	"runtime"
)

type withParallelismType struct{}

var withParallelismTypeKey withParallelismType

// WithParallelism returns a context with the given parallelism.
// Certain functions in this package will use this value to determine
// how many goroutines to run concurrently.
// If the context already has a parallelism value, it will be overwritten.
// If not set, the default parallelism is the number of CPUs.
func WithParallelism(ctx context.Context, n int) context.Context {
	return context.WithValue(ctx, withParallelismTypeKey, n)
}

// parallelismFromCtx returns the parallelism value from the given context.
func parallelismFromCtx(ctx context.Context) int {
	if n, ok := ctx.Value(withParallelismTypeKey).(int); ok {
		return n
	}

	return runtime.NumCPU()
}
