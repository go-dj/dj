package dj_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestWorkerPool_Submit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new pool.
	pool := newTestPool(ctx, 2)

	// Submit 100 jobs to the pool.
	jobs := dj.MapN(100, func(i int) dj.Job[string, int] {
		return pool.Submit(ctx, strconv.Itoa(i))
	})

	// Each job should be processed.
	dj.ForIdx(jobs, func(i int, job dj.Job[string, int]) {
		v, err := job.R()
		require.NoError(t, err)
		require.Equal(t, i, v)
	})
}

func TestWorkerPool_Submit_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new pool.
	pool := newTestPool(ctx, 2)

	// Submit a job which will fail.
	job := pool.Submit(ctx, "foo")

	// The job should fail.
	v, err := job.R()
	require.Error(t, err)
	require.Zero(t, v)
}

func TestWorkerPool_Process(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new pool.
	pool := newTestPool(ctx, 2)

	// Process 100 jobs concurrently.
	res, err := pool.Process(ctx, dj.MapN(100, func(i int) string { return strconv.Itoa(i) })...)
	require.NoError(t, err)

	// Each job should be processed.
	dj.ForIdx(res, func(i int, v int) {
		require.Equal(t, i, v)
	})
}

func TestWorkerPool_Process_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new pool.
	pool := newTestPool(ctx, 2)

	// Attempt to process some jobs which will fail.
	res, err := pool.Process(ctx, dj.MapN(100, func(i int) string { return "foo" + strconv.Itoa(i) })...)
	require.Error(t, err)
	require.Empty(t, res)
}

func TestWorkerPool_JobContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new pool.
	pool := newTestPool(ctx, 2)

	// Submit a job which will fail.
	job := pool.Submit(ctx, "foo")

	// The job should fail.
	v, err := job.R()
	require.Error(t, err)
	require.Zero(t, v)
}

func TestWorkerPool_Close(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new pool.
	pool := newTestPool(ctx, 2)

	// Before the context is cancelled, the pool should accept jobs.
	res, err := pool.Process(context.Background(), "1", "2", "3")
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, res)

	// Cancel the context.
	cancel()

	// After the context is cancelled, the pool should reject jobs.
	res, err = pool.Process(context.Background(), "1", "2", "3")
	require.Error(t, err)
	require.Empty(t, res)
}

// newTestPool creates a new pool with the given number of workers.
// The pool converts strings to ints.
func newTestPool(ctx context.Context, n int) *dj.WorkerPool[string, int] {
	return dj.NewWorkerPool(ctx, n, func(ctx context.Context, in string) (int, error) {
		return strconv.Atoi(in)
	})
}
