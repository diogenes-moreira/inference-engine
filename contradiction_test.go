package inference

import "testing"

func TestContradiction_Detect(t *testing.T) {
	contra := Contradiction{
		Description: "description",
		Facts: []Fact{
			{ID: "fact1", Value: "value1"},
			{ID: "fact2", Value: "value2"},
		},
	}

	facts := map[string]Fact{
		"fact1": {
			ID:    "fact1",
			Value: "value1",
		},
		"fact2": {
			ID:    "fact2",
			Value: "value2",
		},
	}

	if !contra.Detect(facts) {
		t.Errorf("Expected true, got false")
	}
}

func TestContradiction_Detect_FactNotFound(t *testing.T) {
	contra := Contradiction{
		Description: "description",
		Facts: []Fact{
			{ID: "fact1", Value: "value1"},
			{ID: "fact2", Value: "value2"},
		},
	}

	facts := map[string]Fact{
		"fact2": {
			ID:    "fact2",
			Value: "value2",
		},
	}

	if contra.Detect(facts) {
		t.Errorf("Expected false, got true")
	}
}

func TestContradiction_Resolve(t *testing.T) {
	contra := Contradiction{
		Description: "description",
		Facts: []Fact{
			{ID: "fact1", Value: "value1"},
			{ID: "fact2", Value: "value2"},
		},
	}

	kb := &KnowledgeBase{
		Contradictions: []Contradiction{contra},
		Facts: map[string]Fact{
			"fact1": {ID: "fact1", Value: "value1"},
			"fact2": {ID: "fact2", Value: "value2"},
			"fact3": {ID: "fact3", Value: "value3", DerivedFrom: []string{"fact1", "fact2"}},
		},
	}

	if contra.Detect(kb.Facts) {
		contra.Resolve(kb)
	}
	if _, ok := kb.Facts["fact1"]; ok {
		t.Errorf("Expected fact1 to be deleted")
	}
	if _, ok := kb.Facts["fact2"]; ok {
		t.Errorf("Expected fact2 to be deleted")
	}
	if _, ok := kb.Facts["fact3"]; ok {
		t.Errorf("Expected fact3 to not be deleted")
	}
}
