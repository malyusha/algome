package logger

import (
	"github.com/caarlos0/log"
)

func WithField(k string, v any) *log.Entry {
	return log.WithField(k, v)
}

func WithError(err error) *log.Entry {
	return log.WithError(err)
}

func Debug(msg string) {
	log.Debug(msg)
}

func Info(msg string) {
	log.Info(msg)
}
func Error(msg string) {
	log.Error(msg)
}
func Fatal(msg string) {
	log.Fatal(msg)
}

func Warn(msg string) {
	log.Warn(msg)
}
