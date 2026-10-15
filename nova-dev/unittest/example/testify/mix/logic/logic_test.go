package logic

import (
	"nova/framen"
	"nova/redisn"
	"testing"
)

// 分层测试, logic 只测 logic 自身的逻辑

// 测试启动, 做统一的初始化化工作, 同包下每个Test都会自动调用该函数
func TestMain(m *testing.M) {
	framen.TestMainRun(m)
}

func TestDoLogic(t *testing.T) {
	// 无需mock, 测试 logic 自身逻辑
	x := NewLogic()
	t.Logf("standar CallLogic: %d", x.DoLogic(t.Context(), 10, 11))
}

func TestDoRedis(t *testing.T) {
	// 利用 bytedance/mockey mock 三方组件(redis)
	// 接口支持的情况下, 也可使用 testify 进行 mock
	redisn.MockRedis(t)
	x := NewLogic()
	t.Logf("mockey DoRedis: %d", x.DoRedis(t.Context(), 10, 11))
}
