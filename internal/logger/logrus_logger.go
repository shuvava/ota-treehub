package logger

import (
	"context"
	"io"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
)

// LogrusLogger logrus.Logger implementation of interface
type LogrusLogger struct {
	Logger
	entry         *logrus.Entry
	CorrelationID string
}

// SetContext logger context (Operation, CorrelationID)
func (l LogrusLogger) SetContext(operation string) Logger {
	correlationID := l.CorrelationID
	if correlationID == "" {
		correlationID = uuid.New().String()
	}
	return LogrusLogger{entry: l.entry.
		WithField("Operation", operation).
		WithField("CorrelationID", correlationID)}
}

// WithField adds a filed to log entry
func (l LogrusLogger) WithField(key string, value interface{}) Logger {
	return LogrusLogger{entry: l.entry.WithField(key, value)}
}

// WithFields adds a struct of fields to the log entry.
func (l LogrusLogger) WithFields(fields Fields) Logger {
	return LogrusLogger{entry: l.entry.WithFields(logrus.Fields(fields))}
}

// WithError Add an error as single field to the log entry.
func (l LogrusLogger) WithError(err error) Logger {
	return LogrusLogger{entry: l.entry.WithError(err)}
}

// WithContext adds a context to the Entry.
func (l LogrusLogger) WithContext(ctx context.Context) Logger {
	corrID := GetRequestID(ctx)
	log := l.entry.WithContext(ctx)
	if corrID != "" {
		log = log.WithField(ContextKeyRequestID, corrID)
	}
	return LogrusLogger{entry: log}
}

// Trace creates logs entry with Trace level
func (l LogrusLogger) Trace(args ...interface{}) {
	l.entry.Trace(args...)
}

// Debug creates logs entry with Debug level
func (l LogrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Info creates logs entry with Info level
func (l LogrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Warn creates logs entry with Warn level
func (l LogrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Error creates logs entry with Error level
func (l LogrusLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Fatal creates logs entry with Fatal level
func (l LogrusLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Panic creates logs entry with Panic level
func (l LogrusLogger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

// SetOutput sets the output to desired io.Writer like file, stdout, stderr etc
func (l *LogrusLogger) SetOutput(w io.Writer) {
	l.entry.Logger.Out = w
}

// SetLevel sets logger level
func (l LogrusLogger) SetLevel(level Level) error {
	lvl, err := logrus.ParseLevel(ParseLevel(level))
	if err != nil {
		return err
	}

	l.entry.Logger.Level = lvl
	return nil
}

// GetLevel returns current logging level
func (l LogrusLogger) GetLevel() Level {
	return ParseLogrusLevel(l.entry.Logger.Level)
}

// TrackFuncTime creates log record with func execution time
// require debug level or higher
// usage:
// func SomeFunction(list *[]string) {
//    defer TimeTrack(time.Now())
//    // Do whatever you want.
// }
func (l LogrusLogger) TrackFuncTime(start time.Time) {
	lvl := l.GetLevel()
	if lvl > DebugLevel {
		return
	}
	elapsed := time.Since(start)
	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)
	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	//runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	//fname := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	fname := funcObj.Name()

	l.WithField("func", fname).
		WithField("executionTime", elapsed).
		Debug("func execution completed")
}

// NewLogrusLogger creates new instance of LogrusLogger
func NewLogrusLogger(l logrus.Level) LogrusLogger {
	log := logrus.New()
	log.SetLevel(l)

	return LogrusLogger{
		entry: logrus.NewEntry(log),
	}
}

// NewNopLogger returns a logger that discards all log messages.
func NewNopLogger() Logger {
	log := logrus.New()
	log.Out = ioutil.Discard
	return LogrusLogger{
		entry: logrus.NewEntry(log),
	}
}

// ParseLogrusLevel takes a string level and returns the Logrus log level constant.
func ParseLogrusLevel(lvl logrus.Level) Level {
	switch lvl {
	case logrus.PanicLevel:
		return PanicLevel
	case logrus.FatalLevel:
		return FatalLevel
	case logrus.ErrorLevel:
		return ErrorLevel
	case logrus.WarnLevel:
		return WarnLevel
	case logrus.InfoLevel:
		return InfoLevel
	case logrus.DebugLevel:
		return DebugLevel
	case logrus.TraceLevel:
		return TraceLevel
	}
	// Default value
	return InfoLevel
}
