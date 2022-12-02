package dj_test

import (
	"fmt"

	"github.com/go-dj/dj"
)

func ExampleNewPipe() {
	// Create a new pipe for ints.
	in, out := dj.NewPipe[int]()

	// Write some values to the pipe.
	// Writes never block; they are buffered until they are read.
	go func() {
		in <- 1
		in <- 2
		in <- 3

		// Close the pipe to signal that no more values will be written.
		// Buffered values can still be read.
		close(in)
	}()

	// Read the values from the pipe.
	for v := range out {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 2
	// 3
}

func ExampleNewBufPipe() {
	// Create a new pipe for ints with a buffer of size 3.
	in, out := dj.NewBufPipe[int](3)

	// Write some values to the pipe.
	in <- 1
	in <- 2
	in <- 3

	// The buffer is full, so the write blocks.
	select {
	case in <- 4:
		panic("write should block")

	default:
		// ...
	}

	// Read the values from the pipe.
	fmt.Println(<-out) // 1
	fmt.Println(<-out) // 2
	fmt.Println(<-out) // 3

	// The buffer is empty, so the read blocks.
	select {
	case <-out:
		panic("read should block")

	default:
		// ...
	}

	// Output:
	// 1
	// 2
	// 3
}
