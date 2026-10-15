package logger

// LogGNet 提供给 gNet 库使用的日志对象, 固定为实际使用 gLogger 对应的输出
type LogGNet struct {
}

func (x *LogGNet) Debugf(format string, args ...any) {
	gLogger.Debug(format, args...)
}

// Infof logs messages at INFO level.
func (x *LogGNet) Infof(format string, args ...any) {
	gLogger.Info(format, args...)
}

// Warnf logs messages at WARN level.
func (x *LogGNet) Warnf(format string, args ...any) {
	gLogger.Warn(format, args...)
}

// Errorf logs messages at ERROR level.
func (x *LogGNet) Errorf(format string, args ...any) {
	gLogger.Error(format, args...)
}

// Fatalf logs messages at FATAL level.
func (x *LogGNet) Fatalf(format string, args ...any) {
	gLogger.Fatal(format, args...)
}
