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
