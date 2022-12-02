package dj

// Zero returns the Zero value of the given type.
func Zero[T any]() T {
	var zero T
	return zero
}

// Pointer returns a pointer to the given value.
func Pointer[T any](v T) *T {
	return &v
}

// Factorial returns the Factorial of the given number.
func Factorial(n int) int {
	if n < 2 {
		return 1
	}

	return n * Factorial(n-1)
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

		for idx, perm := range permute(Clone(in), size/(len(in)-from), from+1) {
			out[(j-from)*size/(len(in)-from)+idx] = perm
		}

		in[from], in[j] = in[j], in[from]
	}

	return out
}
