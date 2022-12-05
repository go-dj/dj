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
	group.Go(10, func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			require.Fail(t, "the context should not be canceled")

		default:
			// ...
		}
	}).Wait()

	// Now cancel the context.
	cancel()

	// Start 10 more goroutines in the canceled context.
	group.Go(10, func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			// ...

		default:
			require.Fail(t, "the context should be canceled")
		}
	}).Wait()
}

func TestGroup_Child(t *testing.T) {
	// Create a parent group.
	parent := dj.NewGroup(context.Background(), dj.NewSem(1))

	// Create a child group inside the parent group.
	child := parent.Child()

	// Start 10 goroutines in the child group.
	child.Go(10, func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			require.Fail(t, "the context should not be canceled")

		default:
			// ...
		}
	}).Wait()

	// Cancel the child group.
	child.Cancel()

	// Start 10 more goroutines in the child group.
	child.Go(10, func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			// ...

		default:
			require.Fail(t, "the context should be canceled")
		}
	}).Wait()

	// Start 10 more goroutines in the parent group.
	parent.Go(10, func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			require.Fail(t, "the context should not be canceled")

		default:
			// ...
		}
	}).Wait()
}
