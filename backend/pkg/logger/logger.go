package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Lg struct {
	logger *logrus.Logger
}

func New(level string) *Lg {
	var l logrus.Level

	switch strings.ToLower(level) {
	case "error":
		l = logrus.ErrorLevel
	case "warn":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.SetLevel(l)
	logger.SetFormatter(new(logrus.JSONFormatter))
	logger.SetOutput(os.Stdout)

	return &Lg{
		logger: logger,
	}
}

// Debug -.
func (l *Lg) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

// Info -.
func (l *Lg) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

// Warn -.
func (l *Lg) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

// Error -.
func (l *Lg) Error(message interface{}, args ...interface{}) {
	if l.logger.GetLevel() == logrus.DebugLevel {
		l.Debug(message, args...)
		return
	}

	l.msg("error", message, args...)
}

// Fatal -.
func (l *Lg) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}

func (l *Lg) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info(fmt.Sprint(message))
	} else {
		l.logger.Info(fmt.Sprintf(message, args...))
	}
}

func (l *Lg) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
