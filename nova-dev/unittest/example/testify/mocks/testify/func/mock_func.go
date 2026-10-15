package mockrandfunc

import "github.com/stretchr/testify/mock"

type MockRand struct {
	mock.Mock
}

func NewMockRand() *MockRand { return &MockRand{} }

func (m *MockRand) RandomInt(max int) int {
	args := m.Called(max)
	return args.Int(0)
}
