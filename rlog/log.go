package rlog

import (
	"log"
	"os"
)

type LogFormatter func(values ...interface{}) (messages []interface{})
type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
)

func defaultFormatter(values ...interface{}) (messages []interface{}) {
	return values
}

type Logger struct {
	logger    *log.Logger
	level     LogLevel
	formatter LogFormatter
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) SetFormatter(formatter LogFormatter) {
	l.formatter = formatter
}

func (l *Logger) Trace(values ...interface{}) {
	if l.level <= TRACE {

	}
}

func (l *Logger) Debug(values ...interface{}) {
	if l.level <= DEBUG {
		values = append([]interface{}{"DEBUG:"}, values[0:]...)
		l.logger.Println(l.formatter(values...)...)
	}
}

func (l *Logger) Info(values ...interface{}) {
	if l.level <= INFO {

	}
}

func (l *Logger) Warn(values ...interface{}) {
	if l.level <= WARN {

	}
}

func (l *Logger) Error(values ...interface{}) {
	if l.level <= ERROR {

	}
}

func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "", 0), DEBUG, defaultFormatter}
}
