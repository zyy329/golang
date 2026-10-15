package types

type RouterFunc func(s Agent, msg MessagePack)
type HandlerFunc func(s Agent, msg MessagePack)

// MessageCreator 消息工厂函数，用于提供给业务层实现消息的创建逻辑, 避免创建对象时需要反射(低效)
type MessageCreator func() any

type TransportType int
type PayloadType int

const (
	Tcp TransportType = iota + 1
	WebSocket
)

const (
	Protobuf PayloadType = iota + 1
	JSON
)
