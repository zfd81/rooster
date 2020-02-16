package util

import (
	"fmt"
	"reflect"

	"github.com/zfd81/rooster/errors"
)

type IteratorFunc func(index int, key string, value interface{}, field reflect.StructField)

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
		iterator(i, typeOfArg.Field(i).Name, valueOfArg.Field(i).Interface(), typeOfArg.Field(i))
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

func BaseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = Deref(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}
