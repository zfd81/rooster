package log

import (
	"bytes"

	"github.com/spf13/cast"
)

type TextFormatter struct {

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var b = &bytes.Buffer{}
	switch entry.Level {
	case TRACE:
		b.WriteString(FieldKeyTrace)
	case DEBUG:
		b.WriteString(FieldKeyDebug)
	case INFO:
		b.WriteString(FieldKeyInfo)
	case WARN:
		b.WriteString(FieldKeyWarn)
	case ERROR:
		b.WriteString(FieldKeyError)
	default:
		b.WriteString(FieldKeyInfo)
	}
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	b.WriteString("[")
	b.WriteString(entry.Time.Format(timestampFormat))
	b.WriteString("]")
	if len(entry.Data) > 0 {
		b.WriteString("{")
	}
	cnt := 0
	for k, v := range entry.Data {
		if cnt > 0 {
			b.WriteString(",")
		} else {
			cnt++
		}
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(cast.ToString(v))
	}
	if len(entry.Data) > 0 {
		b.WriteString("}")
	}
	b.WriteString(" ")
	b.WriteString(entry.Message)
	b.WriteByte('\n')
	return b.Bytes(), nil
}
