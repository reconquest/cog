package main

import (
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/karma-go"
)

// Logger provides structured logging methods, based on karma package.
type Logger struct {
	output *lorg.Log
}

func (logger *Logger) Tracef(
	context *karma.Context,
	message string,
	args ...interface{},
) {
	logger.log(lorg.LevelTrace, context, message, args...)
}

func (logger *Logger) Debugf(
	context *karma.Context,
	message string,
	args ...interface{},
) {
	logger.log(lorg.LevelDebug, context, message, args...)
}

func (logger *Logger) Infof(
	context *karma.Context,
	message string,
	args ...interface{},
) {
	logger.log(lorg.LevelInfo, context, message, args...)
}

func (logger *Logger) Warningf(
	err error,
	message string,
	args ...interface{},
) {
	logger.log(lorg.LevelWarning, err, message, args...)
}

func (logger *Logger) Errorf(
	err error,
	message string,
	args ...interface{},
) {
	logger.log(lorg.LevelError, err, message, args...)
}

func (logger *Logger) Fatalf(
	err error,
	message string,
	args ...interface{},
) {
	logger.log(lorg.LevelFatal, err, message, args...)
}

func (logger *Logger) log(
	level lorg.Level,
	reason interface{},
	message string,
	args ...interface{},
) {
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

	logger.send(level, hierarchy)
	logger.display(level, hierarchy)
}

func (logger *Logger) display(level lorg.Level, hierarchy karma.Hierarchical) {
	loggers := map[lorg.Level]func(...interface{}){
		lorg.LevelTrace:   logger.output.Trace,
		lorg.LevelDebug:   logger.output.Debug,
		lorg.LevelInfo:    logger.output.Info,
		lorg.LevelWarning: logger.output.Warning,
		lorg.LevelError:   logger.output.Error,
		lorg.LevelFatal:   logger.output.Fatal,
	}

	log := loggers[level]

	log(hierarchy.String())
}

func (logger *Logger) send(level lorg.Level, hierarchy karma.Hierarchical) {
	// TODO: Add ElasticSearch logger here.
}
