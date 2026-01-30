package config

import (
	"fmt"
)

// LLMConfig represents LLM provider configuration
type LLMConfig struct {
	DefaultProvider string                 `json:"default_provider" yaml:"default_provider"`
	Providers       map[string]LLMProvider `json:"providers" yaml:"providers"`
}

// LLMProvider represents a single LLM provider configuration
type LLMProvider struct {
	APIKey      string  `json:"api_key" yaml:"api_key"`
	BaseURL     string  `json:"base_url" yaml:"base_url"`
	Model       string  `json:"model" yaml:"model"`
	Temperature float64 `json:"temperature" yaml:"temperature"`
	MaxTokens   int     `json:"max_tokens" yaml:"max_tokens"`
}

// Validate validates the LLM configuration
func (l *LLMConfig) Validate() error {
	if l.DefaultProvider == "" {
		return fmt.Errorf("llm.default_provider is required")
	}
	if len(l.Providers) == 0 {
		return fmt.Errorf("at least one llm provider is required")
	}
	if _, ok := l.Providers[l.DefaultProvider]; !ok {
		return fmt.Errorf("llm.default_provider '%s' not found in providers", l.DefaultProvider)
	}
	return nil
}
