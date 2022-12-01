package xn

import (
	"context"
)

// CollectChan reads all values from the channel.
func CollectChan[T any](ch <-chan T) []T {
	return CollectChanCtx(context.Background(), ch)
}

// CollectChanCtx reads all values from the channel.
// It stops reading when the context is canceled.
func CollectChanCtx[T any](ctx context.Context, ch <-chan T) []T {
	var out []T

	ForChanCtx(ctx, ch, func(_ context.Context, v T) {
		out = append(out, v)
	})

	return out
}

// ReadChan reads up to n values from the channel.
func ReadChan[T any](ch <-chan T, n int) []T {
	return ReadChanCtx(context.Background(), ch, n)
}

// ReadChanCtx reads up to n values from the channel.
// It stops reading when the context is canceled.
func ReadChanCtx[T any](ctx context.Context, ch <-chan T, n int) []T {
	out := make([]T, 0, n)

	for {
		select {
		case <-ctx.Done():
			return out

		case v, ok := <-ch:
			if !ok {
				return out
			}

			out = append(out, v)
		}

		if len(out) == n {
			return out
		}
	}
}

// ForChan calls the given function for each value in the channel.
func ForChan[T any](ch <-chan T, fn func(T)) {
	ForChanCtx(context.Background(), ch, func(_ context.Context, v T) {
		fn(v)
	})
}

// ForChanCtx calls the given function for each value in the channel.
// It stops iterating when the context is canceled.
func ForChanCtx[T any](ctx context.Context, ch <-chan T, fn func(context.Context, T)) {
	for {
		select {
		case <-ctx.Done():
			return

		case v, ok := <-ch:
			if !ok {
				return
			}

			fn(ctx, v)
		}
	}
}

// ForwardChan forwards values from the src channel to the dst channel.
func ForwardChan[T any](src <-chan T, dst chan<- T) {
	ForwardChanCtx(context.Background(), src, dst)
}

// ForwardChanCtx forwards values from the src channel to the dst channel.
// It stops forwarding when the context is canceled.
func ForwardChanCtx[T any](ctx context.Context, src <-chan T, dst chan<- T) {
	go ForChanCtx(ctx, src, func(ctx context.Context, v T) {
		select {
		case <-ctx.Done():
			// Don't forward the value if the context is canceled.

		case dst <- v:
			// Forward the value.
		}
	})
}

// JoinChan joins the given channels into a single channel.
// That is, the first channel is emptied, then the second, and so on.
// The returned channel is closed when all the given channels are closed.
func JoinChan[T any](chs ...<-chan T) <-chan T {
	return JoinChanCtx(context.Background(), chs...)
}

// JoinChanCtx joins the given channels into a single channel.
// That is, the first channel is emptied, then the second, and so on.
// The returned channel is closed when all the given channels are closed or the context is canceled.
func JoinChanCtx[T any](ctx context.Context, chs ...<-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for _, ch := range chs {
			ForwardChanCtx(ctx, ch, out)
		}
	}()

	return out
}

// MergeChan merges the given channels into a single channel.
// That is, values are read from each channel in parallel.
// The returned channel is closed when all the given channels are closed.
func MergeChan[T any](chs ...<-chan T) <-chan T {
	return MergeChanCtx(context.Background(), chs...)
}

// MergeChanCtx merges the given channels into a single channel.
// That is, values are read from each channel in parallel.
// The returned channel is closed when all the given channels are closed or the context is canceled.
func MergeChanCtx[T any](ctx context.Context, chs ...<-chan T) <-chan T {
	sem := NewSem(ctx, len(chs))
	out := make(chan T)

	go func() {
		defer close(out)
		defer sem.Wait()

		for _, ch := range chs {
			sem.Go(func(ctx context.Context) {
				ForwardChanCtx(ctx, ch, out)
			})
		}
	}()

	return out
}
