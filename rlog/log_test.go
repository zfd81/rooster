package rlog

import (
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	logger := NewLogger()
	logger.Debug("select * from sys_user where name=?", "bb\n", "select * from sys_user where name=?")
}
