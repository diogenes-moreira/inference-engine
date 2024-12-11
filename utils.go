package inference

import (
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
)

func Calculate(sExpression string, params map[string]Fact) (interface{}, []string, error) {
	//TODO: Remove Value from Fact struct
	expression, err := expr.Compile(sExpression, expr.Env(params))
	out := make([]string, 0)
	if err != nil {
		return "", nil, err
	}
	fmt.Println(expression)
	output, err := expr.Run(expression, params)
	if err != nil {
		return "", nil, err
	}
	out = extractFacts(expression.Node(), out)
	return output, out, nil
}

func extractFacts(node interface{}, out []string) []string {
	//TODO: refactor this to use a visitor pattern
	// Add More cases
	member, ok := node.(*ast.MemberNode)
	if ok {
		out = extractFacts(member.Node, out)
	}
	binary, ok := node.(*ast.BinaryNode)
	if ok {
		out = extractFacts(binary.Left, out)
		out = extractFacts(binary.Right, out)
	}
	identifier, ok := node.(*ast.IdentifierNode)
	if ok {
		out = append(out, identifier.Value)
	}
	return out
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
