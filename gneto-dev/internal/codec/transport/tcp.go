package transport

import (
	"encoding/binary"

	"github.com/panjf2000/gnet/v2"
	"gneto/internal/codec/utils"
)

const (
	TCPHeaderSize = 4
)

// TCPTransport 实现 TCP 长度头帧处理
type TCPTransport struct {
	byteOrder binary.ByteOrder
	maxMsgLen int
}

func (x *TCPTransport) ReadRawFrame(c gnet.Conn, _ Context) ([]byte, error) {
	if c.InboundBuffered() < TCPHeaderSize {
		return nil, ErrIncomplete
	}
	lenBuf, _ := c.Peek(TCPHeaderSize)
	bodyLen := x.byteOrder.Uint32(lenBuf)
	totalLen := bodyLen + TCPHeaderSize
	if x.maxMsgLen < int(totalLen) {
		return nil, ErrOverMaxMsgLen
	}

	if c.InboundBuffered() < int(totalLen) {
		return nil, ErrIncomplete
	}

	if _, err := c.Discard(TCPHeaderSize); err != nil {
		return nil, err
	}
	return c.Next(int(bodyLen))
}

func (x *TCPTransport) WriteHeader(w *utils.WriteHelper) error {
	bodyLen := len(w.PacketHeader) + len(w.Body)
	w.Header = w.Header[:TCPHeaderSize]
	x.byteOrder.PutUint32(w.Header, uint32(bodyLen))
	return nil
}

func (x *TCPTransport) GetMaxHeaderSize() int {
	return TCPHeaderSize
}
