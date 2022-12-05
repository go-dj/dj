package dj_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestForN(t *testing.T) {
	var got []int

	dj.ForN(10, func(i int) {
		got = append(got, i)
	})

	require.Equal(t, dj.RangeN(10), got)
}

func TestForNErr(t *testing.T) {
	var got []int

	err := dj.ForNErr(10, func(i int) error {
		got = append(got, i)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, dj.RangeN(10), got)

	require.Error(t, dj.ForNErr(10, func(i int) error {
		return errors.New("oops")
	}))
}

func TestGoForN(t *testing.T) {
	var got dj.Queue[int]

	dj.GoForN(context.Background(), 10, func(_ context.Context, i int) {
		got.PushBack(i)
	})

	require.ElementsMatch(t, dj.RangeN(10), got.Items())
}

func TestGoForNErr(t *testing.T) {
	var got dj.Queue[int]

	err := dj.GoForNErr(context.Background(), 10, func(_ context.Context, i int) error {
		got.PushBack(i)
		return nil
	})
	require.NoError(t, err)
	require.ElementsMatch(t, dj.RangeN(10), got.Items())

	require.Error(t, dj.GoForNErr(context.Background(), 10, func(_ context.Context, i int) error {
		return errors.New("oops")
	}))
}

func TestMapN(t *testing.T) {
	got := dj.MapN(10, func(i int) int {
		return i + 1
	})

	require.Equal(t, dj.Range(1, 11), got)
}

func TestMapNErr(t *testing.T) {
	got, err := dj.MapNErr(10, func(i int) (int, error) {
		return i + 1, nil
	})
	require.NoError(t, err)
	require.Equal(t, dj.Range(1, 11), got)

	bad, err := dj.MapNErr(10, func(i int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

func TestGoMapN(t *testing.T) {
	got, err := dj.GoMapNErr(context.Background(), 10, func(_ context.Context, i int) (int, error) {
		return i + 1, nil
	})
	require.NoError(t, err)
	require.Equal(t, dj.Range(1, 11), got)
}

func TestGoMapNErr(t *testing.T) {
	got, err := dj.GoMapNErr(context.Background(), 10, func(_ context.Context, i int) (int, error) {
		return i + 1, nil
	})
	require.NoError(t, err)
	require.Equal(t, dj.Range(1, 11), got)

	bad, err := dj.GoMapNErr(context.Background(), 10, func(_ context.Context, i int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

func TestForEach(t *testing.T) {
	var got []int

	dj.ForEach(dj.RangeN(10), func(i int) {
		got = append(got, i)
	})

	require.Equal(t, dj.RangeN(10), got)
}

func TestForEachErr(t *testing.T) {
	var got []int

	err := dj.ForEachErr(dj.RangeN(10), func(i int) error {
		got = append(got, i)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, dj.RangeN(10), got)

	require.Error(t, dj.ForEachErr(dj.RangeN(10), func(i int) error {
		return errors.New("oops")
	}))
}

func TestForEachIdx(t *testing.T) {
	got := make(map[int]int)

	dj.ForEachIdx(dj.RangeN(3), func(idx, i int) {
		got[idx] = i
	})

	require.Equal(t, map[int]int{0: 0, 1: 1, 2: 2}, got)
}

func TestForEachIdxErr(t *testing.T) {
	got := make(map[int]int)

	err := dj.ForEachIdxErr(dj.RangeN(3), func(idx, i int) error {
		got[idx] = i
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, map[int]int{0: 0, 1: 1, 2: 2}, got)

	require.Error(t, dj.ForEachIdxErr(dj.RangeN(3), func(idx, i int) error {
		return errors.New("oops")
	}))
}

func TestGoForEach(t *testing.T) {
	var got dj.Queue[int]

	dj.GoForEach(context.Background(), dj.RangeN(10), func(_ context.Context, i int) {
		got.PushBack(i)
	})

	require.ElementsMatch(t, dj.RangeN(10), got.Items())
}

func TestGoForEachErr(t *testing.T) {
	var got dj.Queue[int]

	err := dj.GoForEachErr(context.Background(), dj.RangeN(10), func(_ context.Context, i int) error {
		got.PushBack(i)
		return nil
	})
	require.NoError(t, err)
	require.ElementsMatch(t, dj.RangeN(10), got.Items())

	require.Error(t, dj.GoForEachErr(context.Background(), dj.RangeN(10), func(_ context.Context, i int) error {
		return errors.New("oops")
	}))
}

func TestGoForEachIdx(t *testing.T) {
	got := dj.NewMap[int, int]()

	dj.GoForEachIdx(context.Background(), dj.RangeN(3), func(_ context.Context, idx, i int) {
		got.Set(idx, i)
	})

	require.Equal(t, map[int]int{0: 0, 1: 1, 2: 2}, got.Items())
}

func TestGoForEachIdxErr(t *testing.T) {
	got := dj.NewMap[int, int]()

	err := dj.GoForEachIdxErr(context.Background(), dj.RangeN(3), func(_ context.Context, idx, i int) error {
		got.Set(idx, i)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, map[int]int{0: 0, 1: 1, 2: 2}, got.Items())

	require.Error(t, dj.GoForEachIdxErr(context.Background(), dj.RangeN(3), func(_ context.Context, idx, i int) error {
		return errors.New("oops")
	}))
}

func TestMapEach(t *testing.T) {
	got := dj.MapEach(dj.RangeN(10), func(i int) int {
		return i + 1
	})

	require.Equal(t, dj.Range(1, 11), got)
}

func TestMapEachErr(t *testing.T) {
	got, err := dj.MapEachErr(dj.RangeN(10), func(i int) (int, error) {
		return i + 1, nil
	})
	require.NoError(t, err)
	require.Equal(t, dj.Range(1, 11), got)

	bad, err := dj.MapEachErr(dj.RangeN(10), func(i int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

func TestMapEachIdx(t *testing.T) {
	got := dj.MapEachIdx(dj.RangeN(10), func(idx, i int) int {
		if idx%2 == 0 {
			return i * 2
		}

		return i * -2
	})

	require.Equal(t, []int{0, -2, 4, -6, 8, -10, 12, -14, 16, -18}, got)
}

func TestMapEachIdxErr(t *testing.T) {
	got, err := dj.MapEachIdxErr(dj.RangeN(10), func(idx, i int) (int, error) {
		if idx%2 == 0 {
			return i * 2, nil
		}

		return i * -2, nil
	})
	require.NoError(t, err)
	require.Equal(t, []int{0, -2, 4, -6, 8, -10, 12, -14, 16, -18}, got)

	bad, err := dj.MapEachIdxErr(dj.RangeN(10), func(idx, i int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

func TestGoMap(t *testing.T) {
	got := dj.GoMapEach(context.Background(), dj.RangeN(10), func(_ context.Context, i int) int {
		return i + 1
	})

	require.Equal(t, dj.Range(1, 11), got)
}

func TestGoMapErr(t *testing.T) {
	got, err := dj.GoMapEachErr(context.Background(), dj.RangeN(10), func(_ context.Context, i int) (int, error) {
		return i + 1, nil
	})
	require.NoError(t, err)
	require.Equal(t, dj.Range(1, 11), got)

	bad, err := dj.GoMapEachErr(context.Background(), dj.RangeN(10), func(_ context.Context, i int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

func TestGoMapIdx(t *testing.T) {
	got := dj.GoMapEachIdx(context.Background(), dj.RangeN(10), func(_ context.Context, idx, i int) int {
		if idx%2 == 0 {
			return i * 2
		}

		return i * -2
	})

	require.Equal(t, []int{0, -2, 4, -6, 8, -10, 12, -14, 16, -18}, got)
}

func TestGoMapIdxErr(t *testing.T) {
	got, err := dj.GoMapEachIdxErr(context.Background(), dj.RangeN(10), func(_ context.Context, idx, i int) (int, error) {
		if idx%2 == 0 {
			return i * 2, nil
		}

		return i * -2, nil
	})
	require.NoError(t, err)
	require.Equal(t, []int{0, -2, 4, -6, 8, -10, 12, -14, 16, -18}, got)

	bad, err := dj.GoMapEachIdxErr(context.Background(), dj.RangeN(10), func(_ context.Context, idx, i int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

func TestForWindow(t *testing.T) {
	var got [][]int

	dj.ForWindow(dj.RangeN(5), 3, func(window []int) {
		got = append(got, window)
	})

	require.Equal(t, [][]int{
		{0, 1, 2},
		{1, 2, 3},
		{2, 3, 4},
	}, got)
}

func TestForWindowErr(t *testing.T) {
	var got [][]int

	err := dj.ForWindowErr(dj.RangeN(5), 3, func(window []int) error {
		got = append(got, window)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, [][]int{
		{0, 1, 2},
		{1, 2, 3},
		{2, 3, 4},
	}, got)

	require.Error(t, dj.ForWindowErr(dj.RangeN(5), 3, func(window []int) error {
		return errors.New("oops")
	}))
}

func TestForWindowIdx(t *testing.T) {
	got := make(map[int][]int)

	dj.ForWindowIdx(dj.RangeN(5), 3, func(idx int, window []int) {
		got[idx] = window
	})

	require.Equal(t, map[int][]int{
		0: {0, 1, 2},
		1: {1, 2, 3},
		2: {2, 3, 4},
	}, got)
}

func TestForWindowIdxErr(t *testing.T) {
	got := make(map[int][]int)

	err := dj.ForWindowIdxErr(dj.RangeN(5), 3, func(idx int, window []int) error {
		got[idx] = window
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, map[int][]int{
		0: {0, 1, 2},
		1: {1, 2, 3},
		2: {2, 3, 4},
	}, got)

	require.Error(t, dj.ForWindowIdxErr(dj.RangeN(5), 3, func(idx int, window []int) error {
		return errors.New("oops")
	}))
}

func TestGoForWindow(t *testing.T) {
	var got dj.Queue[[]int]

	dj.GoForWindow(context.Background(), dj.RangeN(5), 3, func(_ context.Context, window []int) {
		got.PushBack(window)
	})

	require.ElementsMatch(t, [][]int{
		{0, 1, 2},
		{1, 2, 3},
		{2, 3, 4},
	}, got.Items())
}

func TestGoForWindowErr(t *testing.T) {
	var got dj.Queue[[]int]

	err := dj.GoForWindowErr(context.Background(), dj.RangeN(5), 3, func(_ context.Context, window []int) error {
		got.PushBack(window)
		return nil
	})
	require.NoError(t, err)
	require.ElementsMatch(t, [][]int{
		{0, 1, 2},
		{1, 2, 3},
		{2, 3, 4},
	}, got.Items())

	require.Error(t, dj.GoForWindowErr(context.Background(), dj.RangeN(5), 3, func(_ context.Context, window []int) error {
		return errors.New("oops")
	}))
}

func TestGoForWindowIdx(t *testing.T) {
	got := make(map[int][]int)

	dj.GoForWindowIdx(context.Background(), dj.RangeN(5), 3, func(_ context.Context, idx int, window []int) {
		got[idx] = window
	})

	require.Equal(t, map[int][]int{
		0: {0, 1, 2},
		1: {1, 2, 3},
		2: {2, 3, 4},
	}, got)
}

func TestGoForWindowIdxErr(t *testing.T) {
	got := make(map[int][]int)

	err := dj.GoForWindowIdxErr(context.Background(), dj.RangeN(5), 3, func(_ context.Context, idx int, window []int) error {
		got[idx] = window
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, map[int][]int{
		0: {0, 1, 2},
		1: {1, 2, 3},
		2: {2, 3, 4},
	}, got)

	require.Error(t, dj.GoForWindowIdxErr(context.Background(), dj.RangeN(5), 3, func(_ context.Context, idx int, window []int) error {
		return errors.New("oops")
	}))
}

func TestMapWindow(t *testing.T) {
	got := dj.MapWindow(dj.RangeN(5), 3, func(window []int) int {
		return window[0] + window[1] + window[2]
	})

	require.Equal(t, []int{3, 6, 9}, got)
}

func TestMapWindowErr(t *testing.T) {
	got, err := dj.MapWindowErr(dj.RangeN(5), 3, func(window []int) (int, error) {
		return window[0] + window[1] + window[2], nil
	})
	require.NoError(t, err)
	require.Equal(t, []int{3, 6, 9}, got)

	bad, err := dj.MapWindowErr(dj.RangeN(5), 3, func(window []int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

func TestMapWindowIdx(t *testing.T) {
	got := dj.MapWindowIdx(dj.RangeN(5), 3, func(idx int, window []int) int {
		return idx + window[0] + window[1] + window[2]
	})

	require.Equal(t, []int{3, 7, 11}, got)
}

func TestMapWindowIdxErr(t *testing.T) {
	got, err := dj.MapWindowIdxErr(dj.RangeN(5), 3, func(idx int, window []int) (int, error) {
		return idx + window[0] + window[1] + window[2], nil
	})
	require.NoError(t, err)
	require.Equal(t, []int{3, 7, 11}, got)

	bad, err := dj.MapWindowIdxErr(dj.RangeN(5), 3, func(idx int, window []int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

func TestGoMapWindow(t *testing.T) {
	got := dj.GoMapWindow(context.Background(), dj.RangeN(5), 3, func(_ context.Context, window []int) int {
		return window[0] + window[1] + window[2]
	})

	require.ElementsMatch(t, []int{3, 6, 9}, got)
}

func TestGoMapWindowErr(t *testing.T) {
	got, err := dj.GoMapWindowErr(context.Background(), dj.RangeN(5), 3, func(_ context.Context, window []int) (int, error) {
		return window[0] + window[1] + window[2], nil
	})
	require.NoError(t, err)
	require.ElementsMatch(t, []int{3, 6, 9}, got)

	bad, err := dj.GoMapWindowErr(context.Background(), dj.RangeN(5), 3, func(_ context.Context, window []int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

func TestGoMapWindowIdx(t *testing.T) {
	got := dj.GoMapWindowIdx(context.Background(), dj.RangeN(5), 3, func(_ context.Context, idx int, window []int) int {
		return idx + window[0] + window[1] + window[2]
	})

	require.ElementsMatch(t, []int{3, 7, 11}, got)
}

func TestGoMapWindowIdxErr(t *testing.T) {
	got, err := dj.GoMapWindowIdxErr(context.Background(), dj.RangeN(5), 3, func(_ context.Context, idx int, window []int) (int, error) {
		return idx + window[0] + window[1] + window[2], nil
	})
	require.NoError(t, err)
	require.ElementsMatch(t, []int{3, 7, 11}, got)

	bad, err := dj.GoMapWindowIdxErr(context.Background(), dj.RangeN(5), 3, func(_ context.Context, idx int, window []int) (int, error) {
		return 0, errors.New("oops")
	})
	require.Error(t, err)
	require.Nil(t, bad)
}

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
			require.Equal(t, tt.want, dj.Concat(tt.in...))
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
			require.ElementsMatch(t, tt.want, dj.Perms(tt.in))
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
			require.ElementsMatch(t, tt.want, dj.PermsIdx(tt.in))
		})
	}
}

func TestShuffle(t *testing.T) {
	for _, slice := range dj.Perms(dj.RangeN(6)) {
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
