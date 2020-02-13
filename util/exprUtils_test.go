package util

import (
	"github.com/spf13/cast"
	"testing"
)

func TestExprParsing(t *testing.T) {
	m := map[string]interface{}{
		"a": 12,
		"b": 10,
	}
	v, err := ExprParsing(m, "a+b")
	if err != nil {
		t.Error(err)
	}
	s := "aa" + cast.ToString(v)
	t.Log(s)
}
