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
	u := &User{55, "insUser5", "pwd515", "5115", 51115, time.Now(), time.Now()}
	cnt, err := db.Save("sys_user", u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}
	rows, _ := db.Query("select * from sys_user", nil)
	users := make([]User, 0)
	err = StructListScan(rows, &users)
	if err != nil {
		t.Error(err)
	}
	for i, u := range users {
		t.Log(i, u)
	}
}
