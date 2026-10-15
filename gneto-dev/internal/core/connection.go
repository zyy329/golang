package core

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/panjf2000/gnet/v2"
	"gneto/internal/codec/packet"
	"gneto/internal/codec/transport"
	"gneto/internal/codec/utils"
	"gneto/internal/logger"
	"gneto/pkg/processor"
)

var (
	connIdCreator int64
)

// Connection 封装 gnet.Conn
type Connection struct {
	conn       gnet.Conn
	processor  *processor.Processor
	transport  transport.Transport
	connId     int64
	bindId     int64
	handshake  bool // 是否已完成握手(Websocket 使用)
	userData   atomic.Pointer[any]
	writerPool *sync.Pool
	localAddr  net.Addr
	remoteAddr net.Addr
}

// NewConnection 创建连接
func NewConnection(c gnet.Conn, processor *processor.Processor, transport transport.Transport) *Connection {
	return &Connection{
		conn:      c,
		processor: processor,
		transport: transport,
		connId:    atomic.AddInt64(&connIdCreator, 1),
		bindId:    0,
		handshake: false,
		writerPool: utils.NewWriteHelperPool(
			transport.GetMaxHeaderSize(),
			packet.GetMaxHeaderSize()),
		localAddr:  c.LocalAddr(),
		remoteAddr: c.RemoteAddr(),
	}
}

// WriteMsg 发送消息 (Payload -> Packet -> Transport -> Network)
func (x *Connection) WriteMsg(msgId uint16, data any) {
	w := x.writerPool.Get().(*utils.WriteHelper)

	if err := x.processor.MarshalWithHelper(msgId, data, w); err != nil {
		return
	}
	// 3. 传输层包装头
	err := x.transport.WriteHeader(w)
	if err != nil {
		logger.Error("FrameHeader error: %v", err)
		return
	}

	// 发送消息; 使用异步方法以支持并发安全发送
	if err := x.conn.AsyncWritev([][]byte{
		w.Header,
		w.PacketHeader,
		w.Body,
	}, func(c gnet.Conn, err error) error {
		w.Reset()
		x.writerPool.Put(w)
		return err
	}); err != nil {
		logger.Error("AsyncWritev error: %v", err)
	}
}

// Close 关闭连接
func (x *Connection) Close() {
	if err := x.conn.Close(); err != nil {
		logger.Error("close conn[connId:%v, bindId:%v] error: %v;", x.connId, x.bindId, err)
	}
}

func (x *Connection) GetConnId() int64 {
	return x.connId
}

func (x *Connection) SetBindId(bindId int64) {
	atomic.StoreInt64(&x.bindId, bindId)
}

func (x *Connection) GetBindId() int64 {
	return atomic.LoadInt64(&x.bindId)
}

func (x *Connection) DispatchKey() int64 {
	idx := x.GetBindId()
	if idx == 0 {
		idx = x.GetConnId()
	}
	return idx
}
func (x *Connection) DispatchFailed() {
	x.Close()
}

func (x *Connection) IsHandshake() bool {
	return x.handshake
}

func (x *Connection) SignHandshake() {
	x.handshake = true
}

func (x *Connection) LocalAddr() net.Addr {
	return x.localAddr
}

func (x *Connection) RemoteAddr() net.Addr {
	return x.remoteAddr
}

func (x *Connection) UserData() any {
	// Load() 操作是原子的，返回 *any，我们需要解引用得到 any
	val := x.userData.Load()
	if val == nil {
		return nil
	}
	return *val
}

func (x *Connection) SetUserData(data any) {
	x.userData.Store(&data)
}
