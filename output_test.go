package inference

import (
	"encoding/json"
	"testing"
)

func TestPipelineResult_JSON(t *testing.T) {
	result := PipelineResult{
		Result: "Patient requires immediate attention",
		Reasoning: Reasoning{
			Signals:     []string{"High temperature detected", "Elevated heart rate"},
			Assumptions: []string{"Vitals are recent"},
			Tradeoffs:   []string{"Speed vs thoroughness"},
		},
		Confidence: ConfidenceHigh,
		FollowUp: FollowUp{
			MissingData: []string{"blood_pressure"},
			NextActions: []string{"Take blood pressure reading"},
		},
		Domain: DomainGeneral,
		Intent: Intent{Type: IntentDecision, Description: "triage"},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded PipelineResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Result != result.Result {
		t.Errorf("Expected result %q, got %q", result.Result, decoded.Result)
	}
	if decoded.Confidence != ConfidenceHigh {
		t.Errorf("Expected high confidence, got %s", decoded.Confidence)
	}
	if len(decoded.Reasoning.Signals) != 2 {
		t.Errorf("Expected 2 signals, got %d", len(decoded.Reasoning.Signals))
	}
}
