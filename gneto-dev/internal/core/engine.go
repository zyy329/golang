package core

import (
	"context"
	"fmt"
	"time"

	"github.com/panjf2000/gnet/v2"
	"gneto/internal/logger"
	"gneto/pkg/config"
	"gneto/pkg/processor"
	"gneto/types"
)

type Engine struct {
	opts       *config.Options
	dispatcher *Dispatcher
	reactor    *Reactor
}

func NewEngine(opts *config.Options, processor *processor.Processor, handler types.RouterFunc) *Engine {
	dispatcher := NewDispatcher(opts.WorkPoolSize, opts.WorkerQueueCap, handler, processor.Payload)
	reactor := NewReactor(opts, dispatcher, processor)
	return &Engine{
		opts:       opts,
		dispatcher: dispatcher,
		reactor:    reactor,
	}
}

// Run 启动网络引擎
func (x *Engine) Run() error {
	addr := fmt.Sprintf("tcp://:%d", x.opts.Port)
	return gnet.Run(x.reactor, addr,
		gnet.WithMulticore(true),
		gnet.WithLogger(&logger.LogGNet{}),
	)
}

func (x *Engine) Stop() error {
	// 1. 设定一个“硬截止”时间，防止清理逻辑卡死导致 Pod 无法删除
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return x.reactor.eng.Stop(ctx)
}

func (x *Engine) GetDispatcher() *Dispatcher {
	return x.dispatcher
}
