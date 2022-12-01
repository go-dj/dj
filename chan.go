package xn

import (
	"context"
	"reflect"
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

// TakeChan reads up to n values from the channel.
func TakeChan[T any](ch <-chan T, n int) []T {
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

		ForEach(chs, func(ch <-chan T) {
			sem.Go(func(ctx context.Context) {
				ForChanCtx(ctx, ch, func(ctx context.Context, v T) {
					select {
					case <-ctx.Done():
					case out <- v:
					}
				})
			})
		})
	}()

	return out
}

// ForwardChan forwards values from the src channel(s) to the dst channel(s).
func ForwardChan[T any](src []<-chan T, dst []chan<- T) {
	ForwardChanCtx(context.Background(), src, dst)
}

// ForwardChanCtx forwards values from the src channel(s) to the dst channel(s).
// It stops forwarding when the context is canceled.
func ForwardChanCtx[T any](ctx context.Context, src []<-chan T, dst []chan<- T) {
	ForChanCtx(ctx, MergeChanCtx(ctx, src...), func(ctx context.Context, v T) {
		SendChanCtx(ctx, v, dst...)
	})
}

// SplitChan splits the given channel into n channels.
// The returned channels are closed when the given channel is closed.
func SplitChan[T any](ch <-chan T, n int) []<-chan T {
	return SplitChanCtx(context.Background(), ch, n)
}

// SplitChanCtx splits the given channel into n channels.
// The returned channels are closed when the given channel is closed or the context is canceled.
func SplitChanCtx[T any](ctx context.Context, ch <-chan T, n int) []<-chan T {
	out := MapN(n, func(int) chan T {
		return make(chan T)
	})

	go func() {
		defer CloseChan(ToSend(out...)...)
		ForwardChanCtx(ctx, []<-chan T{ch}, ToSend(out...))
	}()

	return ToRecv(out...)
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
			ForwardChanCtx(ctx, []<-chan T{ch}, ToSend(out))
		}
	}()

	return out
}

// ZipChan zips the given channels into a single channel.
// The returned channel is closed when all the given channels are closed.
func ZipChan[T any](chs ...<-chan T) <-chan []T {
	return ZipChanCtx(context.Background(), chs...)
}

// ZipChanCtx zips the given channels into a single channel.
// The returned channel is closed when all the given channels are closed or the context is canceled.
func ZipChanCtx[T any](ctx context.Context, chs ...<-chan T) <-chan []T {
	out := make(chan []T)

	go func() {
		defer close(out)

		in := ZipIter(MapEach(chs, func(ch <-chan T) Iter[T] {
			return ChanIterCtx(ctx, ch)
		})...).Chan()

		ForwardChanCtx(ctx, []<-chan []T{in}, ToSend(out))
	}()

	return out
}

// Send sends the given value to one of the given channels.
// It blocks until the value is sent.
func SendChan[T any](v T, chs ...chan<- T) {
	SendChanCtx(context.Background(), v, chs...)
}

// SendCtx sends the given value to one of the given channels.
// It blocks until the value is sent or the context is canceled.
func SendChanCtx[T any](ctx context.Context, v T, chs ...chan<- T) {
	reflect.Select(append(
		MapEach(chs, func(ch chan<- T) reflect.SelectCase {
			return reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(ch),
				Send: reflect.ValueOf(v),
			}
		}),
		reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		}),
	)
}

// Recv receives a value from one of the given channels.
// It blocks until a value is received, which it returns.
// The boolean indicates whether the read was successful; it is false if the channel is closed.
func RecvChan[T any](chs ...<-chan T) (T, bool) {
	_, v, ok := reflect.Select(MapEach(chs, func(ch <-chan T) reflect.SelectCase {
		return reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}))

	return v.Interface().(T), ok
}

// RecvChanCtx receives a value from one of the given channels.
// It blocks until a value is received, which it returns, or the context is canceled.
// The boolean indicates whether the read was successful; it is false if the channel is closed.
func RecvChanCtx[T any](ctx context.Context, chs ...<-chan T) (T, bool) {
	_, v, ok := reflect.Select(append(
		MapEach(chs, func(ch <-chan T) reflect.SelectCase {
			return reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			}
		}),
		reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		}),
	)

	return v.Interface().(T), ok
}

// CloseChan closes the given channels.
func CloseChan[T any](ch ...chan<- T) {
	ForEach(ch, func(ch chan<- T) {
		close(ch)
	})
}

// ToSend converts the given channels to send-only channels.
func ToSend[T any](ch ...chan T) []chan<- T {
	return MapEach(ch, func(ch chan T) chan<- T {
		return ch
	})
}

// ToRecv converts the given channels to receive-only channels.
func ToRecv[T any](ch ...chan T) []<-chan T {
	return MapEach(ch, func(ch chan T) <-chan T {
		return ch
	})
}
