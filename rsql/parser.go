package rsql

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/zfd81/rooster/types/container"

	"github.com/zfd81/rooster/errors"
	"github.com/zfd81/rooster/util"
)

func validCharacter(char byte) bool {
	if (char >= 48 && char <= 57) || (char >= 65 && char <= 90) || (char >= 97 && char <= 122) || char == 95 {
		return true
	}
	return false
}

func foreach(script string, arg *Params) (string, error) {
	end := 0 //切片名的结束位置
	for i, char := range script {
		if i > 0 {
			if !validCharacter(byte(char)) {
				end = i
				break
			}
		}
	}
	name := script[1:end] //要遍历的切片名称
	val := arg.Get(name)
	if val == nil {
		return "", fmt.Errorf("Syntax error, key %s not found", name)
	}
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return "", fmt.Errorf("Syntax error, key %s is not a slice type", name)
	}
	length := v.Len()
	var sql bytes.Buffer
	for index := 0; index < length; index++ {
		item := v.Index(index)
		fragment, err := util.ReplaceByKeyword(script[end:], ':', func(i int, s int, e int, c string) (string, error) {
			if c == "val" {
				arg.Add(fmt.Sprintf("%s%d", c, index), item.Interface())
			} else {
				arg.Add(fmt.Sprintf("%s%d", c, index), item.MapIndex(reflect.ValueOf(c)).Interface())
			}
			return fmt.Sprintf(":%s%d", c, index), nil
		})
		if err != nil {
			return "", err
		}
		if index > 0 {
			sql.WriteString(",")
		}
		sql.WriteString(fragment)
	}
	return sql.String(), nil
}

func empty(script string, arg *Params) (string, error) {
	ignore := false
	fragment, err := util.ReplaceByKeyword(script, ':', func(i int, s int, e int, c string) (string, error) {
		if arg.Get(c) == nil {
			ignore = true
			return "", nil
		}
		return fmt.Sprintf(":%s", c), nil
	})
	if ignore {
		return "", err
	}
	return fragment, err
}

func bindParams(sql string, arg Params) (string, []interface{}, error) {
	newSql, err := util.ReplaceBetween(sql, "{", "}", func(index int, start int, end int, content string) (string, error) {
		if content != "" {
			if content[0] == '@' {
				return foreach(content, &arg)
			} else {
				return empty(content, &arg)
			}
		}
		return "", nil
	})
	if err != nil {
		return "", nil, err
	}
	params := make([]interface{}, 0, 20)
	newSql, err = util.ReplaceByKeyword(newSql, ':', func(index int, start int, end int, content string) (string, error) {
		val := arg.Get(content)
		params = append(params, val)
		return "?", nil
	})
	return newSql, params, err
}

func insert(table string, arg interface{}) (string, []interface{}, error) {
	if table == "" || arg == nil {
		return "", nil, errors.ErrParamNotNil
	}

	typeOfArg := reflect.TypeOf(arg)
	if typeOfArg.Kind() == reflect.Ptr {
		typeOfArg = typeOfArg.Elem()
	}

	var sql bytes.Buffer
	var sql2 bytes.Buffer
	params := make([]interface{}, 0, 20)
	flag := 0 //标识

	sql.WriteString("insert into ")
	sql.WriteString(table)
	sql.WriteString(" (")
	sql2.WriteString(") values (")

	if typeOfArg.Kind() == reflect.Struct {
		p := NewStructParams(arg)
		if p.Size() < 1 {
			return "", nil, errors.ErrParamEmpty
		}
		for i := 0; i < typeOfArg.NumField(); i++ {
			field := typeOfArg.Field(i)
			f := NewField(&field)
			if f.NotIgnore() {
				if flag == 0 {
					flag++
				} else {
					sql.WriteString(",")
					sql2.WriteString(",")
				}
				name := f.AttrName()
				if name == "" {
					name = f.Name
				}
				sql.WriteString(name)
				sql2.WriteString("?")
				params = append(params, p.Get(f.Name))
			}
		}
	}

	if typeOfArg.Kind() == reflect.Map {
		var p Params
		v, ok := arg.(container.JsonMap)
		if ok {
			p = NewMapParams(v.Map())
		} else {
			p = NewMapParams(arg.(map[string]interface{}))
		}
		if p.Size() < 1 {
			return "", nil, errors.ErrParamEmpty
		}
		p.Iterator(func(key string, value interface{}) {
			if flag == 0 {
				flag++
			} else {
				sql.WriteString(",")
				sql2.WriteString(",")
			}
			sql.WriteString(string(key))
			sql2.WriteString("?")
			params = append(params, value)
		})
	}

	sql.WriteString(sql2.String())
	sql.WriteString(")")
	return sql.String(), params, nil
}

func batchInsert(table string, args ...interface{}) (string, []interface{}, error) {
	if table == "" || args == nil || len(args) == 0 {
		return "", nil, errors.ErrParamNotNil
	}

	var sql bytes.Buffer
	var sql2 bytes.Buffer
	params := make([]interface{}, 0, 20)
	columns := make([]string, 0, 20)
	var p Params

	for index, arg := range args {
		typeOfArg := reflect.TypeOf(arg)
		if typeOfArg.Kind() == reflect.Ptr {
			typeOfArg = typeOfArg.Elem()
		}

		if index == 0 {
			sql.WriteString("insert into ")
			sql.WriteString(table)
			sql.WriteString(" (")
			sql2.WriteString(") values ")
		} else {
			sql2.WriteString(",")
		}

		if typeOfArg.Kind() == reflect.Struct {
			p = NewStructParams(arg)
			if p.Size() < 1 {
				return "", nil, errors.ErrParamEmpty
			}
			if index == 0 {
				flag := 0 //标识
				for i := 0; i < typeOfArg.NumField(); i++ {
					field := typeOfArg.Field(i)
					f := NewField(&field)
					if f.NotIgnore() {
						if flag == 0 {
							flag++
						} else {
							sql.WriteString(",")
						}
						name := f.AttrName()
						if name == "" {
							name = f.Name
						}
						sql.WriteString(name)
						columns = append(columns, f.Name)
					}
				}
			}

		}

		if typeOfArg.Kind() == reflect.Map {
			v, ok := arg.(container.JsonMap)
			if ok {
				p = NewMapParams(v.Map())
			} else {
				p = NewMapParams(arg.(map[string]interface{}))
			}
			if p.Size() < 1 {
				return "", nil, errors.ErrParamEmpty
			}
			if index == 0 {
				columns = p.Names()
				for i, v := range columns {
					if i > 0 {
						sql.WriteString(",")
					}
					sql.WriteString(v)
				}
			}
		}

		sql2.WriteString("(")
		for i, v := range columns {
			if i > 0 {
				sql2.WriteString(",")
			}
			sql2.WriteString("?")
			params = append(params, p.Get(v))
		}
		sql2.WriteString(")")
	}
	sql.WriteString(sql2.String())
	return sql.String(), params, nil
}
