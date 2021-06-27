package log

import (
	"sync"
	"sync/atomic"
	"time"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

type Level uint32

const (
	FATAL Level = iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

type Logger struct {
	Out       Appender
	Formatter Formatter
	Level     Level

	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap
}

func (logger *Logger) SetNoLock() {
	logger.mu.Disable()
}

func (logger *Logger) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}

// SetLevel sets the logger level.
func (logger *Logger) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
}

// GetLevel returns the logger level.
func (logger *Logger) GetLevel() Level {
	return logger.level()
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (logger *Logger) IsLevelEnabled(level Level) bool {
	return logger.level() >= level
}

// SetOutput sets the logger output.
func (logger *Logger) SetOutput(output Appender) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Out = output
}

// SetFormatter sets the logger formatter.
func (logger *Logger) SetFormatter(formatter Formatter) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Formatter = formatter
}

func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
}

// WithField allocates a new entry and adds a field to it.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithField(key, value)
}

// Overrides the time of the log entry.
func (logger *Logger) WithTime(t time.Time) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithTime(t)
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.Log(TRACE, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Log(DEBUG, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Log(INFO, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Log(WARN, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Log(ERROR, args...)
}

func (logger *Logger) Log(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Log(level, args...)
	}
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.Logf(TRACE, format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.Logf(DEBUG, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Logf(INFO, format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.Logf(WARN, format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.Logf(ERROR, format, args...)
}

func (logger *Logger) Logf(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Logf(level, format, args...)
	}
}

func (logger *Logger) newEntry() *Entry {
	return NewEntry(logger)
}

func New() *Logger {
	return &Logger{
		Out:       new(ConsoleAppender),
		Formatter: &TextFormatter{TimestampFormat: "2006-01-02 15:04:05"},
		Level:     INFO,
	}
}
