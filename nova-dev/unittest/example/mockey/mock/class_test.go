package mock

import (
	"testing"

	"github.com/bytedance/mockey"
)

// mock 类函数 使用示例

// 类成员函数注入
type A struct {
}

func (x *A) Foo() string {
	return "A.Foo"
}

func TestMockClassSmp(t *testing.T) {
	mockey.Mock((*A).Foo).Return("mock.Foo").Build()

	tests := []struct {
		name string
		want string
	}{
		{"1", "mock.Foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &A{}
			if got := a.Foo(); got != tt.want {
				t.Errorf("Foo() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 类成员函数注入, 传递对象
type B struct {
	n int
}

func (x *B) Add(add int) int {
	x.n += add
	return x.n
}
func Add2Times(x *B, add int) int {
	x.n += add * 2
	return x.n
}

func TestMockClass(t *testing.T) {
	mockey.Mock((*B).Add).To(Add2Times).Build()

	tests := []struct {
		name string
		base int
		add  int
		want int
	}{
		{"1", 0, 1, 2},
		{"1", 3, 2, 7},
		{"1", 10, 4, 18},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &B{n: tt.base}
			if got := x.Add(tt.add); got != tt.want {
				t.Errorf("B.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
