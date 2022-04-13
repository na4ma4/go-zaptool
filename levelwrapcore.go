package zaptool

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type levelWrapCore struct {
	lvl zap.AtomicLevel
	c   zapcore.Core
}

// Enabled returns true if the given level is at or above this level.
func (c *levelWrapCore) Enabled(lvl zapcore.Level) bool {
	return lvl >= c.lvl.Level()
}

// With adds structured context to the Core.
func (c *levelWrapCore) With(fields []zapcore.Field) zapcore.Core {
	return c.c.With(fields)
}

// Check determines whether the supplied Entry should be logged (using the
// embedded LevelEnabler and possibly some extra logic). If the entry
// should be logged, the Core adds itself to the CheckedEntry and returns
// the result.
//
// Callers must use Check before calling Write.
func (c *levelWrapCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.lvl.Enabled(ent.Level) {
		return ce.AddCore(ent, c.c)
	}

	return ce
}

// Write serializes the Entry and any Fields supplied at the log site and
// writes them to their destination.
//
// If called, Write should always log the Entry and Fields; it should not
// replicate the logic of Check.
//nolint:wrapcheck // simple wrapper for a *zap.Logger core.
func (c *levelWrapCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return c.c.Write(ent, fields)
}

// Sync flushes buffered logs (if any).
//nolint:wrapcheck // simple wrapper for a *zap.Logger core.
func (c *levelWrapCore) Sync() error {
	return c.c.Sync()
}
