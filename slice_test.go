package xn_test

import (
	"testing"

	"github.com/jameshoulahan/xn"
	"github.com/stretchr/testify/require"
)

func Test_Equal(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want bool
	}{
		{
			name: "[1]",
			in:   [][]int{{1}},
			want: true,
		},

		{
			name: "[1, 2]",
			in:   [][]int{{1, 2}},
			want: true,
		},

		{
			name: "[1, 2], [1, 2]",
			in:   [][]int{{1, 2}, {1, 2}},
			want: true,
		},

		{
			name: "[1, 2], [1, 3]",
			in:   [][]int{{1, 2}, {1, 3}},
			want: false,
		},

		{
			name: "[1, 2], [1, 2, 3]",
			in:   [][]int{{1, 2}, {1, 2, 3}},
			want: false,
		},

		{
			name: "[1, 2], [1, 2], [1, 2]",
			in:   [][]int{{1, 2}, {1, 2}, {1, 2}},
			want: true,
		},

		{
			name: "[1, 2], [1, 2], [1, 3]",
			in:   [][]int{{1, 2}, {1, 2}, {1, 3}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Equal(tt.in...))
		})
	}
}

func Test_Join(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want []int
	}{
		{
			name: "[1]",
			in:   [][]int{{1}},
			want: []int{1},
		},

		{
			name: "[1, 2]",
			in:   [][]int{{1}, {2}},
			want: []int{1, 2},
		},

		{
			name: "[1, 2], [3, 4]",
			in:   [][]int{{1, 2}, {3, 4}},
			want: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Join(tt.in...))
		})
	}
}

func Test_Zip(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want [][]int
	}{
		{
			name: "[1]",
			in:   [][]int{{1}},
			want: [][]int{{1}},
		},

		{
			name: "[1, 2]",
			in:   [][]int{{1, 2}},
			want: [][]int{{1}, {2}},
		},

		{
			name: "[1, 2], [3, 4]",
			in:   [][]int{{1, 2}, {3, 4}},
			want: [][]int{{1, 3}, {2, 4}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Zip(tt.in...))
		})
	}
}

func Test_Unzip(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want [][]int
	}{
		{
			name: "[1]",
			in:   [][]int{{1}},
			want: [][]int{{1}},
		},

		{
			name: "[1, 2]",
			in:   [][]int{{1}, {2}},
			want: [][]int{{1, 2}},
		},

		{
			name: "[1, 2], [3, 4]",
			in:   [][]int{{1, 3}, {2, 4}},
			want: [][]int{{1, 2}, {3, 4}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Unzip(tt.in))
		})
	}
}

func Test_Chunk(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		size int
		want [][]int
	}{
		{
			name: "[1]",
			in:   []int{1},
			size: 1,
			want: [][]int{{1}},
		},

		{
			name: "[1, 2]",
			in:   []int{1, 2},
			size: 1,
			want: [][]int{{1}, {2}},
		},

		{
			name: "[1, 2, 3]",
			in:   []int{1, 2, 3},
			size: 2,
			want: [][]int{{1, 2}, {3}},
		},

		{
			name: "[1, 2, 3, 4]",
			in:   []int{1, 2, 3, 4},
			size: 2,
			want: [][]int{{1, 2}, {3, 4}},
		},

		{
			name: "[1, 2, 3, 4, 5]",
			in:   []int{1, 2, 3, 4, 5},
			size: 2,
			want: [][]int{{1, 2}, {3, 4}, {5}},
		},

		{
			name: "[1, 2, 3, 4, 5]",
			in:   []int{1, 2, 3, 4, 5},
			size: 3,
			want: [][]int{{1, 2, 3}, {4, 5}},
		},

		{
			name: "[1, 2, 3, 4, 5]",
			in:   []int{1, 2, 3, 4, 5},
			size: 4,
			want: [][]int{{1, 2, 3, 4}, {5}},
		},

		{
			name: "[1, 2, 3, 4, 5]",
			in:   []int{1, 2, 3, 4, 5},
			size: 5,
			want: [][]int{{1, 2, 3, 4, 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Chunk(tt.in, tt.size))
		})
	}
}

func Test_Power(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want [][]int
	}{
		{
			name: "[1]",
			in:   []int{1},
			want: [][]int{{}, {1}},
		},

		{
			name: "[1, 2]",
			in:   []int{1, 2},
			want: [][]int{{}, {1}, {2}, {1, 2}},
		},

		{
			name: "[1, 2, 3]",
			in:   []int{1, 2, 3},
			want: [][]int{{}, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Power(tt.in))
		})
	}
}

func Test_PowerIdx(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want [][]int
	}{
		{
			name: "1",
			in:   1,
			want: [][]int{nil, {0}},
		},

		{
			name: "2",
			in:   2,
			want: [][]int{nil, {0}, {1}, {0, 1}},
		},

		{
			name: "3",
			in:   3,
			want: [][]int{nil, {0}, {1}, {0, 1}, {2}, {0, 2}, {1, 2}, {0, 1, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.PowerIdx(tt.in))
		})
	}
}

func Test_Permutations(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want [][]int
	}{
		{
			name: "[1]",
			in:   []int{1},
			want: [][]int{{1}},
		},

		{
			name: "[1, 2]",
			in:   []int{1, 2},
			want: [][]int{{1, 2}, {2, 1}},
		},

		{
			name: "[1, 2, 3]",
			in:   []int{1, 2, 3},
			want: [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ElementsMatch(t, tt.want, xn.Permutations(tt.in))
		})
	}
}

func Test_PermutationsIdx(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want [][]int
	}{
		{
			name: "1",
			in:   1,
			want: [][]int{{0}},
		},

		{
			name: "2",
			in:   2,
			want: [][]int{{0, 1}, {1, 0}},
		},

		{
			name: "3",
			in:   3,
			want: [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ElementsMatch(t, tt.want, xn.PermutationsIdx(tt.in))
		})
	}
}

func Test_Shuffle(t *testing.T) {
	for _, slice := range xn.Permutations(xn.RangeN(6)) {
		require.Equal(t, xn.Sort(xn.Shuffle(slice)), xn.Sort(slice))
	}
}

func Test_Count(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		val  int
		want int
	}{
		{
			name: "[1] 1",
			in:   []int{1},
			val:  1,
			want: 1,
		},

		{
			name: "[1, 2] 1",
			in:   []int{1, 2},
			val:  1,
			want: 1,
		},

		{
			name: "[1, 2] 2",
			in:   []int{1, 2},
			val:  2,
			want: 1,
		},

		{
			name: "[1, 2] 3",
			in:   []int{1, 2},
			val:  3,
			want: 0,
		},

		{
			name: "[1, 2, 2, 3] 2",
			in:   []int{1, 2, 2, 3},
			val:  2,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Count(tt.in, tt.val))
		})
	}
}

func Test_Min(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{
			name: "[1]",
			in:   []int{1},
			want: 1,
		},

		{
			name: "[1, 2]",
			in:   []int{1, 2},
			want: 1,
		},

		{
			name: "[1, 2, 3]",
			in:   []int{1, 2, 3},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Min(tt.in))
		})
	}
}

func Test_Max(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{
			name: "[1]",
			in:   []int{1},
			want: 1,
		},

		{
			name: "[1, 2]",
			in:   []int{1, 2},
			want: 2,
		},

		{
			name: "[1, 2, 3]",
			in:   []int{1, 2, 3},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Max(tt.in))
		})
	}
}

func Test_Uniq(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{
			name: "[1]",
			in:   []int{1},
			want: []int{1},
		},

		{
			name: "[1, 2]",
			in:   []int{1, 2},
			want: []int{1, 2},
		},

		{
			name: "[1, 2, 3]",
			in:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},

		{
			name: "[1, 2, 1, 3]",
			in:   []int{1, 2, 1, 3},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Uniq(tt.in))
		})
	}
}

func Test_Intersect(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{
			name: "[1, 2] [1, 2]",
			a:    []int{1, 2},
			b:    []int{1, 2},
			want: []int{1, 2},
		},

		{
			name: "[1, 2] [1, 3]",
			a:    []int{1, 2},
			b:    []int{1, 3},
			want: []int{1},
		},

		{
			name: "[1, 2] [3, 4]",
			a:    []int{1, 2},
			b:    []int{3, 4},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Intersect(tt.a, tt.b))
		})
	}
}

func Test_Difference(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{
			name: "[1, 2] [1, 2]",
			a:    []int{1, 2},
			b:    []int{1, 2},
			want: []int{},
		},

		{
			name: "[1, 2] [1, 3]",
			a:    []int{1, 2},
			b:    []int{1, 3},
			want: []int{2},
		},

		{
			name: "[1, 2] [3, 4]",
			a:    []int{1, 2},
			b:    []int{3, 4},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Difference(tt.a, tt.b))
		})
	}
}

func Test_Sort(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{
			name: "[1]",
			in:   []int{1},
			want: []int{1},
		},

		{
			name: "[2, 1]",
			in:   []int{2, 1},
			want: []int{1, 2},
		},

		{
			name: "[3, 2, 1]",
			in:   []int{3, 2, 1},
			want: []int{1, 2, 3},
		},

		{
			name: "[1, 2, 3]",
			in:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},

		{
			name: "[1, 2, 1, 3]",
			in:   []int{1, 2, 1, 3},
			want: []int{1, 1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Sort(tt.in))
		})
	}
}

func Test_Reverse(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{
			name: "[1]",
			in:   []int{1},
			want: []int{1},
		},

		{
			name: "[2, 1]",
			in:   []int{2, 1},
			want: []int{1, 2},
		},

		{
			name: "[3, 2, 1]",
			in:   []int{3, 2, 1},
			want: []int{1, 2, 3},
		},

		{
			name: "[1, 2, 3]",
			in:   []int{1, 2, 3},
			want: []int{3, 2, 1},
		},

		{
			name: "[1, 2, 1, 3]",
			in:   []int{1, 2, 1, 3},
			want: []int{3, 1, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Reverse(tt.in))
		})
	}
}

func Test_Set(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want map[int]struct{}
	}{
		{
			name: "[1]",
			in:   []int{1},
			want: map[int]struct{}{1: {}},
		},

		{
			name: "[2, 1]",
			in:   []int{2, 1},
			want: map[int]struct{}{1: {}, 2: {}},
		},

		{
			name: "[3, 2, 1]",
			in:   []int{3, 2, 1},
			want: map[int]struct{}{1: {}, 2: {}, 3: {}},
		},

		{
			name: "[1, 2, 3]",
			in:   []int{1, 2, 3},
			want: map[int]struct{}{1: {}, 2: {}, 3: {}},
		},

		{
			name: "[1, 2, 1, 3]",
			in:   []int{1, 2, 1, 3},
			want: map[int]struct{}{1: {}, 2: {}, 3: {}},
		},

		{
			name: "[]",
			in:   []int{},
			want: map[int]struct{}{},
		},

		{
			name: "[1, 1, 1, 1]",
			in:   []int{1, 1, 1, 1},
			want: map[int]struct{}{1: {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Set(tt.in))
		})
	}
}

func Test_Contains(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		val  int
		want bool
	}{
		{
			name: "[1] 1",
			in:   []int{1},
			val:  1,
			want: true,
		},

		{
			name: "[2, 1] 1",
			in:   []int{2, 1},
			val:  1,
			want: true,
		},

		{
			name: "[3, 2, 1] 2",
			in:   []int{3, 2, 1},
			val:  2,
			want: true,
		},

		{
			name: "[1, 2, 3] 3",
			in:   []int{1, 2, 3},
			val:  3,
			want: true,
		},

		{
			name: "[1, 2, 3] 4",
			in:   []int{1, 2, 3},
			val:  4,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Contains(tt.in, tt.val))
		})
	}
}

func Test_ContainsAll(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		vals []int
		want bool
	}{
		{
			name: "[1] [1]",
			in:   []int{1},
			vals: []int{1},
			want: true,
		},

		{
			name: "[2, 1] [1]",
			in:   []int{2, 1},
			vals: []int{1},
			want: true,
		},

		{
			name: "[3, 2, 1] [2, 1]",
			in:   []int{3, 2, 1},
			vals: []int{2, 1},
			want: true,
		},

		{
			name: "[1, 2, 3] [3, 2, 1]",
			in:   []int{1, 2, 3},
			vals: []int{3, 2, 1},
			want: true,
		},

		{
			name: "[1, 2, 3] [4, 3, 2]",
			in:   []int{1, 2, 3},
			vals: []int{4, 3, 2},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.ContainsAll(tt.in, tt.vals...))
		})
	}
}

func Test_ContainsAny(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		vals []int
		want bool
	}{
		{
			name: "[1] [1]",
			in:   []int{1},
			vals: []int{1},
			want: true,
		},

		{
			name: "[2, 1] [1]",
			in:   []int{2, 1},
			vals: []int{1},
			want: true,
		},

		{
			name: "[3, 2, 1] [2, 1]",
			in:   []int{3, 2, 1},
			vals: []int{2, 1},
			want: true,
		},

		{
			name: "[1, 2, 3] [3, 2, 1]",
			in:   []int{1, 2, 3},
			vals: []int{3, 2, 1},
			want: true,
		},

		{
			name: "[1, 2, 3] [4, 3, 2]",
			in:   []int{1, 2, 3},
			vals: []int{4, 3, 2},
			want: true,
		},

		{
			name: "[1, 2, 3] [4, 5, 6]",
			in:   []int{1, 2, 3},
			vals: []int{4, 5, 6},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.ContainsAny(tt.in, tt.vals...))
		})
	}
}

func Test_ContainsNone(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		vals []int
		want bool
	}{
		{
			name: "[1] [1]",
			in:   []int{1},
			vals: []int{1},
			want: false,
		},

		{
			name: "[2, 1] [1]",
			in:   []int{2, 1},
			vals: []int{1},
			want: false,
		},

		{
			name: "[3, 2, 1] [2, 1]",
			in:   []int{3, 2, 1},
			vals: []int{2, 1},
			want: false,
		},

		{
			name: "[1, 2, 3] [3, 2, 1]",
			in:   []int{1, 2, 3},
			vals: []int{3, 2, 1},
			want: false,
		},

		{
			name: "[1, 2, 3] [4, 3, 2]",
			in:   []int{1, 2, 3},
			vals: []int{4, 3, 2},
			want: false,
		},

		{
			name: "[1, 2, 3] [4, 5, 6]",
			in:   []int{1, 2, 3},
			vals: []int{4, 5, 6},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.ContainsNone(tt.in, tt.vals...))
		})
	}
}

func Test_Index(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		val  int
		want int
	}{
		{
			name: "[1] 1",
			in:   []int{1},
			val:  1,
			want: 0,
		},

		{
			name: "[2, 1] 1",
			in:   []int{2, 1},
			val:  1,
			want: 1,
		},

		{
			name: "[3, 2, 1] 2",
			in:   []int{3, 2, 1},
			val:  2,
			want: 1,
		},

		{
			name: "[1, 2, 3] 3",
			in:   []int{1, 2, 3},
			val:  3,
			want: 2,
		},

		{
			name: "[1, 2, 3] 4",
			in:   []int{1, 2, 3},
			val:  4,
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Index(tt.in, tt.val))
		})
	}
}

func Test_IndexAll(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		val  int
		want []int
	}{
		{
			name: "[1] 1",
			in:   []int{1},
			val:  1,
			want: []int{0},
		},

		{
			name: "[2, 1] 1",
			in:   []int{2, 1},
			val:  1,
			want: []int{1},
		},

		{
			name: "[3, 2, 1] 2",
			in:   []int{3, 2, 1},
			val:  2,
			want: []int{1},
		},

		{
			name: "[1, 2, 3] 3",
			in:   []int{1, 2, 3},
			val:  3,
			want: []int{2},
		},

		{
			name: "[1, 2, 3] 4",
			in:   []int{1, 2, 3},
			val:  4,
			want: []int{},
		},

		{
			name: "[1, 2, 3, 1] 1",
			in:   []int{1, 2, 3, 1},
			val:  1,
			want: []int{0, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.IndexAll(tt.in, tt.val))
		})
	}
}

func Test_Insert(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		idx  int
		vals []int
		want []int
	}{
		{
			name: "[1] 0 [1]",
			in:   []int{1},
			idx:  0,
			vals: []int{1},
			want: []int{1, 1},
		},

		{
			name: "[2, 1] 0 [1]",
			in:   []int{2, 1},
			idx:  0,
			vals: []int{1},
			want: []int{1, 2, 1},
		},

		{
			name: "[3, 2, 1] 0 [2, 1]",
			in:   []int{3, 2, 1},
			idx:  0,
			vals: []int{2, 1},
			want: []int{2, 1, 3, 2, 1},
		},

		{
			name: "[1, 2, 3] 1 [3, 2, 1]",
			in:   []int{1, 2, 3},
			idx:  1,
			vals: []int{3, 2, 1},
			want: []int{1, 3, 2, 1, 2, 3},
		},

		{
			name: "[1, 2, 3] 2 [4, 3, 2]",
			in:   []int{1, 2, 3},
			idx:  2,
			vals: []int{4, 3, 2},
			want: []int{1, 2, 4, 3, 2, 3},
		},

		{
			name: "[1, 2, 3] 3 [4, 5, 6]",
			in:   []int{1, 2, 3},
			idx:  3,
			vals: []int{4, 5, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Insert(tt.in, tt.idx, tt.vals...))
		})
	}
}

func Test_Remove(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		vals []int
		want []int
	}{
		{
			name: "[1] [1]",
			in:   []int{1},
			vals: []int{1},
			want: []int{},
		},

		{
			name: "[2, 1] [1]",
			in:   []int{2, 1},
			vals: []int{1},
			want: []int{2},
		},

		{
			name: "[3, 2, 1] [2, 1]",
			in:   []int{3, 2, 1},
			vals: []int{2, 1},
			want: []int{3},
		},

		{
			name: "[1, 2, 3] [3, 2, 1]",
			in:   []int{1, 2, 3},
			vals: []int{3, 2, 1},
			want: []int{},
		},

		{
			name: "[1, 2, 3] [4, 5, 6]",
			in:   []int{1, 2, 3},
			vals: []int{4, 5, 6},
			want: []int{1, 2, 3},
		},

		{
			name: "[1, 2, 3, 1] [1]",
			in:   []int{1, 2, 3, 1},
			vals: []int{1},
			want: []int{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.Remove(tt.in, tt.vals...))
		})
	}
}

func Test_RemoveN(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		idx  int
		n    int
		want []int
	}{
		{
			name: "[1] 0 1",
			in:   []int{1},
			idx:  0,
			n:    1,
			want: []int{},
		},

		{
			name: "[2, 1] 0 1",
			in:   []int{2, 1},
			idx:  0,
			n:    1,
			want: []int{1},
		},

		{
			name: "[2, 1] 0 2",
			in:   []int{2, 1},
			idx:  0,
			n:    2,
			want: []int{},
		},

		{
			name: "[3, 2, 1] 0 1",
			in:   []int{3, 2, 1},
			idx:  0,
			n:    1,
			want: []int{2, 1},
		},

		{
			name: "[3, 2, 1] 1 2",
			in:   []int{3, 2, 1},
			idx:  1,
			n:    2,
			want: []int{3},
		},

		{
			name: "[3, 2, 1] 2 1",
			in:   []int{3, 2, 1},
			idx:  2,
			n:    1,
			want: []int{3, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.RemoveN(tt.in, tt.idx, tt.n))
		})
	}
}

func Test_RemoveRange(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		from int
		to   int
		want []int
	}{
		{
			name: "[1] 0 1",
			in:   []int{1},
			from: 0,
			to:   1,
			want: []int{},
		},

		{
			name: "[2, 1] 0 1",
			in:   []int{2, 1},
			from: 0,
			to:   1,
			want: []int{1},
		},

		{
			name: "[2, 1] 0 2",
			in:   []int{2, 1},
			from: 0,
			to:   2,
			want: []int{},
		},

		{
			name: "[3, 2, 1] 0 1",
			in:   []int{3, 2, 1},
			from: 0,
			to:   1,
			want: []int{2, 1},
		},

		{
			name: "[3, 2, 1] 1 2",
			in:   []int{3, 2, 1},
			from: 1,
			to:   2,
			want: []int{3, 1},
		},

		{
			name: "[3, 2, 1] 2 3",
			in:   []int{3, 2, 1},
			from: 2,
			to:   3,
			want: []int{3, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.RemoveRange(tt.in, tt.from, tt.to))
		})
	}
}

func Test_RemoveIdx(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		idxs []int
		want []int
	}{
		{
			name: "[1] [0]",
			in:   []int{1},
			idxs: []int{0},
			want: []int{},
		},

		{
			name: "[2, 1] [0]",
			in:   []int{2, 1},
			idxs: []int{0},
			want: []int{1},
		},

		{
			name: "[3, 2, 1] [0, 1]",
			in:   []int{3, 2, 1},
			idxs: []int{0, 1},
			want: []int{1},
		},

		{
			name: "[1, 2, 3] [0, 1, 2]",
			in:   []int{1, 2, 3},
			idxs: []int{0, 1, 2},
			want: []int{},
		},

		{
			name: "[1, 2, 3, 1] [0, 3]",
			in:   []int{1, 2, 3, 1},
			idxs: []int{0, 3},
			want: []int{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, xn.RemoveIdx(tt.in, tt.idxs...))
		})
	}
}
