package inference

import (
	"github.com/expr-lang/expr"
)

// Entity represents an extracted entity from input facts.
type Entity struct {
	FactID     string      `json:"fact_id"`
	Value      interface{} `json:"value"`
	Source     string      `json:"source"`
	Confidence float64     `json:"confidence"`
}

// ExtractionRule defines how to extract an entity from facts.
type ExtractionRule struct {
	FactID          string  `json:"fact_id"`
	Expression      string  `json:"expression"`
	Source          string  `json:"source"`
	ConfidenceValue float64 `json:"confidence"`
}

// EntityExtractor extracts entities from facts using extraction rules.
type EntityExtractor struct {
	Rules []ExtractionRule `json:"rules"`
}

// Extract evaluates extraction rules against facts and returns discovered entities.
func (ee *EntityExtractor) Extract(facts map[string]Fact) ([]Entity, error) {
	if len(ee.Rules) == 0 {
		return nil, nil
	}

	env := make(map[string]interface{})
	for k, v := range facts {
		env[k] = v.Value
	}

	var entities []Entity
	for _, rule := range ee.Rules {
		program, err := expr.Compile(rule.Expression, expr.Env(env))
		if err != nil {
			continue
		}
		output, err := expr.Run(program, env)
		if err != nil {
			continue
		}
		// Skip nil or false results
		if output == nil {
			continue
		}
		if boolVal, ok := output.(bool); ok && !boolVal {
			continue
		}
		entities = append(entities, Entity{
			FactID:     rule.FactID,
			Value:      output,
			Source:     rule.Source,
			Confidence: rule.ConfidenceValue,
		})
	}
	return entities, nil
}
