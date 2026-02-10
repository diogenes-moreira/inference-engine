package inference

import "testing"

func TestComputeConfidence_High(t *testing.T) {
	if ComputeConfidence(0.8) != ConfidenceHigh {
		t.Error("Expected high confidence for 0.8")
	}
	if ComputeConfidence(1.0) != ConfidenceHigh {
		t.Error("Expected high confidence for 1.0")
	}
}

func TestComputeConfidence_Medium(t *testing.T) {
	if ComputeConfidence(0.5) != ConfidenceMedium {
		t.Error("Expected medium confidence for 0.5")
	}
	if ComputeConfidence(0.79) != ConfidenceMedium {
		t.Error("Expected medium confidence for 0.79")
	}
}

func TestComputeConfidence_Low(t *testing.T) {
	if ComputeConfidence(0.0) != ConfidenceLow {
		t.Error("Expected low confidence for 0.0")
	}
	if ComputeConfidence(0.49) != ConfidenceLow {
		t.Error("Expected low confidence for 0.49")
	}
}
