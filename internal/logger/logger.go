package logger

import (
	"context"
	"strings"
	"time"

	"github.com/shuvava/treehub/internal/apperrors"
)

// Level is logger severity level
type Level uint32

const (
	// PanicLevel level, highest level of severity
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`.
	FatalLevel
	// ErrorLevel level.
	ErrorLevel
	// WarnLevel level.
	WarnLevel
	// InfoLevel level.
	InfoLevel
	// DebugLevel level.
	DebugLevel
	// TraceLevel level.
	TraceLevel
)

// ContextKeyRequestID is correlationId key for context
const ContextKeyRequestID = "correlationId"

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
	// default value
	return "info"
}

// ToLogLevel converts string to log Level
func ToLogLevel(s string) Level {
	lvl := strings.ToLower(s)
	switch lvl {
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	case "trace":
		return TraceLevel
	}

	return InfoLevel
}

// GetRequestID will get reqID from a http request and return it as a string
func GetRequestID(ctx context.Context) string {
	reqID := ctx.Value(ContextKeyRequestID)

	if ret, ok := reqID.(string); ok {
		return ret
	}

	return ""
}

// CreateErrorAndLogIt create log record and throw an error
func CreateErrorAndLogIt(log Logger, code apperrors.AppErrorCode, descr string, err error) error {
	log.WithError(err).
		Error(descr)
	return apperrors.CreateError(code, descr, err)
}
