package dj

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
	return TakeChanCtx(context.Background(), ch, n)
}

// TakeChanCtx reads up to n values from the channel.
// It stops reading when the context is canceled.
func TakeChanCtx[T any](ctx context.Context, ch <-chan T, n int) []T {
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

// FanIn merges the given channels into a single channel.
// That is, values are read from each channel in parallel.
// The returned channel is closed when all the given channels are closed.
func FanIn[T any](chs ...<-chan T) <-chan T {
	return FanInCtx(context.Background(), chs...)
}

// FanInCtx merges the given channels into a single channel.
// That is, values are read from each channel in parallel.
// The returned channel is closed when all the given channels are closed or the context is canceled.
func FanInCtx[T any](ctx context.Context, chs ...<-chan T) <-chan T {
	grp := NewGroup(ctx, NewSem(len(chs)))
	out := make(chan T)

	go func() {
		defer close(out)
		defer grp.Wait()

		grp.GoN(len(chs), func(ctx context.Context, i int) {
			ForChanCtx(ctx, chs[i], func(ctx context.Context, v T) {
				SendCtx(ctx, v, out)
			})
		})
	}()

	return out
}

// FanOut splits the given channel into n channels.
// The returned channels are closed when the given channel is closed.
func FanOut[T any](ch <-chan T, n int) []<-chan T {
	return FanOutCtx(context.Background(), ch, n)
}

// FanOutCtx splits the given channel into n channels.
// The returned channels are closed when the given channel is closed or the context is canceled.
func FanOutCtx[T any](ctx context.Context, ch <-chan T, n int) []<-chan T {
	out := MapN(n, func(int) chan T {
		return make(chan T)
	})

	go func() {
		defer CloseChan(AsSend(out...)...)
		ForwardChanCtx(ctx, []<-chan T{ch}, AsSend(out...))
	}()

	return AsRecv(out...)
}

// ForwardChan forwards values from the src channel(s) to the dst channel(s).
func ForwardChan[T any](src []<-chan T, dst []chan<- T) {
	ForwardChanCtx(context.Background(), src, dst)
}

// ForwardChanCtx forwards values from the src channel(s) to the dst channel(s).
// It stops forwarding when the context is canceled.
func ForwardChanCtx[T any](ctx context.Context, src []<-chan T, dst []chan<- T) {
	ForChanCtx(ctx, FanInCtx(ctx, src...), func(ctx context.Context, v T) {
		SendCtx(ctx, v, dst...)
	})
}

// ConcatChan joins the given channels into a single channel.
// That is, the first channel is emptied, then the second, and so on.
// The returned channel is closed when all the given channels are closed.
func ConcatChan[T any](chs ...<-chan T) <-chan T {
	return ConcatChanCtx(context.Background(), chs...)
}

// ConcatChanCtx joins the given channels into a single channel.
// That is, the first channel is emptied, then the second, and so on.
// The returned channel is closed when all the given channels are closed or the context is canceled.
func ConcatChanCtx[T any](ctx context.Context, chs ...<-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for _, ch := range chs {
			ForwardChanCtx(ctx, []<-chan T{ch}, AsSend(out))
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

		in := ZipIter(Map(chs, func(ch <-chan T) Iter[T] {
			return ChanIterCtx(ctx, ch)
		})...).Recv()

		ForwardChanCtx(ctx, []<-chan []T{in}, AsSend(out))
	}()

	return out
}

// Send sends the given value to one of the given channels.
// It blocks until the value is sent.
func Send[T any](v T, chs ...chan<- T) {
	SendCtx(context.Background(), v, chs...)
}

// SendCtx sends the given value to one of the given channels.
// It blocks until the value is sent or the context is canceled.
func SendCtx[T any](ctx context.Context, v T, chs ...chan<- T) {
	reflect.Select(append(
		Map(chs, func(ch chan<- T) reflect.SelectCase {
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
func Recv[T any](chs ...<-chan T) (T, bool) {
	_, v, ok := reflect.Select(Map(chs, func(ch <-chan T) reflect.SelectCase {
		return reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}))

	return v.Interface().(T), ok
}

// RecvCtx receives a value from one of the given channels.
// It blocks until a value is received, which it returns, or the context is canceled.
// The boolean indicates whether the read was successful; it is false if the channel is closed.
func RecvCtx[T any](ctx context.Context, chs ...<-chan T) (T, bool) {
	_, v, ok := reflect.Select(append(
		Map(chs, func(ch <-chan T) reflect.SelectCase {
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
	For(ch, func(ch chan<- T) {
		close(ch)
	})
}

// AsSend converts the given channels to send-only channels.
func AsSend[T any](ch ...chan T) []chan<- T {
	return Map(ch, func(ch chan T) chan<- T {
		return ch
	})
}

// AsRecv converts the given channels to receive-only channels.
func AsRecv[T any](ch ...chan T) []<-chan T {
	return Map(ch, func(ch chan T) <-chan T {
		return ch
	})
}
