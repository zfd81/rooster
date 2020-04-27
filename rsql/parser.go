package rsql

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/zfd81/rooster/types/container"

	"github.com/zfd81/rooster/errors"
	"github.com/zfd81/rooster/util"
)

func validCharacter(char byte) bool {
	if (char >= 48 && char <= 57) || (char >= 65 && char <= 90) || (char >= 97 && char <= 122) || char == 46 || char == 95 || char == 9 || char == 10 || char == 32 {
		return true
	}
	return false
}

func foreach(script string, arg *Params) (string, error) {
	start := 0      //切片名称开始位置
	end := 0        //切片名称结束位置
	open := 0       //[方括号开始位置
	close := 0      //]方括号结束位置
	var name string //要遍历的切片名称
	separator := "" //分隔符
	for i, char := range script {
		if char != 32 {
			if start == 0 && validCharacter(byte(char)) {
				start = i
				continue
			}
			if !validCharacter(byte(char)) && open == 0 && start > 0 {
				end = i
			}
			if char == 91 {
				if start == 0 || open > 0 {
					return "", fmt.Errorf("Syntax error,near '%s'", script[:i+1])
				}
				open = i
				continue
			}
			if char == 93 {
				if open == 0 {
					return "", fmt.Errorf("Syntax error,near '%s'", script[:i+1])
				}
				close = i
				break
			}
		}
	}
	if close == 0 {
		return "", fmt.Errorf("Syntax error,near '%s'", script)
	}
	name = strings.TrimSpace(script[start:end])
	if close-open > 1 {
		separator = script[open+1 : close]
	}
	content := script[close+1:]
	val := arg.Get(name)
	if val == nil {
		return "", fmt.Errorf("Syntax error, key '%s' not found, near '%s'", name, script[:end])
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
		item := reflect.Indirect(v.Index(index))
		fragment, err := util.ReplaceByKeyword(content, ':', func(i int, s int, e int, c string) (string, error) {
			if strings.HasPrefix(c, "this.") {
				key := c[5:]
				if key == "val" {
					arg.Add(fmt.Sprintf("%s.%s%d", name, c, index), item.Interface())
				} else {
					value := item.MapIndex(reflect.ValueOf(key))
					if value.IsValid() {
						arg.Add(fmt.Sprintf("%s.%s%d", name, c, index), value.Interface())
					} else {
						arg.Add(fmt.Sprintf("%s.%s%d", name, c, index), new(interface{}))
					}
				}
				return fmt.Sprintf(":%s.%s%d", name, c, index), nil
			}
			return fmt.Sprintf(":%s", c), nil
		})
		if err != nil {
			return "", err
		}
		if index > 0 {
			sql.WriteString(" ")
			sql.WriteString(separator)
			sql.WriteString(" ")
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

	var sql bytes.Buffer
	var sql2 bytes.Buffer
	var p Params

	typeOfArg := reflect.TypeOf(arg)
	if typeOfArg.Kind() == reflect.Ptr {
		typeOfArg = typeOfArg.Elem()
	}
	if typeOfArg.Kind() == reflect.Struct {
		p = NewStructParams(arg)
	} else if typeOfArg.Kind() == reflect.Map {
		v, ok := arg.(container.JsonMap)
		if ok {
			p = NewMapParams(v.Map())
		} else {
			p = NewMapParams(arg.(map[string]interface{}))
		}
	}
	if p.Size() < 1 {
		return "", nil, errors.ErrParamEmpty
	}
	params := make([]interface{}, 0, p.Size())
	sql.WriteString("insert into ")
	sql.WriteString(table)
	sql.WriteString(" (")
	sql2.WriteString(") values (")

	flag := 0 //标识
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
	var p Params
	var columns []string
	params := make([]interface{}, 0, 50)

	for index, arg := range args {
		typeOfArg := reflect.TypeOf(arg)
		if typeOfArg.Kind() == reflect.Ptr {
			typeOfArg = typeOfArg.Elem()
		}
		if typeOfArg.Kind() == reflect.Struct {
			p = NewStructParams(arg)
		} else if typeOfArg.Kind() == reflect.Map {
			v, ok := arg.(container.JsonMap)
			if ok {
				p = NewMapParams(v.Map())
			} else {
				p = NewMapParams(arg.(map[string]interface{}))
			}
		}
		if p.Size() < 1 {
			return "", nil, errors.ErrParamEmpty
		}
		if index == 0 {
			columns = p.Names()
			sql.WriteString("insert into ")
			sql.WriteString(table)
			sql.WriteString(" (")
			sql2.WriteString(") values ")
			for i, v := range columns {
				if i > 0 {
					sql.WriteString(",")
				}
				sql.WriteString(v)
			}
		} else {
			sql2.WriteString(",")
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
