package service

import (
	"context"
	"nova/logz"
	"nova/unittest/example/testify/mix/logic"
)

type Service interface {
	DoService(ctx context.Context, p1, p2 int) int
	CallLogic(ctx context.Context, p1, p2 int) int
}

type service struct {
	logic logic.Logic
}

type Option func(*service)

func NewService(opts ...Option) Service {
	x := &service{}
	for _, opt := range opts {
		opt(x)
	}

	if x.logic == nil {
		x.logic = logic.NewLogic()
	}
	return x
}

func WithLogic(l logic.Logic) Option {
	return func(x *service) { x.logic = l }
}

func (x *service) DoService(_ context.Context, p1, p2 int) int {
	logz.Info().Msgf("call service.DoService: %d %d", p1, p2)
	return p1 + p2
}

func (x *service) CallLogic(ctx context.Context, p1, p2 int) int {
	logz.Info().Msgf("call service.CallLogic: %d %d", p1, p2)

	// ....业务逻辑
	// ........

	// 调用
	x.logic.DoLogic(ctx, p1, p2)

	// ....业务逻辑
	// ........

	// 调用
	x.logic.DoLogic(ctx, p1+10, p2+10)

	// ....业务逻辑
	// ........

	// 调用
	return x.logic.DoLogic(ctx, p1, p2)
}
