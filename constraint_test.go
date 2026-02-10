package inference

import "testing"

func TestConstraint_Satisfied(t *testing.T) {
	c := Constraint{
		Description: "age must be >= 18",
		Type:        ConstraintHard,
		Expression:  "age >= 18",
		Weight:      1.0,
	}
	facts := map[string]Fact{"age": {ID: "age", Value: 20}}
	satisfied, err := c.Satisfied(facts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !satisfied {
		t.Error("Expected constraint to be satisfied")
	}
}

func TestConstraint_NotSatisfied(t *testing.T) {
	c := Constraint{
		Description: "age must be >= 18",
		Type:        ConstraintHard,
		Expression:  "age >= 18",
		Weight:      1.0,
	}
	facts := map[string]Fact{"age": {ID: "age", Value: 15}}
	satisfied, err := c.Satisfied(facts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if satisfied {
		t.Error("Expected constraint to not be satisfied")
	}
}

func TestConstraintSet_Identify(t *testing.T) {
	cs := ConstraintSet{
		Constraints: []Constraint{
			{Description: "age check", Type: ConstraintHard, Expression: "age >= 18", Weight: 1.0},
			{Description: "score check", Type: ConstraintSoft, Expression: "score > 50", Weight: 0.5},
		},
	}
	facts := map[string]Fact{
		"age":   {ID: "age", Value: 20},
		"score": {ID: "score", Value: 30},
	}
	active, err := cs.Identify(facts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(active) != 2 {
		t.Errorf("Expected 2 active constraints, got %d", len(active))
	}
}
