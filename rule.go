package inference

import (
	"fmt"
)

// Rule represents a rule in the knowledge base
// A rule is a logical expression that can be evaluated to true, false or unknown
type Rule struct {
	FactTargetID string
	Description  string
	Question     string
	Expression   string
}

func (rule *Rule) evaluate(facts map[string]Fact) (bool, error) {
	output, err := Calculate(rule.Expression, facts)
	if err != nil {
		return false, err
	}
	result, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf("expression did not evaluate to a boolean")
	}
	return result, nil
}

type WeightedRule struct {
	Rule
	Weight float64
}

func (rule *Rule) GenerateFact(values map[string]interface{}) Fact {
	return Fact{
		ID:     rule.FactTargetID,
		Values: values,
	}

}
