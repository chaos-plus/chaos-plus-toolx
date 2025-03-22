package xgrpool

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// GoroutinePool 表示一个协程池
type GoroutinePool struct {
	wg     sync.WaitGroup  // 用于等待所有协程完成
	ctx    context.Context // 上下文
	cancel func()          // 取消函数
	panic  func(any)       // panic处理函数
}

func New() *GoroutinePool {
	return NewWithContext(context.Background())
}

func NewWithContext(ctx context.Context) *GoroutinePool {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	return &GoroutinePool{
		wg:     sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
		panic: func(e any) {
			fmt.Printf("goroutine error: %v\n", e)
		},
	}
}

func (p *GoroutinePool) UncaughtErrorHandler(handler func(error any)) *GoroutinePool {
	p.panic = handler
	return p
}

func (p *GoroutinePool) Wait() {
	p.wg.Wait()
}

func (p *GoroutinePool) Stop() {
	p.cancel()
	p.wg.Wait()
}

func (p *GoroutinePool) Add(run func(context.Context) error) func() {
	return p.AddWithCancelAndRecover(run, nil, nil)
}

func (p *GoroutinePool) AddWithCancel(run func(context.Context) error, cancel func(context.Context) error) func() {
	return p.AddWithCancelAndRecover(run, cancel, nil)
}

func (p *GoroutinePool) AddWithRecover(run func(context.Context) error, caught func(context.Context, interface{})) func() {
	return p.AddWithCancelAndRecover(run, nil, caught)
}

func (p *GoroutinePool) AddWithCancelAndRecover(run func(context.Context) error, cancel func(context.Context) error, caught func(context.Context, interface{})) func() {
	ctx, ctxcancel := context.WithCancel(p.ctx)
	p.wg.Add(1)

	ctxch := make(chan struct{})
	go func() {
		defer p.wg.Done()
		defer p.caught(nil, caught)
        defer close(ctxch)
		if run == nil {
			return
		}
		if err := run(ctx); err != nil {
			p.caught(err, caught)
		}
		
	}()
	go func() {
		if cancel == nil {
			return
		}
		defer p.caught(nil, caught)
		select {
		case <-ctxch:
			return
		case <-ctx.Done():
			if err := cancel(ctx); err != nil {
				p.caught(err, caught)
			}
		}
	}()
	return ctxcancel
}

func (p *GoroutinePool) caught(err any, caught func(context.Context, interface{})) {
	if err == nil {
		err = recover()
	}
	if err != nil {
		stack := make([]byte, 64*1024)
		n := runtime.Stack(stack, false)
		stack = stack[:n]
		msg := fmt.Sprintf("%v\n%s", err, stack)
		if caught != nil {
			caught(p.ctx, msg)
		} else if p.panic != nil {
			p.panic(msg)
		} else {
			fmt.Println("goroutine panic:", msg)
		}
	}
}
