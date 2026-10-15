package service

import (
	"nova/framen"
	mocklogic "nova/unittest/example/testify/mocks/mix/logic"
	"testing"

	"github.com/stretchr/testify/mock"
)

// 分层测试, service 只测 service 自身的逻辑

// 测试启动, 做统一的初始化化工作, 同包下每个Test都会自动调用该函数
func TestMain(m *testing.M) {
	framen.TestMainRun(m)
}

func TestDoService(t *testing.T) {
	// 无需mock, 测试 service 自身逻辑
	x := NewService()
	t.Logf("standar DoService: %d", x.DoService(t.Context(), 10, 11))
}

func TestDoLogic(t *testing.T) {
	// 只测自身逻辑, 外部调用logic 通过mock 方式屏蔽掉
	m := mocklogic.NewMockLogic()
	// 预期目标参数调用会触发两次 (AssertExpectations 会验证该预期)
	m.On("DoLogic", mock.Anything, 10, 11).Return(1).Twice()
	// 预期目标参数调用会触发一次
	m.On("DoLogic", mock.Anything, 20, 21).Return(2).Once()

	x := NewService(WithLogic(m))
	t.Logf("mock call CallLogic ret: %d", x.CallLogic(t.Context(), 10, 11))

	// 断言检查是否符合调用预期
	m.AssertExpectations(t)
}
