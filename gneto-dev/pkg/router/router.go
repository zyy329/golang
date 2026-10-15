package router

import (
	"gneto/types"
)

// Router 路由管理器, 映射查找消息对应处理函数
// 可选, 也可业务层自行实现
type Router struct {
	handlers map[uint16]types.HandlerFunc
}

// NewRouter 创建路由
func NewRouter() *Router {
	return &Router{
		handlers: make(map[uint16]types.HandlerFunc),
	}
}

// Register 注册路由, 必须在启动时就完成全部注册(GetHandler 被调用前), 否则可能出现并发风险
func (x *Router) Register(msgId uint16, handler types.HandlerFunc) {
	x.handlers[msgId] = handler
}

// GetHandler 获取路由处理函数
func (x *Router) GetHandler(msgId uint16) (types.HandlerFunc, bool) {
	h, ok := x.handlers[msgId]
	return h, ok
}

func (x *Router) RouterFunc() types.RouterFunc {
	return func(s types.Agent, msg types.MessagePack) {
		if h, ok := x.GetHandler(msg.GetMsgId()); ok {
			h(s, msg)
		}
	}
}
