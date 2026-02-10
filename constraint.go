package inference

import (
	"github.com/expr-lang/expr"
)

// ConstraintType represents whether a constraint is hard (must satisfy) or soft (should satisfy).
type ConstraintType string

const (
	ConstraintHard ConstraintType = "hard"
	ConstraintSoft ConstraintType = "soft"
)

// Constraint defines a condition that should or must be satisfied.
type Constraint struct {
	Description string         `json:"description"`
	Type        ConstraintType `json:"type"`
	Expression  string         `json:"expression"`
	Weight      float64        `json:"weight"`
}

// Satisfied evaluates whether this constraint is met given the current facts.
func (c *Constraint) Satisfied(facts map[string]Fact) (bool, error) {
	env := make(map[string]interface{})
	for k, v := range facts {
		env[k] = v.Value
	}
	program, err := expr.Compile(c.Expression, expr.Env(env))
	if err != nil {
		return false, err
	}
	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}
	result, ok := output.(bool)
	if !ok {
		return false, nil
	}
	return result, nil
}

// ConstraintSet holds a collection of constraints.
type ConstraintSet struct {
	Constraints []Constraint `json:"constraints"`
}

// Identify returns constraints that are relevant (evaluable) given the current facts.
func (cs *ConstraintSet) Identify(facts map[string]Fact) ([]Constraint, error) {
	var active []Constraint
	for _, c := range cs.Constraints {
		_, err := c.Satisfied(facts)
		if err != nil {
			// Constraint references facts not yet available — skip
			continue
		}
		active = append(active, c)
	}
	return active, nil
}
