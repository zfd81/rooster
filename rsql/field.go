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

type Field struct {
	*reflect.StructField
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
