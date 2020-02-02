package util

import (
	"github.com/zfd81/rooster/errors"
	"reflect"
)

type IteratorFunc func(index int, key string, value interface{})

func StructIterator(arg interface{}, iterator IteratorFunc) error {
	if arg == nil {
		return errors.ErrParamNotNil
	}
	typeOfArg := reflect.TypeOf(arg)
	valueOfArg := reflect.ValueOf(arg)
	if valueOfArg.Kind() == reflect.Ptr {
		typeOfArg = typeOfArg.Elem()
		valueOfArg = valueOfArg.Elem()
	}
	if valueOfArg.Kind() != reflect.Struct || !valueOfArg.IsValid() {
		return errors.ErrParamType
	}
	for i := 0; i < valueOfArg.NumField(); i++ {
		iterator(i, typeOfArg.Field(i).Name, valueOfArg.Field(i).Interface())
	}
	return nil
}

// Deref is Indirect for reflect.Types
func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// FieldByIndexes returns a value for the field given by the struct traversal
// for the given value.
func FieldByIndexes(v reflect.Value, indexes []int) reflect.Value {
	for _, i := range indexes {
		v = reflect.Indirect(v).Field(i)
		// if this is a pointer and it's nil, allocate a new value and set it
		if v.Kind() == reflect.Ptr && v.IsNil() {
			alloc := reflect.New(Deref(v.Type()))
			v.Set(alloc)
		}
		if v.Kind() == reflect.Map && v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
	}
	return v
}
