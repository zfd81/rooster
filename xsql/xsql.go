package xsql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/zfd81/rooster/util"
	"reflect"
)

var _scannerInterface = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
var _valuerInterface = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

type Rows struct {
	*sql.Rows
}

// SliceScan using this Rows.
func (r *Rows) SliceScan() ([]interface{}, error) {
	return SliceScan(r)
}

// MapScan using this Rows.
func (r *Rows) MapScan(dest map[string]interface{}) error {
	return MapScan(r, dest)
}

// StructScan a single Row into dest.
func (r *Rows) StructScan(dest interface{}) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}

	v = v.Elem()

	columns, err := r.Columns()
	if err != nil {
		return err
	}

	t := v.Type()
	values := make([]interface{}, len(columns))
	fieldNum := v.NumField()
	for i := 0; i < fieldNum; i++ {
		fname := t.Field(i).Name
		for index, name := range columns {
			if fname == name {
				valueOfField := v.FieldByName(name)
				values[index] = valueOfField.Interface()
				break
			}
		}
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
	columns, err := r.Columns()
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

	for i := range columns {
		values[i] = *(values[i].(*interface{}))
	}

	return values, r.Err()
}

func MapScan(r *Rows, dest map[string]interface{}) error {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.Columns()
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
		dest[column] = *(values[i].(*interface{}))
	}

	return r.Err()
}

func structOnlyError(t reflect.Type) error {
	isStruct := t.Kind() == reflect.Struct
	isScanner := reflect.PtrTo(t).Implements(_scannerInterface)
	if !isStruct {
		return fmt.Errorf("expected %s but got %s", reflect.Struct, t.Kind())
	}
	if isScanner {
		return fmt.Errorf("structscan expects a struct dest but the provided struct type %s implements scanner", t.Name())
	}
	return fmt.Errorf("expected a struct, but struct %s has no exported fields", t.Name())
}

func isScannable(t reflect.Type) bool {
	if reflect.PtrTo(t).Implements(_scannerInterface) {
		return true
	}
	if t.Kind() != reflect.Struct {
		return true
	}

	return false
}

func fieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(interface{})
			continue
		}
		f := util.FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}
	return nil
}
