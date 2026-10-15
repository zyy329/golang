package transport

import (
	"errors"

	"github.com/panjf2000/gnet/v2"
	"gneto/internal/codec/utils"
	"gneto/internal/logger"
	"gneto/pkg/config"
	"gneto/types"
)

var (
	// ErrIncomplete 数据不完整信号, 表示当前缓冲区（InboundBuffer）中的字节数还不够组成一个完整的业务包, 需等待
	ErrIncomplete = errors.New("incomplete data")
	// ErrHandshaking 协议层内部处理信号, 数据已经被消费了，但它属于协议层消息; 不需要后续进一步处理
	ErrHandshaking   = errors.New("handling protocol handshake/control")
	ErrOverMaxMsgLen = errors.New("message too long")
)

// Context Transport 需要的上下文信息
type Context interface {
	IsHandshake() bool
	SignHandshake() // 标记已完成握手
}

// Transport 处理流协议的拆解 (TCP/WebSocket 等...)
type Transport interface {
	// ReadRawFrame 负责从网络流中切分出一个完整的原始数据块 (Raw Frame)
	ReadRawFrame(c gnet.Conn, ctx Context) ([]byte, error)
	// WriteHeader 基于原始数据块生成协议头 (WS Header 或 TCP Length Header)
	WriteHeader(w *utils.WriteHelper) error
	GetMaxHeaderSize() int
}

func Create(opts *config.Options) Transport {
	switch opts.TransportType {
	case types.Tcp:
		return &TCPTransport{
			byteOrder: opts.ByteOrder,
			maxMsgLen: opts.MaxMsgLen,
		}
	case types.WebSocket:
		return &WSTransport{
			maxMsgLen: opts.MaxMsgLen,
		}
	}
	logger.Fatal("invalid transport type:%v", opts.TransportType)
	return nil
}
