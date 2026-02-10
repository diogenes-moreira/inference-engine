package inference

import "testing"

func TestDomainDetector_Detect(t *testing.T) {
	detector := DomainDetector{
		Signals: map[Domain][]string{
			DomainFinance:   {"revenue", "profit", "cost"},
			DomainEcommerce: {"cart", "product", "price"},
		},
	}
	facts := map[string]Fact{
		"revenue":     {ID: "revenue", Value: 1000},
		"profit":      {ID: "profit", Value: 200},
		"total_cost":  {ID: "total_cost", Value: 800},
	}
	domain := detector.Detect(facts)
	if domain != DomainFinance {
		t.Errorf("Expected finance, got %s", domain)
	}
}

func TestDomainDetector_DetectFromValues(t *testing.T) {
	detector := DomainDetector{
		Signals: map[Domain][]string{
			DomainEcommerce: {"cart", "product"},
		},
	}
	facts := map[string]Fact{
		"item": {ID: "item", Value: "product in cart"},
	}
	domain := detector.Detect(facts)
	if domain != DomainEcommerce {
		t.Errorf("Expected ecommerce, got %s", domain)
	}
}

func TestDomainDetector_DetectGeneral(t *testing.T) {
	detector := DomainDetector{
		Signals: map[Domain][]string{
			DomainFinance: {"revenue"},
		},
	}
	facts := map[string]Fact{
		"temperature": {ID: "temperature", Value: 37.5},
	}
	domain := detector.Detect(facts)
	if domain != DomainGeneral {
		t.Errorf("Expected general, got %s", domain)
	}
}

func TestDomainDetector_EmptySignals(t *testing.T) {
	detector := DomainDetector{}
	facts := map[string]Fact{"x": {ID: "x", Value: 1}}
	domain := detector.Detect(facts)
	if domain != DomainGeneral {
		t.Errorf("Expected general, got %s", domain)
	}
}
