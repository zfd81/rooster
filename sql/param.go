package sql

import (
	"github.com/zfd81/rooster/util"
)

type Params map[string]interface{}

func (p Params) Get(name string) (interface{}, bool) {
	val, ok := p[name]
	return val, ok
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

func NewParams() Params {
	return make(Params)
}

func NewMapParams(params map[string]interface{}) Params {
	if params == nil || len(params) < 1 {
		return make(Params)
	}
	return params
}

func NewStructParams(params interface{}) Params {
	newParams := make(map[string]interface{})
	util.StructIterator(params, func(index int, key string, value interface{}) {
		newParams[key] = value
	})
	return newParams
}
