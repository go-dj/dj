package dj_test

import (
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestRWBuffer_Readable(t *testing.T) {
	buf := dj.NewRWBuffer([]int{1, 2, 3})

	got, err := dj.NewIter[int](buf).Collect()
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, got)
}

func TestRWBuffer_Writable(t *testing.T) {
	buf := dj.NewRWBuffer[int](nil)

	n, err := dj.NewWriter[int](buf).WriteFrom(dj.SliceIter(1, 2, 3))
	require.NoError(t, err)
	require.Equal(t, 3, n)
}

func TestRWBuffer_RW(t *testing.T) {
	buf := dj.NewRWBuffer[int](nil)

	// Create a reader and writer from the buffer.
	r := dj.NewIter[int](buf)
	w := dj.NewWriter[int](buf)

	// Write some data.
	n, err := w.WriteFrom(dj.SliceIter(1, 2, 3))
	require.NoError(t, err)
	require.Equal(t, 3, n)

	// Read it back.
	got, err := r.Collect()
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, got)

	// The buffer should be empty now.
	empty, err := r.Collect()
	require.NoError(t, err)
	require.Equal(t, []int{}, empty)

	// Write some more data.
	n, err = w.WriteFrom(dj.SliceIter(4, 5, 6))
	require.NoError(t, err)
	require.Equal(t, 3, n)

	// Read it back.
	more, err := r.Collect()
	require.NoError(t, err)
	require.Equal(t, []int{4, 5, 6}, more)
}
