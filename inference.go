package inference

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Inference struct {
	Description      string
	Rules            []WeightedRule
	FactID           string
	FactValue        interface{}
	IsValeCalculated bool
	IsIDCalculated   bool
}

func (inf *Inference) infer(facts map[string]Fact) (string, interface{}, error) {
	for _, rule := range inf.Rules {
		result, err := rule.evaluate(facts)
		if err != nil {
			return "", nil, err
		}
		if !result {
			return "", nil, fmt.Errorf("rule %s is false", rule.Description)
		}
	}
	id, err := inf.getFactID(facts)
	if err != nil {
		return "", nil, err
	}
	value, err := inf.getFactValue(facts)
	if err != nil {
		return "", nil, err
	}
	return id, value, nil
}

func (inf *Inference) getFactID(facts map[string]Fact) (string, error) {
	if !inf.IsIDCalculated {
		return inf.FactID, nil
	} else {
		id, err := Calculate(inf.FactID, facts)
		if err != nil {
			return "", err
		}
		return id.(string), nil
	}
}

func (inf *Inference) getFactValue(facts map[string]Fact) (interface{}, error) {
	if !inf.IsValeCalculated {
		return inf.FactValue, nil
	} else {
		value, err := Calculate(inf.FactValue.(string), facts)
		if err != nil {
			return "", err
		}
		return value, nil
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
		_, err := rule.evaluate(facts)
		if err != nil {
			continue
		}
	}
	return accomplished / total
}

func (inf *Inference) PendingRules(facts map[string]Fact) []WeightedRule {
	var pending []WeightedRule
	for _, rule := range inf.Rules {
		_, err := rule.evaluate(facts)
		if err != nil {
			pending = append(pending, rule)
		}
	}
	return pending
}
