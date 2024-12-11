package inference

import (
	log "github.com/sirupsen/logrus"
)

type KnowledgeBase struct {
	Facts      map[string]Fact
	Inferences []Inference
}

func (kb *KnowledgeBase) AddFact(fact Fact) {
	kb.Facts[fact.ID] = fact
	for _, inference := range kb.Inferences {
		if inference.IsNeeded(kb.Facts) {
			id, value, err := inference.infer(kb.Facts)
			if err != nil {
				log.Errorf("Error inferring fact: %s", err)
				continue
			}
			kb.Facts[id] = Fact{ID: "Inference." + id, Values: map[string]interface{}{id: value}}
		}
	}
}

func (kb *KnowledgeBase) Infer() int {
	inferred := 0
	for _, inference := range kb.Inferences {
		if inference.IsNeeded(kb.Facts) {
			id, value, err := inference.infer(kb.Facts)
			if err != nil {
				log.Errorf("Error inferring fact: %s", err)
				continue
			}
			kb.Facts[id] = Fact{ID: "Inference." + id, Values: map[string]interface{}{id: value}}
			inferred++
		}
	}
	return inferred
}

func (kb *KnowledgeBase) GetPendingInference() []string {
	pending := []string{}
	for _, inference := range kb.Inferences {
		if inference.IsNeeded(kb.Facts) {
			pending = append(pending, inference.FactID)
		}
	}
	return pending
}
