package rsql

import (
	"reflect"

	"github.com/zfd81/rooster/types/container"

	"github.com/zfd81/rooster/util"
)

const (
	SingleParameterName string = "val"
)

type Params map[string]interface{}

func (p Params) Get(name string) interface{} {
	return p[name]
}

func (p Params) Add(name string, value interface{}) Params {
	_, ok := p[name]
	if ok {
		delete(p, name)
	}
	p[name] = value
	return p
}

func (p Params) Remove(name string) Params {
	_, ok := p[name]
	if ok {
		delete(p, name)
	}
	return p
}

func (p Params) Names() []string {
	if len(p) < 1 {
		return nil
	}
	names := make([]string, 0, 10)
	for k := range p {
		names = append(names, k)
	}
	return names
}

func (p Params) Size() int {
	return len(p)
}

func (p Params) Iterator(handler func(key string, value interface{})) {
	if len(p) < 1 {
		return
	}
	for k, v := range p {
		handler(k, v)
	}
}

func (p Params) Clone() Params {
	p2 := make(Params, len(p))
	for k, v := range p {
		p2[k] = v
	}
	return p2
}

func NewParams(arg interface{}) Params {
	if arg != nil {
		value := reflect.ValueOf(arg)
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				return make(Params)
			}
			value = value.Elem()
		}
		if value.Kind() == reflect.Map {
			p, ok := value.Interface().(Params)
			if ok {
				return p
			}

			jm, ok := value.Interface().(container.JsonMap)
			if ok {
				return NewMapParams(jm.Map())
			}

			m, ok := value.Interface().(map[string]interface{})
			if ok {
				return NewMapParams(m)
			}
		}
		if value.Kind() == reflect.Struct {
			return NewStructParams(value.Interface())
		}
		if value.Kind() == reflect.String || value.Kind() == reflect.Int || value.Kind() == reflect.Int64 {
			return NewSingleParams(value.Interface())
		}
	}
	return make(Params)
}

func NewMapParams(params map[string]interface{}) Params {
	if params == nil || len(params) < 1 {
		return make(Params)
	}
	return params
}

func NewStructParams(params interface{}) Params {
	p := make(map[string]interface{})
	util.StructIterator(params, func(index int, key string, value interface{}, field reflect.StructField) {
		p[key] = value
	})
	return p
}

func NewSingleParams(param interface{}) Params {
	p := make(Params)
	p.Add(SingleParameterName, param)
	return p
}
