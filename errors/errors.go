package errors

import "errors"

var (
	ErrParamNotNil = errors.New("Parameters cannot be nil")
	ErrParamType   = errors.New("Parameter type error")
	ErrParamEmpty  = errors.New("Parameter size is zero")
)
