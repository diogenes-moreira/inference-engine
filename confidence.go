package inference

// ConfidenceLevel represents the confidence category of a result.
type ConfidenceLevel string

const (
	ConfidenceHigh   ConfidenceLevel = "high"
	ConfidenceMedium ConfidenceLevel = "medium"
	ConfidenceLow    ConfidenceLevel = "low"
)

// ComputeConfidence maps a certainty value (0-1) to a ConfidenceLevel.
func ComputeConfidence(certainty float64) ConfidenceLevel {
	if certainty >= 0.8 {
		return ConfidenceHigh
	}
	if certainty >= 0.5 {
		return ConfidenceMedium
	}
	return ConfidenceLow
}
