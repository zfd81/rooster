package log

import "time"

const (
	defaultTimestampFormat = time.RFC3339
	FieldKeyMsg            = "msg"
	FieldKeyLevel          = "level"
	FieldKeyTime           = "time"
	FieldKeyTrace          = "TRACE"
	FieldKeyDebug          = "DEBUG"
	FieldKeyInfo           = "INFO"
	FieldKeyWarn           = "WARN"
	FieldKeyError          = "ERROR"
)

type Formatter interface {
	Format(*Entry) ([]byte, error)
}
