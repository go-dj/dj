package xn

import (
	"context"
	"sync"
)

type Sem struct {
	ctx context.Context
	ch  chan struct{}
	wg  sync.WaitGroup
}

func NewSem(ctx context.Context, n int) *Sem {
	return &Sem{
		ctx: ctx,
		ch:  make(chan struct{}, n),
	}
}

func (s *Sem) Wait() {
	s.wg.Wait()
}

func (s *Sem) Go(fn func(context.Context)) {
	s.acquire()
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		defer s.release()

		select {
		case <-s.ctx.Done():
			// Context was canceled, don't run the function.

		default:
			fn(s.ctx)
		}
	}()
}

func (s *Sem) acquire() {
	s.ch <- struct{}{}
}

func (s *Sem) release() {
	<-s.ch
}
