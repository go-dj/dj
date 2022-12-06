package dj_test

import (
	"fmt"

	"github.com/go-dj/dj"
)

func ExampleRWBuffer() {
	// Create a new buffer for ints.
	buf := dj.NewRWBuffer[int](nil)

	// Create a reader and writer over the buffer.
	r := dj.NewIter[int](buf)
	w := dj.NewWriter[int](buf)

	// Write some data.
	n, err := w.WriteFrom(dj.SliceIter(1, 2, 3))
	fmt.Println("Wrote", n, "values:", err)

	// Read it back.
	got, err := r.Collect()
	fmt.Println("Collected:", got, err)

	// The buffer should be empty now.
	empty, err := r.Collect()
	fmt.Println("Buffer is empty:", empty, err)

	// Write some more data.
	n, err = w.WriteFrom(dj.SliceIter(4, 5, 6))
	fmt.Println("Wrote", n, "values:", err)

	// Read it back.
	more, err := r.Collect()
	fmt.Println("Collected:", more, err)

	// Output:
	// Wrote 3 values: <nil>
	// Collected: [1 2 3] <nil>
	// Buffer is empty: [] <nil>
	// Wrote 3 values: <nil>
	// Collected: [4 5 6] <nil>
}
