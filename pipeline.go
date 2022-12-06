package dj

import "context"

// NewPipeline returns a pair of connected channels.
// The pipeline calls fn for each value sent to the first channel and sends the result to the second channel.
// If size is greater than zero, it is buffered; writes to the first channel block when the buffer is full.
func NewPipeline[In, Out any](ctx context.Context, size int, fn func(context.Context, In) Result[Out]) (chan<- In, <-chan Result[Out]) {
	var (
		reqCh chan<- In
		jobCh <-chan In
	)

	if size > 0 {
		reqCh, jobCh = NewBufPipe[In](size)
	} else {
		reqCh, jobCh = NewPipe[In]()
	}

	return reqCh, MapChanCtx(ctx, jobCh, fn)
}
