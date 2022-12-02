package dj

import (
	"context"
)

// WorkerPool is a pool of workers that can be used to process work concurrently.
// Jobs of type Req are submitted to the pool and are processed by a worker.
// The result of the job is returned as a value of type Res.
type WorkerPool[Req, Res any] struct {
	reqCh chan<- Job[Req, Res]
	resCh <-chan Job[Req, Res]
	sem   *Group
}

// NewWorkerPool returns a new worker pool with the given number of workers.
// The given function is called for each job submitted to the pool.
// The function must return the result of the job.
func NewWorkerPool[Req, Res any](ctx context.Context, numWorkers int, fn func(context.Context, Req) (Res, error)) *WorkerPool[Req, Res] {
	reqCh, resCh := NewPipe[Job[Req, Res]]()

	grp := NewGroup(ctx, NewSem(numWorkers))

	go func() {
		defer close(reqCh)
		defer grp.Wait()

		grp.GoN(numWorkers, func(ctx context.Context, _ int) {
			ForChanCtx(ctx, resCh, func(ctx context.Context, job Job[Req, Res]) {
				v, err := fn(ctx, job.req)
				if err != nil {
					job.postErr(err)
				} else {
					job.postRes(v)
				}

				job.wait()
			})
		})
	}()

	return &WorkerPool[Req, Res]{reqCh: reqCh, resCh: resCh, sem: grp}
}

// Submit submits the given job to the pool.
func (p *WorkerPool[Req, Res]) Submit(ctx context.Context, req Req) Job[Req, Res] {
	job := newJob[Req, Res](ctx, req)

	p.reqCh <- job

	return job
}

// Process processes a list of jobs concurrently.
func (p *WorkerPool[Req, Res]) Process(ctx context.Context, reqs ...Req) ([]Res, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Submit the jobs to the pool.
	jobs := Map(reqs, func(req Req) Job[Req, Res] {
		return p.Submit(ctx, req)
	})

	res := make([]Res, len(jobs))

	// Collect the results.
	for i, job := range jobs {
		v, err := job.R()
		if err != nil {
			return nil, err
		}

		res[i] = v
	}

	return res, nil
}

// Job is a job that has been submitted to a worker pool.
type Job[Req, Res any] struct {
	ctx  context.Context
	req  Req
	res  chan Res
	err  chan error
	done chan struct{}
}

// newJob returns a new job with the given request.
func newJob[Req, Res any](ctx context.Context, req Req) Job[Req, Res] {
	return Job[Req, Res]{
		ctx:  ctx,
		req:  req,
		res:  make(chan Res),
		err:  make(chan error),
		done: make(chan struct{}),
	}
}

// R returns the result of the job, blocking until the job has finished.
// If the job failed, the error is returned.
func (job *Job[Req, Res]) R() (Res, error) {
	defer close(job.done)

	select {
	case <-job.ctx.Done():
		return Zero[Res](), job.ctx.Err()

	case err := <-job.err:
		return Zero[Res](), err

	case res := <-job.res:
		return res, nil
	}
}

// postRes posts the result of the job (success).
func (job *Job[Req, Res]) postRes(res Res) {
	select {
	case <-job.ctx.Done():
		// Context was canceled, don't send the result.

	case <-job.done:
		// Job must have failed.

	case job.res <- res:
		// ...
	}
}

// postErr posts the error of the job (failure).
func (job *Job[Req, Res]) postErr(err error) {
	select {
	case <-job.ctx.Done():
		// Context was canceled, don't send the error.

	case <-job.done:
		// Job must have succeeded.

	case job.err <- err:
		// ...
	}
}

// wait blocks until the job has either succeeded or failed.
func (job *Job[Req, Res]) wait() {
	select {
	case <-job.ctx.Done():
		// Context was canceled, don't wait for the job to finish.

	case <-job.done:
		// ...
	}
}
