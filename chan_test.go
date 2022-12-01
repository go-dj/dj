package xn_test

import (
	"testing"

	"github.com/jameshoulahan/xn"
	"github.com/stretchr/testify/require"
)

func Test_ForChan(t *testing.T) {
	var got []int

	xn.ForChan(newCh(1, 2, 3), func(v int) {
		got = append(got, v)
	})

	require.Equal(t, []int{1, 2, 3}, got)
}

func Test_CollectChan(t *testing.T) {
	require.Equal(t, []int{1, 2, 3}, xn.CollectChan(newCh(1, 2, 3)))
}

func Test_ReadChan(t *testing.T) {
	ch := newCh(1, 2, 3, 4, 5)

	require.Equal(t, []int{1, 2}, xn.TakeChan(ch, 2))
	require.Equal(t, []int{3, 4}, xn.TakeChan(ch, 2))
	require.Equal(t, []int{5}, xn.TakeChan(ch, 2))
}

func Test_ForwarChan(t *testing.T) {
	// 3 senders of 3 ints each.
	src := xn.MapN(3, func(i int) <-chan int {
		return newCh(xn.Range(3*i, 3*(i+1))...)
	})

	// 3 receivers.
	dst := xn.MapN(3, func(i int) chan int {
		return make(chan int)
	})

	// Forward the src channels to the dst channels.
	go func() {
		defer xn.CloseChan(xn.ToSend(dst...)...)
		xn.ForwardChan(src, xn.ToSend(dst...))
	}()

	// Collect the results.
	require.ElementsMatch(t, xn.RangeN(9), xn.CollectChan(xn.MergeChan(xn.ToRecv(dst...)...)))
}

func Test_JoinChan(t *testing.T) {
	ch1 := newCh(1, 2, 3)
	ch2 := newCh(4, 5, 6)

	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, xn.CollectChan(xn.JoinChan(ch1, ch2)))
}

func Test_ZipChan(t *testing.T) {
	ch1 := newCh(1, 2, 3)
	ch2 := newCh(4, 5, 6)

	require.Equal(t, [][]int{{1, 4}, {2, 5}, {3, 6}}, xn.CollectChan(xn.ZipChan(ch1, ch2)))
}

func Test_MergeChan(t *testing.T) {
	ch1 := newCh(1, 2, 3)
	ch2 := newCh(4, 5, 6)

	require.ElementsMatch(t, []int{1, 2, 3, 4, 5, 6}, xn.CollectChan(xn.MergeChan(ch1, ch2)))
}

func Test_SplitChan(t *testing.T) {
	chs := xn.SplitChan(newCh(xn.RangeN(1000)...), 4)

	xn.ForEachIdx(chs, func(idx int, ch <-chan int) {
		require.Equal(t, xn.Range(250*idx, 250*(idx+1)), xn.TakeChan(ch, 250))
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
