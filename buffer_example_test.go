package dj_test

import (
	"fmt"

	"github.com/jameshoulahan/dj"
)

func ExampleRWBuffer() {
	// Create a new buffer for ints.
	buf := dj.NewRWBuffer[int](nil)

	// Create a reader and writer over the buffer.
	r := dj.NewIter[int](buf)
	w := dj.NewWriter[int](buf)

	// Write some data.
	n, ok := w.WriteFrom(dj.SliceIter(1, 2, 3))
	fmt.Println("Wrote", n, "values:", ok)

	// Read it back.
	fmt.Println("Collected:", r.Collect())

	// The buffer should be empty now.
	fmt.Println("Buffer is empty:", r.Collect())

	// Write some more data.
	n, ok = w.WriteFrom(dj.SliceIter(4, 5, 6))
	fmt.Println("Wrote", n, "values:", ok)

	// Read it back.
	fmt.Println("Collected:", r.Collect())

	// Output:
	// Wrote 3 values: true
	// Collected: [1 2 3]
	// Buffer is empty: []
	// Wrote 3 values: true
	// Collected: [4 5 6]
}
