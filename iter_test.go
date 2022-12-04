package dj_test

import (
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestIter_Collect(t *testing.T) {
	tests := []struct {
		name string
		in   dj.Iter[int]
		want []int
	}{
		{
			name: "single",
			in:   dj.SliceIter(1),
			want: []int{1},
		},

		{
			name: "double",
			in:   dj.SliceIter(1, 2),
			want: []int{1, 2},
		},

		{
			name: "empty",
			in:   dj.SliceIter[int](),
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.in.Collect())
		})
	}
}

func TestIter_Chan(t *testing.T) {
	tests := []struct {
		name string
		in   dj.Iter[int]
		want []int
	}{
		{
			name: "single",
			in:   dj.SliceIter(1),
			want: []int{1},
		},

		{
			name: "double",
			in:   dj.SliceIter(1, 2),
			want: []int{1, 2},
		},

		{
			name: "empty",
			in:   dj.SliceIter[int](),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, dj.CollectChan(tt.in.Recv()))
		})
	}
}

func TestWithPeek(t *testing.T) {
	iter := dj.WithPeek(dj.SliceIter(1, 2, 3))

	// Call read to get the first value.
	{
		next, ok := iter.Read()
		require.True(t, ok)
		require.Equal(t, 1, next)
	}

	// Call peek to get the next value without advancing the iterator.
	{
		peek, ok := iter.Peek()
		require.True(t, ok)
		require.Equal(t, 2, peek)
	}

	// Call peek again, the value should be the same.
	{
		peek, ok := iter.Peek()
		require.True(t, ok)
		require.Equal(t, 2, peek)
	}

	// Call read to get the next value.
	{
		next, ok := iter.Read()
		require.True(t, ok)
		require.Equal(t, 2, next)
	}

	// Call peek to peek at the last value.
	{
		peek, ok := iter.Peek()
		require.True(t, ok)
		require.Equal(t, 3, peek)
	}

	// Call read to get the last value.
	{
		next, ok := iter.Read()
		require.True(t, ok)
		require.Equal(t, 3, next)
	}

	// There should be no more values when calling peek.
	{
		_, ok := iter.Peek()
		require.False(t, ok)
	}

	// There should be no more values when calling next.
	{
		_, ok := iter.Read()
		require.False(t, ok)
	}
}

func TestMapIter(t *testing.T) {
	tests := []struct {
		name string
		in   dj.Iter[int]
		fn   func(int) int
		want []int
	}{
		{
			name: "add 1",
			in:   dj.SliceIter(1, 2, 3),
			fn:   func(i int) int { return i + 1 },
			want: []int{2, 3, 4},
		},

		{
			name: "double",
			in:   dj.SliceIter(1, 2, 3),
			fn:   func(i int) int { return i * 2 },
			want: []int{2, 4, 6},
		},

		{
			name: "empty",
			in:   dj.SliceIter[int](),
			fn:   func(i int) int { return i + 1 },
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, dj.MapIter(tt.in, tt.fn).Collect())
		})
	}
}

func TestChunkIter(t *testing.T) {
	tests := []struct {
		name string
		in   dj.Iter[int]
		size int
		want [][]int
	}{
		{
			name: "[1, 2, 3] size 1",
			in:   dj.SliceIter(1, 2, 3),
			size: 1,
			want: [][]int{{1}, {2}, {3}},
		},

		{
			name: "[1, 2, 3] size 2",
			in:   dj.SliceIter(1, 2, 3),
			size: 2,
			want: [][]int{{1, 2}, {3}},
		},

		{
			name: "[1, 2, 3] size 3",
			in:   dj.SliceIter(1, 2, 3),
			size: 3,
			want: [][]int{{1, 2, 3}},
		},

		{
			name: "[1, 2, 3] size 4",
			in:   dj.SliceIter(1, 2, 3),
			size: 4,
			want: [][]int{{1, 2, 3}},
		},

		{
			name: "empty",
			in:   dj.SliceIter[int](),
			size: 1,
			want: [][]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, dj.ChunkIter(tt.in, tt.size).Collect())
		})
	}
}

func TestFilterIter(t *testing.T) {
	tests := []struct {
		name string
		in   dj.Iter[int]
		fn   func(int) bool
		want []int
	}{
		{
			name: "odd",
			in:   dj.SliceIter(1, 2, 3),
			fn:   func(i int) bool { return i%2 == 1 },
			want: []int{1, 3},
		},

		{
			name: "even",
			in:   dj.SliceIter(1, 2, 3),
			fn:   func(i int) bool { return i%2 == 0 },
			want: []int{2},
		},

		{
			name: "empty",
			in:   dj.SliceIter[int](),
			fn:   func(i int) bool { return i%2 == 0 },
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, dj.FilterIter(tt.in, tt.fn).Collect())
		})
	}
}

func TestFlattenIter(t *testing.T) {
	tests := []struct {
		name string
		in   dj.Iter[dj.Iter[int]]
		want []int
	}{
		{
			name: "[[1, 2], [3, 4]]",
			in: dj.SliceIter(
				dj.SliceIter(1, 2),
				dj.SliceIter(3, 4),
			),
			want: []int{1, 2, 3, 4},
		},

		{
			name: "[[1, 2], [3, 4], []]",
			in: dj.SliceIter(
				dj.SliceIter(1, 2),
				dj.SliceIter(3, 4),
				dj.SliceIter[int](),
			),
			want: []int{1, 2, 3, 4},
		},

		{
			name: "[[], [1, 2], [3, 4]]",
			in: dj.SliceIter(
				dj.SliceIter[int](),
				dj.SliceIter(1, 2),
				dj.SliceIter(3, 4),
			),
			want: []int{1, 2, 3, 4},
		},

		{
			name: "empty",
			in:   dj.SliceIter[dj.Iter[int]](),
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, dj.FlattenIter(tt.in).Collect())
		})
	}
}

func TestJoinIter(t *testing.T) {
	tests := []struct {
		name string
		in   []dj.Iter[int]
		want []int
	}{
		{
			name: "[[1, 2], [3, 4]]",
			in: []dj.Iter[int]{
				dj.SliceIter(1, 2),
				dj.SliceIter(3, 4),
			},
			want: []int{1, 2, 3, 4},
		},

		{
			name: "[[1, 2], [3, 4], []]",
			in: []dj.Iter[int]{
				dj.SliceIter(1, 2),
				dj.SliceIter(3, 4),
				dj.SliceIter[int](),
			},
			want: []int{1, 2, 3, 4},
		},

		{
			name: "[[], [1, 2], [3, 4]]",
			in: []dj.Iter[int]{
				dj.SliceIter[int](),
				dj.SliceIter(1, 2),
				dj.SliceIter(3, 4),
			},
			want: []int{1, 2, 3, 4},
		},

		{
			name: "empty",
			in:   []dj.Iter[int]{},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, dj.JoinIter(tt.in...).Collect())
		})
	}
}

func TestZipIter(t *testing.T) {
	tests := []struct {
		name string
		in   []dj.Iter[int]
		want [][]int
	}{
		{
			name: "[[1, 2], [3, 4]]",
			in: []dj.Iter[int]{
				dj.SliceIter(1, 2),
				dj.SliceIter(3, 4),
			},
			want: [][]int{{1, 3}, {2, 4}},
		},

		{
			name: "[[1, 2], [3, 4], [5, 6]]",
			in: []dj.Iter[int]{
				dj.SliceIter(1, 2),
				dj.SliceIter(3, 4),
				dj.SliceIter(5, 6),
			},
			want: [][]int{{1, 3, 5}, {2, 4, 6}},
		},

		{
			name: "[[1, 2], [3]]",
			in: []dj.Iter[int]{
				dj.SliceIter(1, 2),
				dj.SliceIter(3),
			},
			want: [][]int{{1, 3}},
		},

		{
			name: "empty",
			in:   []dj.Iter[int]{},
			want: [][]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, dj.ZipIter(tt.in...).Collect())
		})
	}
}
