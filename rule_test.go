package inference

import "testing"

func TestRule(t *testing.T) {
	rule := Rule{
		Description: "Be of legal age",
		Expression:  "age >= 18",
	}

	facts := map[string]Fact{"age": {ID: "age", Value: 20}}

	result, _, err := rule.evaluate(facts)
	if err != nil {
		t.Errorf("Error evaluating rule: %s", err)
	}
	if !result {
		t.Errorf("Expected rule to be true")
	}

	facts = map[string]Fact{"age": {ID: "age", Value: 15}}

	result, _, err = rule.evaluate(facts)
	if err != nil {
		t.Errorf("Error evaluating rule: %s", err)
	}
	if result {
		t.Errorf("Expected rule to be false")
	}

	facts = map[string]Fact{"other": {ID: "age", Value: 15}}

	result, _, err = rule.evaluate(facts)
	if err == nil {
		t.Errorf("Expected error evaluating rule")
	}

}

func TestBadRule(t *testing.T) {
	rule := Rule{
		Description: "be of legal age",
		Expression:  "age + 18",
	}

	facts := map[string]Fact{"age": {ID: "age", Value: 20}}

	_, _, err := rule.evaluate(facts)
	if err == nil {
		t.Errorf("Expected error evaluating rule")
	}
	t.Logf("Error evaluating rule: %s", err)
}
