package dj_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {
	inCh, outCh := dj.NewPipeline(context.Background(), 1, 2, func(ctx context.Context, v int) dj.Result[int] {
		return dj.Ok(v * 2)
	})

	inCh <- 1
	inCh <- 2

	// The pipeline is full, so the next write blocks.
	select {
	case inCh <- 3:
		t.Fatal("expected write to block")

	default:
		// ...
	}

	require.Equal(t, 2, (<-outCh).Val())
	require.Equal(t, 4, (<-outCh).Val())
}

func TestPipeline_NoBuffer(t *testing.T) {
	inCh, outCh := dj.NewPipeline(context.Background(), 1, -1, func(ctx context.Context, v int) dj.Result[int] {
		return dj.Ok(v * 2)
	})

	// No buffer, so writes never block.
	dj.ForN(10, func(i int) {
		inCh <- i
	})

	dj.ForN(10, func(i int) {
		require.Equal(t, i*2, (<-outCh).Val())
	})
}

func TestPipeline_Error(t *testing.T) {
	inCh, outCh := dj.NewPipeline(context.Background(), 1, 2, func(ctx context.Context, v string) dj.Result[int] {
		return dj.NewResult(strconv.Atoi(v))
	})

	inCh <- "1"
	inCh <- "two"

	require.NoError(t, (<-outCh).Err())
	require.Error(t, (<-outCh).Err())
}
