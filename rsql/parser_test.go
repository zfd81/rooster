package rsql

import (
	"testing"
)

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

func Test_batchInsert(t *testing.T) {
	user1 := &User{Name: "aaa", Password: "111", FullName: "aaaa", Number: "1111"}
	user2 := &User{Name: "bbb", Password: "222", FullName: "bbbb", Number: "2222"}
	str, params, err := batchInsert("userInfo", user1, user2)
	t.Log(str)
	t.Log(params)
	t.Log(len(params))
	t.Log(err)

	m1 := make(map[string]interface{})
	m1["Name"] = "aa"
	m1["Code"] = "aa11"
	m1["Capitale"] = "aaaa"
	//param := NewMapParams(m1)
	m2 := make(map[string]interface{})
	m2["Name"] = "bb"
	m2["Code"] = "bb11"
	m2["Capitale"] = "bbbb"
	//param := NewMapParams(m1)
	str, params, err = batchInsert("countryInfo", m1, m2)
	t.Log(str)
	t.Log(params)
	t.Log(len(params))
	t.Log(err)
}
