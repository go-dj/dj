package dj_test

import (
	"testing"

	"github.com/jameshoulahan/dj"
	"github.com/stretchr/testify/require"
)

func TestRWBuffer_Iterable(t *testing.T) {
	buf := dj.NewRWBuffer([]int{1, 2, 3})

	require.Equal(t, []int{1, 2, 3}, dj.NewIter[int](buf).Collect())
}

func TestRWBuffer_Writable(t *testing.T) {
	buf := dj.NewRWBuffer[int](nil)

	n, ok := dj.NewWriter[int](buf).WriteFrom(dj.SliceIter(1, 2, 3))
	require.True(t, ok)
	require.Equal(t, 3, n)
}

func TestRWBuffer_RW(t *testing.T) {
	buf := dj.NewRWBuffer[int](nil)

	// Create a reader and writer from the buffer.
	r := dj.NewIter[int](buf)
	w := dj.NewWriter[int](buf)

	// Write some data.
	n, ok := w.WriteFrom(dj.SliceIter(1, 2, 3))
	require.True(t, ok)
	require.Equal(t, 3, n)

	// Read it back.
	require.Equal(t, []int{1, 2, 3}, r.Collect())

	// The buffer should be empty now.
	require.Equal(t, []int{}, r.Collect())

	// Write some more data.
	n, ok = w.WriteFrom(dj.SliceIter(4, 5, 6))
	require.True(t, ok)
	require.Equal(t, 3, n)

	// Read it back.
	require.Equal(t, []int{4, 5, 6}, r.Collect())
}
