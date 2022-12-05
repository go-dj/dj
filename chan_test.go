package dj_test

import (
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestForChan(t *testing.T) {
	var got []int

	dj.ForChan(newCh(1, 2, 3), func(v int) {
		got = append(got, v)
	})

	require.Equal(t, []int{1, 2, 3}, got)
}

func TestCollectChan(t *testing.T) {
	require.Equal(t, []int{1, 2, 3}, dj.CollectChan(newCh(1, 2, 3)))
}

func TestReadChan(t *testing.T) {
	ch := newCh(1, 2, 3, 4, 5)

	require.Equal(t, []int{1, 2}, dj.TakeChan(ch, 2))
	require.Equal(t, []int{3, 4}, dj.TakeChan(ch, 2))
	require.Equal(t, []int{5}, dj.TakeChan(ch, 2))
}

func TestForwarChan(t *testing.T) {
	// 3 senders of 3 ints each.
	src := dj.MapN(3, func(i int) <-chan int {
		return newCh(dj.Range(3*i, 3*(i+1))...)
	})

	// 3 receivers.
	dst := dj.MapN(3, func(i int) chan int {
		return make(chan int)
	})

	// Forward the src channels to the dst channels.
	go func() {
		defer dj.CloseChan(dj.AsSend(dst...)...)
		dj.ForwardChan(src, dj.AsSend(dst...))
	}()

	// Collect the results.
	require.ElementsMatch(t, dj.RangeN(9), dj.CollectChan(dj.FanIn(dj.AsRecv(dst...)...)))
}

func TestJoinChan(t *testing.T) {
	ch1 := newCh(1, 2, 3)
	ch2 := newCh(4, 5, 6)

	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, dj.CollectChan(dj.ConcatChan(ch1, ch2)))
}

func TestZipChan(t *testing.T) {
	ch1 := newCh(1, 2, 3)
	ch2 := newCh(4, 5, 6)

	require.Equal(t, [][]int{{1, 4}, {2, 5}, {3, 6}}, dj.CollectChan(dj.ZipChan(ch1, ch2)))
}

func TestMergeChan(t *testing.T) {
	ch1 := newCh(1, 2, 3)
	ch2 := newCh(4, 5, 6)

	require.ElementsMatch(t, []int{1, 2, 3, 4, 5, 6}, dj.CollectChan(dj.FanIn(ch1, ch2)))
}

func TestSplitChan(t *testing.T) {
	chs := dj.FanOut(newCh(dj.RangeN(1000)...), 4)

	dj.ForEachIdx(chs, func(idx int, ch <-chan int) {
		require.Equal(t, dj.Range(250*idx, 250*(idx+1)), dj.TakeChan(ch, 250))
	})
}

func newCh(vals ...int) <-chan int {
	ch := make(chan int, len(vals))

	for _, v := range vals {
		ch <- v
	}

	close(ch)

	return ch
}
