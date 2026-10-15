package zstdpbn

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestZipUnzip 基础的功能准确性测试
func TestZipUnzip(t *testing.T) {
	// 准备测试数据
	originalMsg := timestamppb.Now()

	// 1. 测试压缩
	compressedData := Zip(originalMsg)
	if len(compressedData) == 0 {
		t.Fatal("压缩失败：返回的字节数组为空")
	}
	t.Logf("原始消息大小估算: %d bytes, 压缩后大小: %d bytes", proto.Size(originalMsg), len(compressedData))

	// 2. 测试解压
	unzippedMsg := &timestamppb.Timestamp{}
	err := Unzip(compressedData, unzippedMsg)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 3. 验证数据一致性
	if !proto.Equal(originalMsg, unzippedMsg) {
		t.Fatalf("数据损坏！\n压缩前: %v\n解压后: %v", originalMsg, unzippedMsg)
	}
}

// TestUnzipEmpty 测试空数据保护
func TestUnzipEmpty(t *testing.T) {
	msg := &timestamppb.Timestamp{}
	err := Unzip([]byte{}, msg)
	if err != nil {
		t.Fatalf("解析空字节切片不应返回错误，但得到了: %v", err)
	}
}

// BenchmarkZipUnzip 基准测试：用于验证内部零分配 (Zero-Allocation) 和高并发性能
func BenchmarkZipUnzip(b *testing.B) {
	// 准备固定测试数据
	msg := timestamppb.Now()
	compressedData := Zip(msg)

	b.Run("Zip", func(b *testing.B) {
		b.ReportAllocs() // 开启内存分配报告
		b.ResetTimer()   // 重置计时器，排除初始化耗时

		for i := 0; i < b.N; i++ {
			// 这里的分配仅来自于 EncodeAll 必须返回的新切片，内部 Buffer 是零分配的
			_ = Zip(msg)
		}
	})

	b.Run("Unzip", func(b *testing.B) {
		outMsg := &timestamppb.Timestamp{}
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// 复用 outMsg 进行反序列化，严格测试 zstdpb 内部的分配情况
			proto.Reset(outMsg)
			_ = Unzip(compressedData, outMsg)
		}
	})
}
