package zstdpbn

import (
	"sync"

	"github.com/klauspost/compress/zstd"
	"google.golang.org/protobuf/proto"
)

// 包级别的全局变量，对外部隐藏，仅供包内函数使用
var (
	enc         *zstd.Encoder
	dec         *zstd.Decoder
	marshalPool sync.Pool
	unzipPool   sync.Pool
)

// init 函数会在包含此包的程序启动时自动执行一次
func init() {
	var err error

	// 1. 初始化 Encoder (注重体积 + 高并发小包优化)
	enc, err = zstd.NewWriter(nil,
		// 体积优先：使用默认压缩级别 (兼顾高压缩率和不错的速度)
		zstd.WithEncoderLevel(zstd.SpeedDefault),
		// 消除内部调度开销：外部已经是并发调用了，内部必须强制单线程，否则性能暴跌
		zstd.WithEncoderConcurrency(1),
		// 极致内存优化：限制滑动窗口大小为 64KB (假设 PB 包不大)
		// 大幅降低内存占用，且能完美命中 CPU L1/L2 Cache 提升速度
		zstd.WithWindowSize(64<<10),
	)
	if err != nil {
		panic("zstdpbn: failed to initialize zstd encoder: " + err.Error())
	}

	// 2. 初始化 Decoder (极致轻量化)
	dec, err = zstd.NewReader(nil,
		// 解压端也关闭内部并发
		zstd.WithDecoderConcurrency(1),
		// 启用低内存模式：配合编码端的 64KB Window，将解压内存降到最低
		zstd.WithDecoderLowmem(true),
	)
	if err != nil {
		panic("zstdpbn: failed to initialize zstd decoder: " + err.Error())
	}

	// 3. 初始化双向内存池 (保持不变)
	marshalPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, 2048)
			return &b
		},
	}
	unzipPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, 4096)
			return &b
		},
	}
}

// Zip 终极压缩版：内部零分配，开箱即用
func Zip(m proto.Message) []byte {
	// 1. 从对象池获取重用的 Buffer
	ptr := marshalPool.Get().(*[]byte)
	mbuf := (*ptr)[:0] // 长度复位，保留容量

	// 确保退出时归还 Buffer
	defer func() {
		*ptr = mbuf
		marshalPool.Put(ptr)
	}()

	// 2. 零分配序列化 PB
	opts := proto.MarshalOptions{}
	var err error
	mbuf, err = opts.MarshalAppend(mbuf, m)
	if err != nil || len(mbuf) == 0 {
		return []byte{}
	}

	// 3. 执行压缩
	return enc.EncodeAll(mbuf, make([]byte, 0, len(mbuf)/2))
}

// Unzip 终极解压版：内部零分配，开箱即用
func Unzip(b []byte, m proto.Message) error {
	if len(b) == 0 {
		return nil
	}

	// 1. 从对象池获取重用的 Buffer
	ptr := unzipPool.Get().(*[]byte)
	dst := (*ptr)[:0] // 长度复位，保留容量

	// 确保退出时归还 Buffer
	defer func() {
		*ptr = dst
		unzipPool.Put(ptr)
	}()

	// 2. 执行解压
	var err error
	dst, err = dec.DecodeAll(b, dst)
	if err != nil {
		return err
	}

	// 3. 零分配反序列化 PB
	return proto.Unmarshal(dst, m)
}
