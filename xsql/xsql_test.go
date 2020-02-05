package xsql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"testing"
	"time"
)

var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", "root", "123456", "localhost", "hdss")

func TestDB_Query(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	p := NewParams()
	p.Add("Name", "tester")
	rows, _ := db.Query("select * from sys_user where name=:Name", p)
	arr, _ := rows.SliceScan()
	for i, v := range arr {
		t.Log(i, v)
	}
	t.Log("-------------------------------------------------------")
	mp := make(map[string]interface{})
	mp["Name"] = "admin"
	rows, _ = db.Query("select * from sys_user where name=:Name", mp)
	m, _ := rows.MapScan()
	for k, v := range m {
		t.Log(k, v, reflect.TypeOf(v).String())
	}
	t.Log("-------------------------------------------------------")
	u := &User{Name: "admin"}
	rows, _ = db.Query("select * from sys_user where name=:Name", u)
	rows.StructScan(u)
	t.Log(u.Name)
	t.Log(u.Id)
	t.Log(u.Password)
	t.Log(u.Department_id)

	t.Log("-------------------------------------------------------")
	rows, _ = db.Query("select * from sys_user", nil)
	l, _ := rows.MapListScan()
	for _, m := range l {
		for k, v := range m {
			t.Log(k, v, reflect.TypeOf(v).String())
		}
	}

	t.Log("-------------------------------------------------------")
	up := &User{Name: "tester"}
	rows, _ = db.Query("select * from sys_user where name=:Name", up)
	users := make([]User, 0)
	rows.StructListScan(&users)
	for i, u := range users {
		t.Log(i, u)
	}
}

func TestMapScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	p := NewParams()
	p.Add("Name", "insUser4")
	rows, _ := db.Query("select * from sys_user where name=:Name", p)
	if rows.Next() {
		m, _ := MapScan(rows)
		for k, v := range m {
			t.Log(k, v, reflect.TypeOf(v).String())
		}
	}
}

func TestSliceScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	p := NewParams()
	p.Add("Name", "insUser4")
	rows, _ := db.Query("select * from sys_user where name=:Name", p)
	if rows.Next() {
		arr, err := SliceScan(rows)
		if err != nil {
			t.Error(err)
		}
		for i, v := range arr {
			t.Log(i, v)
		}
	}
}

func TestStructScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	u := &User{Name: "insUser4"}
	rows, _ := db.Query("select * from sys_user where name=:Name", u)
	if rows.Next() {
		err := StructScan(rows, u)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log("Id: ", u.Id)
	t.Log("Name: ", u.Name)
	t.Log("Number: ", u.Number)
	t.Log("Password: ", u.Password)
	t.Log("Department_id: ", u.Department_id)
	t.Log("Created_date: ", u.Created_date)
	t.Log("Lastmodified_date: ", u.Lastmodified_date)
}

func TestStructListScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	rows, _ := db.Query("select * from sys_user", nil)
	users := make([]User, 0)
	err = StructListScan(rows, &users)
	if err != nil {
		t.Error(err)
	}
	for i, u := range users {
		t.Log(i, "：", u)
	}
}

func TestMapListScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	rows, _ := db.Query("select * from sys_user", nil)
	l, _ := MapListScan(rows)
	for _, m := range l {
		for k, v := range m {
			t.Log(k, ": ", v)
		}
	}
}

func Test_param(t *testing.T) {
	u := User{}
	//users := make([]User, 0)
	//m := make(map[string]interface{})
	param(&u)
}

func TestDB_Save(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	u := &User{55, "用户5", "pwd515", "5115", 51115, time.Now(), time.Now()}
	cnt, err := db.Save("sys_user", u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}
}

func TestDB_Exec(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	u := &User{85, "用户8", "pwd715", "7115", 61115, time.Now(), time.Now()}
	cnt, err := db.Exec("insert into sys_user (id,created_date,lastmodified_date,name,number,password,department_id) values (:Id,:Created_date,:Lastmodified_date,:Name,:Number,:Password,:Department_id)", u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	mp := make(map[string]interface{})
	mp["name"] = "%7"
	cnt, err = db.Exec("delete FROM sys_user where name like :name", mp)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

}

func TestDB_Exec_Ins(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	u := &User{85, "用户8", "pwd715", "7115", 61115, time.Now(), time.Now()}
	sql := "insert into sys_user (id,created_date,lastmodified_date,name,number,password,department_id) values (:Id,:Created_date,:Lastmodified_date,:Name,:Number,:Password,:Department_id)"
	cnt, err := db.Exec(sql, u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}
}

func TestDB_Exec_Del(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	mp := make(map[string]interface{})
	mp["name"] = "%7"
	cnt, err := db.Exec("delete FROM sys_user where name like :name", mp)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	u := &User{Name: "2"}
	cnt, err = db.Exec("delete FROM sys_user where name like CONCAT('%',:Name)", u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	//参数类型为整型
	cnt, err = db.Exec("delete FROM sys_user where id=:val", 55)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	//参数类型为字符串
	cnt, err = db.Exec("delete FROM sys_user where name=:val", "用户8")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

}

func TestDB_Exec_Update(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	p := NewParams()
	p.Add("id", 15)
	p.Add("modifier", 151)
	p.Add("Full_name", "用户1")
	sql := "update sys_user set modifier=:modifier ,full_name=:Full_name where id = :id"
	cnt, err := db.Exec(sql, p)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}
}

func TestDB_QueryForSlice(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	//Map作为参数查询
	mp := make(map[string]interface{})
	mp["Name"] = "admin"
	s, err := db.QueryForSlice("select * from sys_user where name=:Name", mp)
	if err != nil {
		t.Error(err)
	}
	for i, v := range s {
		t.Log(i, v)
	}

	t.Log("---------------------------------------------------------------")

	//类作为参数查询
	u := &User{Name: "tester"}
	s, err = db.QueryForSlice("select * from sys_user where name=:Name", u)
	if err != nil {
		t.Error(err)
	}
	for i, v := range s {
		t.Log(i, v)
	}
}

func TestDB_QueryForMap(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	//Map作为参数查询
	mp := make(map[string]interface{})
	mp["Name"] = "admin"
	m, err := db.QueryForMap("select * from sys_user where name=:Name", mp)
	if err != nil {
		t.Error(err)
	}
	for k, v := range m {
		t.Log(k, v)
	}

	t.Log("---------------------------------------------------------------")

	//单一参数查询
	m, err = db.QueryForMap("select * from sys_user where name=:val", "tester")
	if err != nil {
		t.Error(err)
	}
	for k, v := range m {
		t.Log(k, v)
	}
}

func TestDB_QueryForMapList(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	//空参数查询
	l, err := db.QueryForMapList("select * from sys_user", nil)
	if err != nil {
		t.Error(err)
	}
	for _, m := range l {
		t.Log(m)
	}
}
