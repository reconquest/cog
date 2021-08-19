package cog // import "github.com/reconquest/cog"

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kovetskiy/lorg"
	"github.com/reconquest/karma-go"
)

type (
	Sender    func(lorg.Level, karma.Hierarchical) error
	Displayer func(lorg.Level, karma.Hierarchical)
	Exiter    func(int)
)

// Logger provides structured logging methods, based on karma package.
type Logger struct {
	*lorg.Log
	sender    Sender
	displayer Displayer
	exiter    Exiter
	exitcode  *int
}

func NewLogger(output *lorg.Log) *Logger {
	logger := &Logger{Log: output}
	output.SetExiter(logger.handleExit)
	return logger
}

func (logger *Logger) NewChild() *Logger {
	child := NewLogger(logger.Log.NewChild())
	child.SetSender(logger.sender)
	child.SetDisplayer(logger.displayer)
	return child
}

func (logger *Logger) NewChildWithPrefix(prefix string) *Logger {
	child := NewLogger(logger.Log.NewChildWithPrefix(prefix))
	child.SetSender(logger.sender)
	child.SetDisplayer(logger.displayer)
	child.SetExiter(logger.exiter)
	return child
}

// SetSender sets given function as callack for every log line
func (logger *Logger) SetSender(sender Sender) {
	logger.sender = sender
}

func (logger *Logger) SetExiter(exiter func(int)) {
	logger.exiter = exiter
}

// SetDisplayer sets given function as callack when displays log line
func (logger *Logger) SetDisplayer(display Displayer) {
	logger.displayer = display
}

func (logger *Logger) handleExit(code int) {
	logger.exitcode = &code
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

func (logger *Logger) Trace(
	args ...interface{},
) {
	logger.write(lorg.LevelTrace, args...)
}

func (logger *Logger) Debug(
	args ...interface{},
) {
	logger.write(lorg.LevelDebug, args...)
}

func (logger *Logger) Info(
	args ...interface{},
) {
	logger.write(lorg.LevelInfo, args...)
}

func (logger *Logger) Warning(
	args ...interface{},
) {
	logger.write(lorg.LevelWarning, args...)
}

func (logger *Logger) Error(
	args ...interface{},
) {
	logger.write(lorg.LevelError, args...)
}

func (logger *Logger) Fatal(
	args ...interface{},
) {
	logger.write(lorg.LevelFatal, args...)
}

func (logger *Logger) write(level lorg.Level, args ...interface{}) {
	if len(args) == 0 {
		return
	}

	if logger == nil {
		return
	}

	var hierarchy karma.Hierarchical

	arg0, isKarma := args[0].(karma.Hierarchical)
	if isKarma {
		hierarchy = arg0
	}

	if len(args) > 1 || !isKarma {
		reason := hierarchy

		var message string
		if isKarma {
			message = fmt.Sprint(args[1:]...)
		} else {
			message = fmt.Sprint(args...)
		}

		hierarchy = karma.Karma{
			Reason:  reason,
			Message: message,
		}
	}

	logger.displaySend(level, hierarchy)
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

	var hierarchy karma.Hierarchical

	switch reason := reason.(type) {
	case karma.Hierarchical:
		hierarchy = karma.Format(reason, message, args...)

	case *karma.Context:
		newHierarchy := karma.Format(nil, message, args...)
		newHierarchy.Context = reason

		hierarchy = newHierarchy

	default:
		hierarchy = karma.Format(reason, message, args...)
	}

	logger.displaySend(level, hierarchy)
}

func (logger *Logger) displaySend(
	level lorg.Level,
	hierarchy karma.Hierarchical,
) {
	logger.Display(level, hierarchy)

	err := logger.Send(level, hierarchy)
	if err != nil {
		logger.Display(
			lorg.LevelError,
			karma.Format(err, "error while sending log"),
		)
	}

	if logger.exitcode != nil {
		if logger.exiter != nil {
			logger.exiter(*logger.exitcode)
		} else {
			os.Exit(*logger.exitcode)
		}
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
		lorg.LevelTrace:   logger.Log.Trace,
		lorg.LevelDebug:   logger.Log.Debug,
		lorg.LevelInfo:    logger.Log.Info,
		lorg.LevelWarning: logger.Log.Warning,
		lorg.LevelError:   logger.Log.Error,
		lorg.LevelFatal:   logger.Log.Fatal,
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
