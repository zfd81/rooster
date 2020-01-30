package xsql

import (
	"bytes"
	"fmt"
	"github.com/zfd81/rooster/errors"
	"github.com/zfd81/rooster/util"
)

func bindParams(sql string, arg Params) (string, []interface{}, error) {
	newSql, err := util.ReplaceBetween(sql, "{", "}", func(index int, start int, end int, content string) (string, error) {
		ignore := false
		fragment, err := util.ReplaceByKeyword(content, ':', func(i int, s int, e int, c string) (string, error) {
			_, ok := arg.Get(c)
			if !ok {
				ignore = true
				return "", nil
			}
			return fmt.Sprintf(":%s", c), nil
		})
		if ignore {
			return "", err
		}
		return fragment, err
	})
	if err != nil {
		return "", nil, err
	}
	params := make([]interface{}, 0, 20)
	newSql, err = util.ReplaceByKeyword(newSql, ':', func(index int, start int, end int, content string) (string, error) {
		val, ok := arg.Get(content)
		if !ok {
			return "?", fmt.Errorf("could not find name %s in %#v", content, arg)
		}
		params = append(params, val)
		return "?", nil
	})
	return newSql, params, err
}

func insert(table string, arg Params) (string, []interface{}, error) {
	if table == "" || arg == nil {
		return "", nil, errors.ErrParamNotNil
	}
	if arg.Size() < 1 {
		return "", nil, errors.ErrParamEmpty
	}
	var sql bytes.Buffer
	var sql2 bytes.Buffer
	sql.WriteString("insert into ")
	sql.WriteString(table)
	sql.WriteString(" (")
	sql2.WriteString(") values (")
	params := make([]interface{}, 0, 20)
	index := 0
	arg.Iterator(func(key string, value interface{}) {
		if index == 0 {
			index++
		} else {
			sql.WriteString(",")
			sql2.WriteString(",")
		}
		sql.WriteString(string(key))
		sql2.WriteString("?")
		params = append(params, value)
	})
	sql.WriteString(sql2.String())
	sql.WriteString(")")
	return sql.String(), params, nil
}
