package xn

import (
	"context"
	"math/rand"
	"runtime"
	"sort"

	"golang.org/x/exp/constraints"
)

// MaxParallel sets the default maximum number of goroutines that can be run in parallel.
var MaxParallel = runtime.NumCPU()

// ForN calls the given function for each index in the given range.
func ForN(n int, fn func(int)) {
	for i := 0; i < n; i++ {
		fn(i)
	}
}

// GoForN calls the given function with n unique values in parallel.
func GoForN(ctx context.Context, n int, fn func(int)) {
	sem := NewSem(ctx, MaxParallel)
	defer sem.Wait()

	ForN(n, func(idx int) {
		sem.Go(func(context.Context) {
			fn(idx)
		})
	})
}

// MapN returns a slice of the results of the given function applied to each index in the given range.
func MapN[T any](n int, fn func(int) T) []T {
	out := make([]T, n)

	ForN(n, func(i int) {
		out[i] = fn(i)
	})

	return out
}

// GoMapN returns a slice of the results of the given function applied in parallel to each index in the given range.
func GoMapN[T any](ctx context.Context, n int, fn func(int) T) []T {
	out := make([]T, n)

	GoForN(ctx, n, func(i int) {
		out[i] = fn(i)
	})

	return out
}

// ForEach calls the given function for each element in the given slice.
func ForEach[T any](slice []T, fn func(T)) {
	ForEachIdx(slice, func(_ int, v T) {
		fn(v)
	})
}

// ForEachIdx calls the given function for each index of the given slice.
func ForEachIdx[T any](slice []T, fn func(int, T)) {
	ForN(len(slice), func(i int) {
		fn(i, slice[i])
	})
}

// GoForEach calls the given function with each value in the given slice in parallel.
func GoForEach[T any](ctx context.Context, slice []T, fn func(T)) {
	GoForEachIdx(ctx, slice, func(_ int, v T) {
		fn(v)
	})
}

// GoForEachIdx calls the given function with each value in the given slice in parallel.
func GoForEachIdx[T any](ctx context.Context, slice []T, fn func(int, T)) {
	GoForN(ctx, len(slice), func(i int) {
		fn(i, slice[i])
	})
}

// MapEach returns a slice of the results of the given function applied to each element in the given slice.
func MapEach[T, U any](slice []T, fn func(T) U) []U {
	return MapEachIdx(slice, func(_ int, v T) U {
		return fn(v)
	})
}

// MapEachIdx returns a slice of the results of the given function applied to each index of the given slice.
func MapEachIdx[T, U any](slice []T, fn func(int, T) U) []U {
	return MapN(len(slice), func(i int) U {
		return fn(i, slice[i])
	})
}

// GoMapEach returns a slice of the results of the given function applied to each element in the given slice
// in parallel.
func GoMapEach[T, U any](ctx context.Context, slice []T, fn func(T) U) []U {
	return GoMapEachIdx(ctx, slice, func(_ int, v T) U {
		return fn(v)
	})
}

// GoMapEachIdx returns a slice of the results of the given function applied to each index of the given slice
// in parallel.
func GoMapEachIdx[T, U any](ctx context.Context, slice []T, fn func(int, T) U) []U {
	return GoMapN(ctx, len(slice), func(i int) U {
		return fn(i, slice[i])
	})
}

// ForWindow calls the given function for each window of the given size in the given slice.
func ForWindow[T any](slice []T, size int, fn func([]T)) {
	ForWindowIdx(slice, size, func(_ int, window []T) {
		fn(window)
	})
}

// ForWindowIdx calls the given function for each window of the given size in the given slice.
func ForWindowIdx[T any](slice []T, size int, fn func(int, []T)) {
	ForN(len(slice)-size+1, func(idx int) {
		fn(idx, slice[idx:idx+size])
	})
}

// GoForWindow calls the given function for each window of the given size in the given slice in parallel.
func GoForWindow[T any](ctx context.Context, slice []T, size int, fn func([]T)) {
	GoForWindowIdx(ctx, slice, size, func(_ int, window []T) {
		fn(window)
	})
}

// GoForWindowIdx calls the given function for each window of the given size in the given slice in parallel.
func GoForWindowIdx[T any](ctx context.Context, slice []T, size int, fn func(int, []T)) {
	GoForN(ctx, len(slice)-size+1, func(idx int) {
		fn(idx, slice[idx:idx+size])
	})
}

// MapWindow returns a slice of the results of the given function applied to each window of the given size
// in the given slice.
func MapWindow[T, U any](slice []T, size int, fn func([]T) U) []U {
	return MapWindowIdx(slice, size, func(_ int, window []T) U {
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

// GoMapWindow returns a slice of the results of the given function applied to each window of the given size
// in the given slice in parallel.
func GoMapWindow[T, U any](ctx context.Context, slice []T, size int, fn func([]T) U) []U {
	return GoMapWindowIdx(ctx, slice, size, func(_ int, window []T) U {
		return fn(window)
	})
}

// GoMapWindowIdx returns a slice of the results of the given function applied to each window of the given size
// in the given slice in parallel.
func GoMapWindowIdx[T, U any](ctx context.Context, slice []T, size int, fn func(int, []T) U) []U {
	return GoMapN(ctx, len(slice)-size+1, func(idx int) U {
		return fn(idx, slice[idx:idx+size])
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

// Any returns true if the given predicate is true for any element in the given slice.
func Any[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if fn(v) {
			return true
		}
	}

	return false
}

// All returns true if the given predicate is true for all elements in the given slice.
func All[T any](slice []T, fn func(T) bool) bool {
	return !Any(slice, func(v T) bool {
		return !fn(v)
	})
}

// None returns true if the given predicate is false for all elements in the given slice.
func None[T any](slice []T, fn func(T) bool) bool {
	return !Any(slice, fn)
}

// Same returns true if all the given values are equal.
func Same[T comparable](in ...T) bool {
	return SameFn(in, func(a, b T) bool {
		return a == b
	})
}

// SameFn returns true if all the given values are equal, according to the given function.
func SameFn[T any](in []T, fn func(T, T) bool) bool {
	for i := 1; i < len(in); i++ {
		if !fn(in[0], in[i]) {
			return false
		}
	}

	return true
}

// Equal returns true if the given slices are equal.
func Equal[T comparable](in ...[]T) bool {
	return EqualFn(in, func(a, b T) bool {
		return a == b
	})
}

// EqualFn returns true if the given slices are equal, according to the given function.
func EqualFn[T any](in [][]T, eq func(T, T) bool) bool {
	if !Same(MapEach(in, func(slice []T) int { return len(slice) })...) {
		return false
	}

	for _, tuple := range Zip(in...) {
		if !SameFn(tuple, eq) {
			return false
		}
	}

	return true
}

// Join returns a slice of the elements in the given slices.
func Join[T any](in ...[]T) []T {
	out := make([]T, 0, Sum(MapEach(in, func(slice []T) int { return len(slice) })))

	for _, slice := range in {
		out = append(out, slice...)
	}

	return out
}

// Zip returns a slice of tuples, where each tuple contains the elements at the same index in the given slices.
func Zip[T any](in ...[]T) [][]T {
	return MapN(len(in[0]), func(idx int) []T {
		return MapEach(in, func(slice []T) T {
			return slice[idx]
		})
	})
}

// Unzip returns a slice of slices, where each slice contains the elements at the same index in the given tuples.
func Unzip[T any](in [][]T) [][]T {
	return MapN(len(in[0]), func(idx int) []T {
		return MapEach(in, func(tuple []T) T {
			return tuple[idx]
		})
	})
}

// Chunk chunks the given slice into slices of the given size.
func Chunk[T any](slice []T, size int) [][]T {
	var buf int

	if len(slice)%size != 0 {
		buf = 1
	}

	out := make([][]T, len(slice)/size+buf)

	ForEachIdx(slice, func(i int, v T) {
		out[i/size] = append(out[i/size], v)
	})

	return out
}

// Reduce returns the result of applying the given function to each element in the given slice,
// starting with the given initial value.
func Reduce[T, U any](slice []T, init U, fn func(U, T) U) U {
	return ReduceIdx(slice, init, func(u U, idx int) U {
		return fn(u, slice[idx])
	})
}

// ReduceIdx returns the result of applying the given function to each index of the given slice,
// starting with the given initial value.
func ReduceIdx[T, U any](slice []T, init U, fn func(U, int) U) U {
	ForN(len(slice), func(idx int) {
		init = fn(init, idx)
	})

	return init
}

// Min returns the minimum value in the given slice.
func Min[T constraints.Ordered](in []T) T {
	return in[MinIdx(in)]
}

// MinFn returns the minimum value in the given slice, according to the given function.
func MinFn[T any](in []T, fn func(T, T) bool) T {
	return in[MinIdxFn(in, fn)]
}

// MinIdx returns the index of the minimum element in the given slice.
func MinIdx[T constraints.Ordered](in []T) int {
	return MinIdxFn(in, func(min, b T) bool {
		return min < b
	})
}

// MinIdxFn returns the index of the minimum element in the given slice, according to the given function.
func MinIdxFn[T any](in []T, fn func(T, T) bool) int {
	return ReduceIdx(in, 0, func(min int, idx int) int {
		if fn(in[idx], in[min]) {
			return idx
		}

		return min
	})
}

// Max returns the maximum value in the given slice.
func Max[T constraints.Ordered](in []T) T {
	return in[MaxIdx(in)]
}

// MaxFn returns the maximum value in the given slice, according to the given function.
func MaxFn[T any](in []T, fn func(T, T) bool) T {
	return in[MaxIdxFn(in, fn)]
}

// MaxIdx returns the index of the maximum element in the given slice.
func MaxIdx[T constraints.Ordered](in []T) int {
	return MaxIdxFn(in, func(max, b T) bool {
		return max > b
	})
}

// MaxIdxFn returns the index of the maximum element in the given slice, according to the given function.
func MaxIdxFn[T any](in []T, fn func(T, T) bool) int {
	return ReduceIdx(in, 0, func(max int, idx int) int {
		if fn(in[idx], in[max]) {
			return idx
		}

		return max
	})
}

// Sum returns the sum of the given slice.
func Sum[T constraints.Ordered](in []T) T {
	return Reduce(in, zero[T](), func(a, b T) T {
		return a + b
	})
}

// Count returns the number of elements in the given slice that equal the given value.
func Count[T comparable](in []T, val T) int {
	return CountFn(in, func(a T) bool {
		return a == val
	})
}

// CountFn returns the number of elements in the given slice that satisfy the given predicate.
func CountFn[T any](in []T, fn func(T) bool) int {
	return Reduce(in, 0, func(a int, b T) int {
		if fn(b) {
			return a + 1
		}

		return a
	})
}

// Filter returns a slice of the elements in the given slice that satisfy the given predicate.
func Filter[T any](slice []T, fn func(T) bool) []T {
	return FilterIdx(slice, func(_ int, v T) bool {
		return fn(v)
	})
}

// FilterIdx returns a slice of the elements in the given slice that satisfy the given predicate.
func FilterIdx[T any](slice []T, fn func(int, T) bool) []T {
	out := make([]T, 0, len(slice))

	ForEachIdx(slice, func(i int, v T) {
		if fn(i, v) {
			out = append(out, v)
		}
	})

	return out
}

// Uniq returns a slice of the unique elements in the given slice.
func Uniq[T comparable](in []T) []T {
	return UniqFn(in, func(a, b T) bool {
		return a == b
	})
}

// UniqFn returns a slice of the unique elements in the given slice, according to the given function.
func UniqFn[T any](in []T, eq func(T, T) bool) []T {
	uniq := make([]T, 0, len(in))

	ForEach(in, func(v T) {
		if None(uniq, func(other T) bool { return eq(other, v) }) {
			uniq = append(uniq, v)
		}
	})

	return uniq
}

// Intersect returns a slice of the elements that are in both the given slices.
func Intersect[T comparable](in ...[]T) []T {
	return IntersectFn(in, func(a, b T) bool {
		return a == b
	})
}

// IntersectFn returns a slice of the elements that are in both the given slices, according to the given function.
func IntersectFn[T any](in [][]T, eq func(T, T) bool) []T {
	if len(in) == 0 {
		return nil
	}

	if len(in) == 1 {
		return in[0]
	}

	return Reduce(in[1:], in[0], func(a, b []T) []T {
		return Filter(a, func(v T) bool {
			return Any(b, func(u T) bool {
				return eq(u, v)
			})
		})
	})
}

// Difference returns a slice of the elements that are in the first slice but not the second.
func Difference[T comparable](in ...[]T) []T {
	return DifferenceFn(in, func(a, b T) bool {
		return a == b
	})
}

// DifferenceFn returns a slice of the elements that are in the first slice but not the second, according to the given function.
func DifferenceFn[T any](in [][]T, eq func(T, T) bool) []T {
	if len(in) == 0 {
		return nil
	}

	if len(in) == 1 {
		return in[0]
	}

	return Reduce(in[1:], in[0], func(a, b []T) []T {
		return Filter(a, func(v T) bool {
			return !Any(b, func(u T) bool {
				return eq(u, v)
			})
		})
	})
}

// Power returns a slice of all the possible combinations of the given slice.
func Power[T any](in []T) [][]T {
	return MapEach(PowerIdx(len(in)), func(idx []int) []T {
		return MapEach(idx, func(i int) T {
			return in[i]
		})
	})
}

// PowerIdx returns a slice containing the indices of all the possible combinations of a slice of the given length.
func PowerIdx(n int) [][]int {
	return MapN(1<<n, func(i int) []int {
		var idx []int

		for j := 0; j < n; j++ {
			if i>>j&1 == 1 {
				idx = append(idx, j)
			}
		}

		return idx
	})
}

// Permutations returns a slice of all the possible permutations of the given slice.
func Permutations[T any](in []T) [][]T {
	return permute(in, factorial(len(in)), 0)
}

// PermutationsIdx returns a slice containing the indices of all the possible permutations of a slice of the given length.
func PermutationsIdx(n int) [][]int {
	return permute(RangeN(n), factorial(n), 0)
}

// Shuffle returns a shuffled slice of the given slice.
func Shuffle[T any](in []T) []T {
	return MapEach(rand.Perm(len(in)), func(idx int) T {
		return in[idx]
	})
}

// Reverse returns a reversed slice of the given slice.
func Reverse[T any](in []T) []T {
	out := make([]T, len(in))

	ForEachIdx(in, func(i int, v T) {
		out[len(in)-i-1] = v
	})

	return out
}

// Sort returns a sorted slice of the given slice.
func Sort[T constraints.Ordered](in []T) []T {
	return SortFn(in, func(a, b T) bool {
		return a < b
	})
}

// SortFn returns a sorted slice of the given slice, according to the given function.
func SortFn[T any](in []T, fn func(T, T) bool) []T {
	out := make([]T, len(in))

	copy(out, in)

	sort.Slice(out, func(i, j int) bool {
		return fn(out[i], out[j])
	})

	return out
}

// SetFn returns a map set of the elements in the given slice, with duplicates removed.
func Set[T comparable](in []T) map[T]struct{} {
	set := make(map[T]struct{}, len(in))

	ForEach(in, func(v T) {
		set[v] = struct{}{}
	})

	return set
}

// Contains returns true if the given slice contains the given element.
func Contains[T comparable](in []T, elem T) bool {
	return Any(in, func(v T) bool {
		return v == elem
	})
}

// ContainsAll returns true if the given slice contains all of the given elements.
func ContainsAll[T comparable](in []T, elems ...T) bool {
	return All(elems, func(elem T) bool {
		return Contains(in, elem)
	})
}

// ContainsAny returns true if the given slice contains any of the given elements.
func ContainsAny[T comparable](in []T, elems ...T) bool {
	return Any(elems, func(elem T) bool {
		return Contains(in, elem)
	})
}

// ContainsNone returns true if the given slice contains none of the given elements.
func ContainsNone[T comparable](in []T, elems ...T) bool {
	return !ContainsAny(in, elems...)
}

// Index returns the index of the given element in the given slice, or -1 if it is not found.
func Index[T comparable](in []T, elem T) int {
	return IndexFn(in, func(v T) bool {
		return v == elem
	})
}

// IndexFn returns the index of the first element that satisfies the given function in the given slice, or -1 if it is not found.
func IndexFn[T any](in []T, fn func(T) bool) int {
	for i, v := range in {
		if fn(v) {
			return i
		}
	}

	return -1
}

// IndexAll returns the indices of all elements that are equal to the given element in the given slice.
func IndexAll[T comparable](in []T, elem T) []int {
	return IndexAllFn(in, func(v T) bool {
		return v == elem
	})
}

// IndexAllFn returns the indices of all elements that satisfy the given function in the given slice.
func IndexAllFn[T any](in []T, fn func(T) bool) []int {
	indices := make([]int, 0, len(in))

	ForEachIdx(in, func(i int, v T) {
		if fn(v) {
			indices = append(indices, i)
		}
	})

	return indices
}

// Insert returns a slice with the given elements inserted at the given index.
func Insert[T any](in []T, idx int, elems ...T) []T {
	out := make([]T, len(in)+len(elems))

	copy(out, in[:idx])
	copy(out[idx:], elems)
	copy(out[idx+len(elems):], in[idx:])

	return out
}

// Remove returns a slice with the given elements removed.
func Remove[T comparable](in []T, elems ...T) []T {
	return RemoveFn(in, func(v T) bool {
		return Contains(elems, v)
	})
}

// RemoveFn returns a slice with the elements that satisfy the given function removed.
func RemoveFn[T any](in []T, fn func(T) bool) []T {
	out := make([]T, 0, len(in))

	ForEach(in, func(v T) {
		if !fn(v) {
			out = append(out, v)
		}
	})

	return out
}

// RemoveN returns a slice with n elements removed at the given index.
func RemoveN[T any](in []T, idx, n int) []T {
	return RemoveRange(in, idx, idx+n)
}

// RemoveRange returns a slice with the elements in the given range removed.
func RemoveRange[T any](in []T, start, end int) []T {
	out := make([]T, len(in)-(end-start))

	copy(out, in[:start])
	copy(out[start:], in[end:])

	return out
}

// RemoveIdx returns a slice with the elements at the given indices removed.
func RemoveIdx[T any](in []T, indices ...int) []T {
	out := make([]T, 0, len(in)-len(indices))

	ForEachIdx(in, func(i int, v T) {
		if !Contains(indices, i) {
			out = append(out, v)
		}
	})

	return out
}
