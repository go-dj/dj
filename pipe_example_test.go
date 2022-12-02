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
