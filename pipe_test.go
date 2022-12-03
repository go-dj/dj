package dj_test

import (
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {
	in, out := dj.NewPipe[int]()
	defer close(in)

	in <- 1
	in <- 2
	in <- 3

	require.Equal(t, 1, <-out)
	require.Equal(t, 2, <-out)
	require.Equal(t, 3, <-out)
}

func TestPipe_Large(t *testing.T) {
	in, out := dj.NewPipe[int]()
	defer close(in)

	dj.ForN(100, func(i int) {
		in <- i
	})

	dj.ForN(100, func(i int) {
		require.Equal(t, i, <-out)
	})
}

func TestPipe_Close(t *testing.T) {
	in, out := dj.NewPipe[int]()

	dj.ForN(100, func(i int) {
		in <- i
	})

	close(in)

	require.Equal(t, dj.RangeN(100), dj.CollectChan(out))
}

func TestPipe_Iterator(t *testing.T) {
	in, out := dj.NewPipe[int]()
	defer close(in)

	dj.ForN(100, func(i int) {
		in <- i
	})

	require.Equal(t, dj.RangeN(100), dj.ChanIter(out).Take(100))
}

func TestPipe_Forward(t *testing.T) {
	type pipe struct {
		in  chan<- int
		out <-chan int
	}

	// Create 100 pipes.
	pipes := dj.MapN(100, func(int) pipe {
		in, out := dj.NewPipe[int]()
		return pipe{in, out}
	})

	// Chain the pipes together.
	dj.ForWindow(pipes, 2, func(pair []pipe) {
		go func() {
			defer close(pair[1].in)
			dj.ForwardChan([]<-chan int{pair[0].out}, []chan<- int{pair[1].in})
		}()
	})

	// Write data into the first pipe's input channel.
	pipes[0].in <- 1
	pipes[0].in <- 2
	pipes[0].in <- 3

	// It should be available in the last pipe's output channel.
	require.Equal(t, 1, <-dj.Last(pipes).out)
	require.Equal(t, 2, <-dj.Last(pipes).out)
	require.Equal(t, 3, <-dj.Last(pipes).out)

	// Add some more data.
	pipes[0].in <- 4
	pipes[0].in <- 5
	pipes[0].in <- 6

	// Close the first pipe's input channel.
	close(pipes[0].in)

	// Read data from the last pipe's output channel.
	require.Equal(t, []int{4, 5, 6}, dj.CollectChan(dj.Last(pipes).out))
}

func TestPipe_Read_Block(t *testing.T) {
	in, out := dj.NewPipe[int]()
	defer close(in)

	go func() {
		in <- 1
		in <- 2
		in <- 3
	}()

	require.Equal(t, 1, <-out)
	require.Equal(t, 2, <-out)
	require.Equal(t, 3, <-out)
}

func TestPipe_Write_Block(t *testing.T) {
	in, out := dj.NewPipe[int]()
	defer close(in)

	go func() {
		require.Equal(t, 1, <-out)
		require.Equal(t, 2, <-out)
		require.Equal(t, 3, <-out)
	}()

	in <- 1
	in <- 2
	in <- 3
}

func TestBufPipe_ReadWrite_Block(t *testing.T) {
	in, out := dj.NewBufPipe[int](3)
	defer close(in)

	in <- 1
	in <- 2
	in <- 3

	// The buffer is full, so the write blocks.
	select {
	case in <- 4:
		t.Fatal("write should block")

	default:
		// ...
	}

	require.Equal(t, 1, <-out)
	require.Equal(t, 2, <-out)
	require.Equal(t, 3, <-out)

	// The buffer is empty, so the read blocks.
	select {
	case <-out:
		t.Fatal("read should block")

	default:
		// ...
	}
}
