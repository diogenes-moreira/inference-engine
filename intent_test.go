package inference

import "testing"

func TestIntentClassifier_Classify(t *testing.T) {
	classifier := IntentClassifier{
		Rules: []IntentRule{
			{Expression: "action == true", IntentType: IntentAction, Weight: 1.0},
			{Expression: "query == true", IntentType: IntentQuery, Weight: 0.5},
		},
	}
	facts := map[string]Fact{
		"action": {ID: "action", Value: true},
		"query":  {ID: "query", Value: true},
	}
	intent, err := classifier.Classify(facts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if intent.Type != IntentAction {
		t.Errorf("Expected action intent, got %s", intent.Type)
	}
}

func TestIntentClassifier_ClassifyDefault(t *testing.T) {
	classifier := IntentClassifier{}
	facts := map[string]Fact{"x": {ID: "x", Value: 1}}
	intent, err := classifier.Classify(facts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if intent.Type != IntentQuery {
		t.Errorf("Expected default query intent, got %s", intent.Type)
	}
}

func TestIntentClassifier_NoMatch(t *testing.T) {
	classifier := IntentClassifier{
		Rules: []IntentRule{
			{Expression: "action == true", IntentType: IntentAction, Weight: 1.0},
		},
	}
	facts := map[string]Fact{
		"action": {ID: "action", Value: false},
	}
	intent, err := classifier.Classify(facts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if intent.Type != IntentQuery {
		t.Errorf("Expected default query intent when no match, got %s", intent.Type)
	}
}
