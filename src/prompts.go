package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type promptsConfig struct {
	Default   string            `json:"default"`
	Overrides map[string]string `json:"overrides,omitempty"`
	Annual    map[string]string `json:"annual,omitempty"`
}

func loadPrompt(location *time.Location) string {
	data, err := os.ReadFile("prompts.json")
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf(warnReadPrompt+"\n", err)
		}
		return ""
	}

	var cfg promptsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Printf(warnParsePrompt+"\n", err)
		return ""
	}

	today := time.Now().In(location)
	dateKey := today.Format("2006-01-02")
	annualKey := today.Format("01-02")

	if cfg.Overrides != nil {
		if p, ok := cfg.Overrides[dateKey]; ok {
			fmt.Printf(infoPromptOverride+"\n", dateKey)
			return p
		}
	}
	if cfg.Annual != nil {
		if p, ok := cfg.Annual[annualKey]; ok {
			fmt.Printf(infoPromptAnnual+"\n", annualKey)
			return p
		}
	}
	if cfg.Default != "" {
		fmt.Println(infoPromptDefault)
		return cfg.Default
	}
	return ""
}
