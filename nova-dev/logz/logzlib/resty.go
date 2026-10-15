package logzlib

import "nova/logz"

type RestyLogger struct {
}

func (x *RestyLogger) Errorf(format string, v ...interface{}) {
	logz.Error().Msgf(format, v...)
}

func (x *RestyLogger) Warnf(format string, v ...interface{}) {
	logz.Warn().Msgf(format, v...)
}
func (x *RestyLogger) Debugf(format string, v ...interface{}) {
	logz.Debug().Msgf(format, v...)
}
