package xsql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", "root", "123456", "localhost", "tcm")

func TestDB_Query(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	rows, _ := db.Query("select * from userInfo", nil)
	m := make(map[string]interface{})
	rows.MapScan(m)
	t.Log(m)
}
