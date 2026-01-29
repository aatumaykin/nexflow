package config

import (
	"fmt"
)

// SkillsConfig represents skills configuration
type SkillsConfig struct {
	Directory      string `json:"directory" yaml:"directory"`
	TimeoutSec     int    `json:"timeout_sec" yaml:"timeout_sec"`
	SandboxEnabled bool   `json:"sandbox_enabled" yaml:"sandbox_enabled"`
}

// Validate validates the skills configuration
func (s *SkillsConfig) Validate() error {
	if s.Directory == "" {
		return fmt.Errorf("skills.directory is required")
	}
	if s.TimeoutSec <= 0 {
		return fmt.Errorf("skills.timeout_sec must be positive")
	}
	return nil
}
