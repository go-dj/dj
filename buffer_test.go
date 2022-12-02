package xn_test

import (
	"testing"

	"github.com/jameshoulahan/xn"
	"github.com/stretchr/testify/require"
)

func TestRWBuffer_Iterable(t *testing.T) {
	buf := xn.NewRWBuffer(1, 2, 3)

	require.Equal(t, []int{1, 2, 3}, xn.NewIter[int](buf).Collect())
}

func TestRWBuffer_Writable(t *testing.T) {
	buf := xn.NewRWBuffer[int]()

	n, ok := xn.NewWriter[int](buf).WriteFrom(xn.SliceIter(1, 2, 3))
	require.True(t, ok)
	require.Equal(t, 3, n)
}

func TestRWBuffer_RW(t *testing.T) {
	buf := xn.NewRWBuffer[int]()

	// Create a reader and writer from the buffer.
	r := xn.NewIter[int](buf)
	w := xn.NewWriter[int](buf)

	// Write some data.
	n, ok := w.WriteFrom(xn.SliceIter(1, 2, 3))
	require.True(t, ok)
	require.Equal(t, 3, n)

	// Read it back.
	require.Equal(t, []int{1, 2, 3}, r.Collect())

	// The buffer should be empty now.
	require.Equal(t, []int{}, r.Collect())

	// Write some more data.
	n, ok = w.WriteFrom(xn.SliceIter(4, 5, 6))
	require.True(t, ok)
	require.Equal(t, 3, n)

	// Read it back.
	require.Equal(t, []int{4, 5, 6}, r.Collect())
}
