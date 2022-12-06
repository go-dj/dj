package dj

import "context"

// NewPipeline returns a pair of connected channels.
// Values sent to the first channel are processed by fn and the results are sent to the second channel.
// The pipeline uses n goroutines to process the values.
// If buf is greater than or equal to zero, up to buf values are buffered before the pipeline blocks.
// Otherwise, writes never block.
func NewPipeline[In, Out any](ctx context.Context, n, buf int, fn func(context.Context, In) Result[Out]) (chan<- In, <-chan Result[Out]) {
	var (
		reqCh chan<- In
		jobCh <-chan In
	)

	if buf >= 0 {
		reqCh, jobCh = NewBufPipe[In](buf)
	} else {
		reqCh, jobCh = NewPipe[In]()
	}

	return reqCh, GoMapChanCtx(WithParallelism(ctx, n), jobCh, fn)
}
