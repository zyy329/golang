package mocklogic

import (
	"context"

	"github.com/stretchr/testify/mock"
)

func NewMockLogic() *MockLogic {
	return &MockLogic{}
}

type MockLogic struct {
	mock.Mock
}

func (x *MockLogic) DoLogic(ctx context.Context, p1, p2 int) int {
	args := x.Called(ctx, p1, p2)
	return args.Int(0)
}

func (x *MockLogic) DoRedis(ctx context.Context, p1, p2 int) int {
	args := x.Called(ctx, p1, p2)
	return args.Int(0)
}
