package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	FilePath     string `json:"filePath"`
	MaxSizeMB    int    `json:"maxSizeMB"`
	OutputFolder string `json:"outputFolder"`
}

func ReadConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
