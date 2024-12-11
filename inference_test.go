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
		Facts:      map[string]Fact{"age": {ID: "age", Value: 20}},
	}

	if !inference.IsNeeded(kb.Facts) {
		t.Errorf("Expected inference to be needed")
	}

	kb.Dump("test.json")
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
		Facts:      map[string]Fact{"can_vote": {ID: "can_vote", Value: false}},
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
		FactValue:        "age.Value >= 18 && age.Value < 50",
		IsIDCalculated:   false,
		IsValeCalculated: true,
	}

	kb := KnowledgeBase{
		Inferences: []Inference{inference},
		Facts:      map[string]Fact{},
	}

	kb.Start()
	kb.AddFact(Fact{ID: "age", Value: 18})
	if value, ok := kb.Facts["can_vote"]; !ok || value.Value != true {
		t.Errorf("Expected fact to be inferred")
	}
	kb.AddFact(Fact{ID: "age", Value: 15})
	if value, ok := kb.Facts["can_vote"]; !ok || value.Value != false {
		t.Errorf("Expected fact to be inferred")
	}
}
