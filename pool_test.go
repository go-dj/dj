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

	// Create a pool with 2 workers which converts strings to ints.
	pool := dj.NewWorkerPool(ctx, 2, func(ctx context.Context, in string) (int, error) {
		return strconv.Atoi(in)
	})

	// Submit 100 jobs to the pool.
	jobs := dj.MapN(100, func(i int) dj.Job[string, int] {
		return pool.Submit(ctx, strconv.Itoa(i))
	})

	// Each job should be processed.
	dj.ForEachIdx(jobs, func(i int, job dj.Job[string, int]) {
		v, err := job.R()
		require.NoError(t, err)
		require.Equal(t, i, v)
	})
}

func TestWorkerPool_Submit_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a pool with 2 workers which converts strings to ints.
	pool := dj.NewWorkerPool(ctx, 2, func(ctx context.Context, in string) (int, error) {
		return strconv.Atoi(in)
	})

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

	// Create a pool with 2 workers which converts strings to ints.
	pool := dj.NewWorkerPool(ctx, 2, func(ctx context.Context, in string) (int, error) {
		return strconv.Atoi(in)
	})

	// Process 100 jobs concurrently.
	res, err := pool.Process(ctx, dj.MapN(100, func(i int) string { return strconv.Itoa(i) })...)
	require.NoError(t, err)

	// Each job should be processed.
	dj.ForEachIdx(res, func(i int, v int) {
		require.Equal(t, i, v)
	})
}

func TestWorkerPool_Process_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a pool with 2 workers which converts strings to ints.
	pool := dj.NewWorkerPool(ctx, 2, func(ctx context.Context, in string) (int, error) {
		return strconv.Atoi(in)
	})

	// Attempt to process some jobs which will fail.
	res, err := pool.Process(ctx, dj.MapN(100, func(i int) string { return "foo" + strconv.Itoa(i) })...)
	require.Error(t, err)
	require.Empty(t, res)
}
