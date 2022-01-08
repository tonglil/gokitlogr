// Package gokitlogr defines an implementation of the github.com/go-logr/logr
// interfaces built on top of go-kit/log (https://github.com/go-kit/log).
//
// Usage
//
// A new logr.Logger can be constructed from an existing log.Logger using
// the New function:
//
//  log := gokitlogr.New(someGoKitLogger)
//
// Implementation Details
//
// For the most part, concepts in go-kit/log correspond directly with those in
// logr.
//
// Levels in logr correspond to custom debug levels in go-kit/log.
// V(0) and V(1) are equivalent to go-kit/log's Info level, while V(2) is
// equvalent to go-kit/log's Debug level. The Warn level is unused.
package gokitlogr

import (
	"fmt"

	kitlog "github.com/go-kit/log"
	kitlevel "github.com/go-kit/log/level"

	"github.com/go-logr/logr"
)

// TODO: as options, see:
// https://github.com/go-logr/zapr/blob/master/zapr.go#L245
var (
	// NameFieldKey is the field key for logr.WithName
	NameFieldKey = "logger"
	// NameSeparator separates names for logr.WithName
	NameSeparator = "/"
	// ErrorFieldKey is the field key for logr.Error
	ErrorFieldKey = "error"
	// CallerFieldKey is the field key for call site information
	// When using LogfmtLogger, this should be set to a different key
	// than configured in go-kit/log.Logger to represent the true value
	// of WithCallDepth call sites.
	// When using JSONLogger, this should be set to the same key
	// configured in go-kit/log.Logger to overwrite it with the true value
	// of WithCallDepth call sites.
	CallerFieldKey = "caller"
)

var (
	_ logr.LogSink          = &kitlogger{}
	_ logr.CallDepthLogSink = &kitlogger{}
)

// New returns a logr.Logger with logr.LogSink implemented by go-kit/log.
func New(l *kitlog.Logger) logr.Logger {
	ls := newKitLogger(l)
	return logr.New(ls)
}

// kitlogger implements the LogSink interface.
type kitlogger struct {
	kl     *kitlog.Logger
	name   string
	values []interface{}
	depth  int
}

// newKitLogger returns a logr.LogSink implemented by go-kit/log.
func newKitLogger(l *kitlog.Logger) *kitlogger {
	return &kitlogger{kl: l}
}

// Enabled tests whether this LogSink is enabled at the specified V-level.
func (l kitlogger) Enabled(level int) bool {
	// Optimization: Info() will check level internally.
	const debugLevel = 2
	return level <= debugLevel
}

// WithName returns a new LogSink with the specified name appended in NameFieldName.
// Name elements are separated by NameSeparator.
func (l kitlogger) WithName(name string) logr.LogSink {
	if l.name != "" {
		l.name += NameSeparator + name
	} else {
		l.name = name
	}
	return &l
}

// WithValues returns a new LogSink with additional key/value pairs.
// NOTE: look at "github.com/go-logr/logr/funcr".Formatter.AddValues for a more exhaustive implementation.
func (l kitlogger) WithValues(keysAndValues ...interface{}) logr.LogSink {
	l.values = append(l.values, keysAndValues...)
	return &l
}

// Info logs a non-error message at specified V-level with the given key/value pairs as context.
// Duplicate key/values are not allowed for JSONLogger, and last key overwrites previous values.
func (l *kitlogger) Info(level int, msg string, keysAndValues ...interface{}) {
	kvs := append(l.values, keysAndValues...)
	kvs = append(kvs, "msg", msg)
	if l.name != "" {
		kvs = append(kvs, NameFieldKey, l.name)
	}
	kvs = defaultRender(kvs)
	if level > 1 {
		kitlevel.Debug(*l.kl).Log(kvs...)
		// NOTE: WON'T DO
		// } else if level == 0 {
		// 	kitlevel.Warn(*l.kl).Log(kvs...)
	} else {
		kitlevel.Info(*l.kl).Log(kvs...)
	}
}

// Error logs an error, with the given message and key/value pairs as context.
// Duplicate key/values are not allowed for JSONLogger, and last key overwrites previous values.
func (l *kitlogger) Error(err error, msg string, keysAndValues ...interface{}) {
	kvs := append(l.values, keysAndValues...)
	kvs = append(kvs, "msg", msg, ErrorFieldKey, err)
	if l.name != "" {
		kvs = append(kvs, NameFieldKey, l.name)
	}
	kvs = defaultRender(kvs)
	kitlevel.Error(*l.kl).Log(kvs...)
}

// defaultRender supports logr.Marshaler and fmt.Stringer.
// From: https://github.com/go-logr/zerologr/blob/33354eecabe37c0eacbba4df530534fed6d8a3f3/zerologr.go#L150-L162
func defaultRender(keysAndValues []interface{}) []interface{} {
	for i, n := 1, len(keysAndValues); i < n; i += 2 {
		value := keysAndValues[i]
		switch v := value.(type) {
		case logr.Marshaler:
			keysAndValues[i] = v.MarshalLog()
		case fmt.Stringer:
			keysAndValues[i] = v.String()
		}
	}
	return keysAndValues
}

// Init receives runtime info about the logr library.
func (l *kitlogger) Init(info logr.RuntimeInfo) {
	l.depth = info.CallDepth + 4
}

// WithCallDepth returns a new LogSink that offsets the call
// stack by the specified number of frames when logging call
// site information.
//
// If depth is 0, the LogSink should skip exactly the number
// of call frames defined in RuntimeInfo.CallDepth when Info
// or Error are called, i.e. the attribution should be to the
// direct caller of Logger.Info or Logger.Error.
//
// If depth is 1 the attribution should skip 1 call frame, and so on.
// Successive calls to this are additive.
//
// WithCallDepth will be duplicated for LogfmtLogger if go-kit/log.Caller
// or go-kit/log.DefaultCaller is used, and present two key/value pairs,
// one of which will be incorrect.
func (l kitlogger) WithCallDepth(depth int) logr.LogSink {
	newLogger := kitlog.With(*l.kl, CallerFieldKey, kitlog.Caller(l.depth+depth))
	l.kl = &newLogger
	return &l
}

// Underlier exposes access to the underlying logging implementation.  Since
// callers only have a logr.Logger, they have to know which implementation is
// in use, so this interface is less of an abstraction and more of way to test
// type conversion.
type Underlier interface {
	GetUnderlying() kitlog.Logger
}

// GetUnderlying returns the go-kit/log.Logger underneath this logSink.
func (l *kitlogger) GetUnderlying() kitlog.Logger {
	return *l.kl
}
