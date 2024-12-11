package inference

import (
	"fmt"
)

// Rule represents a rule in the knowledge base
// A rule is a logical expression that can be evaluated to true, false or unknown
type Rule struct {
	FactTargetID string `json:"fact_target_id"`
	Description  string `json:"description"`
	Question     string `json:"question"`
	Expression   string `json:"expression"`
}

func (rule *Rule) evaluate(facts map[string]Fact) (bool, []string, error) {
	output, derived, err := Calculate(rule.Expression, facts)
	if err != nil {
		return false, nil, err
	}
	result, ok := output.(bool)
	if !ok {
		return false, nil, fmt.Errorf("expression did not evaluate to a boolean")
	}
	return result, derived, nil
}

type WeightedRule struct {
	Rule
	Weight float64
}
