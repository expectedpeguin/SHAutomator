package executor

import (
	"context"
)

type Pool struct {
	ctx     context.Context
	workers chan struct{}
}

func NewPool(ctx context.Context, maxWorkers int) *Pool {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	return &Pool{
		ctx:     ctx,
		workers: make(chan struct{}, maxWorkers),
	}
}

func (p *Pool) Submit(task func()) {
	select {
	case p.workers <- struct{}{}:
		go func() {
			defer func() { <-p.workers }()
			task()
		}()
	case <-p.ctx.Done():
		return
	}
}
