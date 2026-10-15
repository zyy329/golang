package config

import (
	"encoding/binary"

	"gneto/types"
)

// 函数式配置选项
// Functional configuration options

// Options 服务器配置
type Options struct {
	Port           int
	WorkPoolSize   int
	WorkerQueueCap int
	MaxConn        int
	MaxMsgLen      int
	ByteOrder      binary.ByteOrder
	TransportType  types.TransportType
	PayloadType    types.PayloadType
	AgentHood      types.AgentHook
}

// Option 配置函数
type Option func(*Options)

func NewDefaultOptions() *Options {
	return &Options{
		Port:           8888,
		WorkPoolSize:   1024,
		WorkerQueueCap: 1024,
		MaxConn:        1000,
		MaxMsgLen:      4096,
		ByteOrder:      binary.BigEndian, // 在网络编程中，通常固定为 BigEndian
		TransportType:  types.WebSocket,
		PayloadType:    types.Protobuf,
		AgentHood:      nil,
	}
}

// WithPort 设置端口
func WithPort(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

// WithWorkPoolSize 设置 工作协程池大小
func WithWorkPoolSize(size int) Option {
	return func(o *Options) {
		o.WorkPoolSize = size
	}
}

// WithWorkerQueueCap
// 单个工作协程的任务队列缓冲大小。
// 当业务处理速度跟不上网络接收速度时，消息会在此排队。
// 如果队列溢出，认为负载过高, 将触发熔断, 丢弃消息并踢掉当前链接, 避免消息持续积压
func WithWorkerQueueCap(len int) Option {
	return func(o *Options) {
		o.WorkerQueueCap = len
	}
}

// WithMaxConn 最大链接数
func WithMaxConn(max int) Option {
	return func(o *Options) {
		o.MaxConn = max
	}
}

// WithMaxMsgLen 单个消息最大长度
//
//goland:noinspection GoUnusedExportedFunction
func WithMaxMsgLen(max int) Option {
	return func(o *Options) {
		o.MaxMsgLen = max
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithByteOrder(binaryByteOrder binary.ByteOrder) Option {
	return func(o *Options) {
		o.ByteOrder = binaryByteOrder
	}
}

// WithTransportType 设置 传输协议类型
func WithTransportType(transportType types.TransportType) Option {
	return func(o *Options) {
		o.TransportType = transportType
	}
}

// WithPayloadType 设置 负载类型
func WithPayloadType(payloadType types.PayloadType) Option {
	return func(o *Options) {
		o.PayloadType = payloadType
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithAgentHook(hook types.AgentHook) Option {
	return func(o *Options) {
		o.AgentHood = hook
	}
}
