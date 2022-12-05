package dj_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestWithParallelism(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var n atomic.Int32

	dj.GoForN(dj.WithParallelism(ctx, 10), 20, func(context.Context, int) {
		n.Add(1)
		defer n.Add(-1)

		require.LessOrEqual(t, n.Load(), int32(10))

		time.Sleep(time.Millisecond)
	})
}
