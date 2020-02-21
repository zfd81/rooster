package rsql

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zfd81/rooster/conf"
	"github.com/zfd81/rooster/util"
)

var (
	config = conf.GetGlobalConfig()
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
func (r *Rows) MapScan() (map[string]interface{}, error) {
	r.Next()
	return MapScan(r)
}

func (r *Rows) MapListScan() ([]map[string]interface{}, error) {
	return MapListScan(r)
}

// StructScan a single Row into dest.
func (r *Rows) StructScan(dest interface{}) error {
	if r.Next() {
		return StructScan(r, dest)
	}
	return nil
}

func (r *Rows) StructListScan(list interface{}) error {
	return StructListScan(r, list)
}

type DB struct {
	*sql.DB
	driverName     string
	dataSourceName string
	unsafe         bool
}

// Open is the same as rsql.Open, but returns an *rooster.rsql.DB instead.
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

func (db *DB) Query(query string, arg interface{}) (*Rows, error) {
	sql, params, err := bindParams(query, NewParams(arg))
	if err != nil {
		return nil, err
	}
	r, err := db.DB.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r}, err
}

func (db *DB) QuerySlice(query string, arg interface{}) ([]interface{}, error) {
	rows, err := db.Query(query, arg)
	if err != nil {
		return nil, err
	}
	return rows.SliceScan()
}

func (db *DB) QueryMap(query string, arg interface{}) (map[string]interface{}, error) {
	rows, err := db.Query(query, arg)
	if err != nil {
		return nil, err
	}
	return rows.MapScan()
}

func (db *DB) QueryMapList(query string, arg interface{}, pageNumber int, pageSize int) ([]map[string]interface{}, error) {
	sql, err := pagesql(db.driverName, query, pageNumber, pageSize)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(sql, arg)
	if err != nil {
		return nil, err
	}
	return rows.MapListScan()
}

func (db *DB) QueryStruct(dest interface{}, query string, arg interface{}) error {
	rows, err := db.Query(query, arg)
	if err != nil {
		return err
	}
	return rows.StructScan(dest)
}

func (db *DB) QueryStructList(list interface{}, query string, arg interface{}, pageNumber int, pageSize int) error {
	sql, err := pagesql(db.driverName, query, pageNumber, pageSize)
	if err != nil {
		return err
	}
	rows, err := db.Query(sql, arg)
	if err != nil {
		return err
	}
	return rows.StructListScan(list)
}

func (db *DB) Save(arg interface{}, table ...string) (int64, error) {
	tableName := ""
	if len(table) > 0 {
		tableName = table[0]
	} else {
		model, ok := arg.(Modeler)
		if ok {
			tableName = model.TableName()
		} else {
			return 0, errors.New("must pass a Model to Save destination")
		}
	}
	sql, params, err := insert(tableName, arg)
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

func (db *DB) Exec(query string, arg interface{}) (int64, error) {
	sql, params, err := bindParams(query, NewParams(arg))
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
//	rsql, args, err := bindParams(query, param)
//	if err != nil {
//		return -1, err
//	}
//	stmt, err := db.Prepare(rsql)
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

// func (db *DB) Execute(query string, arg interface{}) (rsql.Result, error) {
// 	var rsql string
// 	var arglist []interface{}
// 	var err error
// 	if maparg, ok := arg.(map[string]interface{}); ok {
// 		rsql, arglist, err = bindMap(query, maparg)
// 	} else {
// 		// rsql, arglist, err = bindStruct(query, maparg)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db.Exec(rsql, arglist...)
// }

func SliceScan(r *Rows) ([]interface{}, error) {
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
		values[i] = value(column.ScanType(), values[i])
	}
	return values, r.Err()
}

func MapScan(r *Rows) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	columns, err := r.ColumnTypes()
	if err != nil {
		return m, err
	}
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}
	err = r.Scan(values...)
	if err != nil {
		return m, err
	}
	for i, column := range columns {
		m[column.Name()] = value(column.ScanType(), values[i])
	}
	return m, r.Err()
}

func MapListScan(r *Rows) ([]map[string]interface{}, error) {
	l := make([]map[string]interface{}, 0, 10)
	columns, err := r.ColumnTypes()
	if err != nil {
		return l, err
	}
	for r.Next() {
		m := make(map[string]interface{})
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}
		err = r.Scan(values...)
		if err != nil {
			return l, err
		}
		for i, column := range columns {
			m[column.Name()] = value(column.ScanType(), values[i])
		}
		l = append(l, m)
	}
	return l, r.Err()
}

func StructScan(r *Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	v = v.Elem()

	nameMapping := map[string]int{}
	t := v.Type()
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		f := NewField(&field)
		if f.NotIgnore() {
			name := f.AttrName()
			if name == "" {
				name = f.Name
			}
			nameMapping[strings.ToLower(name)] = i
		}
	}

	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}
	err = r.Scan(values...)

	for i, column := range columns {
		index := nameMapping[strings.ToLower(column)]
		f := v.Field(index)
		if !reflect.ValueOf(values[i]).Elem().IsZero() {
			switch f.Kind() {
			case reflect.String:
				f.SetString(cast.ToString(values[i]))
			case reflect.Int:
				f.SetInt(cast.ToInt64(values[i]))
			case reflect.Bool:
				f.SetBool(cast.ToBool(values[i]))
			case reflect.Struct:
				f.Set(reflect.ValueOf(cast.ToTime(values[i])))
			}
		}
	}
	return err
}

func StructListScan(r *Rows, list interface{}) error {
	var v, vp reflect.Value

	value := reflect.ValueOf(list)

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
			field := t.Field(i)
			f := NewField(&field)
			if f.NotIgnore() {
				aname := f.AttrName()
				if aname == "" {
					aname = f.Name
				}
				if strings.ToLower(aname) == name {
					valueOfField := v.FieldByName(f.Name)
					values[index] = valueOfField.Addr().Interface()
					flag = false
					break
				}
			}
		}
		if flag {
			values[index] = new(interface{})
		}
	}
	return values
}

func value(t reflect.Type, v interface{}) interface{} {
	switch t.String() {
	case "sql.RawBytes":
		if reflect.ValueOf(v).Elem().IsZero() {
			return ""
		}
		return string((*(v.(*interface{}))).([]uint8))
	case "int64", "sql.NullInt64":
		if reflect.ValueOf(v).Elem().IsZero() {
			return 0
		}
		if "int64" == reflect.TypeOf(*(v.(*interface{}))).String() {
			return *(v.(*interface{}))
		} else {
			return cast.ToInt(string((*(v.(*interface{}))).([]uint8)))
		}
	case "mysql.NullTime":
		if reflect.ValueOf(v).Elem().IsZero() {
			return time.Time{}
		}
		return *(v.(*interface{}))
	default:
		return *(v.(*interface{}))
	}
}

func pagesql(driverName string, sql string, pageNumber int, pageSize int) (string, error) {
	sql = config.PageSql(driverName, sql)
	env := map[string]interface{}{
		"_pageNumber": pageNumber,
		"_pageSize":   pageSize,
	}
	newSql, err := util.ReplaceBetween(sql, "${", "}", func(index int, start int, end int, content string) (string, error) {
		val, err := util.ExprParsing(env, strings.TrimSpace(content))
		if err != nil {
			return content, err
		}
		return cast.ToString(val), nil
	})
	if err != nil {
		return sql, err
	}
	return newSql, nil
}
