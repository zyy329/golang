package payload

import (
	"gneto/internal/logger"
	"gneto/types"
)

// Payload 序列化接口 (JSON/Protobuf 等...)
type Payload interface {
	NewMessage(msgID uint16) (msg any, err error)
	PopMessage(msgId uint16) (msg any, err error)
	PushMessage(msgId uint16, msg any) // 未能成功还回并不会造成泄漏而是会被GC处理, 主动还回可优化性能
	Marshal(msg any) ([]byte, error)
	Unmarshal(data []byte, msg any) error
}

func Create(pool *MessagePool, payloadType types.PayloadType) Payload {
	switch payloadType {
	case types.Protobuf:
		return NewProtobufPayload(pool)
	case types.JSON:
		return NewJsonPayload(pool)
	}
	logger.Fatal("invalid payload type:%v", payloadType)
	return nil
}
