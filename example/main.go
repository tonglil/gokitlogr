package main

import (
	"os"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-logr/logr"
	"github.com/tonglil/gokitlogr"
)

type e struct {
	str string
}

func (e e) Error() string {
	return e.str
}

func helper(log logr.Logger, msg string) {
	helper2(log, msg)
}

func helper2(log logr.Logger, msg string) {
	log.WithCallDepth(2).Info(msg)
}

func main() {
	kl := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	kl = level.NewFilter(kl, level.AllowInfo())
	kl = kitlog.With(kl, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.Caller(5))

	gokitlogr.NameFieldKey = "logger"
	gokitlogr.NameSeparator = "/"
	var log logr.Logger = gokitlogr.New(&kl)

	log = log.WithName("MyName")
	example(log.WithValues("module", "example"))
}

// example only depends on logr except when explicitly breaking the
// abstraction. Even that part is written so that it works with non-zap
// loggers.
func example(log logr.Logger) {
	log.Info("hello", "val1", 1, "val2", map[string]int{"k": 1})
	log.V(1).Info("you should see this")
	log.V(1).V(1).Info("you should NOT see this")
	log.Error(nil, "uh oh", "trouble", true, "reasons", []float64{0.1, 0.11, 3.14})
	log.Error(e{"an error occurred"}, "goodbye", "code", -1)
	helper(log, "thru a helper")

	if gokitLogger, ok := log.GetSink().(gokitlogr.Underlier); ok {
		kl := gokitLogger.GetUnderlying()
		kl = kitlog.With(kl, "ts", kitlog.DefaultTimestamp, "caller", kitlog.DefaultCaller)
		kl.Log("msg", "go-kit/log now")
	}
}
