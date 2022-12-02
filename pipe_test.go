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

func TestPipe_Iterator(t *testing.T) {
	in, out := xn.NewPipe[int]()

	in <- 1
	in <- 2
	in <- 3

	close(in)

	require.Equal(t, []int{1, 2, 3}, xn.ChanIter(out).Collect())
}

func TestPipe_Forward(t *testing.T) {
	in1, out1 := xn.NewPipe[int]()
	in2, out2 := xn.NewPipe[int]()

	// Write data into the first pipe's input channel.
	xn.ChanWriter(in1).WriteFrom(xn.SliceIter(1, 2, 3))

	// close the first pipe's input channel;
	// the written data is still available in the output channel.
	close(in1)

	// Forward the data from the first pipe to the second.
	// Close the second pipe's input channel when finished.
	go func() {
		defer close(in2)
		xn.ForwardChan([]<-chan int{out1}, []chan<- int{in2})
	}()

	// Read data from the second pipe's output channel.
	require.Equal(t, []int{1, 2, 3}, xn.ChanIter(out2).Collect())
}
