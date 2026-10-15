package logger

import (
	"gneto/types"
)

var gLogger types.Logger = &DefaultLogger{}

func Init(logger types.Logger) {
	gLogger = logger
}

//goland:noinspection GoUnusedExportedFunction
func Debug(format string, a ...interface{}) {
	gLogger.Debug(format, a...)
}

func Info(format string, a ...interface{}) {
	gLogger.Info(format, a...)
}

func Warn(format string, a ...interface{}) {
	gLogger.Warn(format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.Error(format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.Fatal(format, a...)
}
