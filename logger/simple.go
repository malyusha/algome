package logger

import (
	"fmt"
	"io"
	"os"
)

type logLevel int8

const (
	LevelDebug logLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l logLevel) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	}
	return ""
}

type simpleLogger struct {
	level  logLevel
	output io.StringWriter
}

func NewDefaultSimpleLogger(level logLevel) *simpleLogger {
	return &simpleLogger{
		level:  level,
		output: os.Stdout,
	}
}

func NewSimpleLogger(level logLevel, output io.StringWriter) *simpleLogger {
	return &simpleLogger{
		level:  level,
		output: output,
	}
}

func (s *simpleLogger) Debug(msg string, args ...any) {
	s.write(msg, LevelDebug, args...)
}

func (s *simpleLogger) Info(msg string, args ...any) {
	s.write(msg, LevelInfo, args...)
}

func (s *simpleLogger) Error(msg string, args ...any) {
	s.write(msg, LevelError, args...)
}

func (s *simpleLogger) Warn(msg string, args ...any) {
	s.write(msg, LevelWarn, args...)
}

func (s *simpleLogger) Fatal(msg string, args ...any) {
	s.write(msg, LevelFatal, args...)
	os.Exit(1)
}

func (s *simpleLogger) write(msg string, level logLevel, args ...any) {
	if s.level > level {
		return
	}

	s.output.WriteString(fmt.Sprintf(msg, args...))
	s.output.WriteString("\n")
}
