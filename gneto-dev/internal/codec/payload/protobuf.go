package payload

import (
	"fmt"

	"gneto/internal/logger"
	"google.golang.org/protobuf/proto"
)

const (
	// MaxPoolObjectSize 定义对象归还池中的最大字节数 (1MB)
	// 超过此大小的对象在 Reset 后依然可能占用大量底层内存，建议直接丢弃
	MaxPoolObjectSize = 1024 * 1024
)

type ProtobufPayload struct {
	*MessagePool
}

func NewProtobufPayload(pool *MessagePool) *ProtobufPayload {
	p := &ProtobufPayload{MessagePool: pool}
	return p
}

func (x *ProtobufPayload) PushMessage(msgId uint16, msg any) {
	// 校验是否为 proto.Message，非该类型的对象禁止放入池中
	pm, ok := msg.(proto.Message)
	if !ok {
		logger.Warn("PushMessage failed: object is not a proto.Message, msgId:%d", msgId)
		return
	}

	// 拦截大对象：如果消息内容过大，不归还到池中，防止内存空洞
	// proto.Reset 会清除数据，但通常会保留底层切片的容量（Capacity）
	// 这里判断 序列化后的大小, 避免巨大的切片残留在内存中
	if proto.Size(pm) > MaxPoolObjectSize {
		logger.Debug("PushMessage: object too large, skipping pool return; msgId: %d, size: %d",
			msgId, proto.Size(pm))
		return
	}

	// 潜在问题: proto.Reset 只会清除 .proto 文件中定义的字段, 不会清理 自定义字段
	proto.Reset(pm)
	x.pushToPool(msgId, msg)
}

func (x *ProtobufPayload) Marshal(v any) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("not a proto.Message; type: %T", v)
	}
	return proto.Marshal(msg)
}

func (x *ProtobufPayload) Unmarshal(data []byte, v any) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("not a proto.Message; type: %T", v)
	}
	return proto.Unmarshal(data, msg)
}
