package randfunc

import (
	"nova/framen"
	mock_randfunc "nova/unittest/example/testify/mocks/testify/func"
	"testing"
)

// 测试启动, 做统一的初始化化工作, 同包下每个Test都会自动调用该函数
func TestMain(m *testing.M) {
	framen.TestMainRun(m)
}

func TestMock(t *testing.T) {
	// 原始调用(对比用,实际测试中不会有这一部分)
	r := newStandarRand()
	t.Logf("standar: %d", AddByRand(r, 10))
	t.Logf("standar: %d", AddByRand(r, 20))

	// Mock
	m := mock_randfunc.NewMockRand()
	m.On("RandomInt", 10).Return(6)
	m.On("RandomInt", 20).Return(7)

	// Mock 调用
	t.Logf("mock: %d", AddByRand(m, 10))
	t.Logf("mock: %d", AddByRand(m, 20))
}
