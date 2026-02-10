package inference

import (
	"github.com/expr-lang/expr"
)

// IntentType represents the type of user intent.
type IntentType string

const (
	IntentQuery    IntentType = "query"
	IntentDecision IntentType = "decision"
	IntentAnalysis IntentType = "analysis"
	IntentAction   IntentType = "action"
)

// Intent represents a classified user intent.
type Intent struct {
	Type        IntentType        `json:"type"`
	Description string            `json:"description"`
	Context     map[string]string `json:"context,omitempty"`
}

// IntentRule maps an expr-lang expression to an intent type.
type IntentRule struct {
	Expression string     `json:"expression"`
	IntentType IntentType `json:"intent_type"`
	Weight     float64    `json:"weight"`
}

// IntentClassifier classifies facts into an intent using rules.
type IntentClassifier struct {
	Rules []IntentRule `json:"rules"`
}

// Classify evaluates intent rules against facts and returns the best matching intent.
func (ic *IntentClassifier) Classify(facts map[string]Fact) (Intent, error) {
	if len(ic.Rules) == 0 {
		return Intent{Type: IntentQuery, Description: "default"}, nil
	}

	env := make(map[string]interface{})
	for k, v := range facts {
		env[k] = v.Value
	}

	bestWeight := 0.0
	bestIntent := Intent{Type: IntentQuery, Description: "default"}

	for _, rule := range ic.Rules {
		program, err := expr.Compile(rule.Expression, expr.Env(env))
		if err != nil {
			continue
		}
		output, err := expr.Run(program, env)
		if err != nil {
			continue
		}
		result, ok := output.(bool)
		if !ok || !result {
			continue
		}
		if rule.Weight > bestWeight {
			bestWeight = rule.Weight
			bestIntent = Intent{
				Type:        rule.IntentType,
				Description: rule.Expression,
			}
		}
	}
	return bestIntent, nil
}
