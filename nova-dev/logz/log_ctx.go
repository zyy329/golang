package logz

import (
	"context"
	"sync/atomic"

	"nova/gn"
	"nova/logz/zerologutil"

	"github.com/rs/zerolog"
)

// 附带上下文支持的 业务日志输出工具

var logId = int64(0)

type Ctx struct {
	fields zerologutil.FieldArr
}

func NewCtx() *Ctx {
	return &Ctx{}
}

//goland:noinspection GoUnusedExportedFunction
func CtxWith(k string, v interface{}) *Ctx {
	l := NewCtx()
	l.AddFiled(k, v)
	return l
}

func WithTrace(ctx context.Context) *Ctx {
	return NewCtx().WithTrace(ctx)
}

func (x *Ctx) WithTrace(ctx context.Context) *Ctx {
	traceID := TraceIDFromContext(ctx)
	if len(traceID) > 0 {
		x.AddFiled(gn.TraceKey, traceID)
	}
	x.AddFiled(gn.LogIdKey, atomic.AddInt64(&logId, int64(1)))
	return x
}

func (x *Ctx) AddFiled(k string, v interface{}) *Ctx {
	if x.fields == nil {
		x.fields = make([]interface{}, 4)
	}
	x.fields = append(x.fields, k, v)
	return x
}

func (x *Ctx) withCtx(out *zerolog.Event) *zerolog.Event {
	if len(x.fields) != 0 {
		out = out.Fields(x.fields)
	}
	return out
}

func (x *Ctx) Debug() *zerolog.Event {
	return x.withCtx(Debug())
}

func (x *Ctx) Info() *zerolog.Event {
	return x.withCtx(Info())
}

func (x *Ctx) Warn() *zerolog.Event {
	return x.withCtx(Warn())
}

func (x *Ctx) Error() *zerolog.Event {
	return x.withCtx(Error())
}

func (x *Ctx) Fatal() *zerolog.Event {
	return x.withCtx(Fatal())
}

func (x *Ctx) Panic() *zerolog.Event {
	return x.withCtx(Panic())
}
