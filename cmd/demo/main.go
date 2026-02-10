package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	inference "github.com/diogenes-moreira/inference-engine"
)

//go:embed static
var staticFiles embed.FS

//go:embed examples
var exampleFiles embed.FS

var (
	mu             sync.Mutex
	currentConfig  *inference.PipelineConfig
	currentPipeline *inference.Pipeline
)

type ExampleInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(staticFS)))
	http.HandleFunc("/api/examples", handleExamples)
	http.HandleFunc("/api/pipeline/run", handlePipelineRun)
	http.HandleFunc("/api/pipeline/pending", handlePipelinePending)
	http.HandleFunc("/api/pipeline/reset", handlePipelineReset)
	http.HandleFunc("/api/pipeline/load", handlePipelineLoad)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	fmt.Printf("Demo server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleExamples(w http.ResponseWriter, r *http.Request) {
	examples := []ExampleInfo{
		{ID: "triage", Name: "Hospital Triage", Description: "Emergency room patient triage based on vital signs"},
		{ID: "gamification", Name: "Pizza Gamification", Description: "Pizza loyalty rewards based on accumulated sales"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(examples)
}

func handlePipelineLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ExampleID string `json:"example_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var configPath string
	switch req.ExampleID {
	case "triage":
		configPath = "examples/triage/definition.json"
	case "gamification":
		configPath = "examples/gamification/pipeline.json"
	default:
		http.Error(w, "Unknown example: "+req.ExampleID, http.StatusBadRequest)
		return
	}

	data, err := exampleFiles.ReadFile(configPath)
	if err != nil {
		http.Error(w, "Failed to read config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var config inference.PipelineConfig
	if err := json.Unmarshal(data, &config); err != nil {
		http.Error(w, "Failed to parse config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	mu.Lock()
	currentConfig = &config
	currentConfig.KnowledgeBase.Start()
	currentPipeline = inference.NewPipeline(*currentConfig)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "loaded", "example": req.ExampleID})
}

func handlePipelineRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var inputFacts map[string]inference.Fact
	if err := json.NewDecoder(r.Body).Decode(&inputFacts); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if currentPipeline == nil {
		http.Error(w, "No pipeline loaded. Call /api/pipeline/load first.", http.StatusBadRequest)
		return
	}

	result, err := currentPipeline.Run(inputFacts)
	if err != nil {
		http.Error(w, "Pipeline error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handlePipelinePending(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if currentConfig == nil || currentConfig.KnowledgeBase == nil {
		http.Error(w, "No pipeline loaded", http.StatusBadRequest)
		return
	}

	pending := currentConfig.KnowledgeBase.GetPendingInference()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pending)
}

func handlePipelineReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if currentConfig == nil || currentConfig.KnowledgeBase == nil {
		http.Error(w, "No pipeline loaded", http.StatusBadRequest)
		return
	}

	currentConfig.KnowledgeBase.Start()
	currentPipeline = inference.NewPipeline(*currentConfig)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "reset"})
}
