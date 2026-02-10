package inference

import "strings"

// Domain represents a business domain category.
type Domain string

const (
	DomainFinance        Domain = "finance"
	DomainEcommerce      Domain = "ecommerce"
	DomainInfrastructure Domain = "infrastructure"
	DomainData           Domain = "data"
	DomainAIML           Domain = "aiml"
	DomainGeneral        Domain = "general"
)

// DomainDetector detects the domain from a set of facts using keyword signals.
type DomainDetector struct {
	Signals map[Domain][]string `json:"signals"`
}

// Detect returns the domain with the most keyword matches in fact IDs and string values.
func (d *DomainDetector) Detect(facts map[string]Fact) Domain {
	if len(d.Signals) == 0 {
		return DomainGeneral
	}
	scores := make(map[Domain]int)
	for domain, keywords := range d.Signals {
		for _, keyword := range keywords {
			kw := strings.ToLower(keyword)
			for id, fact := range facts {
				if strings.Contains(strings.ToLower(id), kw) {
					scores[domain]++
				}
				if sv, ok := fact.Value.(string); ok {
					if strings.Contains(strings.ToLower(sv), kw) {
						scores[domain]++
					}
				}
			}
		}
	}
	best := DomainGeneral
	bestScore := 0
	for domain, score := range scores {
		if score > bestScore {
			best = domain
			bestScore = score
		}
	}
	return best
}

// DefaultDomainWeights returns domain-specific scoring weight presets.
func DefaultDomainWeights() map[Domain]SolutionScore {
	return map[Domain]SolutionScore{
		DomainFinance:        {BusinessImpact: 0.4, ImplementationComplexity: 0.2, RiskLevel: 0.3, TimeToValue: 0.1},
		DomainEcommerce:      {BusinessImpact: 0.3, ImplementationComplexity: 0.2, RiskLevel: 0.2, TimeToValue: 0.3},
		DomainInfrastructure: {BusinessImpact: 0.2, ImplementationComplexity: 0.3, RiskLevel: 0.3, TimeToValue: 0.2},
		DomainData:           {BusinessImpact: 0.3, ImplementationComplexity: 0.3, RiskLevel: 0.2, TimeToValue: 0.2},
		DomainAIML:           {BusinessImpact: 0.3, ImplementationComplexity: 0.3, RiskLevel: 0.2, TimeToValue: 0.2},
		DomainGeneral:        {BusinessImpact: 0.25, ImplementationComplexity: 0.25, RiskLevel: 0.25, TimeToValue: 0.25},
	}
}
