package xsql

import (
	"testing"
	"time"
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

func Test_bindParams(t *testing.T) {
	str := "select name from tbale where name=:Name and Password =:PWD and{} {age>1}"
	user := &User{Name: "zfd", Password: "456"}
	param := NewStructParams(user)
	str, params, err := bindParams(str, param)
	t.Log(str)
	t.Log(params)
	t.Log(err)
	t.Log(param.Size())
	t.Log(param.Names())
	t.Log(len(param.Names()))
}

func Test_insert(t *testing.T) {
	user := &User{Name: "zfd", Password: "4568"}
	str, params, err := insert("userInfo", NewStructParams(user))
	t.Log(str)
	t.Log(params)
	t.Log(len(params))
	t.Log(err)

	countrylMap := make(map[string]interface{})
	countrylMap["Name"] = "China"
	countrylMap["Code"] = "86"
	countrylMap["Capitale"] = "BeiJing"
	param := NewMapParams(countrylMap)
	str, params, err = insert("countryInfo", param)
	t.Log(str)
	t.Log(params)
	t.Log(len(params))
	t.Log(err)
}
