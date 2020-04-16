package rsql

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/zfd81/rooster/types/container"
)

type MM struct {
	TesName string `rsql:"-"` //忽略这个字段
}

type Model struct {
	Creator          int
	CreatedDate      time.Time `rsql:"name:created_date"` //标注数据表中对应的列名称
	Modifier         int
	LastmodifiedDate time.Time `rsql:"name:lastmodified_date"` //标注数据表中对应的列名称
	MM
}

type User struct {
	Id           int //没有注解，表示名称和数据库中的列名相同
	Name         string
	Password     string
	FullName     string `rsql:"name:full_name"` //标注数据表中对应的列名称
	Number       string
	DepartmentId int       `rsql:"name:department_id"` //标注数据表中对应的列名称
	Field1       string    `rsql:"-"`                  //忽略这个字段
	Field2       int       `rsql:"-"`                  //忽略这个字段
	Field3       time.Time `rsql:"-"`                  //忽略这个字段
	Model
}

type User1 struct {
	Id           int //没有注解，表示名称和数据库中的列名相同
	Name         string
	Password     string
	FullName     string `rsql:"name:full_name"` //标注数据表中对应的列名称
	Number       string
	DepartmentId int `rsql:"name:department_id"` //标注数据表中对应的列名称
	Model
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
	u := &User1{Name: "user23"}
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
		Id:       22,
		Name:     "user22",
		Password: "pwd22",
		FullName: "用户22",
		Number:   "num22",
		//DepartmentId: 1022,
		Model: Model{
			Creator:          11,
			CreatedDate:      time.Now(),
			Modifier:         12,
			LastmodifiedDate: time.Now(),
		},
		Field1: "test",
		Field2: 999,
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

	t.Log("参数为JsonMap类型>>>>>>>>>>>>>>>>>>>>>>>")
	jm := container.JsonMap{}
	jm.Put("Id", 26)
	jm.Put("Name", "user26")
	jm.Put("Password", "pwd26")

	jm.Put("full_name", "用户26")
	jm.Put("Number", "num26")
	jm.Put("department_id", 1026)
	jm.Put("Creator", "1")
	jm.Put("created_date", time.Now())
	jm.Put("Modifier", "1")
	jm.Put("lastmodified_date", time.Now())

	cnt, err = db.Save(jm, "sys_user")
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
		Id:           26,
		Name:         "user26",
		Password:     "pwd26",
		FullName:     "用户26",
		Number:       "num26",
		DepartmentId: 1026,
		Model: Model{
			Creator:          1,
			CreatedDate:      time.Now(),
			Modifier:         1,
			LastmodifiedDate: time.Now(),
		},
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
		Id:           23,
		Name:         "user23",
		Password:     "pwd23",
		Number:       "num23",
		DepartmentId: 1023,
		Model: Model{
			Creator:          1,
			CreatedDate:      time.Now(),
			Modifier:         1,
			LastmodifiedDate: time.Now(),
		},
	}
	sql := "insert into sys_user (id,created_date,lastmodified_date,name,number,password,department_id) values (:Id,:CreatedDate,:LastmodifiedDate,:Name,:Number,:Password,:DepartmentId)"
	cnt, err := db.Exec(sql, u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	mps := []map[string]interface{}{
		{
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
		},
		{
			"Id":                26,
			"Name":              "user26",
			"Password":          "pwd26",
			"full_name":         "用户26",
			"Number":            "num26",
			"department_id":     1026,
			"Creator":           1,
			"created_date":      time.Now(),
			"Modifier":          1,
			"lastmodified_date": time.Now(),
		},
		{
			"Id":                27,
			"Name":              "user27",
			"Password":          "pwd27",
			"full_name":         "用户27",
			"Number":            "num27",
			"department_id":     1027,
			"Creator":           1,
			"created_date":      time.Now(),
			"Modifier":          1,
			"lastmodified_date": time.Now(),
		}}
	sql = "insert into sys_user (id,created_date,lastmodified_date,name,number,password,department_id) values {@vals[,] (:this.Id,:this.created_date,:this.lastmodified_date,:this.Name,:this.Number,:this.Password,:this.department_id)}"
	cnt, err = db.Exec(sql, mps)
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

	//SQL中使用遍历实现in的用法
	ps := []int{25, 26, 27}
	sql := "select * from sys_user where id in({@vals[,] :this.val})"
	l, err = db.QueryMapList(sql, ps, 1, 10)
	if err != nil {
		t.Error(err)
	}
	for _, m := range l {
		t.Log(m)
	}

	//SQL中使用遍历实现in的用法
	mp := map[string]interface{}{
		"ids": []int{25, 26, 27},
	}
	sql = "select * from sys_user where {@ids[OR] id=:this.val}"
	l, err = db.QueryMapList(sql, mp, 1, 10)
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

func TestDB_BatchSave(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	t.Log("参数为实现Modeler接口的struct类型>>>>>>>>>>>>>>>>>>>>>>>")
	u1 := &User{
		Id:           222,
		Name:         "user222",
		Password:     "pwd222",
		FullName:     "用户222",
		Number:       "num222",
		DepartmentId: 10222,
		Model: Model{
			Creator:          1,
			CreatedDate:      time.Now(),
			Modifier:         1,
			LastmodifiedDate: time.Now(),
		},
		Field1: "test",
		Field2: 9999,
	}
	u2 := &User{
		Id:           333,
		Name:         "user333",
		Password:     "pwd333",
		FullName:     "用户333",
		Number:       "num333",
		DepartmentId: 10333,
		Model: Model{
			Creator:          1,
			CreatedDate:      time.Now(),
			Modifier:         1,
			LastmodifiedDate: time.Now(),
		},
		Field1: "test",
		Field2: 0000,
	}
	cnt, err := db.BatchSave([](interface{}){u1, u2})
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	t.Log("参数为struct类型>>>>>>>>>>>>>>>>>>>>>>>")
	ui1 := &UserInfo{
		Id:               233,
		Name:             "user233",
		Password:         "pwd233",
		FullName:         "用户233",
		Number:           "num233",
		DepartmentId:     10233,
		Creator:          1,
		CreatedDate:      time.Now(),
		Modifier:         1,
		LastmodifiedDate: time.Now(),
		Field1:           "test",
		Field2:           9993,
	}
	ui2 := &UserInfo{
		Id:               234,
		Name:             "user234",
		Password:         "pwd234",
		FullName:         "用户234",
		Number:           "num234",
		DepartmentId:     10234,
		Creator:          1,
		CreatedDate:      time.Now(),
		Modifier:         1,
		LastmodifiedDate: time.Now(),
		Field1:           "test",
		Field2:           9994,
	}
	cnt, err = db.BatchSave([]interface{}{ui1, ui2}, "sys_user")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}

	t.Log("参数为map类型>>>>>>>>>>>>>>>>>>>>>>>")
	mp1 := map[string]interface{}{
		"Id":                255,
		"Name":              "user255",
		"Password":          "pwd255",
		"full_name":         "用户255",
		"Number":            "num255",
		"department_id":     10255,
		"Creator":           1,
		"created_date":      time.Now(),
		"Modifier":          1,
		"lastmodified_date": time.Now(),
	}
	mp2 := map[string]interface{}{
		"Id":                256,
		"Name":              "user256",
		"Password":          "pwd256",
		"full_name":         "用户256",
		"Number":            "num256",
		"department_id":     10256,
		"Creator":           1,
		"created_date":      time.Now(),
		"Modifier":          1,
		"lastmodified_date": time.Now(),
	}
	cnt, err = db.BatchSave([]interface{}{mp1, mp2}, "sys_user")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cnt)
	}
}

func TestDB_QueryCount(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}

	cnt, err := db.QueryCount("select * from sys_user where name like :val", "%2")
	if err != nil {
		t.Error(err)
	}
	t.Log(cnt)

	cnt, err = db.QueryCount("select * from sys_user", nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(cnt)
}

func TestDB_BeginTx(t *testing.T) {
	db, err := Open("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	u := &User{
		Id:           29,
		Name:         "user23",
		Password:     "pwd23",
		Number:       "num23",
		DepartmentId: 1023,
		Model: Model{
			Creator:          1,
			CreatedDate:      time.Now(),
			Modifier:         1,
			LastmodifiedDate: time.Now(),
		},
	}
	u1 := &User{
		Id:           239,
		Name:         "user23",
		Password:     "pwd23",
		Number:       "num23",
		DepartmentId: 1023,
		Model: Model{
			Creator:          1,
			CreatedDate:      time.Now(),
			Modifier:         1,
			LastmodifiedDate: time.Now(),
		},
	}
	sql := "insert into sys_user (id,created_date,lastmodified_date,name,number,password,department_id) values (:Id,:CreatedDate,:LastmodifiedDate,:Name,:Number,:Password,:DepartmentId)"
	tx, err := db.BeginTx()

	cnt, err := tx.ExecTx(sql, u)
	if err != nil {
		t.Log(err)
		tx.Rollback()
	} else {
		t.Log(cnt)
	}
	cnt, err = tx.ExecTx(sql, u1)
	if err != nil {
		t.Log(err)
		tx.Rollback()
	} else {
		t.Log(cnt)

	}
	tx.Commit()
	//tx.Rollback()
}
