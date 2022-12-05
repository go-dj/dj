package dj

import (
	"context"
	"sync"
)

// ForN calls the given function for each index in the given range.
func ForN(n int, fn func(int)) {
	_ = ForNErr(n, func(i int) error {
		fn(i)
		return nil
	})
}

// ForNErr calls the given function for each index in the given range.
func ForNErr(n int, fn func(int) error) error {
	for i := 0; i < n; i++ {
		if err := fn(i); err != nil {
			return err
		}
	}

	return nil
}

// GoForN calls the given function with n unique values in parallel.
func GoForN(ctx context.Context, n int, fn func(context.Context, int)) {
	_ = GoForNErr(ctx, n, func(ctx context.Context, i int) error {
		fn(ctx, i)
		return nil
	})
}

// GoForNErr calls the given function with n unique values in parallel.
func GoForNErr(ctx context.Context, n int, fn func(context.Context, int) error) error {
	group := NewGroup(ctx, NewSem(parallelismFromCtx(ctx)))
	defer group.Wait()

	var err error
	var mu sync.Mutex

	group.Go(n, func(ctx context.Context, n int) {
		if fnErr := fn(ctx, n); fnErr != nil {
			mu.Lock()
			defer mu.Unlock()

			group.Cancel()

			if err == nil {
				err = fnErr
			}
		}
	}).Wait()

	return err
}

// MapN returns a slice of the results of the given function applied to each index in the given range.
func MapN[T any](n int, fn func(int) T) []T {
	out, _ := MapNErr(n, func(i int) (T, error) {
		return fn(i), nil
	})

	return out
}

// MapNErr returns a slice of the results of the given function applied to each index in the given range.
func MapNErr[T any](n int, fn func(int) (T, error)) ([]T, error) {
	out := make([]T, n)

	if err := ForNErr(n, func(i int) error {
		v, err := fn(i)
		if err != nil {
			return err
		}

		out[i] = v

		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

// GoMapN returns a slice of the results of the given function applied in parallel to each index in the given range.
func GoMapN[T any](ctx context.Context, n int, fn func(context.Context, int) T) []T {
	out, _ := GoMapNErr(ctx, n, func(ctx context.Context, i int) (T, error) {
		return fn(ctx, i), nil
	})

	return out
}

// GoMapNErr returns a slice of the results of the given function applied in parallel to each index in the given range.
func GoMapNErr[T any](ctx context.Context, n int, fn func(context.Context, int) (T, error)) ([]T, error) {
	out := make([]T, n)

	if err := GoForNErr(ctx, n, func(ctx context.Context, i int) error {
		v, err := fn(ctx, i)
		if err != nil {
			return err
		}

		out[i] = v

		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

// ForEach calls the given function for each element in the given slice.
func ForEach[T any](slice []T, fn func(T)) {
	ForEachIdx(slice, func(_ int, v T) {
		fn(v)
	})
}

// ForEachErr calls the given function for each element in the given slice.
func ForEachErr[T any](slice []T, fn func(T) error) error {
	return ForEachIdxErr(slice, func(_ int, v T) error {
		return fn(v)
	})
}

// ForEachIdx calls the given function for each index of the given slice.
func ForEachIdx[T any](slice []T, fn func(int, T)) {
	ForN(len(slice), func(i int) {
		fn(i, slice[i])
	})
}

// ForEachIdxErr calls the given function for each index of the given slice.
func ForEachIdxErr[T any](slice []T, fn func(int, T) error) error {
	return ForNErr(len(slice), func(i int) error {
		return fn(i, slice[i])
	})
}

// GoForEach calls the given function with each value in the given slice in parallel.
func GoForEach[T any](ctx context.Context, slice []T, fn func(context.Context, T)) {
	GoForEachIdx(ctx, slice, func(ctx context.Context, _ int, v T) {
		fn(ctx, v)
	})
}

// GoForEachErr calls the given function with each value in the given slice in parallel.
func GoForEachErr[T any](ctx context.Context, slice []T, fn func(context.Context, T) error) error {
	return GoForEachIdxErr(ctx, slice, func(ctx context.Context, _ int, v T) error {
		return fn(ctx, v)
	})
}

// GoForEachIdx calls the given function with each value in the given slice in parallel.
func GoForEachIdx[T any](ctx context.Context, slice []T, fn func(context.Context, int, T)) {
	GoForN(ctx, len(slice), func(ctx context.Context, i int) {
		fn(ctx, i, slice[i])
	})
}

// GoForEachIdxErr calls the given function with each value in the given slice in parallel.
func GoForEachIdxErr[T any](ctx context.Context, slice []T, fn func(context.Context, int, T) error) error {
	return GoForNErr(ctx, len(slice), func(ctx context.Context, i int) error {
		return fn(ctx, i, slice[i])
	})
}

// MapEach returns a slice of the results of the given function applied to each element in the given slice.
func MapEach[T, U any](slice []T, fn func(T) U) []U {
	return MapEachIdx(slice, func(_ int, v T) U {
		return fn(v)
	})
}

// MapEachErr returns a slice of the results of the given function applied to each element in the given slice.
func MapEachErr[T, U any](slice []T, fn func(T) (U, error)) ([]U, error) {
	return MapEachIdxErr(slice, func(_ int, v T) (U, error) {
		return fn(v)
	})
}

// MapEachIdx returns a slice of the results of the given function applied to each index of the given slice.
func MapEachIdx[T, U any](slice []T, fn func(int, T) U) []U {
	return MapN(len(slice), func(i int) U {
		return fn(i, slice[i])
	})
}

// MapEachIdxErr returns a slice of the results of the given function applied to each index of the given slice.
func MapEachIdxErr[T, U any](slice []T, fn func(int, T) (U, error)) ([]U, error) {
	return MapNErr(len(slice), func(i int) (U, error) {
		return fn(i, slice[i])
	})
}

// GoMapEach returns a slice of the results of the given function applied to each element in the given slice
// in parallel.
func GoMapEach[T, U any](ctx context.Context, slice []T, fn func(context.Context, T) U) []U {
	return GoMapEachIdx(ctx, slice, func(ctx context.Context, _ int, v T) U {
		return fn(ctx, v)
	})
}

// GoMapEachErr returns a slice of the results of the given function applied to each element in the given slice
// in parallel.
func GoMapEachErr[T, U any](ctx context.Context, slice []T, fn func(context.Context, T) (U, error)) ([]U, error) {
	return GoMapEachIdxErr(ctx, slice, func(ctx context.Context, _ int, v T) (U, error) {
		return fn(ctx, v)
	})
}

// GoMapEachIdx returns a slice of the results of the given function applied to each index of the given slice
// in parallel.
func GoMapEachIdx[T, U any](ctx context.Context, slice []T, fn func(context.Context, int, T) U) []U {
	return GoMapN(ctx, len(slice), func(ctx context.Context, i int) U {
		return fn(ctx, i, slice[i])
	})
}

// GoMapEachIdxErr returns a slice of the results of the given function applied to each index of the given slice
// in parallel.
func GoMapEachIdxErr[T, U any](ctx context.Context, slice []T, fn func(context.Context, int, T) (U, error)) ([]U, error) {
	return GoMapNErr(ctx, len(slice), func(ctx context.Context, i int) (U, error) {
		return fn(ctx, i, slice[i])
	})
}

// ForWindow calls the given function for each window of the given size in the given slice.
func ForWindow[T any](slice []T, size int, fn func([]T)) {
	ForWindowIdx(slice, size, func(_ int, window []T) {
		fn(window)
	})
}

// ForWindowErr calls the given function for each window of the given size in the given slice.
func ForWindowErr[T any](slice []T, size int, fn func([]T) error) error {
	return ForWindowIdxErr(slice, size, func(_ int, window []T) error {
		return fn(window)
	})
}

// ForWindowIdx calls the given function for each window of the given size in the given slice.
func ForWindowIdx[T any](slice []T, size int, fn func(int, []T)) {
	ForN(len(slice)-size+1, func(idx int) {
		fn(idx, slice[idx:idx+size])
	})
}

// ForWindowIdxErr calls the given function for each window of the given size in the given slice.
func ForWindowIdxErr[T any](slice []T, size int, fn func(int, []T) error) error {
	return ForNErr(len(slice)-size+1, func(idx int) error {
		return fn(idx, slice[idx:idx+size])
	})
}

// GoForWindow calls the given function for each window of the given size in the given slice in parallel.
func GoForWindow[T any](ctx context.Context, slice []T, size int, fn func(context.Context, []T)) {
	GoForWindowIdx(ctx, slice, size, func(ctx context.Context, _ int, window []T) {
		fn(ctx, window)
	})
}

// GoForWindowErr calls the given function for each window of the given size in the given slice in parallel.
func GoForWindowErr[T any](ctx context.Context, slice []T, size int, fn func(context.Context, []T) error) error {
	return GoForWindowIdxErr(ctx, slice, size, func(ctx context.Context, _ int, window []T) error {
		return fn(ctx, window)
	})
}

// GoForWindowIdx calls the given function for each window of the given size in the given slice in parallel.
func GoForWindowIdx[T any](ctx context.Context, slice []T, size int, fn func(context.Context, int, []T)) {
	GoForN(ctx, len(slice)-size+1, func(ctx context.Context, idx int) {
		fn(ctx, idx, slice[idx:idx+size])
	})
}

// GoForWindowIdxErr calls the given function for each window of the given size in the given slice in parallel.
func GoForWindowIdxErr[T any](ctx context.Context, slice []T, size int, fn func(context.Context, int, []T) error) error {
	return GoForNErr(ctx, len(slice)-size+1, func(ctx context.Context, idx int) error {
		return fn(ctx, idx, slice[idx:idx+size])
	})
}

// MapWindow returns a slice of the results of the given function applied to each window of the given size
// in the given slice.
func MapWindow[T, U any](slice []T, size int, fn func([]T) U) []U {
	return MapWindowIdx(slice, size, func(_ int, window []T) U {
		return fn(window)
	})
}

// MapWindowErr returns a slice of the results of the given function applied to each window of the given size
// in the given slice.
func MapWindowErr[T, U any](slice []T, size int, fn func([]T) (U, error)) ([]U, error) {
	return MapWindowIdxErr(slice, size, func(_ int, window []T) (U, error) {
		return fn(window)
	})
}

// MapWindowIdx returns a slice of the results of the given function applied to each window of the given size
// in the given slice.
func MapWindowIdx[T, U any](slice []T, size int, fn func(int, []T) U) []U {
	return MapN(len(slice)-size+1, func(idx int) U {
		return fn(idx, slice[idx:idx+size])
	})
}

// MapWindowIdxErr returns a slice of the results of the given function applied to each window of the given size
// in the given slice.
func MapWindowIdxErr[T, U any](slice []T, size int, fn func(int, []T) (U, error)) ([]U, error) {
	return MapNErr(len(slice)-size+1, func(idx int) (U, error) {
		return fn(idx, slice[idx:idx+size])
	})
}

// GoMapWindow returns a slice of the results of the given function applied to each window of the given size
// in the given slice in parallel.
func GoMapWindow[T, U any](ctx context.Context, slice []T, size int, fn func(context.Context, []T) U) []U {
	return GoMapWindowIdx(ctx, slice, size, func(ctx context.Context, _ int, window []T) U {
		return fn(ctx, window)
	})
}

// GoMapWindowErr returns a slice of the results of the given function applied to each window of the given size
// in the given slice in parallel.
func GoMapWindowErr[T, U any](ctx context.Context, slice []T, size int, fn func(context.Context, []T) (U, error)) ([]U, error) {
	return GoMapWindowIdxErr(ctx, slice, size, func(ctx context.Context, _ int, window []T) (U, error) {
		return fn(ctx, window)
	})
}

// GoMapWindowIdx returns a slice of the results of the given function applied to each window of the given size
// in the given slice in parallel.
func GoMapWindowIdx[T, U any](ctx context.Context, slice []T, size int, fn func(context.Context, int, []T) U) []U {
	return GoMapN(ctx, len(slice)-size+1, func(ctx context.Context, idx int) U {
		return fn(ctx, idx, slice[idx:idx+size])
	})
}

// GoMapWindowIdxErr returns a slice of the results of the given function applied to each window of the given size
// in the given slice in parallel.
func GoMapWindowIdxErr[T, U any](ctx context.Context, slice []T, size int, fn func(context.Context, int, []T) (U, error)) ([]U, error) {
	return GoMapNErr(ctx, len(slice)-size+1, func(ctx context.Context, idx int) (U, error) {
		return fn(ctx, idx, slice[idx:idx+size])
	})
}

// RangeN returns a slice of integers from 0 to n-1.
func RangeN(n int) []int {
	return MapN(n, func(i int) int {
		return i
	})
}

// Range returns a slice of integers from start to end-1.
func Range(start, end int) []int {
	return MapN(end-start, func(i int) int {
		return start + i
	})
}
