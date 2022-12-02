package dj

import (
	"context"
)

// WorkerPool is a pool of workers that can be used to process work concurrently.
// Jobs of type In are submitted to the pool and are processed by a worker.
// The result of the job is returned as a value of type Out.
type WorkerPool[In, Out any] struct {
	inCh  chan<- Job[In, Out]
	outCh <-chan Job[In, Out]
	sem   *Group
}

// NewWorkerPool returns a new worker pool with the given number of workers.
// The given function is called for each job submitted to the pool.
// The function must return the result of the job.
func NewWorkerPool[In, Out any](ctx context.Context, numWorkers int, fn func(context.Context, In) (Out, error)) *WorkerPool[In, Out] {
	inCh, outCh := NewPipe[Job[In, Out]]()

	grp := NewGroup(ctx, NewSem(numWorkers))

	go func() {
		defer close(inCh)
		defer grp.Wait()

		grp.GoN(numWorkers, func(ctx context.Context, _ int) {
			ForChanCtx(ctx, outCh, func(ctx context.Context, job Job[In, Out]) {
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

	return &WorkerPool[In, Out]{inCh: inCh, outCh: outCh, sem: grp}
}

// Submit submits the given job to the pool.
func (p *WorkerPool[In, Out]) Submit(ctx context.Context, req In) Job[In, Out] {
	job := newJob[In, Out](ctx, req)

	p.inCh <- job

	return job
}

// Process processes a list of jobs concurrently.
func (p *WorkerPool[In, Out]) Process(ctx context.Context, reqs ...In) ([]Out, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Submit the jobs to the pool.
	jobs := MapEach(reqs, func(req In) Job[In, Out] {
		return p.Submit(ctx, req)
	})

	res := make([]Out, len(jobs))

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
type Job[In, Out any] struct {
	ctx  context.Context
	req  In
	res  chan Out
	err  chan error
	done chan struct{}
}

// newJob returns a new job with the given request.
func newJob[In, Out any](ctx context.Context, req In) Job[In, Out] {
	return Job[In, Out]{
		ctx:  ctx,
		req:  req,
		res:  make(chan Out),
		err:  make(chan error),
		done: make(chan struct{}),
	}
}

// R returns the result of the job, blocking until the job has finished.
// If the job failed, the error is returned.
func (job *Job[In, Out]) R() (Out, error) {
	defer close(job.done)

	select {
	case <-job.ctx.Done():
		return Zero[Out](), job.ctx.Err()

	case err := <-job.err:
		return Zero[Out](), err

	case res := <-job.res:
		return res, nil
	}
}

// postRes posts the result of the job (success).
func (job *Job[In, Out]) postRes(res Out) {
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
func (job *Job[In, Out]) postErr(err error) {
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
func (job *Job[In, Out]) wait() {
	select {
	case <-job.ctx.Done():
		// Context was canceled, don't wait for the job to finish.

	case <-job.done:
		// ...
	}
}
