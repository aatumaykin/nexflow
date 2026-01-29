package dto

// SkillDTO represents a skill data transfer object
type SkillDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`        // Unique skill name
	Version     string `json:"version"`     // Skill version
	Location    string `json:"location"`    // Path to skill directory
	Permissions string `json:"permissions"` // JSON array of required permissions
	Metadata    string `json:"metadata"`    // JSON metadata (timeout, etc.)
	CreatedAt   string `json:"created_at"`  // ISO 8601 format
}

// CreateSkillRequest represents a request to create a skill
type CreateSkillRequest struct {
	Name        string                 `json:"name" yaml:"name"`
	Version     string                 `json:"version" yaml:"version"`
	Location    string                 `json:"location" yaml:"location"`
	Permissions []string               `json:"permissions" yaml:"permissions"`
	Metadata    map[string]interface{} `json:"metadata" yaml:"metadata"`
}

// UpdateSkillRequest represents a request to update a skill
type UpdateSkillRequest struct {
	Version     string                 `json:"version,omitempty" yaml:"version,omitempty"`
	Location    string                 `json:"location,omitempty" yaml:"location,omitempty"`
	Permissions []string               `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// SkillResponse represents a skill response
type SkillResponse struct {
	Success bool      `json:"success"`
	Skill   *SkillDTO `json:"skill,omitempty"`
	Error   string    `json:"error,omitempty"`
}

// SkillsResponse represents a list of skills response
type SkillsResponse struct {
	Success bool        `json:"success"`
	Skills  []*SkillDTO `json:"skills,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SkillExecutionRequest represents a request to execute a skill
type SkillExecutionRequest struct {
	Skill string                 `json:"skill" yaml:"skill"`
	Input map[string]interface{} `json:"input" yaml:"input"`
}

// SkillExecutionResponse represents a skill execution response
type SkillExecutionResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output,omitempty"`
	Error   string `json:"error,omitempty"`
}
