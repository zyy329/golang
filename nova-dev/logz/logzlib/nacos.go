package logzlib

import (
	"fmt"

	"nova/logz"
	"nova/logz/logzdef"
	"nova/logz/zerologutil"

	"github.com/rs/zerolog"
)

// 重载 nacos 底层的日志输出

var nacosLog = nacos{}

const skipAdjust = 1

type nacos struct {
	log *zerolog.Logger
}

//goland:noinspection GoExportedFuncWithUnexportedType
func GetNacos(level logzdef.LogLevel) *nacos {
	logz.ListenLog(func(l *zerolog.Logger, setLevel bool) {
		if setLevel {
			nacosLog.log = l
		} else {
			nacosLog.log = zerologutil.ChildByLevel(l, level)
		}
	})

	return &nacosLog
}

// 实现 nacos 的log接口, 以供覆盖日志输出

func (x *nacos) Info(args ...interface{}) {
	x.log.Info().Timestamp().CallerSkipFrame(skipAdjust).Msg(fmt.Sprint(args...))
}

func (x *nacos) Warn(args ...interface{}) {
	x.log.Warn().Timestamp().CallerSkipFrame(skipAdjust).Msg(fmt.Sprint(args...))
}

func (x *nacos) Error(args ...interface{}) {
	x.log.Error().Timestamp().CallerSkipFrame(skipAdjust).Msg(fmt.Sprint(args...))
}

func (x *nacos) Debug(args ...interface{}) {
	x.log.Debug().Timestamp().CallerSkipFrame(skipAdjust).Msg(fmt.Sprint(args...))
}

func (x *nacos) Infof(fmt string, args ...interface{}) {
	x.log.Info().Timestamp().CallerSkipFrame(skipAdjust).Msgf(fmt, args...)
}

func (x *nacos) Warnf(fmt string, args ...interface{}) {
	x.log.Warn().Timestamp().CallerSkipFrame(skipAdjust).Msgf(fmt, args...)
}

func (x *nacos) Errorf(fmt string, args ...interface{}) {
	x.log.Error().Timestamp().CallerSkipFrame(skipAdjust).Msgf(fmt, args...)
}

func (x *nacos) Debugf(fmt string, args ...interface{}) {
	x.log.Debug().Timestamp().CallerSkipFrame(skipAdjust).Msgf(fmt, args...)
}

func (x *nacos) Close() error {
	return nil
}
