package zaptool

import (
	"fmt"

	"go.uber.org/zap"
)

type SubLogLevels struct {
	prefix string
	logmgr LogManager
}

func NewSubLogLevels(name string, logmgr LogManager) *SubLogLevels {
	return &SubLogLevels{
		prefix: name,
		logmgr: logmgr,
	}
}

func (s *SubLogLevels) levelName(name string) string {
	return fmt.Sprintf("%s.%s", s.prefix, name)
}

func (s *SubLogLevels) NewLevel(name string) *zap.AtomicLevel {
	return s.logmgr.NewLevel(s.levelName(name))
}

func (s *SubLogLevels) Named(name string, opts ...interface{}) *zap.Logger {
	return s.logmgr.Named(s.levelName(name), opts...)
}

func (s *SubLogLevels) Iterator(f func(string, *zap.AtomicLevel) error) error {
	return s.logmgr.Iterator(f)
}

func (s *SubLogLevels) IsLogger(name string) bool {
	return s.logmgr.IsLogger(s.levelName(name))
}

func (s *SubLogLevels) SetLevel(name string, lvl interface{}) bool {
	return s.logmgr.SetLevel(s.levelName(name), lvl)
}

func (s *SubLogLevels) DeleteLevel(name string) {
	s.logmgr.DeleteLevel(s.levelName(name))
}

func (s *SubLogLevels) String() string {
	return s.logmgr.String()
}
