package inference

import (
	log "github.com/sirupsen/logrus"
	"slices"
)

type KnowledgeBase struct {
	RunningCount   int             `json:"runnings_count"`
	Facts          map[string]Fact `json:"facts"`
	Inferences     []Inference     `json:"inferences"`
	Contradictions []Contradiction `json:"contradictions"`
	Conclusions    []Conclusion    `json:"conclusions"`
}

// Start the knowledge base session
// this is used to keep track of the number of times the knowledge base has
// been used and change the probability of the inferences
// also it clears the facts
func (kb *KnowledgeBase) Start() {
	kb.RunningCount++
	kb.Facts = make(map[string]Fact)
}

// AddFact adds a fact to the knowledge base
func (kb *KnowledgeBase) AddFact(fact Fact) {
	kb.Facts[fact.ID] = fact
	kb.RemoveDerivedFrom(fact.ID)
	kb.Infer()
	kb.ResolveContradictions()
}

// Infer runs all the inferences in the knowledge base
func (kb *KnowledgeBase) Infer() {
	slices.SortFunc(kb.Inferences, func(i, j Inference) int {
		return i.Order - j.Order
	})
	for _, inference := range kb.Inferences {
		if inference.IsNeeded(kb.Facts) {
			id, value, derived, err := inference.infer(kb.Facts)
			if err != nil {
				log.Infof("Error inferring fact: %s", err)
				continue
			} else {
				inference.CountOfTrue++
				inference.Probability = (inference.Probability + float64(inference.CountOfTrue)/float64(kb.RunningCount)) / 2
			}
			var accumulative bool
			f, ok := kb.Facts[id]
			if ok {
				accumulative = f.Accumulative
			}
			kb.Facts[id] = Fact{ID: id, Value: value, DerivedFrom: derived, Accumulative: accumulative}
		}
	}
}

// GetPendingInference returns all the inferences that are needed
func (kb *KnowledgeBase) GetPendingInference() []Inference {
	var pending []Inference
	for _, inference := range kb.Inferences {
		if inference.IsNeeded(kb.Facts) {
			pending = append(pending, inference)
		}
	}
	slices.SortFunc(pending, func(i, j Inference) int {
		return int(i.Probability*100 - j.Probability*100)
	})
	return pending
}

// RemoveDerivedFrom removes all the facts that are derived from a given fact
func (kb *KnowledgeBase) RemoveDerivedFrom(id string) {
	for key, fact := range kb.Facts {
		for _, derived := range fact.DerivedFrom {
			if derived == id && !fact.Accumulative {
				delete(kb.Facts, key)
			}
		}
	}
}

// ResolveContradictions resolves all the contradictions in the knowledge base
func (kb *KnowledgeBase) ResolveContradictions() {
	for _, contradiction := range kb.Contradictions {
		if contradiction.Detect(kb.Facts) {
			contradiction.Resolve(kb)
		}
	}
}

// GetTrueConclusions returns all the conclusions that are true
func (kb *KnowledgeBase) GetTrueConclusions() []Conclusion {
	var conclusions []Conclusion
	for _, conclusion := range kb.Conclusions {
		if conclusion.Assert(kb.Facts) {
			conclusions = append(conclusions, conclusion)
		}
	}
	return conclusions
}

// GetFalseConclusions returns all the conclusions that are false
func (kb *KnowledgeBase) GetFalseConclusions() []Conclusion {
	var conclusions []Conclusion
	for _, conclusion := range kb.Conclusions {
		if !conclusion.Assert(kb.Facts) {
			conclusions = append(conclusions, conclusion)
		}
	}
	return conclusions
}

// CertaintyForConclusion returns the certainty of a conclusion
func (kb *KnowledgeBase) CertaintyForConclusion(conclusion Conclusion) float64 {
	return conclusion.Certainty(kb.Facts)
}
