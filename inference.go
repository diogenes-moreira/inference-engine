package inference

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"slices"
	"strings"
)

type Inference struct {
	Description      string         `json:"description"`
	Rules            []WeightedRule `json:"rules"`
	FactID           string         `json:"fact_id"`
	FactValue        interface{}    `json:"fact_value"`
	IsValeCalculated bool           `json:"is_value_calculated"`
	IsIDCalculated   bool           `json:"is_id_calculated"`
	CountOfTrue      int            `json:"count_of_true"`
	Probability      float64        `json:"probability"`
	// The Probability of the inference being true,
	// this is used to sort inferences
}

func (inf *Inference) infer(facts map[string]Fact) (string, interface{}, []string, error) {
	var derived []string
	for _, rule := range inf.Rules {
		result, d, err := rule.evaluate(facts)
		if err != nil {
			return "", nil, nil, err
		}
		if !result {
			return "", nil, nil, fmt.Errorf("rule %s is false", rule.Description)
		}
		derived = append(derived, d...)
	}
	id, err := inf.getFactID(facts)
	if err != nil {
		return "", nil, nil, err
	}
	value, d, err := inf.getFactValue(facts)
	if err != nil {
		return "", nil, nil, err
	}
	derived = append(derived, d...)
	derived = unique(derived)
	return id, value, derived, nil
}

func (inf *Inference) getFactID(facts map[string]Fact) (string, error) {
	if !inf.IsIDCalculated {
		return inf.FactID, nil
	} else {
		id, _, err := Calculate(inf.FactID, facts)
		if err != nil {
			return "", err
		}
		return id.(string), nil
	}
}

func (inf *Inference) getFactValue(facts map[string]Fact) (interface{}, []string, error) {
	var empty []string
	if !inf.IsValeCalculated {
		return inf.FactValue, empty, nil
	} else {
		value, derived, err := Calculate(inf.FactValue.(string), facts)
		if err != nil {
			return "", empty, err
		}
		return value, derived, nil
	}

}

func (inf *Inference) IsNeeded(facts map[string]Fact) bool {
	id, err := inf.getFactID(facts)
	if err != nil {
		log.Errorf("Error getting fact id: %s", err)
		return false
	}
	for key := range facts {
		if strings.HasPrefix(key, id) {
			return false
		}
	}
	return true
}

func (inf *Inference) Certainty(facts map[string]Fact) float64 {
	total := 0.0
	accomplished := 0.0
	for _, rule := range inf.Rules {
		total += rule.Weight
		_, _, err := rule.evaluate(facts)
		if err != nil {
			continue
		}
	}
	return accomplished / total
}

// PendingRules returns the rules that have not been achieved yet
// sorted by weight
func (inf *Inference) PendingRules(facts map[string]Fact) []WeightedRule {
	var pending []WeightedRule
	for _, rule := range inf.Rules {
		_, _, err := rule.evaluate(facts)
		if err != nil {
			pending = append(pending, rule)
		}
	}
	slices.SortFunc(pending, func(a, b WeightedRule) int {
		return int(a.Weight - b.Weight)
	})
	return pending
}
