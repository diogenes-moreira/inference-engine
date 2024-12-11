package inference

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"slices"
)

type KnowledgeBase struct {
	RunningCount int             `json:"runnings_count"`
	Facts        map[string]Fact `json:"facts"`
	Inferences   []Inference     `json:"inferences"`
}

// Start the knowledge base session
// this is used to keep track of the number of times the knowledge base has
// been used and change the probability of the inferences
func (kb *KnowledgeBase) Start() {
	kb.RunningCount++
}

// AddFact adds a fact to the knowledge base
func (kb *KnowledgeBase) AddFact(fact Fact) {
	kb.Facts[fact.ID] = fact
	kb.RemoveDerivedFrom(fact.ID)
	kb.Infer()
}

// Infer runs all the inferences in the knowledge base
func (kb *KnowledgeBase) Infer() {
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
			kb.Facts[id] = Fact{ID: id, Value: value, DerivedFrom: derived}
		}
	}
}

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

func (kb *KnowledgeBase) DumpInferences(filename string) {
	data, err := json.Marshal(kb.Inferences)
	if err != nil {
		log.Errorf("Error dumping rules: %s", err)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Errorf("Error creating file: %s", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Error closing file: %s", err)
		}
	}(file)

	_, err = file.Write(data)
	if err != nil {
		log.Errorf("Error writing to file: %s", err)
		return
	}

}

func (kb *KnowledgeBase) Dump(filename string) {
	data, err := json.Marshal(kb)
	if err != nil {
		log.Errorf("Error dumping knowledge base: %s", err)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Errorf("Error creating file: %s", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Error closing file: %s", err)
		}
	}(file)

	_, err = file.Write(data)
	if err != nil {
		log.Errorf("Error writing to file: %s", err)
		return
	}
}

func (kb *KnowledgeBase) RemoveDerivedFrom(id string) {
	for key, fact := range kb.Facts {
		for _, derived := range fact.DerivedFrom {
			if derived == id {
				delete(kb.Facts, key)
			}
		}
	}
}
