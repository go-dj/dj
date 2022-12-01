package xn_test

import (
	"testing"

	"github.com/jameshoulahan/xn"
	"github.com/stretchr/testify/require"
)

func Test_ForChan(t *testing.T) {
	var got []int

	xn.ForChan(newTestChan(1, 2, 3), func(v int) {
		got = append(got, v)
	})

	require.Equal(t, []int{1, 2, 3}, got)
}

func Test_CollectChan(t *testing.T) {
	require.Equal(t, []int{1, 2, 3}, xn.CollectChan(newTestChan(1, 2, 3)))
}

func Test_ReadChan(t *testing.T) {
	ch := newTestChan(1, 2, 3, 4, 5)

	require.Equal(t, []int{1, 2}, xn.ReadChan(ch, 2))
	require.Equal(t, []int{3, 4}, xn.ReadChan(ch, 2))
	require.Equal(t, []int{5}, xn.ReadChan(ch, 2))
}

func Test_ForwardChan(t *testing.T) {
	ch := make(chan int, 3)

	xn.ForwardChan(newTestChan(1, 2, 3), ch)

	require.Equal(t, []int{1, 2, 3}, xn.ReadChan(ch, 3))
}

func newTestChan(vals ...int) <-chan int {
	ch := make(chan int, len(vals))

	for _, v := range vals {
		ch <- v
	}

	close(ch)

	return ch
}
