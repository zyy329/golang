package packet

import (
	"encoding/binary"
	"errors"

	"gneto/internal/codec/utils"
)

// 数据包定义

// Packet 协议包结构
// 数据包定义 (处理自定义包结构: MsgId, Body)
// 后续可扩展 魔树(快速过滤非法连接),序列号,加密,签名压缩等功能

type Info struct {
	MsgId uint16 // 业务消息 Id
	Body  []byte // 原始负载 (二进制)
}

const (
	// HeaderSize 包头固定长度: MsgId(2)
	HeaderSize = 2
)

type Packet struct {
	byteOrder binary.ByteOrder
}

func New(byteOrder binary.ByteOrder) *Packet {
	return &Packet{
		byteOrder: byteOrder,
	}
}

func GetMaxHeaderSize() int {
	return HeaderSize
}

func (x *Packet) Unmarshal(data []byte) (*Info, error) {
	if len(data) < HeaderSize {
		return nil, errors.New("packet short than header")
	}
	msgId := x.byteOrder.Uint16(data[:HeaderSize])
	return &Info{
		MsgId: msgId,
		Body:  data[HeaderSize:],
	}, nil
}

func (x *Packet) Marshal(msgId uint16, w *utils.WriteHelper) {
	w.PacketHeader = w.PacketHeader[:HeaderSize]
	x.byteOrder.PutUint16(w.PacketHeader, msgId)
}
