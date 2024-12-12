package inference

type Conclusion struct {
	Description string `json:"description"`
	Facts       []Fact `json:"facts"`
}

func (c *Conclusion) Assert(facts map[string]Fact) bool {
	for _, fact := range c.Facts {
		f, ok := facts[fact.ID]
		if !ok {
			return false
		}
		if f.Value != fact.Value {
			return false
		}
	}
	return true
}

// Certainty returns the certainty of the conclusion given the facts even if
// some facts are missing.
func (c *Conclusion) Certainty(facts map[string]Fact) float64 {
	if len(c.Facts) == 0 {
		return 0
	}
	trueCount := 0
	for _, fact := range c.Facts {
		f, ok := facts[fact.ID]
		if !ok {
			continue
		}
		if f.Value != fact.Value {
			return 0
		}
		trueCount++
	}
	return float64(trueCount) / float64(len(c.Facts))
}
