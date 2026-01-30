package skills

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Config represents local skill runtime configuration
type Config struct {
	Directory      string
	TimeoutSeconds int
	SandboxEnabled bool
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Directory == "" {
		return fmt.Errorf("local runtime: directory is required")
	}
	if c.TimeoutSeconds <= 0 {
		return fmt.Errorf("local runtime: timeout_seconds must be positive")
	}
	return nil
}

// LocalRuntime is a local skill runtime that executes skills on the host system
type LocalRuntime struct {
	config *Config
	logger *slog.Logger
}

// NewLocalRuntime creates a new local skill runtime
func NewLocalRuntime(config *Config, logger *slog.Logger) (*LocalRuntime, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Ensure directory exists
	if err := os.MkdirAll(config.Directory, 0755); err != nil {
		return nil, fmt.Errorf("local runtime: failed to create skills directory: %w", err)
	}

	return &LocalRuntime{
		config: config,
		logger: logger.With("runtime", "local"),
	}, nil
}

// Name returns the runtime name
func (r *LocalRuntime) Name() string {
	return "local"
}

// Execute runs a skill with the given input
func (r *LocalRuntime) Execute(ctx context.Context, req *ExecutionRequest) (*ExecutionResult, error) {
	startTime := time.Now()

	r.logger.Debug("Executing skill",
		"skill_name", req.SkillName,
		"input", req.Input)

	// Find skill executable
	skillPath, err := r.findSkillPath(req.SkillName)
	if err != nil {
		return &ExecutionResult{
			Success: false,
			Output:  nil,
			Error:   fmt.Errorf("skill not found: %w", err),
		}, nil
	}

	// Create context with timeout
	timeout := time.Duration(r.config.TimeoutSeconds) * time.Second
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Prepare command
	cmd := exec.CommandContext(execCtx, skillPath)

	// Set environment variables for input
	for key, value := range req.Input {
		// Convert value to string
		var envValue string
		switch v := value.(type) {
		case string:
			envValue = v
		case []byte:
			envValue = string(v)
		default:
			// JSON encode complex values
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				r.logger.Warn("Failed to marshal input value",
					"key", key,
					"value", value,
					"error", err)
				continue
			}
			envValue = string(jsonBytes)
		}
		cmd.Env = append(cmd.Env, fmt.Sprintf("NEXFLOW_%s=%s", strings.ToUpper(key), envValue))
	}

	// Capture stdout and stderr
	output, err := cmd.CombinedOutput()

	executionTime := time.Since(startTime)

	if err != nil {
		r.logger.Error("Skill execution failed",
			"skill_name", req.SkillName,
			"execution_time_ms", executionTime.Milliseconds(),
			"error", err,
			"output", string(output))

		return &ExecutionResult{
			Success: false,
			Output:  nil,
			Error:   fmt.Errorf("execution failed: %w: %s", err, string(output)),
		}, nil
	}

	// Parse output as JSON
	var parsedOutput map[string]interface{}
	if len(output) > 0 {
		if err := json.Unmarshal(output, &parsedOutput); err != nil {
			r.logger.Warn("Failed to parse skill output as JSON, using raw output",
				"skill_name", req.SkillName,
				"error", err)

			// Use raw output if JSON parsing fails
			parsedOutput = map[string]interface{}{
				"output": string(output),
			}
		}
	}

	r.logger.Info("Skill execution completed",
		"skill_name", req.SkillName,
		"execution_time_ms", executionTime.Milliseconds())

	return &ExecutionResult{
		Success: true,
		Output:  parsedOutput,
		Error:   nil,
		Metadata: map[string]interface{}{
			"execution_time_ms": executionTime.Milliseconds(),
			"skill_path":        skillPath,
		},
	}, nil
}

// List returns a list of available skills
func (r *LocalRuntime) List(ctx context.Context) ([]string, error) {
	r.logger.Debug("Listing available skills")

	entries, err := os.ReadDir(r.config.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	var skills []string
	for _, entry := range entries {
		// Skip directories and hidden files
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Check if file is executable
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Check if file has execute permission
		mode := info.Mode()
		if mode&0111 != 0 {
			skills = append(skills, entry.Name())
		}
	}

	r.logger.Debug("Found skills", "count", len(skills))

	return skills, nil
}

// GetMetadata returns metadata for a specific skill
func (r *LocalRuntime) GetMetadata(ctx context.Context, skillName string) (map[string]interface{}, error) {
	skillPath, err := r.findSkillPath(skillName)
	if err != nil {
		return nil, err
	}

	// Get file info
	info, err := os.Stat(skillPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill info: %w", err)
	}

	metadata := map[string]interface{}{
		"name":         skillName,
		"path":         skillPath,
		"size":         info.Size(),
		"modified":     info.ModTime(),
		"runtime":      r.Name(),
		"sandbox_mode": r.config.SandboxEnabled,
	}

	// Try to read SKILL.md for additional metadata
	skillDir := filepath.Dir(skillPath)
	skillMDPath := filepath.Join(skillDir, fmt.Sprintf("%s.SKILL.md", skillName))
	if _, err := os.Stat(skillMDPath); err == nil {
		metadata["has_documentation"] = true
	}

	return metadata, nil
}

// IsAvailable checks if a skill is available for execution
func (r *LocalRuntime) IsAvailable(ctx context.Context, skillName string) bool {
	skillPath, err := r.findSkillPath(skillName)
	if err != nil {
		return false
	}

	// Check if file exists and is executable
	info, err := os.Stat(skillPath)
	if err != nil {
		return false
	}

	mode := info.Mode()
	return mode.IsRegular() && mode&0111 != 0
}

// findSkillPath finds the full path to a skill executable
func (r *LocalRuntime) findSkillPath(skillName string) (string, error) {
	// First, check if skillName is an absolute path
	if filepath.IsAbs(skillName) {
		if _, err := os.Stat(skillName); err != nil {
			return "", fmt.Errorf("skill file not found: %s", skillName)
		}
		return skillName, nil
	}

	// Check in skills directory
	skillPath := filepath.Join(r.config.Directory, skillName)
	if _, err := os.Stat(skillPath); err == nil {
		return skillPath, nil
	}

	// Try with common extensions
	for _, ext := range []string{"", ".sh", ".py", ".js", ".rb", ".go"} {
		pathWithExt := filepath.Join(r.config.Directory, skillName+ext)
		if _, err := os.Stat(pathWithExt); err == nil {
			return pathWithExt, nil
		}
	}

	// Check if skillName is in PATH
	path, err := exec.LookPath(skillName)
	if err == nil {
		return path, nil
	}

	return "", fmt.Errorf("skill '%s' not found in directory: %s", skillName, r.config.Directory)
}
