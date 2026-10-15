package logzdef

import "github.com/rs/zerolog"

type LogLevel = string

const (
	Debug    LogLevel = "debug"
	Info     LogLevel = "info"
	Warn     LogLevel = "warn"
	Error    LogLevel = "error"
	Fatal    LogLevel = "fatal"
	Panic    LogLevel = "panic"
	NoLevel  LogLevel = "noLevel"
	Disabled LogLevel = "disabled"
	Trace    LogLevel = "trace"
)

func ToZeroLevel(level LogLevel) zerolog.Level {
	switch level {
	case Debug:
		return zerolog.DebugLevel
	case Info:
		return zerolog.InfoLevel
	case Warn:
		return zerolog.WarnLevel
	case Error:
		return zerolog.ErrorLevel
	case Fatal:
		return zerolog.FatalLevel
	case Panic:
		return zerolog.PanicLevel
	case NoLevel:
		return zerolog.NoLevel
	case Disabled:
		return zerolog.Disabled
	case Trace:
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}

//goland:noinspection GoUnusedExportedFunction
func ToLevel(level zerolog.Level) LogLevel {
	switch level {
	case zerolog.DebugLevel:
		return Debug
	case zerolog.InfoLevel:
		return Info
	case zerolog.WarnLevel:
		return Warn
	case zerolog.ErrorLevel:
		return Error
	case zerolog.FatalLevel:
		return Fatal
	case zerolog.PanicLevel:
		return Panic
	case zerolog.NoLevel:
		return NoLevel
	case zerolog.Disabled:
		return Disabled
	case zerolog.TraceLevel:
		return Trace
	default:
		return Info
	}
}
