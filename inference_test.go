package inference

import "testing"

func TestInference_IsNeeded(t *testing.T) {
	inference := Inference{
		Description: "Can vote",
		Rules: []WeightedRule{{
			Weight: 1,
			Rule: Rule{
				Description: "Be of legal age",
				Expression:  "age >= 18",
			}},
		},
		FactID:           "can_vote",
		FactValue:        "age >= 18",
		IsIDCalculated:   false,
		IsValeCalculated: true,
	}

	kb := KnowledgeBase{
		Inferences: []Inference{inference},
		Facts:      map[string]Fact{"age": {ID: "age", Values: map[string]interface{}{"age": 20}}},
	}

	if !inference.IsNeeded(kb.Facts) {
		t.Errorf("Expected inference to be needed")
	}
}

func TestInference_IsNeededFalse(t *testing.T) {
	inference := Inference{
		Description: "Can vote",
		Rules: []WeightedRule{{
			Weight: 1,
			Rule: Rule{
				Description: "Be of legal age",
				Expression:  "age >= 18",
			}},
		},
		FactID:           "can_vote",
		FactValue:        "age >= 18",
		IsIDCalculated:   false,
		IsValeCalculated: true,
	}

	kb := KnowledgeBase{
		Inferences: []Inference{inference},
		Facts:      map[string]Fact{"can_vote": {ID: "can_vote", Values: map[string]interface{}{"can_vote": false}}},
	}

	if inference.IsNeeded(kb.Facts) {
		t.Errorf("Expected inference to not be needed")
	}
}

func TestInference_Facts(t *testing.T) {
	inference := Inference{
		Description: "Can vote",
		Rules: []WeightedRule{{
			Weight: 1,
			Rule: Rule{
				Description: "Be of legal age",
				Expression:  "true",
			}},
		},
		FactID:           "can_vote",
		FactValue:        "age.Values[\"value\"] >= 18",
		IsIDCalculated:   false,
		IsValeCalculated: true,
	}

	kb := KnowledgeBase{
		Inferences: []Inference{inference},
		Facts:      map[string]Fact{},
	}

	kb.AddFact(Fact{ID: "age", Values: map[string]interface{}{"value": 15}})
	if _, ok := kb.Facts["can_vote"]; !ok {
		t.Errorf("Expected fact to be inferred")
	}
}
