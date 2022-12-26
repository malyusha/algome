package logger

import (
	"sync"
)

var globalLogger Logger

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}
func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}
func Fatal(msg string, args ...any) {
	GetLogger().Fatal(msg, args...)
}

func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

func GetLogger() Logger {
	var once sync.Once
	once.Do(func() {
		if globalLogger == nil {
			globalLogger = new(simpleLogger)
		}
	})

	return globalLogger
}

func SetGlobalLogger(l Logger) {
	var once sync.Once
	once.Do(func() {
		if globalLogger == nil {
			globalLogger = l
		}
	})
}
