package zaptool_test

import (
	"strings"
	"testing"

	"github.com/na4ma4/go-zaptool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogLevels_ChangeLevel(t *testing.T) {
	fac, observedLogs := observer.New(zapcore.DebugLevel)
	logger := zap.New(fac)

	loglvls := zaptool.NewLogLevels(logger)
	if loglvls.String() != "Internal.LogLevels:debug" {
		t.Error("Initial log level for Internal.LogLevels is not debug")
	}

	loglvls.SetLevel("Internal.LogLevels", zapcore.InfoLevel)
	if loglvls.String() != "Internal.LogLevels:info" {
		t.Error("Log level for Internal.LogLevels should be debug after SetLevel")
	}

	if len(observedLogs.Filter(func(le observer.LoggedEntry) bool {
		return (le.Level.String() != "debug")
	}).All()) > 0 {
		t.Error("should not contain any non debug logs")
	}
}

func TestLogLevels_LogAtLevel(t *testing.T) {
	fac, observedLogs := observer.New(zapcore.InfoLevel)
	logger := zap.New(fac)

	loglvls := zaptool.NewLogLevels(logger)
	if loglvls.String() != "Internal.LogLevels:info" {
		t.Error("Initial log level for Internal.LogLevels is not info")
	}

	testLogger := loglvls.Named("TestLogger")

	testLogger.Debug("[info] should not log")
	testLogger.Info("[info] should log")
	testLogger.Warn("[info] should log")

	loglvls.SetLevel("TestLogger", zapcore.WarnLevel)

	testLogger.Debug("[warn] should not log")
	testLogger.Info("[warn] should not log")
	testLogger.Warn("[warn] should log")

	loglvls.SetLevel("TestLogger", zapcore.DebugLevel)

	testLogger.Debug("[debug] should log")
	testLogger.Info("[debug] should log")
	testLogger.Warn("[debug] should log")

	if len(observedLogs.All()) != 6 {
		t.Errorf("should contain 6 log messages, instead contained %d messages", len(observedLogs.All()))
	}

	if len(observedLogs.Filter(func(le observer.LoggedEntry) bool {
		return strings.Contains(le.Message, "not")
	}).All()) > 0 {
		for _, le := range observedLogs.All() {
			if strings.Contains(le.Message, "not") {
				t.Logf("this message should not have been logged: %s", le.Message)
			}
		}
		t.Error("should not contain a 'should not log' message")
	}
}

func TestLogLevels_NotALevel(t *testing.T) {
	fac, observedLogs := observer.New(zapcore.InfoLevel)
	logger := zap.New(fac)

	loglvls := zaptool.NewLogLevels(logger)
	if loglvls.String() != "Internal.LogLevels:info" {
		t.Error("Initial log level for Internal.LogLevels is not info")
	}

	testLogger := loglvls.Named("TestLogger")

	if loglvls.String() != "Internal.LogLevels:info,TestLogger:info" {
		t.Error("Initial log level for TestLogger is not info")
	}

	testLogger.Debug("[info] should not log")
	testLogger.Info("[warn] should log")
	testLogger.Warn("[warn] should log")

	loglvls.SetLevel("TestLogger", zapcore.InvalidLevel)

	testLogger.Debug("[warn] should not log")
	testLogger.Info("[warn] should not log")
	testLogger.Warn("[warn] should not log")

	t.Logf("Levels: %s", loglvls.String())

	if loglvls.String() != "Internal.LogLevels:info,TestLogger:invalid" {
		t.Error("Updated log level for TestLogger is not invalid")
	}

	if len(observedLogs.All()) != 2 {
		t.Errorf("should contain 2 log messages, instead contained %d messages", len(observedLogs.All()))
	}

	if len(observedLogs.Filter(func(le observer.LoggedEntry) bool {
		return strings.Contains(le.Message, "not")
	}).All()) > 0 {
		for _, le := range observedLogs.All() {
			if strings.Contains(le.Message, "not") {
				t.Logf("this message should not have been logged: %s", le.Message)
			}
		}
		t.Error("should not contain a 'should not log' message")
	}
}
