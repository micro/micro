package logger

import (
	"os"

	"github.com/micro/go-micro/v3/logger"
)

type Level = logger.Level

var (
	TraceLevel = logger.TraceLevel
	DebugLevel = logger.DebugLevel
	InfoLevel  = logger.InfoLevel
	WarnLevel  = logger.WarnLevel
	ErrorLevel = logger.ErrorLevel
	FatalLevel = logger.FatalLevel
)

var (
	// DefaultLogger to use for logging
	DefaultLogger logger.Logger = logger.NewLogger()

	// fields to use when logging
	fields map[string]interface{}
)

func Info(args ...interface{}) {
	if !V(logger.InfoLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Log(logger.InfoLevel, args...)
}

func Infof(template string, args ...interface{}) {
	if !V(logger.InfoLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Logf(logger.InfoLevel, template, args...)
}

func Trace(args ...interface{}) {
	if !V(logger.TraceLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Log(logger.TraceLevel, args...)
}

func Tracef(template string, args ...interface{}) {
	if !V(logger.TraceLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Logf(logger.TraceLevel, template, args...)
}

func Debug(args ...interface{}) {
	if !V(logger.DebugLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Log(logger.DebugLevel, args...)
}

func Debugf(template string, args ...interface{}) {
	if !V(logger.DebugLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Logf(logger.DebugLevel, template, args...)
}

func Warn(args ...interface{}) {
	if !V(logger.WarnLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Log(logger.WarnLevel, args...)
}

func Warnf(template string, args ...interface{}) {
	if !V(logger.WarnLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Logf(logger.WarnLevel, template, args...)
}

func Error(args ...interface{}) {
	if !V(logger.ErrorLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Log(logger.ErrorLevel, args...)
}

func Errorf(template string, args ...interface{}) {
	if !V(logger.ErrorLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Logf(logger.ErrorLevel, template, args...)
}

func Fatal(args ...interface{}) {
	if !V(logger.FatalLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Log(logger.FatalLevel, args...)
	os.Exit(1)
}

func Fatalf(template string, args ...interface{}) {
	if !V(logger.FatalLevel, DefaultLogger) {
		return
	}
	DefaultLogger.Fields(fields).Logf(logger.FatalLevel, template, args...)
	os.Exit(1)
}

// V returns true if the given level is at or lower the current logger level
func V(lvl Level, logger logger.Logger) bool {
	return logger.Options().Level <= lvl
}
