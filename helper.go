package dj

// zero returns the zero value of the given type.
func zero[T any]() T {
	var zero T
	return zero
}

// factorial returns the factorial of the given number.
func factorial(n int) int {
	if n < 2 {
		return 1
	}

	return n * factorial(n-1)
}

// permute returns all permutations of the given slice.
func permute[T any](in []T, size, from int) [][]T {
	if from == len(in)-1 {
		return [][]T{in}
	}

	out := make([][]T, size)

	for idx, perm := range permute(in, size/(len(in)-from), from+1) {
		out[idx] = perm
	}

	for j := from + 1; j < len(in); j++ {
		in[from], in[j] = in[j], in[from]

		for idx, perm := range permute(clone(in), size/(len(in)-from), from+1) {
			out[(j-from)*size/(len(in)-from)+idx] = perm
		}

		in[from], in[j] = in[j], in[from]
	}

	return out
}

// clone returns a copy of the given slice.
func clone[T any](v []T) []T {
	out := make([]T, len(v))

	copy(out, v)

	return out
}
