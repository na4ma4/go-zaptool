package zaptool

import "go.uber.org/zap/zapcore"

// ErrorLevel returns ErrorLevel if err is non-nil; otherwise, it returns
// InfoLevel.
func ErrorLevel(err error) zapcore.Level {
	if err == nil {
		return zapcore.InfoLevel
	}
	return zapcore.ErrorLevel
}
