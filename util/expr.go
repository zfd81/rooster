package util

import (
	"github.com/antonmedv/expr"
	"github.com/zfd81/rooster/errors"
)

func ExprParsing(env map[string]interface{}, expression string) (interface{}, error) {
	if expression == "" {
		return nil, errors.ErrParamNotNil
	}
	output, err := expr.Eval(expression, env)
	if err != nil {
		return nil, err
	}
	return output, nil
}
