package logzlib

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nova/gn/tracen"
	"nova/logz"

	"gorm.io/gorm/logger"
)

const (
	TraceStr     = "[%.3fms] [rows:%v] %s"
	TraceWarnStr = "%s [%.3fms] [rows:%v] %s"
	TraceErrStr  = "%s [%.3fms] [rows:%v] %s"
)

type GormLogger struct {
	LogLevel      logger.LogLevel
	ShowSQL       bool
	SlowThreshold time.Duration
}

func (x *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	x.LogLevel = level
	return x
}

func (x *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if x.LogLevel >= logger.Info {
		logz.WithTrace(ctx).Info().Msgf(msg, data...)
	}
}

func (x *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if x.LogLevel >= logger.Warn {
		logz.WithTrace(ctx).Warn().Msgf(msg, data...)
	}
}

func (x *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if x.LogLevel >= logger.Error {
		logz.WithTrace(ctx).Error().Msgf(msg, data...)
	}
}

func (x *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rows int64), err error) {
	if x.LogLevel == logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && x.LogLevel >= logger.Error && !errors.Is(err, logger.ErrRecordNotFound):
		logz.WithTrace(ctx).AddFiled(tracen.Duration, elapsed.Milliseconds()).Error().Msgf(TraceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case x.SlowThreshold > 0 && elapsed > x.SlowThreshold:
		slowLog := fmt.Sprintf("SLOW SQL >= %d", x.SlowThreshold)
		logz.WithTrace(ctx).AddFiled(tracen.Duration, elapsed.Milliseconds()).Warn().Msgf(TraceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case x.ShowSQL:
		logz.WithTrace(ctx).AddFiled(tracen.Duration, elapsed.Milliseconds()).Info().Timestamp().Msgf(TraceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
