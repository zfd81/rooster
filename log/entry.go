package log

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

type Entry struct {
	Logger *Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level Level

	// Message passed to Trace, Debug, Info, Warn, Error
	Message string

	// err may contain a field formatting error
	err string
}

// Returns the bytes representation of this entry from the formatter.
func (entry *Entry) Bytes() ([]byte, error) {
	return entry.Logger.Formatter.Format(entry)
}

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	serialized, err := entry.Bytes()
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	fieldErr := entry.err
	for k, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch t.Kind() {
			case reflect.Func:
				isErrField = true
			case reflect.Ptr:
				isErrField = t.Elem().Kind() == reflect.Func
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", k)
			if fieldErr != "" {
				fieldErr = entry.err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			data[k] = v
		}
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time, err: fieldErr}
}

// Overrides the time of the Entry.
func (entry *Entry) WithTime(t time.Time) *Entry {
	dataCopy := make(Fields, len(entry.Data))
	for k, v := range entry.Data {
		dataCopy[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: dataCopy, Time: t, err: entry.err}
}

func (entry *Entry) write() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	if _, err := entry.Logger.Out.Append(entry); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}

func (entry Entry) log(level Level, msg string) {
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}
	entry.Level = level
	entry.Message = msg
	entry.write()
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprint(args...))
	}
}

func (entry *Entry) Trace(args ...interface{}) {
	entry.Log(TRACE, args...)
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.Log(DEBUG, args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Log(INFO, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.Log(WARN, args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.Log(ERROR, args...)
}

func (entry *Entry) Logf(level Level, format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.Log(level, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Tracef(format string, args ...interface{}) {
	entry.Logf(TRACE, format, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Logf(DEBUG, format, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Logf(INFO, format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Logf(WARN, format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Logf(ERROR, format, args...)
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
	}
}
