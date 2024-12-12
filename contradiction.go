package inference

type Contradiction struct {
	Description string `json:"description"`
	Facts       []Fact `json:"facts"`
}

func (c *Contradiction) Detect(facts map[string]Fact) bool {
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

func (c *Contradiction) Resolve(base *KnowledgeBase) {
	for _, fact := range c.Facts {
		delete(base.Facts, fact.ID)
		base.RemoveDerivedFrom(fact.ID)
	}
}
