package mock

import (
	"testing"

	"nova/framen"

	"github.com/bytedance/mockey"
)

// mock 函数 使用示例

// 测试启动, 做统一的初始化化工作, 同包下每个Test都会自动调用该函数
func TestMain(m *testing.M) {
	framen.TestMainRun(m)
}

func RetSelf(a int) int {
	return a
}
func Ret10Times(a int) int {
	return a * 10
}

// 函数注入, 简单替换结果
func TestMockFuncRet(t *testing.T) {
	mockey.Mock(RetSelf).Return(10).Build()

	tests := []struct {
		name string
		args int
		want int
	}{
		{"1", 1, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetSelf(tt.args); got != tt.want {
				t.Errorf("RetSelf() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 函数注入, 顺序替换结果
func TestMockFuncRetSeq(t *testing.T) {
	mockey.Mock(RetSelf).Return(mockey.Sequence(10).Times(2).Then(50).Then(90)).Build()

	tests := []struct {
		name string
		args int
		want int
	}{
		{"1", 1, 10},
		{"2", 2, 10},
		{"3", 3, 50},
		{"4", 4, 90},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetSelf(tt.args); got != tt.want {
				t.Errorf("RetSelf() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 函数注入, 替换为另一函数进行调用
func TestMockFunc(t *testing.T) {
	mockey.Mock(RetSelf).To(Ret10Times).Build()

	tests := []struct {
		name string
		args int
		want int
	}{
		{"1", 1, 10},
		{"2", 2, 20},
		{"3", 3, 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetSelf(tt.args); got != tt.want {
				t.Errorf("RetSelf() = %v, want %v", got, tt.want)
			}
		})
	}
}
