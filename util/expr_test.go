package util

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cast"
)

type User struct {
	Id                int
	Name              string
	Password          string
	Number            string
	Department_id     int
	Created_date      time.Time
	Lastmodified_date time.Time
}

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

	//mm := make(map[string]string)
	//var mm map[string]string
	var us User = User{}
	vv := reflect.ValueOf(us)
	if vv.Kind() != reflect.Ptr {
		t.Error("err")
	}

	if vv.IsNil() {
		t.Error("==========")
	}
	//t.Log(len(mm))
	//aa(&mm)
	//t.Log(len(mm))
}

func aa(p *map[string]string) {
	(*p)["aa"] = "aa"
	(*p)["bb"] = "bb"
}
