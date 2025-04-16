package config

import (
	"encoding/json"
	"os"
	"regexp"
)

type FilterMode string

const (
	FilterModeNone    FilterMode = "none"
	FilterModeInclude FilterMode = "include"
	FilterModeExclude FilterMode = "exclude"
)

type Config struct {
	FilePath     string    `json:"filePath"`
	MaxSizeMB    int       `json:"maxSizeMB"`
	OutputFolder string    `json:"outputFolder"`
	Analysis     Analysis  `json:"analysis"`
	Filtering    Filtering `json:"filtering"`
}

type Analysis struct {
	Enabled        bool     `json:"enabled"`
	ErrorPatterns  []string `json:"errorPatterns"`
	CustomPatterns []string `json:"customPatterns"`
	OutputFile     string   `json:"outputFile"`
}

type Filtering struct {
	Mode     FilterMode `json:"mode"`
	Patterns []string   `json:"patterns"`
}

type RuntimeConfig struct {
	Config
	ErrorPatternRegexps  []*regexp.Regexp
	CustomPatternRegexps []*regexp.Regexp
	FilterPatternRegexps []*regexp.Regexp
}

func ReadConfig(filePath string) (RuntimeConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return RuntimeConfig{}, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return RuntimeConfig{}, err
	}

	// Prepare runtime config with compiled regexps
	rtCfg := RuntimeConfig{Config: cfg}

	// Compile error patterns
	for _, pattern := range cfg.Analysis.ErrorPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue // Skip invalid patterns
		}
		rtCfg.ErrorPatternRegexps = append(rtCfg.ErrorPatternRegexps, re)
	}

	// Compile custom patterns
	for _, pattern := range cfg.Analysis.CustomPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue // Skip invalid patterns
		}
		rtCfg.CustomPatternRegexps = append(rtCfg.CustomPatternRegexps, re)
	}

	// Compile filter patterns
	for _, pattern := range cfg.Filtering.Patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue // Skip invalid patterns
		}
		rtCfg.FilterPatternRegexps = append(rtCfg.FilterPatternRegexps, re)
	}

	return rtCfg, nil
}
