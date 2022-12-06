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
		if v, ok := RecvFromCtx(ctx, ch); !ok {
			return out
		} else if out = append(out, v); len(out) == n {
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
		v, ok := RecvFromCtx(ctx, ch)
		if !ok {
			return
		}

		fn(ctx, v)
	}
}

// MapChan returns a channel that applies the given function to each value in the input channel.
func MapChan[T, R any](ch <-chan T, fn func(T) R) <-chan R {
	return MapChanCtx(context.Background(), ch, func(_ context.Context, v T) R {
		return fn(v)
	})
}

// MapChanCtx returns a channel that applies the given function to each value in the input channel.
// It stops iterating when the context is canceled.
func MapChanCtx[T, R any](ctx context.Context, ch <-chan T, fn func(context.Context, T) R) <-chan R {
	out := make(chan R)

	go func() {
		defer close(out)

		ForChanCtx(ctx, ch, func(ctx context.Context, v T) {
			out <- fn(ctx, v)
		})
	}()

	return out
}

// GoMapChan returns a channel that applies the given function to each value in the input channel in parallel.
func GoMapChan[T, R any](ch <-chan T, fn func(T) R) <-chan R {
	return GoMapChanCtx(context.Background(), ch, func(_ context.Context, v T) R {
		return fn(v)
	})
}

// GoMapChanCtx returns a channel that applies the given function to each value in the input channel in parallel.
// It stops iterating when the context is canceled.
func GoMapChanCtx[T, R any](ctx context.Context, ch <-chan T, fn func(context.Context, T) R) <-chan R {
	out := make(chan R)

	go func() {
		defer close(out)

		group := NewGroup(ctx, NewSem(parallelismFromCtx(ctx)))
		defer group.Wait()

		group.Go(parallelismFromCtx(ctx), func(ctx context.Context, _ int) {
			ForChanCtx(ctx, ch, func(ctx context.Context, v T) {
				out <- fn(ctx, v)
			})
		})
	}()

	return out
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

		grp.Go(len(chs), func(ctx context.Context, i int) {
			ForChanCtx(ctx, chs[i], func(ctx context.Context, v T) {
				SendToCtx(ctx, v, out)
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
		ForwardCtx(ctx, []<-chan T{ch}, AsSend(out...))
	}()

	return AsRecv(out...)
}

// Forward forwards values from the src channel(s) to the dst channel(s).
func Forward[T any](src []<-chan T, dst []chan<- T) {
	ForwardCtx(context.Background(), src, dst)
}

// ForwardCtx forwards values from the src channel(s) to the dst channel(s).
// It stops forwarding when the context is canceled.
func ForwardCtx[T any](ctx context.Context, src []<-chan T, dst []chan<- T) {
	ForChanCtx(ctx, FanInCtx(ctx, src...), func(ctx context.Context, v T) {
		SendToCtx(ctx, v, dst...)
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
			ForwardCtx(ctx, []<-chan T{ch}, AsSend(out))
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
		})...).Recv()

		ForwardCtx(ctx, []<-chan []T{in}, AsSend(out))
	}()

	return out
}

// SendTo sends the given value to one of the given channels.
// It blocks until the value is sent.
func SendTo[T any](v T, chs ...chan<- T) {
	SendToCtx(context.Background(), v, chs...)
}

// SendToCtx sends the given value to one of the given channels.
// It blocks until the value is sent or the context is canceled.
func SendToCtx[T any](ctx context.Context, v T, chs ...chan<- T) {
	reflect.Select(append(
		MapEach(chs, func(ch chan<- T) reflect.SelectCase {
			return sendCase(ch, v)
		}),
		recvCase(ctx.Done()),
	))
}

// RecvFrom receives a value from one of the given channels.
// It blocks until a value is received, which it returns.
// The boolean indicates whether the read was successful; it is false if the channel is closed.
func RecvFrom[T any](chs ...<-chan T) (T, bool) {
	return RecvFromCtx(context.Background(), chs...)
}

// RecvFromCtx receives a value from one of the given channels.
// It blocks until a value is received, which it returns, or the context is canceled.
// The boolean indicates whether the read was successful; it is false if the channel is closed.
func RecvFromCtx[T any](ctx context.Context, chs ...<-chan T) (T, bool) {
	if _, v, ok := reflect.Select(append(
		MapEach(chs, func(ch <-chan T) reflect.SelectCase {
			return recvCase(ch)
		}),
		recvCase(ctx.Done()),
	)); ok {
		return v.Interface().(T), true
	}

	return Zero[T](), false
}

// CloseChan closes the given channels.
func CloseChan[T any](ch ...chan<- T) {
	ForEach(ch, func(ch chan<- T) {
		close(ch)
	})
}

// AsSend converts the given channels to send-only channels.
func AsSend[T any](ch ...chan T) []chan<- T {
	return MapEach(ch, func(ch chan T) chan<- T {
		return ch
	})
}

// AsRecv converts the given channels to receive-only channels.
func AsRecv[T any](ch ...chan T) []<-chan T {
	return MapEach(ch, func(ch chan T) <-chan T {
		return ch
	})
}

func sendCase[T any](ch chan<- T, v T) reflect.SelectCase {
	return reflect.SelectCase{
		Dir:  reflect.SelectSend,
		Chan: reflect.ValueOf(ch),
		Send: reflect.ValueOf(v),
	}
}

func recvCase[T any](ch <-chan T) reflect.SelectCase {
	return reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ch),
	}
}
