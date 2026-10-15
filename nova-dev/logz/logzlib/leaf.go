package logzlib

import "nova/logz"

const leafCallerSkipFrame = 2

var gLeaf = new(LogLeaf)

type LogLeaf struct {
}

//goland:noinspection GoUnusedExportedFunction
func GetLeaf() *LogLeaf {
	return gLeaf
}

func (x *LogLeaf) Debug(format string, args ...interface{}) {
	logz.Debug().CallerSkipFrame(leafCallerSkipFrame).Msgf(format, args...)
}

func (x *LogLeaf) Info(format string, args ...interface{}) {
	logz.Info().CallerSkipFrame(leafCallerSkipFrame).Msgf(format, args...)
}

func (x *LogLeaf) Warn(format string, args ...interface{}) {
	logz.Warn().CallerSkipFrame(leafCallerSkipFrame).Msgf(format, args...)
}

func (x *LogLeaf) Error(format string, args ...interface{}) {
	logz.Error().CallerSkipFrame(leafCallerSkipFrame).Msgf(format, args...)
}

func (x *LogLeaf) Fatal(format string, args ...interface{}) {
	logz.Fatal().CallerSkipFrame(leafCallerSkipFrame).Msgf(format, args...)
}
