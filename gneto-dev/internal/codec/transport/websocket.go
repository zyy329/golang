package transport

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"hash"
	"io"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"gneto/internal/codec/utils"
)

type WSTransport struct {
	maxMsgLen int
}

func (x *WSTransport) ReadRawFrame(c gnet.Conn, ctx Context) ([]byte, error) {
	// 1. 检查是否已完成 WebSocket 握手
	if !ctx.IsHandshake() {
		return nil, x.handleHandshake(c, ctx)
	}

	// 2. 检查缓冲区数据是否足够读取最小 WS 帧头
	size := c.InboundBuffered()
	if size < ws.MinHeaderSize {
		return nil, ErrIncomplete
	}

	// 3. 预读取 (Peek) 当前所有数据，尝试解析帧头
	// 注意：此处不消耗缓冲区，因为不确定数据是否足够一整帧
	buf, _ := c.Peek(size)
	br := bytes.NewReader(buf)
	hdr, err := ws.ReadHeader(br)
	if err != nil {
		// 如果由于 EOF 导致读取失败，说明数据包不完整，等待更多数据
		if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, ErrIncomplete
		}
		return nil, err
	}

	// 4. 计算帧头的实际字节长度
	headerLen := size - br.Len()
	// 5. 计算这一帧的总长度 (头 + 载荷)
	totalLen := headerLen + int(hdr.Length)

	// 6. 如果缓冲区总长度不足一整帧，返回不完整错误，等待下一次流量触发
	if size < totalLen {
		return nil, ErrIncomplete
	}

	if x.maxMsgLen < totalLen {
		return nil, ErrOverMaxMsgLen
	}

	// 7. 确认有一封完整的帧，开始消耗缓冲区：先丢弃已解析的帧头部分
	if _, err := c.Discard(headerLen); err != nil {
		return nil, err
	}
	// 8. 提取完整的载荷部分 (Payload)
	payload, err := c.Next(int(hdr.Length))
	if err != nil {
		return nil, err
	}

	// 9. 处理 WebSocket 控制帧 (如 Ping, Pong, Close)
	if hdr.OpCode.IsControl() {
		if hdr.OpCode == ws.OpPing {
			// 如果是 Ping 帧，自动回复 Pong
			_ = wsutil.WriteServerMessage(c, ws.OpPong, payload)
		}
		// 控制帧处理完后返回特定错误，以便外层继续循环处理后续业务帧
		return nil, ErrHandshaking
	}

	// 10. 如果客户端数据带有掩码 (RFC 规定客户端发给服务端必须带掩码)，进行解码
	if hdr.Masked {
		ws.Cipher(payload, hdr.Mask, 0)
	}

	// 11. 返回解码后的原始载荷，供业务层 UnmarshalPacket 使用
	return payload, nil
}

func (x *WSTransport) WriteHeader(w *utils.WriteHelper) error {
	// 1. 只构造 WebSocket 帧头
	header := ws.Header{
		Fin:    true,
		OpCode: ws.OpBinary,
		Length: int64(len(w.PacketHeader) + len(w.Body)),
	}

	// 2. 将 Header 写入 buffer
	err := ws.WriteHeader(w, header)
	if err != nil {
		return err
	}

	return nil
}

func (x *WSTransport) GetMaxHeaderSize() int {
	return ws.MaxHeaderSize
}

var (
	// 预定义字节数组，避免运行时分配
	binUpgradeKey    = []byte("Upgrade")
	binWebsocketVal  = []byte("websocket")
	binConnectionKey = []byte("Connection")
	binUpgradeVal    = []byte("upgrade")
	binSecKey        = []byte("Sec-WebSocket-Key")
	wsGUID           = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

	// 预置响应模板
	httpResponsePre = []byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: ")
	crlf            = []byte("\r\n")

	sha1Pool = sync.Pool{New: func() interface{} { return sha1.New() }}
)

func (x *WSTransport) handleHandshake(c gnet.Conn, ctx Context) error {
	// 1. 检查报文完整性
	size := c.InboundBuffered()
	data, _ := c.Peek(size)
	headerEnd := bytes.Index(data, []byte("\r\n\r\n"))
	if headerEnd == -1 {
		return ErrIncomplete
	}
	headerData := data[:headerEnd]

	// 2. 核心校验变量
	var (
		hasUpgrade    bool
		hasConnection bool
		secKey        []byte
	)

	// 3. 单次扫描解析 (Single-pass Byte Scanning)
	lines := bytes.Split(headerData, crlf)
	for i, line := range lines {
		if i == 0 {
			// 第一行校验：必须是 GET
			if !bytes.HasPrefix(line, []byte("GET ")) {
				return errors.New("invalid method")
			}
			continue
		}

		// 寻找冒号分隔符
		pos := bytes.IndexByte(line, ':')
		if pos == -1 {
			continue
		}

		key := bytes.TrimSpace(line[:pos])
		val := bytes.TrimSpace(line[pos+1:])

		// 使用 EqualFold 进行无分配的大小写无关比较
		switch {
		case bytes.EqualFold(key, binUpgradeKey):
			if bytes.EqualFold(val, binWebsocketVal) {
				hasUpgrade = true
			}
		case bytes.EqualFold(key, binConnectionKey):
			// Connection 可能包含多个 token (如 "keep-alive, upgrade")
			if containsToken(val, binUpgradeVal) {
				hasConnection = true
			}
		case bytes.EqualFold(key, binSecKey):
			secKey = val
		}
	}

	// 4. 合规性判定
	if !hasUpgrade || !hasConnection || len(secKey) == 0 {
		return errors.New("invalid websocket handshake")
	}

	// ---------------------------------------------------------
	// 5. 剩下的 SHA-1 计算与发送逻辑
	// ---------------------------------------------------------

	// 从对象池获取 SHA1 实例，减少 GC 压力
	h := sha1Pool.Get().(hash.Hash)
	h.Reset()
	defer sha1Pool.Put(h)

	// 计算签名: SHA1(Key + GUID)
	h.Write(secKey)
	h.Write(wsGUID)
	hashSum := h.Sum(nil)

	// 构造最终响应报文
	// 响应长度 = 模板前缀 + Key + 结尾换行
	encodedLen := base64.StdEncoding.EncodedLen(len(hashSum))
	res := make([]byte, len(httpResponsePre)+encodedLen+4) // 4 = \r\n\r\n
	// 按照偏移填充数据
	offset := copy(res, httpResponsePre)
	// 在正确的位置执行 Base64 编码（直接写入 res，无需额外分配）
	base64.StdEncoding.Encode(res[offset:offset+encodedLen], hashSum)
	offset += encodedLen
	// 填充结尾
	copy(res[offset:], "\r\n\r\n")

	// 一次性写入缓冲区
	if _, err := c.Write(res); err != nil {
		return err
	}

	// 消耗掉已处理的 HTTP 头部数据，保留后续可能的 Data Frame
	fullHeaderLen := headerEnd + 4 // 4 = \r\n\r\n
	if _, err := c.Discard(fullHeaderLen); err != nil {
		return err
	}

	ctx.SignHandshake()
	return ErrHandshaking
}

// 辅助函数：在不分配内存的情况下检查逗号分隔的 Token
func containsToken(header []byte, token []byte) bool {
	parts := bytes.Split(header, []byte(","))
	for _, p := range parts {
		if bytes.EqualFold(bytes.TrimSpace(p), token) {
			return true
		}
	}
	return false
}
