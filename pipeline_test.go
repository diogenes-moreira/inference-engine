package inference

import "testing"

func TestPipeline_EndToEnd(t *testing.T) {
	kb := &KnowledgeBase{
		Facts: map[string]Fact{},
		Inferences: []Inference{
			{
				Description: "High urgency",
				Rules: []WeightedRule{
					{Rule: Rule{Expression: "temperature > 39", FactTargetID: "temperature"}, Weight: 1},
				},
				FactID:    "urgency",
				FactValue: "red",
				Order:     0,
			},
		},
		Conclusions: []Conclusion{
			{
				Description: "Critical patient",
				Facts:       []Fact{{ID: "urgency", Value: "red"}},
			},
		},
	}
	kb.Start()

	config := PipelineConfig{
		KnowledgeBase: kb,
		DomainDetector: &DomainDetector{
			Signals: map[Domain][]string{
				DomainGeneral: {"temperature", "symptom"},
			},
		},
		IntentClassifier: &IntentClassifier{
			Rules: []IntentRule{
				{Expression: "temperature > 0", IntentType: IntentDecision, Weight: 1.0},
			},
		},
		EntityExtractor: &EntityExtractor{
			Rules: []ExtractionRule{
				{FactID: "fever", Expression: "temperature > 38", Source: "vitals", ConfidenceValue: 0.9},
			},
		},
		RiskAnalyzer: &RiskAnalyzer{
			Risks: []Risk{
				{Description: "Critical temperature", Level: RiskHigh, Expression: "temperature > 40", Mitigation: "Immediate cooling"},
			},
		},
	}

	pipeline := NewPipeline(config)
	inputFacts := map[string]Fact{
		"temperature": {ID: "temperature", Value: 41.0},
	}

	result, err := pipeline.Run(inputFacts)
	if err != nil {
		t.Fatalf("Pipeline failed: %v", err)
	}

	if result.Result != "Critical patient" {
		t.Errorf("Expected 'Critical patient', got %q", result.Result)
	}
	if result.Confidence != ConfidenceHigh {
		t.Errorf("Expected high confidence, got %s", result.Confidence)
	}
	if result.Intent.Type != IntentDecision {
		t.Errorf("Expected decision intent, got %s", result.Intent.Type)
	}
	if result.Domain != DomainGeneral {
		t.Errorf("Expected general domain, got %s", result.Domain)
	}

	// Check risks
	hasHighRisk := false
	for _, r := range result.Risks {
		if r.Level == RiskHigh {
			hasHighRisk = true
			break
		}
	}
	if !hasHighRisk {
		t.Error("Expected high risk for temperature > 40")
	}
}

func TestPipeline_NoKnowledgeBase(t *testing.T) {
	config := PipelineConfig{}
	pipeline := NewPipeline(config)
	_, err := pipeline.Run(map[string]Fact{})
	if err == nil {
		t.Error("Expected error when KB is nil")
	}
}

func TestPipeline_MinimalConfig(t *testing.T) {
	kb := &KnowledgeBase{
		Facts: map[string]Fact{},
	}
	kb.Start()

	config := PipelineConfig{KnowledgeBase: kb}
	pipeline := NewPipeline(config)

	result, err := pipeline.Run(map[string]Fact{
		"x": {ID: "x", Value: 1},
	})
	if err != nil {
		t.Fatalf("Pipeline failed: %v", err)
	}
	if result.Domain != DomainGeneral {
		t.Errorf("Expected general domain, got %s", result.Domain)
	}
	if result.Confidence != ConfidenceLow {
		t.Errorf("Expected low confidence with no conclusions, got %s", result.Confidence)
	}
}
