package logzlib

import "nova/logz"

type FxLogger struct {
}

//goland:noinspection GoUnusedExportedFunction
func NewFxLogger() *FxLogger {
	return &FxLogger{}
}

func (l *FxLogger) Printf(msg string, args ...interface{}) {
	logz.Info().Msgf(msg, args...)
}
