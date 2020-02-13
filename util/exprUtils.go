package util

import (
	"github.com/antonmedv/expr"
)

func ExprParsing(env map[string]interface{}, expression string) (interface{}, error) {
	output, err := expr.Eval(expression, env)
	if err != nil {
		return nil, err
	}
	return output, nil
}
