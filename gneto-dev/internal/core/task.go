package core

import (
	"fmt"

	"gneto/types"
)

// WorkerTask 其中 Msg, CustomFn 用于分别支持两种不同的处理流程, 两者只能一个有效
type WorkerTask interface {
	types.DispatchAble
	LogInfo() string
}

type MsgTask struct {
	types.Agent
	msgId uint16 // 消息Id
	msg   any    // 已解码的消息体; 注意:该消息对象不可持久化持有或传入其他协程,否则可能出现数据错误(因为直接调用逻辑结束后,该消息会被清理和还回重用池)
}

func NewMsgTask(agent types.Agent, msgId uint16, msg any) *MsgTask {
	return &MsgTask{
		Agent: agent,
		msgId: msgId,
		msg:   msg,
	}
}

func (x *MsgTask) GetAgent() types.Agent {
	return x.Agent
}

func (x *MsgTask) LogInfo() string {
	return fmt.Sprintf("[connId:%d, bindId:%d, msgId:%d]", x.GetConnId(), x.GetBindId(), x.msgId)
}

func (x *MsgTask) GetMsgId() uint16 {
	return x.msgId
}
func (x *MsgTask) GetMsg() any {
	return x.msg
}

func (x *MsgTask) AcquireMsg() any {
	msg := x.msg
	// 标记所有权已转移, 后续不要再放入重用池中
	x.msg = nil
	return msg
}

// ------------------------------------

type FuncTask struct {
	types.DispatchAble
	customFn func()
}

func NewFuncTask(d types.DispatchAble, customFn func()) *FuncTask {
	return &FuncTask{
		DispatchAble: d,
		customFn:     customFn,
	}
}

func (x *FuncTask) LogInfo() string {
	return fmt.Sprintf("[DispatchAble:%d]", x.DispatchKey())
}

func (x *FuncTask) GetCustomFn() func() {
	return x.customFn
}
