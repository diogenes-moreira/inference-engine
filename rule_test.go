package inference

import "testing"

func TestRule(t *testing.T) {
	rule := Rule{
		Description: "Be of legal age",
		Expression:  "age.Values[\"value\"] >= 18",
	}

	facts := map[string]Fact{"age": {ID: "age", Values: map[string]interface{}{"value": 20}}}

	result, err := rule.evaluate(facts)
	if err != nil {
		t.Errorf("Error evaluating rule: %s", err)
	}
	if !result {
		t.Errorf("Expected rule to be true")
	}

	facts = map[string]Fact{"age": {ID: "age", Values: map[string]interface{}{"value": 15}}}

	result, err = rule.evaluate(facts)
	if err != nil {
		t.Errorf("Error evaluating rule: %s", err)
	}
	if result {
		t.Errorf("Expected rule to be false")
	}

	facts = map[string]Fact{"other": {ID: "age", Values: map[string]interface{}{"value": 15}}}

	result, err = rule.evaluate(facts)
	if err == nil {
		t.Errorf("Expected error evaluating rule")
	}
	t.Logf("Error evaluating rule: %s", err)

}

func TestBadRule(t *testing.T) {
	rule := Rule{
		Description: "be of legal age",
		Expression:  "age + 18",
	}

	facts := map[string]Fact{"age": {ID: "age", Values: map[string]interface{}{"value": 20}}}

	_, err := rule.evaluate(facts)
	if err == nil {
		t.Errorf("Expected error evaluating rule")
	}
	t.Logf("Error evaluating rule: %s", err)
}
