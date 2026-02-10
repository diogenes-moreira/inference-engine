package inference

// Reasoning captures the explanation chain for a pipeline result.
type Reasoning struct {
	Signals     []string `json:"signals"`
	Assumptions []string `json:"assumptions"`
	Tradeoffs   []string `json:"tradeoffs"`
}

// FollowUp captures what's still needed after pipeline execution.
type FollowUp struct {
	MissingData []string `json:"missing_data"`
	NextActions []string `json:"next_actions"`
}

// PipelineResult is the structured output of the 6-step pipeline.
type PipelineResult struct {
	Result     string          `json:"result"`
	Reasoning  Reasoning       `json:"reasoning"`
	Confidence ConfidenceLevel `json:"confidence"`
	FollowUp   FollowUp       `json:"follow_up"`
	Domain     Domain          `json:"domain"`
	Intent     Intent          `json:"intent"`
	Entities   []Entity        `json:"entities"`
	Constraints []Constraint   `json:"constraints"`
	Risks      []Risk          `json:"risks"`
	Solutions  []RankedSolution `json:"solutions,omitempty"`
}
