package rsql

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id               int //没有注解，表示名称和数据库中的列名相同
	Name             string
	Password         string
	FullName         string `rsql:"name:full_name"` //标注数据表中对应的列名称
	Number           string
	DepartmentId     int `rsql:"name:department_id"` //标注数据表中对应的列名称
	Creator          int
	CreatedDate      time.Time `rsql:"name:created_date"` //标注数据表中对应的列名称
	Modifier         int
	LastmodifiedDate time.Time `rsql:"name:lastmodified_date"` //标注数据表中对应的列名称
	Field1           string    `rsql:"-"`                      //忽略这个字段
	Field2           int       `rsql:"-"`                      //忽略这个字段
	Field3           time.Time `rsql:"-"`                      //忽略这个字段
}

func (u User) TableName() string {
	return "sys_user"
}

type UserInfo struct {
	Id               int //没有注解，表示名称和数据库中的列名相同
	Name             string
	Password         string
	FullName         string `rsql:"name:full_name"` //标注数据表中对应的列名称
	Number           string
	DepartmentId     int `rsql:"name:department_id"` //标注数据表中对应的列名称
	Creator          int
	CreatedDate      time.Time `rsql:"name:created_date"` //标注数据表中对应的列名称
	Modifier         int
	LastmodifiedDate time.Time `rsql:"name:lastmodified_date"` //标注数据表中对应的列名称
	Field1           string    `rsql:"-"`                      //忽略这个字段
	Field2           int       `rsql:"-"`                      //忽略这个字段
	Field3           time.Time `rsql:"-"`                      //忽略这个字段
}

var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", "root", "123456", "localhost", "hdss")

func TestDB_Query(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	t.Log("查询参数为Params类型>>>>>>>>>>>>>>>>>>>>>>>")
	p := make(Params)
	p.Add("Name", "tester")
	rows, _ := db.Query("select * from sys_user where name=:Name", p)
	arr, err := rows.SliceScan()
	if err != nil {
		t.Error(err)
	}
	for i, v := range arr {
		t.Log(i, v)
	}

	t.Log("查询参数为map类型>>>>>>>>>>>>>>>>>>>>>>>")
	mp := make(map[string]interface{})
	mp["Name"] = "admin"
	rows, _ = db.Query("select * from sys_user where name=:Name", mp)
	m, err := rows.MapScan()
	if err != nil {
		t.Error(err)
	}
	for i, k := range m.Keys() {
		v, _ := m.Get(k)
		t.Log(i, k, v)
	}

	t.Log("查询参数为struct类型>>>>>>>>>>>>>>>>>>>>>>>")
	u := &User{Name: "admin"}
	rows, _ = db.Query("select * from sys_user where name=:Name", u)
	err = rows.StructScan(u)
	if err != nil {
		t.Error(err)
	}
	t.Log(u)

	t.Log("查询参数只有一个，而且为string或int类型，变量名为val>>>>>>>>>>>>>>>>>>>>>>>")
	rows, _ = db.Query("select * from sys_user where name=:val", "admin")
	err = rows.StructScan(u)
	if err != nil {
		t.Error(err)
	}
	t.Log(u)

	t.Log("查询参数为空值nil>>>>>>>>>>>>>>>>>>>>>>>")
	rows, _ = db.Query("select * from sys_user", nil)
	l, err := rows.MapListScan()
	if err != nil {
		t.Error(err)
	}
	for i, m := range l {
		t.Log(i, m)
	}

	t.Log("查询参数为struct类型>>>>>>>>>>>>>>>>>>>>>>>")
	up := &User{Name: "admin"}
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
	p := make(Params)
	p.Add("Name", "admin")
	rows, _ := db.Query("select * from sys_user where name=:Name", p)
	if rows.Next() {
		m, _ := MapScan(rows)
		for i, k := range m.Keys() {
			v, _ := m.Get(k)
			t.Log(i, k, v)
		}
	}
}

func TestSliceScan(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	p := NewParams(nil)
	p.Add("Name", "user25")
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
	u := &User{Name: "user23"}
	rows, _ := db.Query("select * from sys_user where name=:Name", u)
	if rows.Next() {
		err := StructScan(rows, u)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(u)
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
		for i, k := range m.Keys() {
			v, _ := m.Get(k)
			t.Log(i, k, v)
		}
	}
}

func TestDB_Save(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	t.Log("参数为实现Modeler接口的struct类型>>>>>>>>>>>>>>>>>>>>>>>")
	u := &User{
		Id:               22,
		Name:             "user22",
		Password:         "pwd22",
		FullName:         "用户22",
		Number:           "num22",
		DepartmentId:     1022,
		Creator:          1,
		CreatedDate:      time.Now(),
		Modifier:         1,
		LastmodifiedDate: time.Now(),
		Field1:           "test",
		Field2:           999,
	}
	cnt, err := db.Save(u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	t.Log("参数为struct类型>>>>>>>>>>>>>>>>>>>>>>>")
	ui := &UserInfo{
		Id:               23,
		Name:             "user23",
		Password:         "pwd23",
		FullName:         "用户23",
		Number:           "num23",
		DepartmentId:     1023,
		Creator:          1,
		CreatedDate:      time.Now(),
		Modifier:         1,
		LastmodifiedDate: time.Now(),
		Field1:           "test",
		Field2:           999,
	}
	cnt, err = db.Save(ui, "sys_user")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	t.Log("参数为map类型>>>>>>>>>>>>>>>>>>>>>>>")
	mp := map[string]interface{}{
		"Id":                25,
		"Name":              "user25",
		"Password":          "pwd25",
		"full_name":         "用户25",
		"Number":            "num25",
		"department_id":     1025,
		"Creator":           1,
		"created_date":      time.Now(),
		"Modifier":          1,
		"lastmodified_date": time.Now(),
	}
	cnt, err = db.Save(mp, "sys_user")
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
	u := &User{
		Id:               26,
		Name:             "user26",
		Password:         "pwd26",
		FullName:         "用户26",
		Number:           "num26",
		DepartmentId:     1026,
		Creator:          1,
		CreatedDate:      time.Now(),
		Modifier:         1,
		LastmodifiedDate: time.Now(),
	}
	cnt, err := db.Exec("insert into sys_user (id,created_date,lastmodified_date,name,number,password,department_id) values (:Id,:CreatedDate,:LastmodifiedDate,:Name,:Number,:Password,:DepartmentId)", u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	mp := make(map[string]interface{})
	mp["name"] = "%26"
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
	u := &User{
		Id:               23,
		Name:             "user23",
		Password:         "pwd23",
		Number:           "num23",
		DepartmentId:     1023,
		CreatedDate:      time.Now(),
		LastmodifiedDate: time.Now(),
	}
	sql := "insert into sys_user (id,created_date,lastmodified_date,name,number,password,department_id) values (:Id,:CreatedDate,:LastmodifiedDate,:Name,:Number,:Password,:DepartmentId)"
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

	//模糊查询用法一：将百分号%写在变量值中
	mp := make(map[string]interface{})
	mp["name"] = "%7"
	cnt, err := db.Exec("delete FROM sys_user where name like :name", mp)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	//模糊查询用法二：用字符串连接函数将百分号%和变量进行连接
	u := &User{Name: "8"}
	cnt, err = db.Exec("delete FROM sys_user where name like CONCAT('%',:Name)", u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	//参数类型为整型
	cnt, err = db.Exec("delete FROM sys_user where id=:val", 23)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	//参数类型为字符串
	cnt, err = db.Exec("delete FROM sys_user where full_name=:val", "用户25")
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
	p := make(Params)
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
	s, err := db.QuerySlice("select * from sys_user where name=:Name", mp)
	if err != nil {
		t.Error(err)
	}
	for i, v := range s {
		t.Log(i, v)
	}

	t.Log("---------------------------------------------------------------")

	//类作为参数查询
	u := &User{Name: "tester"}
	s, err = db.QuerySlice("select * from sys_user where name=:Name", u)
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
	m, err := db.QueryMap("select * from sys_user where name=:Name", mp)
	if err != nil {
		t.Error(err)
	}
	for i, k := range m.Keys() {
		v, _ := m.Get(k)
		t.Log(i, k, v)
	}

	t.Log("---------------------------------------------------------------")

	//单一参数查询
	m, err = db.QueryMap("select * from sys_user where name=:val", "tester")
	if err != nil {
		t.Error(err)
	}
	for i, k := range m.Keys() {
		t.Log(i, k)
	}
}

func TestDB_QueryForMapList(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	//空参数查询
	l, err := db.QueryMapList("select * from sys_user", nil, 1, 4)
	if err != nil {
		t.Error(err)
	}
	for _, m := range l {
		t.Log(m)
	}
}

func TestDB_QueryForStruct(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	//Map作为参数查询
	mp := make(map[string]interface{})
	mp["Name"] = "admin"
	u := &User{}
	err = db.QueryStruct(u, "select * from sys_user where name=:Name", mp)
	if err != nil {
		t.Error(err)
	}
	t.Log(u)

	t.Log("---------------------------------------------------------------")

	//对象作为参数查询
	u1 := &User{Name: "tester"}
	err = db.QueryStruct(u1, "select * from sys_user where name=:Name", u1)
	if err != nil {
		t.Error(err)
	}
	t.Log(u1)

}

func TestDB_QueryForStructList(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	var users []User
	err = db.QueryStructList(&users, "select * from sys_user", nil, 1, 4)
	if err != nil {
		t.Error(err)
	}
	for i, u := range users {
		t.Log(i, u)
	}
}
