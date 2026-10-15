package logger

import "log"

type DefaultLogger struct{}

func (x *DefaultLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
func (x *DefaultLogger) Debug(format string, args ...interface{}) {
	x.Printf("[Debug]"+format, args...)
}
func (x *DefaultLogger) Info(format string, args ...interface{}) {
	x.Printf("[Info]"+format, args...)
}
func (x *DefaultLogger) Warn(format string, args ...interface{}) {
	x.Printf("[Warn]"+format, args...)
}
func (x *DefaultLogger) Error(format string, args ...interface{}) {
	x.Printf("[Error]"+format, args...)
}
func (x *DefaultLogger) Fatal(format string, args ...interface{}) {
	log.Fatalf("[Fatal]"+format, args...)
}
