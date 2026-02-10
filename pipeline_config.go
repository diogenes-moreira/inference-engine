package inference

import (
	"encoding/json"
	"os"
)

// LoadPipelineConfig loads a PipelineConfig from a JSON file.
func LoadPipelineConfig(filename string) (*PipelineConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config PipelineConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SavePipelineConfig writes a PipelineConfig to a JSON file.
func SavePipelineConfig(config *PipelineConfig, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
