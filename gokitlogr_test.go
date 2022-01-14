package gokitlogr_test

import (
	"os"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-logr/logr"
	"github.com/tonglil/gokitlogr"
)

type E struct {
	str string
}

func (e E) Error() string {
	return e.str
}

func Helper(log logr.Logger, msg string) {
	helper2(log, msg)
}

func helper2(log logr.Logger, msg string) {
	log.WithCallDepth(2).Info(msg)
}

func ExampleNew() {
	kl := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	kl = level.NewFilter(kl, level.AllowInfo())
	kl = kitlog.With(kl, "caller", kitlog.Caller(5))

	var log logr.Logger = gokitlogr.New(&kl)
	log = log.WithName("MyName")
	log = log.WithValues("module", "example")

	log.Info("hello", "val1", 1, "val2", map[string]int{"k": 1})
	log.V(1).Info("you should see this")
	log.V(1).V(1).Info("you should NOT see this")
	log.Error(nil, "uh oh", "trouble", true, "reasons", []float64{0.1, 0.11, 3.14})
	log.Error(E{"an error occurred"}, "goodbye", "code", -1)
	Helper(log, "thru a helper")

	// Output:
	// {"caller":"gokitlogr_test.go:37","level":"info","logger":"MyName","module":"example","msg":"hello","val1":1,"val2":{"k":1}}
	// {"caller":"gokitlogr_test.go:38","level":"info","logger":"MyName","module":"example","msg":"you should see this"}
	// {"caller":"gokitlogr_test.go:40","error":null,"level":"error","logger":"MyName","module":"example","msg":"uh oh","reasons":[0.1,0.11,3.14],"trouble":true}
	// {"caller":"gokitlogr_test.go:41","code":-1,"error":"an error occurred","level":"error","logger":"MyName","module":"example","msg":"goodbye"}
	// {"caller":"gokitlogr_test.go:42","level":"info","logger":"MyName","module":"example","msg":"thru a helper"}
}
