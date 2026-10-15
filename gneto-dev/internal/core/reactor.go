package core

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/panjf2000/gnet/v2"
	"gneto/internal/codec/transport"
	"gneto/internal/logger"
	"gneto/pkg/config"
	"gneto/pkg/processor"
	"gneto/types"
)

// Reactor 网络事件处理器
type Reactor struct {
	gnet.BuiltinEventEngine

	eng        gnet.Engine
	dispatcher *Dispatcher
	processor  *processor.Processor
	transport  transport.Transport
	agentHook  types.AgentHook

	maxConn int32
	connNum atomic.Int32
}

// NewReactor 创建 Reactor
func NewReactor(opts *config.Options, dispatcher *Dispatcher, processor *processor.Processor) *Reactor {
	return &Reactor{
		dispatcher: dispatcher,
		processor:  processor,
		transport:  transport.Create(opts),
		agentHook:  opts.AgentHood,
		maxConn:    int32(opts.MaxConn),
		connNum:    atomic.Int32{},
	}
}

func (x *Reactor) OnBoot(eng gnet.Engine) (action gnet.Action) {
	x.eng = eng
	logger.Info("Network engine started")
	return
}

func (x *Reactor) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	newNum := x.connNum.Add(1)
	if newNum > x.maxConn {
		x.connNum.Add(-1)
		logger.Warn("Max connection number exceeded")
		return nil, gnet.Close
	}

	conn := NewConnection(c, x.processor, x.transport)
	c.SetContext(conn)
	if x.agentHook != nil {
		x.agentHook.OnConnect(conn)
	}
	logger.Info("New connection opened: curNum=%d", newNum)
	return
}

func (x *Reactor) OnClose(c gnet.Conn, err error) gnet.Action {
	ctx := c.Context()
	if ctx == nil {
		return gnet.None
	}
	conn, ok := ctx.(*Connection)
	if !ok {
		return gnet.None
	}

	if x.agentHook != nil {
		x.agentHook.OnClose(conn)
	}
	c.SetContext(nil)
	newNum := x.connNum.Add(-1)
	if err != nil {
		logger.Error("Network connection closed: curNum=%d, by Err:%v", newNum, err)
	} else {
		logger.Info("Network connection closed: curNum=%d", newNum)
	}
	return gnet.None
}

func (x *Reactor) OnTraffic(c gnet.Conn) (action gnet.Action) {
	// 循环消费掉已到达的所有网络包
	for {
		conn, ok := c.Context().(*Connection)
		if !ok {
			logger.Error("Network connection closed: connection lost")
			return gnet.Close
		}

		// 传输层切帧 (Raw Frame)
		rawData, err := x.transport.ReadRawFrame(c, conn)
		if err != nil {
			if errors.Is(err, transport.ErrIncomplete) || errors.Is(err, transport.ErrHandshaking) {
				return gnet.None
			}
			return gnet.Close
		}

		// 业务层解包
		msgId, msg, err := x.processor.Unmarshal(rawData, true)
		if err != nil {
			// 非法消息, 认为可能不是合法的客户端
			logger.Error("Network connection closed; unmarshal error: %v", err)
			return gnet.Close
		}

		// 4. 业务分发
		x.dispatcher.Dispatch(NewMsgTask(conn, msgId, msg))
	}
}

func (x *Reactor) OnTick() (delay time.Duration, action gnet.Action) {
	logger.Info("Network connected num:%d", x.connNum.Load())
	return 10 * time.Second, gnet.None
}
