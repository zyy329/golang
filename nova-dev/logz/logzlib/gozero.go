package logzlib

import (
	"fmt"

	"nova/logz"
	"nova/logz/logzdef"
	"nova/logz/zerologutil"

	"github.com/rs/zerolog"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	goZeroLevelStr        = "gozero_level"
	goZeroCallerStr       = "caller"
	goZeroCallerSkipFrame = 4 // Alert是3
)

var gGoZero = new(LogGoZero)

type LogGoZero struct {
	log *zerolog.Logger
}

//goland:noinspection GoUnusedExportedFunction
func GetGoZero(levelS ...logzdef.LogLevel) *LogGoZero {
	logz.ListenLog(func(l *zerolog.Logger, setLevel bool) {
		if setLevel || len(levelS) == 0 {
			gGoZero.log = l
		} else {
			gGoZero.log = zerologutil.ChildByLevel(l, levelS[0])
		}
	})

	return gGoZero
}

func (x *LogGoZero) Alert(v any) {
	e := x.log.Warn().Timestamp().CallerSkipFrame(3).Str(goZeroLevelStr, "alert")
	x.msg(e, v)
}

func (x *LogGoZero) Close() error {
	return nil
}

func (x *LogGoZero) Debug(v any, fields ...logx.LogField) {
	e := x.log.Debug().Timestamp().CallerSkipFrame(goZeroCallerSkipFrame).Str(goZeroLevelStr, "debug")
	x.msg(e, v, fields...)
}

func (x *LogGoZero) Error(v any, fields ...logx.LogField) {
	e := x.log.Error().Timestamp().CallerSkipFrame(goZeroCallerSkipFrame).Str(goZeroLevelStr, "error")
	x.msg(e, v, fields...)
}

func (x *LogGoZero) Info(v any, fields ...logx.LogField) {
	e := x.log.Info().Timestamp().CallerSkipFrame(goZeroCallerSkipFrame).Str(goZeroLevelStr, "info")
	x.msg(e, v, fields...)
}

func (x *LogGoZero) Severe(v any) {
	e := x.log.Error().Timestamp().CallerSkipFrame(goZeroCallerSkipFrame).Str(goZeroLevelStr, "severe")
	x.msg(e, v)
}

func (x *LogGoZero) Slow(v any, fields ...logx.LogField) {
	e := x.log.Error().Timestamp().CallerSkipFrame(goZeroCallerSkipFrame).Str(goZeroLevelStr, "slow")
	x.msg(e, v, fields...)
}

func (x *LogGoZero) Stack(v any) {
	e := x.log.Error().Timestamp().CallerSkipFrame(goZeroCallerSkipFrame).Str(goZeroLevelStr, "stack")
	x.msg(e, v)
}

func (x *LogGoZero) Stat(v any, fields ...logx.LogField) {
	e := x.log.Info().Timestamp().CallerSkipFrame(goZeroCallerSkipFrame).Str(goZeroLevelStr, "stat")
	x.msg(e, v, fields...)
}

func (x *LogGoZero) msg(e *zerolog.Event, v any, fields ...logx.LogField) {
	for _, field := range fields {
		if field.Key == goZeroCallerStr {
			continue
		}

		e = e.Any(field.Key, field.Value)
	}
	e.Msg(fmt.Sprint(v))
}
