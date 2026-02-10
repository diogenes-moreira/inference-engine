package inference

import "sort"

// SolutionScore holds scoring dimensions for a solution (all 0-1).
type SolutionScore struct {
	BusinessImpact           float64 `json:"business_impact"`
	ImplementationComplexity float64 `json:"implementation_complexity"`
	RiskLevel                float64 `json:"risk_level"`
	TimeToValue              float64 `json:"time_to_value"`
}

// Composite returns the weighted composite score.
func (s SolutionScore) Composite(weights SolutionScore) float64 {
	return s.BusinessImpact*weights.BusinessImpact +
		(1-s.ImplementationComplexity)*weights.ImplementationComplexity +
		(1-s.RiskLevel)*weights.RiskLevel +
		s.TimeToValue*weights.TimeToValue
}

// RankedSolution pairs a Conclusion with its scoring.
type RankedSolution struct {
	Conclusion Conclusion    `json:"conclusion"`
	Score      SolutionScore `json:"score"`
	CompositeScore float64   `json:"composite_score"`
}

// RankSolutions sorts solutions by weighted composite score (descending).
func RankSolutions(solutions []RankedSolution, weights SolutionScore) []RankedSolution {
	for i := range solutions {
		solutions[i].CompositeScore = solutions[i].Score.Composite(weights)
	}
	sort.Slice(solutions, func(i, j int) bool {
		return solutions[i].CompositeScore > solutions[j].CompositeScore
	})
	return solutions
}
