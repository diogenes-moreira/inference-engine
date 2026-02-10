package inference

import "testing"

func TestSolutionScore_Composite(t *testing.T) {
	score := SolutionScore{
		BusinessImpact:           0.8,
		ImplementationComplexity: 0.2,
		RiskLevel:                0.3,
		TimeToValue:              0.9,
	}
	weights := SolutionScore{
		BusinessImpact:           0.4,
		ImplementationComplexity: 0.2,
		RiskLevel:                0.2,
		TimeToValue:              0.2,
	}
	composite := score.Composite(weights)
	// 0.8*0.4 + (1-0.2)*0.2 + (1-0.3)*0.2 + 0.9*0.2 = 0.32 + 0.16 + 0.14 + 0.18 = 0.80
	if composite < 0.79 || composite > 0.81 {
		t.Errorf("Expected ~0.80, got %f", composite)
	}
}

func TestRankSolutions(t *testing.T) {
	solutions := []RankedSolution{
		{
			Conclusion: Conclusion{Description: "low"},
			Score:      SolutionScore{BusinessImpact: 0.3, ImplementationComplexity: 0.8, RiskLevel: 0.9, TimeToValue: 0.2},
		},
		{
			Conclusion: Conclusion{Description: "high"},
			Score:      SolutionScore{BusinessImpact: 0.9, ImplementationComplexity: 0.1, RiskLevel: 0.1, TimeToValue: 0.9},
		},
	}
	weights := SolutionScore{BusinessImpact: 0.25, ImplementationComplexity: 0.25, RiskLevel: 0.25, TimeToValue: 0.25}
	ranked := RankSolutions(solutions, weights)
	if ranked[0].Conclusion.Description != "high" {
		t.Errorf("Expected 'high' first, got '%s'", ranked[0].Conclusion.Description)
	}
	if ranked[0].CompositeScore <= ranked[1].CompositeScore {
		t.Error("Expected first solution to have higher composite score")
	}
}
