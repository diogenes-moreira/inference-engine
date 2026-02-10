package inference

import (
	"github.com/expr-lang/expr"
)

// RiskLevel represents the severity of a risk.
type RiskLevel string

const (
	RiskHigh   RiskLevel = "high"
	RiskMedium RiskLevel = "medium"
	RiskLow    RiskLevel = "low"
)

// Risk represents a potential risk identified during analysis.
type Risk struct {
	Description string    `json:"description"`
	Level       RiskLevel `json:"level"`
	Expression  string    `json:"expression"`
	Mitigation  string    `json:"mitigation"`
}

// RiskAnalyzer evaluates risk expressions and checks for contradictions and low-certainty conclusions.
type RiskAnalyzer struct {
	Risks []Risk `json:"risks"`
}

// Analyze evaluates risk expressions against the KB state and returns triggered risks.
func (ra *RiskAnalyzer) Analyze(kb *KnowledgeBase) ([]Risk, error) {
	var triggered []Risk

	env := make(map[string]interface{})
	for k, v := range kb.Facts {
		env[k] = v.Value
	}

	// Evaluate explicit risk rules
	for _, risk := range ra.Risks {
		program, err := expr.Compile(risk.Expression, expr.Env(env))
		if err != nil {
			continue
		}
		output, err := expr.Run(program, env)
		if err != nil {
			continue
		}
		result, ok := output.(bool)
		if !ok || !result {
			continue
		}
		triggered = append(triggered, risk)
	}

	// Check for contradictions as high-risk signals
	for _, contradiction := range kb.Contradictions {
		if contradiction.Detect(kb.Facts) {
			triggered = append(triggered, Risk{
				Description: "Active contradiction: " + contradiction.Description,
				Level:       RiskHigh,
				Mitigation:  "Resolve conflicting facts",
			})
		}
	}

	// Check for low-certainty conclusions as medium-risk signals
	for _, conclusion := range kb.Conclusions {
		certainty := conclusion.Certainty(kb.Facts)
		if certainty > 0 && certainty < 0.5 {
			triggered = append(triggered, Risk{
				Description: "Low certainty conclusion: " + conclusion.Description,
				Level:       RiskMedium,
				Mitigation:  "Gather additional evidence",
			})
		}
	}

	return triggered, nil
}
