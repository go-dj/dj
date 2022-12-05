package dj_test

import (
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := dj.NewQueue(2, 3, 4)

	q.PushBack(5, 6)
	q.PushFront(0, 1)

	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6}, q.Items())
}

func TestQueue_PopFront(t *testing.T) {
	q := dj.NewQueue(1, 2, 3)

	v, ok := q.PopFront()
	require.True(t, ok)
	require.Equal(t, 1, v)

	v, ok = q.PopFront()
	require.True(t, ok)
	require.Equal(t, 2, v)

	v, ok = q.PopFront()
	require.True(t, ok)
	require.Equal(t, 3, v)

	v, ok = q.PopFront()
	require.False(t, ok)
	require.Equal(t, 0, v)
}

func TestQueue_PopBack(t *testing.T) {
	q := dj.NewQueue(1, 2, 3)

	v, ok := q.PopBack()
	require.True(t, ok)
	require.Equal(t, 3, v)

	v, ok = q.PopBack()
	require.True(t, ok)
	require.Equal(t, 2, v)

	v, ok = q.PopBack()
	require.True(t, ok)
	require.Equal(t, 1, v)

	v, ok = q.PopBack()
	require.False(t, ok)
	require.Equal(t, 0, v)
}

func TestQueue_Len(t *testing.T) {
	q := dj.NewQueue(1, 2, 3)

	q.PopBack()
	require.Equal(t, 2, q.Len())

	q.PopBack()
	require.Equal(t, 1, q.Len())

	q.PopBack()
	require.Equal(t, 0, q.Len())
}

func TestMap(t *testing.T) {
	m := dj.NewMap[string, int]()

	m.Set("foo", 1)
	m.Set("bar", 2)

	v1, ok := m.Get("foo")
	require.True(t, ok)
	require.Equal(t, 1, v1)

	v2, ok := m.Get("bar")
	require.True(t, ok)
	require.Equal(t, 2, v2)
}

func TestMap_Delete(t *testing.T) {
	m := dj.NewMap[string, int]()

	m.Set("foo", 1)
	m.Set("bar", 2)
	m.Delete("foo")

	_, ok := m.Get("foo")
	require.False(t, ok)

	v, ok := m.Get("bar")
	require.True(t, ok)
	require.Equal(t, 2, v)
}

func TestMap_Items(t *testing.T) {
	m := dj.NewMap[string, int]()

	m.Set("foo", 1)
	m.Set("bar", 2)

	require.Equal(t, map[string]int{"foo": 1, "bar": 2}, m.Items())
}
