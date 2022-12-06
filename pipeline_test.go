package dj_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {
	inCh, outCh := dj.NewPipeline(context.Background(), 2, func(ctx context.Context, v int) dj.Result[int] {
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

	require.Equal(t, 2, (<-outCh).Value())
	require.Equal(t, 4, (<-outCh).Value())
}

func TestPipeline_Error(t *testing.T) {
	inCh, outCh := dj.NewPipeline(context.Background(), 2, func(ctx context.Context, v string) dj.Result[int] {
		return dj.NewResult(strconv.Atoi(v))
	})

	inCh <- "1"
	inCh <- "two"

	require.NoError(t, (<-outCh).Error())
	require.Error(t, (<-outCh).Error())
}
