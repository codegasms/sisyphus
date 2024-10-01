package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Servers  []ServerAddr `json:"servers"`
	Weights  []float32    `json:"weights,omitempty"`
	Strategy StrategyKind `json:"strategy"`
}

func LoadConfig(configPath string) (*Config, error) {
	rawJson, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(rawJson, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
