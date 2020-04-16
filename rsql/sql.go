package rsql

import (
	"bytes"
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/zfd81/rooster/rlog"

	"github.com/zfd81/rooster/types/container"

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
func (r *Rows) MapScan() (container.Map, error) {
	r.Next()
	return MapScan(r)
}

func (r *Rows) MapListScan() ([]container.Map, error) {
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

type Tx struct {
	*sql.Tx
	driverName string
	logger     *rlog.Logger
}

func (t *Tx) DriverName() string {
	return t.driverName
}

func (t *Tx) ExecTx(query string, arg interface{}) (int64, error) {
	sql, params, err := bindParams(query, NewParams(arg))
	if err != nil {
		return -1, err
	}
	t.logger.Debug(log(sql, params)...)
	res, err := t.Exec(sql, params...)
	if err != nil {
		return -1, err
	}
	num, err := res.RowsAffected() //影响行数
	return num, err
}

type DB struct {
	*sql.DB
	driverName     string
	dataSourceName string
	unsafe         bool
	logger         *rlog.Logger
}

// Open is the same as rsql.Open, but returns an *rooster.rsql.DB instead.
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, driverName: driverName, dataSourceName: dataSourceName, logger: rlog.NewLogger()}, err
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

func (db *DB) SetLogLevel(level rlog.LogLevel) {
	db.logger.SetLevel(level)
}

func (db *DB) Query(query string, arg interface{}) (*Rows, error) {
	sql, params, err := bindParams(query, NewParams(arg))
	if err != nil {
		return nil, err
	}
	db.logger.Debug(log(sql, params)...)
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

func (db *DB) QueryMap(query string, arg interface{}) (container.Map, error) {
	rows, err := db.Query(query, arg)
	if err != nil {
		return nil, err
	}
	return rows.MapScan()
}

func (db *DB) QueryMapList(query string, arg interface{}, pageNumber int, pageSize int) ([]container.Map, error) {
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

func (db *DB) QueryCount(query string, arg interface{}) (int, error) {
	var sql bytes.Buffer
	sql.WriteString("select count(1) recordCount from (")
	sql.WriteString(query)
	sql.WriteString(") roosterCountTable")
	rows, err := db.Query(sql.String(), arg)
	if err != nil {
		return 0, err
	}
	var cnt int
	if rows.Next() {
		err = rows.Scan(&cnt)
	}
	return cnt, err
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
			return 0, errors.New("Please enter the table name")
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
	db.logger.Debug(log(sql, params)...)
	res, err := stmt.Exec(params...)
	if err != nil {
		return -1, err
	}
	num, err := res.RowsAffected() //影响行数
	return num, err
}

func (db *DB) BatchSave(arg []interface{}, table ...string) (int64, error) {
	tableName := ""
	if len(table) > 0 {
		tableName = table[0]
	} else {
		model, ok := arg[0].(Modeler)
		if ok {
			tableName = model.TableName()
		} else {
			return 0, errors.New("Please enter the table name")
		}
	}
	sql, params, err := batchInsert(tableName, arg...)
	if err != nil {
		return -1, err
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return -1, err
	}
	db.logger.Debug(log(sql, params)...)
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
	db.logger.Debug(log(sql, params)...)
	res, err := stmt.Exec(params...)
	if err != nil {
		return -1, err
	}
	num, err := res.RowsAffected() //影响行数
	return num, err
}

func (db *DB) BeginTx() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx, driverName: db.driverName, logger: db.logger}, err
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

func MapScan(r *Rows) (container.Map, error) {
	m := container.JsonMap{}
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
		m.Put(column.Name(), value(column.ScanType(), values[i]))
	}
	return m, r.Err()
}

func MapListScan(r *Rows) ([]container.Map, error) {
	l := make([]container.Map, 0, 10)
	columns, err := r.ColumnTypes()
	if err != nil {
		return l, err
	}
	for r.Next() {
		m := container.JsonMap{}
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}
		err = r.Scan(values...)
		if err != nil {
			return l, err
		}
		for i, column := range columns {
			m.Put(column.Name(), value(column.ScanType(), values[i]))
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

	nameMapping := GetNameMapping(v.Type())

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
		indexes := nameMapping[strings.ToLower(column)]
		f := FieldByIndexes(v, indexes)
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
	mapping := GetNameMapping(v.Type())
	for index, name := range names {
		indexes, ok := mapping[strings.ToLower(name)]
		if ok {
			f := FieldByIndexes(v, indexes)
			values[index] = f.Addr().Interface()
		} else {
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

func log(sql string, params []interface{}) (messages []interface{}) {
	messages = append([]interface{}{"\r\n", sql, "\r\n", "\tparams:"}, params)
	return
}
