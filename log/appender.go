package log

type Appender interface {
	Append(*Entry) (n int, err error) //打印日志
}
