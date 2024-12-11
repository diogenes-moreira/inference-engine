package inference

import "github.com/expr-lang/expr"

func Calculate(sExpression string, params map[string]Fact) (interface{}, error) {
	expression, err := expr.Compile(sExpression, expr.Env(params))
	if err != nil {
		return "", err
	}
	output, err := expr.Run(expression, params)
	if err != nil {
		return "", err
	}
	return output, nil
}
