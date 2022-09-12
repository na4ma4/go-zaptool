package zaptool

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevels provides a wrapper for multiple *zap.Logger levels,
// the individual loggers are not kept, but levels are kept
// indexed by name.
type LogLevels struct {
	logger  *zap.Logger
	iLogger *zap.Logger
	levels  map[string]*zap.AtomicLevel
	lock    sync.RWMutex
}

// NewLogLevels returns a new LogLevels ready for use.
func NewLogLevels(logger *zap.Logger) *LogLevels {
	out := &LogLevels{
		logger:  logger,
		iLogger: logger,
		levels:  map[string]*zap.AtomicLevel{},
		lock:    sync.RWMutex{},
	}
	out.iLogger = out.Named("Internal.LogLevels")

	return out
}

// NewLevel returns a zap.AtomicLevel reference to the stored named level.
func (a *LogLevels) NewLevel(name string) *zap.AtomicLevel {
	a.lock.Lock()
	defer a.lock.Unlock()

	if v, ok := a.levels[name]; ok {
		return v
	}

	atom := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	a.levels[name] = &atom

	return &atom
}

func (a *LogLevels) parseLevel(v interface{}) (zapcore.Level, bool) {
	switch lvl := v.(type) {
	case zapcore.Level:
		return lvl, true
	case *zap.AtomicLevel:
		return lvl.Level(), true
	case zap.AtomicLevel:
		return lvl.Level(), true
	case string:
		if level, err := zap.ParseAtomicLevel(lvl); err == nil {
			return level.Level(), true
		}

		return zapcore.InfoLevel, false
	}

	return zapcore.InfoLevel, false
}

// can't remember what this was for, but it's unused code now.
// func (a *LogLevels) isWildcard(c string) bool {
// 	if strings.HasSuffix(c, "*") {
// 		return true
// 	}

// 	if strings.HasPrefix(c, "*") {
// 		return true
// 	}

// 	return false
// }

// doesKeyMatch tests if a key matches.
func (a *LogLevels) doesKeyMatch(key, check string) bool {
	if strings.EqualFold(key, check) {
		return true
	}

	if len(check) == 1 && check == "*" {
		return true
	}

	switch {
	case strings.HasPrefix(check, "*") && strings.HasSuffix(check, "*"):
		// wildcard both ends (can't be a single * otherwise initial check would fail)
		return strings.Contains(key, check[1:len(check)-1])
	case strings.HasPrefix(check, "*"):
		// wildcard at start
		return strings.HasSuffix(key, check[1:])
	case strings.HasSuffix(check, "*"):
		// wildcard at end
		return strings.HasPrefix(key, check[:len(check)-1])
	}

	return false
}

// SetLevel attempts to set the level supplied, it will attempt to typecast the value
// against string, zapcore.Level and *zap.AtomicLevel.
func (a *LogLevels) SetLevel(name string, lvl interface{}) bool {
	a.iLogger.Debug("SetLevel", zap.String("name", name))

	found := false

	a.lock.Lock()
	defer a.lock.Unlock()

	for itemKey, val := range a.levels {
		if a.doesKeyMatch(itemKey, name) {
			if level, ok := a.parseLevel(lvl); ok {
				a.iLogger.Debug(
					"setting level for name (wildcard)",
					zap.String("name", name),
					zap.String("match", itemKey),
					zap.String("level", level.String()),
				)
				val.SetLevel(level)

				found = true
			}
		}
	}
	// }

	return found
}

// String returns a string representation of the currently stored loggers and their levels.
func (a *LogLevels) String() string {
	a.lock.RLock()
	defer a.lock.RUnlock()

	out := []string{}

	for k, v := range a.levels {
		out = append(out, fmt.Sprintf("%s:%s", k, v.String()))
	}

	sort.Strings(out)

	return strings.Join(out, ",")
}

// Named returns a named *zap.Logger if any additional parameters are specified it will
// try to determine if they represent a log level (by string, zapcore.Level or *zap.AtomicLevel).
func (a *LogLevels) Named(name string, opts ...interface{}) *zap.Logger {
	lvl := a.NewLevel(name)

	if a.logger.Core().Enabled(zapcore.DebugLevel) {
		lvl.SetLevel(zapcore.DebugLevel)
	}

	for _, opt := range opts {
		switch v := opt.(type) {
		case zapcore.Level:
			a.SetLevel(name, v.String())
		case zap.AtomicLevel, *zap.AtomicLevel:
			a.SetLevel(name, v)
		}
	}

	return a.logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return &levelWrapCore{
			lvl: *lvl,
			c:   c,
		}
	})).Named(name)
}
