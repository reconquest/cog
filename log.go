package logger

import (
	"encoding/json"
	"fmt"

	"github.com/kovetskiy/lorg"
	"github.com/reconquest/karma-go"
)

type Sender func(lorg.Level, karma.Hierarchical) error
type Displayer func(lorg.Level, karma.Hierarchical)

// Logger provides structured logging methods, based on karma package.
type Logger struct {
	*lorg.Log
	sender    Sender
	displayer Displayer
}

func NewLogger(output *lorg.Log) *Logger {
	return &Logger{output, nil, nil}
}

// SetSender sets given function as callack for every log line
func (logger *Logger) SetSender(sender Sender) {
	logger.sender = sender
}

// SetDisplayer sets given function as callack when displays log line
func (logger *Logger) SetDisplayer(display Displayer) {
	logger.displayer = display
}

func (logger *Logger) Tracef(
	context *karma.Context,
	message string,
	args ...interface{},
) {
	logger.Write(lorg.LevelTrace, context, message, args...)
}

func (logger *Logger) Debugf(
	context *karma.Context,
	message string,
	args ...interface{},
) {
	logger.Write(lorg.LevelDebug, context, message, args...)
}

func (logger *Logger) Infof(
	context *karma.Context,
	message string,
	args ...interface{},
) {
	logger.Write(lorg.LevelInfo, context, message, args...)
}

func (logger *Logger) Warningf(
	err error,
	message string,
	args ...interface{},
) {
	logger.Write(lorg.LevelWarning, err, message, args...)
}

func (logger *Logger) Errorf(
	err error,
	message string,
	args ...interface{},
) {
	logger.Write(lorg.LevelError, err, message, args...)
}

func (logger *Logger) Fatalf(
	err error,
	message string,
	args ...interface{},
) {
	logger.Write(lorg.LevelFatal, err, message, args...)
}

func (logger *Logger) Write(
	level lorg.Level,
	reason interface{},
	message string,
	args ...interface{},
) {
	if logger == nil {
		return
	}

	var hierarchy karma.Karma

	switch reason := reason.(type) {
	case karma.Hierarchical:
		hierarchy = karma.Format(reason, message, args...)

	case *karma.Context:
		hierarchy = karma.Format(nil, message, args...)
		hierarchy.Context = reason

	default:
		hierarchy = karma.Format(reason, message, args...)
	}

	logger.Display(level, hierarchy)
	err := logger.Send(level, hierarchy)
	if err != nil {
		logger.Display(
			lorg.LevelError,
			karma.Format(err, "error while sending log"),
		)
	}
}

func (logger *Logger) Display(level lorg.Level, hierarchy karma.Hierarchical) {
	if logger.displayer != nil {
		logger.displayer(level, hierarchy)
	} else {
		Display(logger, level, hierarchy)
	}
}

func (logger *Logger) Send(
	level lorg.Level,
	hierarchy karma.Hierarchical,
) error {
	if logger.sender != nil {
		return logger.sender(level, hierarchy)
	}

	return nil
}

func Display(logger *Logger, level lorg.Level, hierarchy karma.Hierarchical) {
	loggers := map[lorg.Level]func(...interface{}){
		lorg.LevelTrace:   logger.Trace,
		lorg.LevelDebug:   logger.Debug,
		lorg.LevelInfo:    logger.Info,
		lorg.LevelWarning: logger.Warning,
		lorg.LevelError:   logger.Error,
		lorg.LevelFatal:   logger.Fatal,
	}

	log := loggers[level]

	log(hierarchy.String())
}

func (logger *Logger) TraceJSON(obj interface{}) (encoded string) {
	if logger.GetLevel() != lorg.LevelTrace {
		return ""
	}

	defer func() {
		err := recover()
		if err != nil {
			encoded = fmt.Sprintf(
				"%#v (unable to encode to json: %s)",
				obj, err,
			)
		}
	}()

	contents, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		return fmt.Sprintf(
			"%#v (unable to encode to json: %s)",
			obj, err,
		)
	}

	return string(contents)
}
