package xsql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"testing"
)

var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", "root", "123456", "localhost", "hdss")

func TestDB_Query(t *testing.T) {
	u := User{Id: 12}
	v := reflect.ValueOf(u)
	t.Log(v.Field(0))

}

func TestMapScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	p := NewParams()
	p.Add("Name", "tester")
	rows, _ := db.Query("select * from sys_user where name=:Name", p)
	m := make(map[string]interface{})
	rows.MapScan(m)
	for k, v := range m {
		t.Log(k, v, reflect.TypeOf(v).String())
	}
}

func TestSliceScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	p := NewParams()
	p.Add("Name", "admin")
	rows, _ := db.Query("select * from sys_user where name=:Name", p)
	arr, _ := rows.SliceScan()
	for i, v := range arr {
		t.Log(i, v)
	}
}

func TestStructScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	u := &User{}
	p := NewParams()
	p.Add("Name", "admin")
	rows, _ := db.Query("select * from sys_user where name=:Name", p)
	rows.StructScan(u)
	t.Log(u.Name)
	t.Log(u.Id)
	t.Log(u.Password)
	t.Log(u.Department_id)
}

func TestStructListScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	p := NewParams()
	p.Add("Name", "admin")
	rows, _ := db.Query("select * from sys_user", p)
	users := make([]User, 0)
	StructListScan(rows, &users)
	for i, u := range users {
		t.Log(i, u)
	}
}
