package utils

import (
	"sync"
)

// 写入中间对象池化, 减少GC

type WriteHelper struct {
	Header       []byte // 传输头
	PacketHeader []byte // 业务头
	Body         []byte // 业务数据
}

func NewWriteHelperPool(headerMax, packetHeaderMax int) *sync.Pool {
	return &sync.Pool{New: func() interface{} {
		w := &WriteHelper{
			Header:       make([]byte, 0, headerMax),
			PacketHeader: make([]byte, 0, packetHeaderMax),
			Body:         nil,
		}
		return w
	}}
}

func (x *WriteHelper) Reset() {
	x.Header = x.Header[:0]
	x.PacketHeader = x.PacketHeader[:0]
	x.Body = nil
}

// 仅服务于对 Header 的 io.Writer 接口实现
func (x *WriteHelper) Write(p []byte) (n int, err error) {
	x.Header = append(x.Header, p...)

	// 返回实际写入的字节数（即 p 的长度）
	return len(p), nil
}
