package logic

import (
	"context"
	"nova/logz"
	"nova/redisn"

	"github.com/redis/go-redis/v9"
)

type Logic interface {
	DoLogic(ctx context.Context, p1, p2 int) int
	DoRedis(ctx context.Context, p1, p2 int) int
}

type logic struct {
	redisC *redis.Client
}
type Option func(*logic)

func NewLogic(opts ...Option) Logic {
	x := &logic{}
	for _, opt := range opts {
		opt(x)
	}

	if x.redisC == nil {
		x.redisC = redisn.NewClient()
	}
	return x
}

//func WithRedis(redisC *redis.Client) Option {
//	return func(x *logic) { x.redisC = redisC }
//}

func (x *logic) DoLogic(_ context.Context, p1, p2 int) int {
	logz.Info().Msgf("call logic.CallLogic: %d %d", p1, p2)
	return p1 + p2
}

func (x *logic) DoRedis(ctx context.Context, p1, p2 int) int {
	logz.Info().Msgf("call logic.DoRedis: %d %d", p1, p2)
	ret := x.redisC.IncrBy(ctx, "testInc", int64(p1+p2))
	return int(ret.Val())
}
