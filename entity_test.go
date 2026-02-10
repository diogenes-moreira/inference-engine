package inference

import "testing"

func TestEntityExtractor_Extract(t *testing.T) {
	extractor := EntityExtractor{
		Rules: []ExtractionRule{
			{FactID: "temperature_high", Expression: "temperature > 38", Source: "vitals", ConfidenceValue: 0.9},
			{FactID: "heart_rate_high", Expression: "heart_rate > 100", Source: "vitals", ConfidenceValue: 0.85},
		},
	}
	facts := map[string]Fact{
		"temperature": {ID: "temperature", Value: 39.5},
		"heart_rate":  {ID: "heart_rate", Value: 110},
	}
	entities, err := extractor.Extract(facts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(entities) != 2 {
		t.Fatalf("Expected 2 entities, got %d", len(entities))
	}
	if entities[0].FactID != "temperature_high" {
		t.Errorf("Expected temperature_high, got %s", entities[0].FactID)
	}
}

func TestEntityExtractor_NoMatch(t *testing.T) {
	extractor := EntityExtractor{
		Rules: []ExtractionRule{
			{FactID: "fever", Expression: "temperature > 38", Source: "vitals", ConfidenceValue: 0.9},
		},
	}
	facts := map[string]Fact{
		"temperature": {ID: "temperature", Value: 36.5},
	}
	entities, err := extractor.Extract(facts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(entities) != 0 {
		t.Errorf("Expected 0 entities, got %d", len(entities))
	}
}

func TestEntityExtractor_Empty(t *testing.T) {
	extractor := EntityExtractor{}
	facts := map[string]Fact{"x": {ID: "x", Value: 1}}
	entities, err := extractor.Extract(facts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if entities != nil {
		t.Errorf("Expected nil entities, got %v", entities)
	}
}
