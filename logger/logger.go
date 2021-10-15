package logger

import (
	"context"
	"time"
)

// Level is logger severity level
type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
	UnknownLevel
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// Logger is Generic logger interface
type Logger interface {
	// SetLevel logging severity level
	SetLevel(level Level) error
	// GetLevel returns current logging level
	GetLevel() Level
	// SetContext set logger context
	SetContext(operation string) Logger

	// WithField adds a filed to log entry
	WithField(key string, value interface{}) Logger
	// WithFields adds a struct of fields to the log entry.
	WithFields(fields Fields) Logger
	// WithError Add an error as single field to the log entry.
	WithError(err error) Logger
	WithContext(ctx context.Context) Logger

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	// TrackFuncTime creates log record with func execution time
	TrackFuncTime(start time.Time)
}

// ParseLevel parse Level to string value
func ParseLevel(lvl Level) string {
	switch lvl {
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case TraceLevel:
		return "trace"
	}

	return "unknown"
}
