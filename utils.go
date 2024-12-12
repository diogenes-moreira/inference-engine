package inference

import (
	"encoding/json"
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	log "github.com/sirupsen/logrus"
	"os"
)

func Calculate(sExpression string, params map[string]Fact) (interface{}, []string, error) {
	par := make(map[string]interface{})
	for k, v := range params {
		par[k] = v.Value
	}
	expression, err := expr.Compile(sExpression, expr.Env(par))
	out := make([]string, 0)
	if err != nil {
		return "", nil, err
	}
	fmt.Println(expression)
	output, err := expr.Run(expression, par)
	if err != nil {
		return "", nil, err
	}
	out = extractFacts(expression.Node(), out)
	return output, out, nil
}

func extractFacts(node ast.Node, out []string) []string {
	visitor := &NodeVisitor{Dependencies: out}
	ast.Walk(&node, visitor)
	return visitor.Dependencies
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

type NodeVisitor struct {
	Dependencies []string
}

func (v *NodeVisitor) Visit(node *ast.Node) {
	identifierNode, ok := (*node).(*ast.IdentifierNode)
	if ok {
		v.Dependencies = append(v.Dependencies, identifierNode.Value)
	}
}

func (kb *KnowledgeBase) Dump(filename string) {
	data, err := json.Marshal(kb)
	if err != nil {
		log.Errorf("Error dumping knowledge base: %s", err)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Errorf("Error creating file: %s", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Error closing file: %s", err)
		}
	}(file)

	_, err = file.Write(data)
	if err != nil {
		log.Errorf("Error writing to file: %s", err)
		return
	}
}
