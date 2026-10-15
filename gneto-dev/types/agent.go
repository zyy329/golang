package types

import "net"

// Agent 代表网络层的一个逻辑会话。
// 它是 网络层 与业务逻辑之间的中间桥梁，不直接持有复杂的业务对象（如 Player,Room），
// 旨在简化逻辑和降低内存泄漏风险 并原生支持“断线重连”场景下的状态迁移。
//
// 并发保证:
//  1. BindId 用于提供会话粘性, 若指定了 BindId，系统保证具有相同 BindId 的消息会被路由至同一个协程/Worker 处理。
//     这对于处理具有时序性要求的玩家逻辑至关重要，可避免业务层频繁加锁。
//  2. 若未指定 BindId(默认值为 0)，则退化为基于 ConnId 的路由模式。
//
// 最佳实践: 业务层应通过 BindID 在自身的 Manager 中检索真正的业务对象。
// 风险警示: 设置bindId后, 后续到达的包会立刻丢到新协程中; 此时可能有并发风险; 请注意设置时机(比如在进入房间处理逻辑最后再设置)
type Agent interface {
	DispatchAble
	GetConnId() int64
	SetBindId(bindId int64) // Bind ID 通常指定为 playerId 或 roomId, 用于确保处理协程相同
	GetBindId() int64
	Close()
	WriteMsg(msgId uint16, data any) // 发送消息
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	UserData() any
	SetUserData(data any)
}

type AgentHook interface {
	OnConnect(agent Agent)
	OnClose(agent Agent)
}
