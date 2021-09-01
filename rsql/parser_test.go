package rsql

import (
	"testing"
	"time"

	"github.com/zfd81/rooster/types/container"
)

func Test_bindParams(t *testing.T) {
	str := "select name from tbale where name=:Name and Password =\\:PWD and{} {age>1}=:aa"
	user := &User{Name: "zfd", Password: "456"}
	param := NewStructParams(user)
	str, params, err := bindParams(str, param)
	if err != nil {
		t.Error(err)
	}
	t.Log(str)
	t.Log(params)
	str1 := "$index"
	if str1 == "$index" {
		t.Log("==========")
	}
	ms1 := []map[string]interface{}{{"aa": 111, "bb": 222, "cc": 333}, {"aa": 444, "bb": 555, "cc": 666}}
	param.Add("msa", ms1)
	param.Add("msa1", 12)
	str = "insert into tbale (name,pwd,age,seq ) values {@msa[,] (:Name,:this.bb,:this.cc,:this.$index)}"
	sql, params, err := bindParams(str, param)
	if err != nil {
		t.Error(err)
	}
	t.Log(sql)
	t.Log(params)
	//t.Log(params)
	//
	//ms2 := []string{"11", "22", "33"}
	//param.Add("msb", ms2)
	//str = "select * from tbale where name in ({@msb[,] :this.val})"
	//sql, params, err = bindParams(str, param)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(sql)
	//t.Log(params)
	//
	//p := NewParams(ms2)
	//str = "select * from tbale where name in ({@vals[,] :this.val})"
	//sql, params, err = bindParams(str, p)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(sql)
	//t.Log(params)
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

func Test_foreach(t *testing.T) {
	//p := Params{}
	//ms1 := []map[string]interface{}{{"aa": 111, "bb": 222, "cc": 333}, {"aa": 111, "bb": 222, "cc": 333}}
	//p.Add("msa", ms1)
	//p.Add("msa1", 12)
	//str := "@msa[,] (:this.aa,:this.bb,:this.cc1)"
	//sql, err := foreach(str, &p)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(sql)
	//t.Log(p)
	//
	//ms2 := []string{"11", "22", "33"}
	//p.Add("msb", ms2)
	//str = "@msb[,] :this.val"
	//sql, err = foreach(str, &p)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(sql)
	//t.Log(p)
	//
	ms3 := []container.JsonMap{}
	ms3 = append(ms3, container.JsonMap{
		"id":   123456,
		"code": "bbb",
		"uid":  "ccc",
		"t":    time.Now(),
	})
	sql := `
@vals [,] (
				:this.id,
				:this.code ,
				:this.uid ,
				:this.t
			)

`
	pp := NewParams(ms3)
	sql, err := foreach(sql, &pp)
	if err != nil {
		t.Error("err:", err)
	} else {
		t.Log(sql)
	}

}

func Test_validCharacter(t *testing.T) {
	var b byte = ' '
	t.Log(b)
	str := "he llo[]"
	for _, v := range str {
		t.Log(v)
	}
}
