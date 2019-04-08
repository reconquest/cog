# cog

<img src="https://liquipedia.net/commons/images/7/71/Clockwerk_power_cog.jpg" width="200px" />

Package provides to structured log both for displaying it sanely in the stderr
logs as well as sending key-valued logs into any log storage such as ElasticSearch.

To setup basic stderr logging, use following snippet in your `main.go`:

```go
import "github.com/reconquest/cog"

// ...

stderr := lorg.NewLog()
stderr.SetIndentLines(true)
stderr.SetFormat(
    lorg.NewFormat("${time} ${level:[%s]:right:short} ${prefix}%s"),
)

if args["--debug"].(bool) {
    stderr.SetLevel(lorg.LevelDebug)
}

log = cog.NewLogger(stderr)

// use like that
log.Infof(nil, "message to log: %d", 1)
log.Infof(karma.Describe("key", "value"), "message to log: %d", 1)

err := errors.New("some error")
log.Fatalf(
    karma.Describe("context", "testing error").Reason(err),
    "message to log: %d",
    1,
)
```

To see more examples of how to use `karma` for structured logging and error
reporting, consider [looking at tests and examples][1].

## Motivation

Following package offers significant improvements above `logrus` and similar
structured loggers:

* Readable tree-like log entries in stderr, which makes easy to debug program,
  because log is more readable.
* Allows to use context errors, that used to describe errors on all call-stack
  levels to ease finding problems and fixing them.
* Sends logs to ES in key-value format like other structured loggers.
* Does not change ordering of key-values in the context.

[1]: https://github.com/reconquest/karma-go/blob/f802f635edd15c647995280b90f7de3e84ca8999/karma_test.go

# License

This project is licensed under the terms of the MIT license.
