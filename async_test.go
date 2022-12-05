package dj_test

import (
	"context"
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestGroup_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := dj.NewGroup(ctx, dj.NewSem(1))
	defer group.Wait()

	// Start 10 goroutines.
	group.GoN(10, func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			require.Fail(t, "should not run")

		default:
			// ...
		}
	})

	// Wait for the goroutines to execute.
	group.Wait()

	// Now cancel the context.
	cancel()

	// Start 10 more goroutines in the canceled context.
	group.GoN(10, func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			// ...

		default:
			require.Fail(t, "should not run")
		}
	})
}
