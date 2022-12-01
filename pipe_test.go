package xn_test

import (
	"testing"

	"github.com/jameshoulahan/xn"
	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {
	in, out := xn.NewPipe[int]()

	in <- 1
	in <- 2
	in <- 3

	require.Equal(t, 1, <-out)
	require.Equal(t, 2, <-out)
	require.Equal(t, 3, <-out)
}

func TestPipe_Large(t *testing.T) {
	in, out := xn.NewPipe[int]()

	xn.ForN(1000, func(i int) {
		in <- i
	})

	for i := 0; i < 1000; i++ {
		require.Equal(t, i, <-out)
	}
}

func TestPipe_Close(t *testing.T) {
	in, out := xn.NewPipe[int]()

	in <- 1
	in <- 2
	in <- 3

	close(in)

	require.Equal(t, []int{1, 2, 3}, xn.CollectChan(out))
}
