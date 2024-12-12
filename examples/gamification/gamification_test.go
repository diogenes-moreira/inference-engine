package main

import (
	"encoding/json"
	"github.com/diogenes-moreira/inference-engine"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

type Sale struct {
	Product string `json:"product"`
	Price   int    `json:"price"`
}

func Test_Gamification(t *testing.T) {
	kb, _ := LoadKnowledgeBase("definition.json")
	kb.Start()
	sale := Sale{Product: "pizza", Price: 100}
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	pending := kb.GetPendingInference()
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	log.Infof("Pending inferences: %v", pending)
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	kb.AddFact(inference.Fact{ID: "sale", Value: sale, DerivedFrom: []string{}})
	log.Infof("Pending inferences: %v", pending)
}

func LoadKnowledgeBase(filename string) (*inference.KnowledgeBase, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Errorf("Error opening file: %s", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Error closing file: %s", err)
		}
	}(file)

	var kb inference.KnowledgeBase
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&kb)
	if err != nil {
		log.Errorf("Error decoding JSON: %s", err)
		return nil, err
	}

	return &kb, nil
}
