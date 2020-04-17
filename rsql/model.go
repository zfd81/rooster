package rsql

import (
	"reflect"
	"strings"

	"github.com/spf13/cast"
	"github.com/zfd81/rooster/util"
)

const (
	TagName                 = "rsql"
	AttributeSeparator byte = ';'
	KvSeparator             = ":"
	AttrName                = "name"
)

type Modeler interface {
	TableName() string
}

type Field struct {
	*reflect.StructField
	Index      []int
	Path       string
	ignore     bool
	attributes map[string]string
}

func (f *Field) NotIgnore() bool {
	return !f.ignore
}

func (f *Field) GetAttr(name string) string {
	if f.attributes == nil {
		return ""
	}
	return f.attributes[name]
}

func (f *Field) AttrName() string {
	return f.GetAttr(AttrName)
}

func NewField(field *reflect.StructField) *Field {
	f := &Field{StructField: field}
	content := field.Tag.Get(TagName)
	if content != "" {
		if strings.TrimSpace(content) == "-" {
			f.ignore = true
		} else {
			attrs := make(map[string]string)
			if !strings.Contains(content, cast.ToString(AttributeSeparator)) && strings.TrimSpace(content) != "" {
				kv := strings.Split(content, KvSeparator)
				attrs[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			} else {
				util.ReplaceByKeyword(content, AttributeSeparator, func(i int, s int, e int, c string) (string, error) {
					if strings.TrimSpace(c) != "" {
						kv := strings.Split(c, KvSeparator)
						attrs[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
					}
					return c, nil
				})
			}
			f.attributes = attrs
		}
	}
	return f
}

func FieldByIndexes(v reflect.Value, indexes []int) reflect.Value {
	for _, i := range indexes {
		v = reflect.Indirect(v).Field(i)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			alloc := reflect.New(util.Deref(v.Type()))
			v.Set(alloc)
		}
		if v.Kind() == reflect.Map && v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
	}
	return v
}

func GetNameMapping(t reflect.Type) map[string][]int {
	nameMapping := map[string][]int{}
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		if field.Anonymous {
			nm := GetNameMapping(field.Type)
			for k, v := range nm {
				nameMapping[k] = util.InsertIntSlice(v, []int{i}, 0)
			}
		} else {
			f := NewField(&field)
			if f.NotIgnore() {
				name := f.AttrName()
				if name == "" {
					name = f.Name
				}
				nameMapping[strings.ToLower(name)] = []int{i}
			}
		}

	}
	return nameMapping
}

func FieldMapping(t reflect.Type) map[string]*Field {
	t = util.Deref(t)
	mapping := map[string]*Field{}
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		if field.Anonymous {
			fm := FieldMapping(field.Type)
			for k, f := range fm {
				f.Index = util.InsertIntSlice(f.Index, []int{i}, 0)
				f.Path = field.Name + "." + f.Path
				mapping[k] = f
			}
		} else {
			f := NewField(&field)
			f.Index = []int{i}
			f.Path = f.Name
			if f.NotIgnore() {
				name := f.AttrName()
				if name == "" {
					name = f.Name
				}
				//nameMapping[strings.ToLower(name)] = f
				mapping[name] = f
			}
		}
	}
	return mapping
}
