package inference

import "testing"

func TestRiskAnalyzer_Analyze(t *testing.T) {
	ra := RiskAnalyzer{
		Risks: []Risk{
			{
				Description: "High temperature risk",
				Level:       RiskHigh,
				Expression:  "temperature > 40",
				Mitigation:  "Immediate cooling",
			},
		},
	}
	kb := &KnowledgeBase{
		Facts: map[string]Fact{
			"temperature": {ID: "temperature", Value: 41},
		},
	}
	risks, err := ra.Analyze(kb)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(risks) != 1 {
		t.Fatalf("Expected 1 risk, got %d", len(risks))
	}
	if risks[0].Level != RiskHigh {
		t.Errorf("Expected high risk, got %s", risks[0].Level)
	}
}

func TestRiskAnalyzer_NoRisks(t *testing.T) {
	ra := RiskAnalyzer{
		Risks: []Risk{
			{
				Description: "High temperature risk",
				Level:       RiskHigh,
				Expression:  "temperature > 40",
				Mitigation:  "Immediate cooling",
			},
		},
	}
	kb := &KnowledgeBase{
		Facts: map[string]Fact{
			"temperature": {ID: "temperature", Value: 37},
		},
	}
	risks, err := ra.Analyze(kb)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(risks) != 0 {
		t.Errorf("Expected 0 risks, got %d", len(risks))
	}
}

func TestRiskAnalyzer_ContradictionRisk(t *testing.T) {
	ra := RiskAnalyzer{}
	kb := &KnowledgeBase{
		Facts: map[string]Fact{
			"status_ok":    {ID: "status_ok", Value: true},
			"status_error": {ID: "status_error", Value: true},
		},
		Contradictions: []Contradiction{
			{
				Description: "status conflict",
				Facts: []Fact{
					{ID: "status_ok", Value: true},
					{ID: "status_error", Value: true},
				},
			},
		},
	}
	risks, err := ra.Analyze(kb)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(risks) != 1 {
		t.Fatalf("Expected 1 contradiction risk, got %d", len(risks))
	}
	if risks[0].Level != RiskHigh {
		t.Errorf("Expected high risk from contradiction, got %s", risks[0].Level)
	}
}

func TestRiskAnalyzer_LowCertaintyRisk(t *testing.T) {
	ra := RiskAnalyzer{}
	kb := &KnowledgeBase{
		Facts: map[string]Fact{
			"symptom1": {ID: "symptom1", Value: true},
		},
		Conclusions: []Conclusion{
			{
				Description: "diagnosis",
				Facts: []Fact{
					{ID: "symptom1", Value: true},
					{ID: "symptom2", Value: true},
					{ID: "symptom3", Value: true},
				},
			},
		},
	}
	risks, err := ra.Analyze(kb)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(risks) != 1 {
		t.Fatalf("Expected 1 low certainty risk, got %d", len(risks))
	}
	if risks[0].Level != RiskMedium {
		t.Errorf("Expected medium risk, got %s", risks[0].Level)
	}
}
