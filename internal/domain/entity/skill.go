package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// Skill represents a registered skill that can be executed by the AI.
// Skills are tools with specific permissions and metadata.
type Skill struct {
	ID          string                 `json:"id"`          // Unique identifier for the skill
	Name        string                 `json:"name"`        // Unique skill name
	Version     string                 `json:"version"`     // Skill version (e.g., "1.0.0")
	Location    string                 `json:"location"`    // Path to skill directory
	Permissions string                 `json:"permissions"` // JSON array of required permissions
	Metadata    string                 `json:"metadata"`    // JSON metadata (timeout, description, etc.)
	CreatedAt   time.Time              `json:"created_at"`  // Timestamp when the skill was registered
	MetadataMap map[string]interface{} `json:"-"`           // Parsed metadata (not persisted)
}

// NewSkill creates a new skill with the specified name, version, location, permissions, and metadata.
func NewSkill(name, version, location string, permissions []string, metadata map[string]interface{}) *Skill {
	return &Skill{
		ID:          utils.GenerateID(),
		Name:        name,
		Version:     version,
		Location:    location,
		Permissions: utils.MarshalJSON(permissions),
		Metadata:    utils.MarshalJSON(metadata),
		CreatedAt:   utils.Now(),
		MetadataMap: metadata,
	}
}

// GetPermissions parses and returns the list of permissions.
// Returns nil if parsing fails or permissions is empty.
func (s *Skill) GetPermissions() []string {
	return utils.UnmarshalJSONToSlice(s.Permissions)
}

// RequiresPermission checks if the skill requires a specific permission.
func (s *Skill) RequiresPermission(permission string) bool {
	permissions := s.GetPermissions()
	if permissions == nil {
		return false
	}
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// RequiresSandbox checks if the skill needs sandbox execution.
// Returns true if the skill has dangerous permissions (shell, filesystem, network).
func (s *Skill) RequiresSandbox() bool {
	permissions := s.GetPermissions()
	if permissions == nil {
		return false
	}
	dangerousPerms := []string{"shell", "filesystem", "network", "system"}
	for _, perm := range permissions {
		for _, dangerous := range dangerousPerms {
			if perm == dangerous {
				return true
			}
		}
	}
	return false
}

// GetTimeout returns the execution timeout in seconds.
// Returns 30 seconds (default) if not specified in metadata.
func (s *Skill) GetTimeout() int {
	metadata := s.GetMetadata()
	if metadata == nil {
		return 30
	}
	if timeout, ok := metadata["timeout"].(float64); ok {
		return int(timeout)
	}
	return 30
}

// HasPermission checks if the skill has a specific permission.
func (s *Skill) HasPermission(perm string) bool {
	return s.RequiresPermission(perm)
}

// GetMetadata parses and returns the metadata as a map.
// Returns nil if parsing fails or metadata is empty.
func (s *Skill) GetMetadata() map[string]interface{} {
	if s.MetadataMap != nil {
		return s.MetadataMap
	}
	s.MetadataMap = utils.UnmarshalJSONToMap(s.Metadata)
	return s.MetadataMap
}
