package zaptool_test

import (
	"strings"
	"testing"

	"github.com/na4ma4/go-zaptool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestSubLogLevels_Names(t *testing.T) {
	fac, observedLogs := observer.New(zapcore.InfoLevel)
	logger := zap.New(fac)

	loglvls := zaptool.NewLogLevels(logger)
	if loglvls.String() != "Internal.LogLevels:info" {
		t.Errorf("Initial log level for Internal.LogLevels is not info: %s", loglvls.String())
	}

	sublvls := zaptool.NewSubLogLevels("Childlog", loglvls)
	sublvls.Named("Core").Debug("should not log")
	sublvls.Named("Core").Info("should log")

	if sublvls.String() != "Childlog.Core:info,Internal.LogLevels:info" {
		t.Errorf("Initial log level for Childlog.Core is not info: %s", sublvls.String())
	}

	sublvls.SetLevel("Core", zapcore.DebugLevel)

	if sublvls.String() != "Childlog.Core:debug,Internal.LogLevels:info" {
		t.Errorf("Updated log level for Childlog.Core is not debug: %s", sublvls.String())
	}

	sublvls.Named("Core").Debug("should log")
	sublvls.Named("Core").Info("should log")

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

	if len(observedLogs.All()) != 3 {
		t.Errorf("should contain 3 log messages, instead contained %d messages", len(observedLogs.All()))
	}
}
