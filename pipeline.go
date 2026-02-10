package inference

import (
	"fmt"
	"strings"
)

// PipelineConfig holds all components needed for the 6-step pipeline.
type PipelineConfig struct {
	IntentClassifier *IntentClassifier `json:"intent_classifier,omitempty"`
	EntityExtractor  *EntityExtractor  `json:"entity_extractor,omitempty"`
	ConstraintSet    *ConstraintSet    `json:"constraint_set,omitempty"`
	KnowledgeBase    *KnowledgeBase    `json:"knowledge_base"`
	RiskAnalyzer     *RiskAnalyzer     `json:"risk_analyzer,omitempty"`
	DomainDetector   *DomainDetector   `json:"domain_detector,omitempty"`
	ScoringWeights   *SolutionScore    `json:"scoring_weights,omitempty"`
	DomainWeights    map[Domain]SolutionScore `json:"domain_weights,omitempty"`
}

// PipelineState tracks intermediate results through pipeline steps.
type PipelineState struct {
	Domain      Domain
	Intent      Intent
	Entities    []Entity
	Constraints []Constraint
	Risks       []Risk
	Signals     []string
	Assumptions []string
	Tradeoffs   []string
}

// Pipeline orchestrates the 6-step deterministic pipeline.
type Pipeline struct {
	Config PipelineConfig
}

// NewPipeline creates a pipeline from a config.
func NewPipeline(config PipelineConfig) *Pipeline {
	return &Pipeline{Config: config}
}

// Run executes the full 6-step pipeline on input facts.
func (p *Pipeline) Run(inputFacts map[string]Fact) (*PipelineResult, error) {
	if p.Config.KnowledgeBase == nil {
		return nil, fmt.Errorf("knowledge base is required")
	}

	state := &PipelineState{}
	kb := p.Config.KnowledgeBase

	// Step 1: Domain detection
	state.Signals = append(state.Signals, "Pipeline started with "+fmt.Sprintf("%d", len(inputFacts))+" input facts")
	if p.Config.DomainDetector != nil {
		state.Domain = p.Config.DomainDetector.Detect(inputFacts)
	} else {
		state.Domain = DomainGeneral
	}
	state.Signals = append(state.Signals, "Detected domain: "+string(state.Domain))

	// Select scoring weights based on domain
	weights := SolutionScore{BusinessImpact: 0.25, ImplementationComplexity: 0.25, RiskLevel: 0.25, TimeToValue: 0.25}
	if p.Config.ScoringWeights != nil {
		weights = *p.Config.ScoringWeights
	} else if p.Config.DomainWeights != nil {
		if dw, ok := p.Config.DomainWeights[state.Domain]; ok {
			weights = dw
		}
	}

	// Step 2: Intent classification
	if p.Config.IntentClassifier != nil {
		intent, err := p.Config.IntentClassifier.Classify(inputFacts)
		if err != nil {
			return nil, fmt.Errorf("intent classification failed: %w", err)
		}
		state.Intent = intent
	} else {
		state.Intent = Intent{Type: IntentQuery, Description: "default"}
	}
	state.Signals = append(state.Signals, "Classified intent: "+string(state.Intent.Type))

	// Step 3: Entity extraction — extract entities and add as facts to KB
	if p.Config.EntityExtractor != nil {
		entities, err := p.Config.EntityExtractor.Extract(inputFacts)
		if err != nil {
			return nil, fmt.Errorf("entity extraction failed: %w", err)
		}
		state.Entities = entities
		for _, entity := range entities {
			kb.AddFact(Fact{
				ID:     entity.FactID,
				Value:  entity.Value,
				Source: "extracted",
			})
		}
		if len(entities) > 0 {
			state.Signals = append(state.Signals, fmt.Sprintf("Extracted %d entities", len(entities)))
		}
	}

	// Add input facts to KB
	for _, fact := range inputFacts {
		if fact.Source == "" {
			fact.Source = "input"
		}
		kb.AddFact(fact)
	}

	// Step 4: Constraint identification
	if p.Config.ConstraintSet != nil {
		constraints, err := p.Config.ConstraintSet.Identify(kb.Facts)
		if err != nil {
			return nil, fmt.Errorf("constraint identification failed: %w", err)
		}
		state.Constraints = constraints

		// Check constraint satisfaction
		for _, c := range constraints {
			satisfied, err := c.Satisfied(kb.Facts)
			if err != nil {
				continue
			}
			if !satisfied && c.Type == ConstraintHard {
				state.Tradeoffs = append(state.Tradeoffs, "Hard constraint not met: "+c.Description)
			} else if !satisfied && c.Type == ConstraintSoft {
				state.Assumptions = append(state.Assumptions, "Soft constraint relaxed: "+c.Description)
			}
		}
	}

	// Step 5: Knowledge application — inference + contradiction resolution
	kb.Infer()
	kb.ResolveContradictions()
	state.Signals = append(state.Signals, fmt.Sprintf("Knowledge base has %d facts after inference", len(kb.Facts)))

	// Step 6: Risk analysis
	if p.Config.RiskAnalyzer != nil {
		risks, err := p.Config.RiskAnalyzer.Analyze(kb)
		if err != nil {
			return nil, fmt.Errorf("risk analysis failed: %w", err)
		}
		state.Risks = risks
		for _, r := range risks {
			if r.Level == RiskHigh {
				state.Tradeoffs = append(state.Tradeoffs, "High risk: "+r.Description)
			}
		}
	}

	// Build structured output
	return p.buildResult(state, kb, weights), nil
}

func (p *Pipeline) buildResult(state *PipelineState, kb *KnowledgeBase, weights SolutionScore) *PipelineResult {
	// Determine result from conclusions
	var resultParts []string
	trueConclusions := kb.GetTrueConclusions()
	for _, c := range trueConclusions {
		resultParts = append(resultParts, c.Description)
	}

	// Rank solutions from conclusions
	var solutions []RankedSolution
	for _, c := range kb.Conclusions {
		certainty := kb.CertaintyForConclusion(c)
		if certainty > 0 {
			riskScore := 0.5
			for _, r := range state.Risks {
				if r.Level == RiskHigh {
					riskScore = 0.9
					break
				} else if r.Level == RiskMedium {
					riskScore = 0.7
				}
			}
			solutions = append(solutions, RankedSolution{
				Conclusion: c,
				Score: SolutionScore{
					BusinessImpact:           certainty,
					ImplementationComplexity: 0.5,
					RiskLevel:                riskScore,
					TimeToValue:              0.5,
				},
			})
		}
	}
	if len(solutions) > 0 {
		solutions = RankSolutions(solutions, weights)
	}

	result := "No conclusions reached"
	if len(resultParts) > 0 {
		result = strings.Join(resultParts, "; ")
	}

	// Compute overall confidence from true conclusions
	maxCertainty := 0.0
	for _, c := range trueConclusions {
		cert := kb.CertaintyForConclusion(c)
		if cert > maxCertainty {
			maxCertainty = cert
		}
	}
	if maxCertainty == 0 && len(kb.Facts) > 0 {
		maxCertainty = 0.3
	}

	// Missing data
	missingData := kb.GetMissingFactIDs()
	pending := kb.GetPendingInference()
	var nextActions []string
	for _, inf := range pending {
		rules := inf.PendingRules(kb.Facts)
		for _, r := range rules {
			if r.Question != "" {
				nextActions = append(nextActions, r.Question)
			}
		}
	}

	return &PipelineResult{
		Result: result,
		Reasoning: Reasoning{
			Signals:     state.Signals,
			Assumptions: state.Assumptions,
			Tradeoffs:   state.Tradeoffs,
		},
		Confidence:  ComputeConfidence(maxCertainty),
		FollowUp: FollowUp{
			MissingData: missingData,
			NextActions: nextActions,
		},
		Domain:      state.Domain,
		Intent:      state.Intent,
		Entities:    state.Entities,
		Constraints: state.Constraints,
		Risks:       state.Risks,
		Solutions:   solutions,
	}
}
