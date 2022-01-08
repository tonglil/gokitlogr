# gokitlogr

[![Go Reference](https://pkg.go.dev/badge/github.com/tonglil/gokitlogr.svg)](https://pkg.go.dev/github.com/tonglil/gokitlogr)
<!-- ![test](https://github.com/tonglil/gokitlogr/workflows/test/badge.svg) -->
[![Go Report Card](https://goreportcard.com/badge/github.com/tonglil/gokitlogr)](https://goreportcard.com/report/github.com/tonglil/gokitlogr)

A [logr](https://github.com/go-logr/logr) LogSink implementation using [go-kit/log](https://github.com/go-kit/log).

## Usage

```go
import (
    "os"

    "github.com/go-logr/logr"
    "github.com/tonglil/gokitlogr"
    kitlog "github.com/go-kit/log"
)

func main() {
    kl := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stderr))
    kl = kitlog.With(kl, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.Caller(5))

    gokitlogr.NameFieldKey = "logger"
    gokitlogr.NameSeparator = "/"
    var log logr.Logger = gokitlogr.New(&kl)

    log = log.WithName("my app")
    log = log.WithValues("format", "json")

    log.Info("Logr in action!", "the answer", 42)
}
```

## Implementation Details

For the most part, concepts in go-kit/log correspond directly with those in logr.

Levels in logr correspond to custom debug levels in go-kit/log.
V(0) and V(1) are equivalent to go-kit/log's Info level, while V(2) is
equvalent to go-kit/log's Debug level. The Warn level is unused.
