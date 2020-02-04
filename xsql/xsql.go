package xsql

import (
	"database/sql"
	"errors"
	"github.com/spf13/cast"
	"github.com/zfd81/rooster/util"
	"reflect"
	"strings"
)

type Rows struct {
	*sql.Rows
}

// SliceScan using this Rows.
func (r *Rows) SliceScan() ([]interface{}, error) {
	if r.Next() {
		return SliceScan(r)
	}
	return nil, nil
}

// MapScan using this Rows.
func (r *Rows) MapScan(dest map[string]interface{}) error {
	if r.Next() {
		return MapScan(r, dest)
	}
	return nil
}

// StructScan a single Row into dest.
func (r *Rows) StructScan(dest interface{}) error {
	if r.Next() {
		return StructScan(r, dest)
	}
	return nil
}

type DB struct {
	*sql.DB
	driverName     string
	dataSourceName string
	unsafe         bool
}

// Open is the same as sql.Open, but returns an *rooster.xsql.DB instead.
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, driverName: driverName, dataSourceName: dataSourceName}, err
}

// Connect to a database and verify with a ping.
func Connect(driverName, dataSourceName string) (*DB, error) {
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (db *DB) Query(query string, arg Params) (*Rows, error) {
	sql, params, err := bindParams(query, arg)
	if err != nil {
		return nil, err
	}
	r, err := db.DB.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r}, err
}

func (db *DB) Save(table string, arg Params) (int64, error) {
	sql, params, err := insert(table, arg)
	if err != nil {
		return -1, err
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(params...)
	if err != nil {
		return -1, err
	}
	num, err := res.RowsAffected() //影响行数
	return num, err
}

func (db *DB) Exec(query string, arg Params) (int64, error) {
	sql, params, err := bindParams(query, arg)
	if err != nil {
		return -1, err
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(params...)
	if err != nil {
		return -1, err
	}
	num, err := res.RowsAffected() //影响行数
	return num, err
}

//func (db *DB) XExec(query string, param Paramer) (int64, error) {
//	sql, args, err := bindParams(query, param)
//	if err != nil {
//		return -1, err
//	}
//	stmt, err := db.Prepare(sql)
//	defer stmt.Close()
//	if err != nil {
//		return -1, err
//	}
//	res, err := stmt.Exec(args...)
//	if err != nil {
//		return -1, err
//	}
//	num, err := res.RowsAffected()
//	return num, err
//}

// func (db *DB) Execute(query string, arg interface{}) (sql.Result, error) {
// 	var sql string
// 	var arglist []interface{}
// 	var err error
// 	if maparg, ok := arg.(map[string]interface{}); ok {
// 		sql, arglist, err = bindMap(query, maparg)
// 	} else {
// 		// sql, arglist, err = bindStruct(query, maparg)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db.Exec(sql, arglist...)
// }

func SliceScan(r *Rows) ([]interface{}, error) {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.ColumnTypes()
	if err != nil {
		return []interface{}{}, err
	}
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}
	err = r.Scan(values...)
	if err != nil {
		return values, err
	}
	for i, column := range columns {
		switch column.ScanType().String() {
		case "sql.RawBytes", "mysql.NullTime":
			values[i] = string((*(values[i].(*interface{}))).([]uint8))
		case "int64", "sql.NullInt64":
			if "int64" == reflect.TypeOf(*(values[i].(*interface{}))).String() {
				values[i] = *(values[i].(*interface{}))
			} else {
				values[i] = cast.ToInt(string((*(values[i].(*interface{}))).([]uint8)))
			}
		default:
			values[i] = *(values[i].(*interface{}))
		}
	}
	return values, r.Err()
}

func MapScan(r *Rows, dest map[string]interface{}) error {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.ColumnTypes()
	if err != nil {
		return err
	}
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}
	err = r.Scan(values...)
	if err != nil {
		return err
	}
	for i, column := range columns {
		switch column.ScanType().String() {
		case "sql.RawBytes", "mysql.NullTime":
			dest[column.Name()] = string((*(values[i].(*interface{}))).([]uint8))
		case "int64", "sql.NullInt64":
			if "int64" == reflect.TypeOf(*(values[i].(*interface{}))).String() {
				dest[column.Name()] = *(values[i].(*interface{}))
			} else {
				dest[column.Name()] = cast.ToInt(string((*(values[i].(*interface{}))).([]uint8)))
			}
		default:
			dest[column.Name()] = *(values[i].(*interface{}))
		}
	}
	return r.Err()
}

func StructScan(r *Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	v = v.Elem()
	columns, err := r.Columns()
	if err != nil {
		return err
	}
	values := wrapFields(v, columns)
	err = r.Scan(values...)
	return err
}

func StructListScan(r *Rows, dest interface{}) error {
	var v, vp reflect.Value

	value := reflect.ValueOf(dest)

	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructListScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructListScan destination")
	}
	direct := reflect.Indirect(value)

	slice, err := util.BaseType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}
	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := util.Deref(slice.Elem())

	columns, err := r.Columns()
	if err != nil {
		return err
	}

	for r.Next() {
		vp = reflect.New(base)
		v = reflect.Indirect(vp)
		values := wrapFields(v, columns)
		err = r.Scan(values...)
		if err != nil {
			return err
		}
		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, v))
		}
	}
	return nil
}

func wrapFields(v reflect.Value, names []string) []interface{} {
	v = reflect.Indirect(v)
	values := make([]interface{}, len(names))
	t := v.Type()
	fieldNum := v.NumField()
	for index, name := range names {
		flag := true
		name = strings.ToLower(name)
		for i := 0; i < fieldNum; i++ {
			fname := t.Field(i).Name
			if strings.ToLower(fname) == name {
				valueOfField := v.FieldByName(fname)
				values[index] = valueOfField.Addr().Interface()
				flag = false
				break
			}
		}
		if flag {
			values[index] = new(interface{})
		}
	}
	return values
}
