package xsql

import (
	"database/sql"
)

type DB struct {
	*sql.DB
	driverName     string
	dataSourceName string
	unsafe         bool
}

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, driverName: driverName, dataSourceName: dataSourceName}, err
}

func (db *DB) Save(table string, arg Params) (int64, error) {
	sql, params, err := insert(table, arg)
	if err != nil {
		return -1, err
	}
	stmt, err := db.DB.Prepare(sql)
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
	stmt, err := db.DB.Prepare(sql)
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
