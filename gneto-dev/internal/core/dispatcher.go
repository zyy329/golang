package core

import (
	"runtime"

	"gneto/internal/codec/payload"
	"gneto/internal/logger"
	"gneto/types"
)

const LenStackBuf = 4096

// Dispatcher 消息分发器, 决定分发到哪个处理协程
type Dispatcher struct {
	payload        payload.Payload
	workPoolSize   int64
	workerQueueCap int
	workers        []chan WorkerTask
	handler        types.RouterFunc
}

// NewDispatcher 创建分发器
func NewDispatcher(
	workPoolSize, workerQueueCap int,
	handler types.RouterFunc,
	payload payload.Payload) *Dispatcher {

	if workPoolSize <= 0 {
		workPoolSize = 1
	}
	d := &Dispatcher{
		payload:        payload,
		workPoolSize:   int64(workPoolSize),
		workerQueueCap: workerQueueCap,
		workers:        make([]chan WorkerTask, workPoolSize),
		handler:        handler,
	}
	for i := 0; i < workPoolSize; i++ {
		d.workers[i] = make(chan WorkerTask, d.workerQueueCap)
		go d.workerLoop(d.workers[i])
	}
	return d
}

func (x *Dispatcher) workerLoop(ch chan WorkerTask) {
	for task := range ch {
		func() {
			defer func() {
				if r := recover(); r != nil {
					buf := make([]byte, LenStackBuf)
					l := runtime.Stack(buf, false)
					logger.Error("%v: %s", r, buf[:l])
				}
			}()

			switch t := task.(type) {
			case *FuncTask:
				// 调用业务层自定义处理函数
				t.GetCustomFn()()
			case *MsgTask:
				// 调用通用消息处理流程
				if x.handler != nil {
					defer func() {
						msg := t.AcquireMsg()
						if msg != nil {
							x.payload.PushMessage(t.GetMsgId(), msg)
						}
					}()

					x.handler(t.GetAgent(), t)
				}
			}
		}()
	}
}

// Dispatch 任务分发
func (x *Dispatcher) Dispatch(task WorkerTask) {
	idx := task.DispatchKey()

	// 简单的负载均衡策略
	idx = idx % x.workPoolSize
	if idx < 0 {
		idx = -idx
	}

	// 分发(避免阻塞gnet IO线程)
	select {
	case x.workers[idx] <- task:
		// 投递成功
	default:
		// 队列已满, 认为负载过高,
		// 丢弃消息, 断开该链接, 避免持续堆积
		logger.Error("Worker queue overflow, idx: %d, len:%d, log:%ag",
			idx, len(x.workers[idx]), task.LogInfo())
		task.DispatchFailed()
	}
}
