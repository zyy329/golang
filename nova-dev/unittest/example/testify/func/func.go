package randfunc

import "math/rand/v2"

// RandGenerator 待Mock接口
type RandGenerator interface {
	RandomInt(max int) int
}

// 原始实现
type standarRand struct{}

func newStandarRand() RandGenerator { return &standarRand{} }

func (s *standarRand) RandomInt(max int) int {
	return rand.N(max)
}

// AddByRand 测试调用入口
func AddByRand(r RandGenerator, num int) int {
	return num + r.RandomInt(num)
}
