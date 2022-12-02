package dj_test

import (
	"testing"

	"github.com/jameshoulahan/dj"
	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Equal(tt.in...))
		})
	}
}

func TestJoin(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Join(tt.in...))
		})
	}
}

func TestZip(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Zip(tt.in...))
		})
	}
}

func TestUnzip(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Unzip(tt.in))
		})
	}
}

func TestChunk(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Chunk(tt.in, tt.size))
		})
	}
}

func TestPower(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Power(tt.in))
		})
	}
}

func TestPowerIdx(t *testing.T) {
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
			require.Equal(t, tt.want, dj.PowerIdx(tt.in))
		})
	}
}

func TestPermutations(t *testing.T) {
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
			require.ElementsMatch(t, tt.want, dj.Permutations(tt.in))
		})
	}
}

func TestPermutationsIdx(t *testing.T) {
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
			require.ElementsMatch(t, tt.want, dj.PermutationsIdx(tt.in))
		})
	}
}

func TestShuffle(t *testing.T) {
	for _, slice := range dj.Permutations(dj.RangeN(6)) {
		require.Equal(t, dj.Sort(dj.Shuffle(slice)), dj.Sort(slice))
	}
}

func TestCount(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Count(tt.in, tt.val))
		})
	}
}

func TestMin(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Min(tt.in))
		})
	}
}

func TestMax(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Max(tt.in))
		})
	}
}

func TestUniq(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Uniq(tt.in))
		})
	}
}

func TestIntersect(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Intersect(tt.a, tt.b))
		})
	}
}

func TestDifference(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Difference(tt.a, tt.b))
		})
	}
}

func TestSort(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Sort(tt.in))
		})
	}
}

func TestReverse(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Reverse(tt.in))
		})
	}
}

func TestSet(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Set(tt.in))
		})
	}
}

func TestContains(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Contains(tt.in, tt.val))
		})
	}
}

func TestContainsAll(t *testing.T) {
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
			require.Equal(t, tt.want, dj.ContainsAll(tt.in, tt.vals...))
		})
	}
}

func TestContainsAny(t *testing.T) {
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
			require.Equal(t, tt.want, dj.ContainsAny(tt.in, tt.vals...))
		})
	}
}

func TestContainsNone(t *testing.T) {
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
			require.Equal(t, tt.want, dj.ContainsNone(tt.in, tt.vals...))
		})
	}
}

func TestIndex(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Index(tt.in, tt.val))
		})
	}
}

func TestIndexAll(t *testing.T) {
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
			require.Equal(t, tt.want, dj.IndexAll(tt.in, tt.val))
		})
	}
}

func TestInsert(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Insert(tt.in, tt.idx, tt.vals...))
		})
	}
}

func TestRemove(t *testing.T) {
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
			require.Equal(t, tt.want, dj.Remove(tt.in, tt.vals...))
		})
	}
}

func TestRemoveN(t *testing.T) {
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
			require.Equal(t, tt.want, dj.RemoveN(tt.in, tt.idx, tt.n))
		})
	}
}

func TestRemoveRange(t *testing.T) {
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
			require.Equal(t, tt.want, dj.RemoveRange(tt.in, tt.from, tt.to))
		})
	}
}

func TestRemoveIdx(t *testing.T) {
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
			require.Equal(t, tt.want, dj.RemoveIdx(tt.in, tt.idxs...))
		})
	}
}
