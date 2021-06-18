package log

import "fmt"

type ConsoleAppender struct {
}

func (a *ConsoleAppender) Append(entry *Entry) (int, error) {
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return 0, err
	}
	fmt.Print(string(serialized))
	return 1, nil
}
