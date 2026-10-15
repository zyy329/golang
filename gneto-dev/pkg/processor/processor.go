package processor

import (
	"encoding/binary"

	"gneto/internal/codec/packet"
	"gneto/internal/codec/payload"
	"gneto/internal/codec/utils"
	"gneto/internal/logger"
	"gneto/pkg/config"
	"gneto/types"
)

type Processor struct {
	Packet    *packet.Packet
	Payload   payload.Payload
	msgPool   *payload.MessagePool
	byteOrder binary.ByteOrder
}

func NewProcessor(opts ...config.Option) *Processor {
	options := config.NewDefaultOptions()
	for _, o := range opts {
		o(options)
	}

	msgPool := payload.NewMessagePool()
	return &Processor{
		msgPool:   msgPool,
		Packet:    packet.New(options.ByteOrder),
		Payload:   payload.Create(msgPool, options.PayloadType),
		byteOrder: options.ByteOrder,
	}
}

// RegisterMessage 注册消息结构体, 用于消息反序列化
func (x *Processor) RegisterMessage(msgId uint16, factory types.MessageCreator) {
	x.msgPool.Register(msgId, factory)
}

// Marshal 业务层数据的编码
func (x *Processor) Marshal(msgId uint16, data any) (header []byte, body []byte, err1 error) {
	w := &utils.WriteHelper{
		PacketHeader: make([]byte, 0, packet.GetMaxHeaderSize()),
	}

	err := x.MarshalWithHelper(msgId, data, w)
	if err != nil {
		return nil, nil, err
	}
	return w.PacketHeader, w.Body, nil
}

func (x *Processor) MarshalWithHelper(msgId uint16, data any, w *utils.WriteHelper) (err error) {
	// 1. 消息序列化
	w.Body, err = x.Payload.Marshal(data)
	if err != nil {
		logger.Error("marshal err", err)
		return
	}

	// 2. 数据包组包处理
	x.Packet.Marshal(msgId, w)
	return
}

// Unmarshal 业务层数据的解码
func (x *Processor) Unmarshal(rawData []byte, pool bool) (msgId uint16, msg any, err error) {
	// 2. 数据包结构解析
	var p *packet.Info
	p, err = x.Packet.Unmarshal(rawData)
	if err != nil {
		return
	}
	msgId = p.MsgId

	// 消息对象生成
	if pool {
		// 重用池分配空消息对象, Dispatch中会执行对应的 PushMessage
		msg, err = x.Payload.PopMessage(p.MsgId)
	} else {
		msg, err = x.Payload.NewMessage(p.MsgId)
	}
	if err != nil {
		return
	}

	// 3. 消息反序列化;
	err = x.Payload.Unmarshal(p.Body, msg)

	return
}

func (x *Processor) GetByteOrder() binary.ByteOrder {
	return x.byteOrder
}
