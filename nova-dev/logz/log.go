package logz

import (
	"io"

	"nova/errorn"
	"nova/logz/logzdef"
	"nova/logz/zerologutil"

	"github.com/rs/zerolog"
)

// 业务日志输出工具
// 以链式风格输出日志
// 使用示例 请参看 test/log/example logPrint()

var (
	// 日志对象 和 关闭对象
	// DefaultLog 是完成初始化前 默认输出, 完成 Init 后会被替换掉
	log, closer = zerologutil.DefaultLog()

	listener = make([]func(l *zerolog.Logger, setLevel bool), 0)
)

//goland:noinspection GoUnusedGlobalVariable
var (
	Dict = zerolog.Dict
)

type Event = zerolog.Event

//goland:noinspection GoUnusedExportedFunction
func Close() {
	_ = closer.Close()
}

// SetLogger 更改logger; 不负责 oldCloser 的close, 由调用方负责确定是否执行
func SetLogger(l *zerolog.Logger, c io.Closer) {
	// 更新
	log = l
	closer = c

	// 触发回调
	for _, f := range listener {
		f(log, false)
	}
}

func GetLogger() *zerolog.Logger {
	return log
}

//goland:noinspection GoUnusedExportedFunction
func SetLevel(level logzdef.LogLevel) {
	newLog := zerologutil.ChildByLevel(log, level)
	log = newLog

	// 触发回调
	for _, f := range listener {
		f(log, true)
	}
}

// ListenLog 监听日志对象变化; 用于支持 logzlib
func ListenLog(f func(l *zerolog.Logger, setLevel bool)) {
	// 通常在启动时调用; 多线程问题不大, 不做多线程加锁
	f(log, false)
	listener = append(listener, f)
}

func Debug() *zerolog.Event {
	return log.Debug().Timestamp()
}

func Info() *zerolog.Event {
	return log.Info().Timestamp()
}

func Warn() *zerolog.Event {
	return log.Warn().Timestamp()
}

func Error() *zerolog.Event {
	return log.Error().Timestamp()
}

func Fatal() *zerolog.Event {
	return log.Fatal().Timestamp()
}

func Panic() *zerolog.Event {
	return log.Panic().Timestamp()
}

func Stack(l *zerolog.Event) *zerolog.Event {
	stS := errorn.GetStackS(1, errorn.DefaultStackLimit)
	zerologutil.LogPrintStack(l, stS)
	return l
}

//goland:noinspection GoUnusedExportedFunction
func StackSkip(l *zerolog.Event, skip int, limit int) *zerolog.Event {
	stS := errorn.GetStackS(skip+1, limit)
	zerologutil.LogPrintStack(l, stS)
	return l
}
