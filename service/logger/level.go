// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/logger/level.go

package logger

import (
	"fmt"
	"os"
)

type Level int8

const (
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel Level = iota - 2
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// InfoLevel is the default logging priority.
	// General operational entries about what's going on inside the application.
	InfoLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	ErrorLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. highest level of severity.
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	}
	return ""
}

// Enabled returns true if the given level is at or above this level.
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// GetLevel converts a level string into a logger Level value.
// returns an error if the input string does not match known values.
func GetLevel(levelStr string) (Level, error) {
	switch levelStr {
	case TraceLevel.String():
		return TraceLevel, nil
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	case FatalLevel.String():
		return FatalLevel, nil
	}
	return InfoLevel, fmt.Errorf("Unknown Level String: '%s', defaulting to InfoLevel", levelStr)
}

func Info(args ...interface{}) {
	DefaultLogger.Log(InfoLevel, args...)
}

func Infof(template string, args ...interface{}) {
	DefaultLogger.Logf(InfoLevel, template, args...)
}

func Trace(args ...interface{}) {
	DefaultLogger.Log(TraceLevel, args...)
}

func Tracef(template string, args ...interface{}) {
	DefaultLogger.Logf(TraceLevel, template, args...)
}

func Debug(args ...interface{}) {
	DefaultLogger.Log(DebugLevel, args...)
}

func Debugf(template string, args ...interface{}) {
	DefaultLogger.Logf(DebugLevel, template, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Log(WarnLevel, args...)
}

func Warnf(template string, args ...interface{}) {
	DefaultLogger.Logf(WarnLevel, template, args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Log(ErrorLevel, args...)
}

func Errorf(template string, args ...interface{}) {
	DefaultLogger.Logf(ErrorLevel, template, args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Log(FatalLevel, args...)
	os.Exit(1)
}

func Fatalf(template string, args ...interface{}) {
	DefaultLogger.Logf(FatalLevel, template, args...)
	os.Exit(1)
}

// Returns true if the given level is at or lower the current logger level
func V(lvl Level, log Logger) bool {
	l := DefaultLogger
	if log != nil {
		l = log
	}
	return l.Options().Level <= lvl
}
