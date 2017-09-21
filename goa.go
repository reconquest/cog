package main

import (
	"fmt"

	"github.com/goadesign/goa"
	"github.com/kovetskiy/lorg"

	karma "github.com/reconquest/karma-go"
)

// GoaLogger is required for using hierarchical logging inside Goa services
// itself.
type GoaLogger struct {
	Logger

	keyvals []interface{}
}

func (logger GoaLogger) Info(message string, keyvals ...interface{}) {
	logger.log(
		lorg.LevelDebug,
		logger.context(message, keyvals...),
		"{goa} "+message,
	)
}

func (logger GoaLogger) Error(message string, keyvals ...interface{}) {
	logger.log(
		lorg.LevelError,
		logger.context(message, keyvals...),
		"{goa} "+message,
	)
}

func (logger GoaLogger) New(keyvals ...interface{}) goa.LogAdapter {
	return &GoaLogger{
		Logger:  logger.Logger,
		keyvals: keyvals,
	}
}

func (logger GoaLogger) context(
	message string,
	keyvals ...interface{},
) *karma.Context {
	var pairs []interface{}

	pairs = append(pairs, logger.keyvals...)
	pairs = append(pairs, keyvals...)

	var context *karma.Context

	for i := 0; i < len(pairs); i += 2 {
		var (
			key   = fmt.Sprint(pairs[i])
			value = "MISSING"
		)

		if i+1 < len(pairs) {
			value = fmt.Sprint(pairs[i+1])
		}

		// without that goa will duplicate messages in log
		if message == "uncaught error" {
			if key == "msg" {
				continue
			}
		}

		context = context.Describe(key, value)
	}

	return context
}
