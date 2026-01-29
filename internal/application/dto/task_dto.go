package dto

// TaskDTO represents a task data transfer object
type TaskDTO struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
	Skill     string `json:"skill"`      // Name of the skill to execute
	Input     string `json:"input"`      // Input parameters (JSON)
	Output    string `json:"output"`     // Output result (JSON)
	Status    string `json:"status"`     // "pending", "running", "completed", "failed"
	Error     string `json:"error"`      // Error message if failed
	CreatedAt string `json:"created_at"` // ISO 8601 format
	UpdatedAt string `json:"updated_at"` // ISO 8601 format
}

// CreateTaskRequest represents a request to create a task
type CreateTaskRequest struct {
	SessionID string                 `json:"session_id" yaml:"session_id"`
	Skill     string                 `json:"skill" yaml:"skill"`
	Input     map[string]interface{} `json:"input" yaml:"input"`
}

// UpdateTaskRequest represents a request to update a task
type UpdateTaskRequest struct {
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
	Output string `json:"output,omitempty" yaml:"output,omitempty"`
	Error  string `json:"error,omitempty"  yaml:"error,omitempty"`
}

// TaskResponse represents a task response
type TaskResponse struct {
	Success bool     `json:"success"`
	Task    *TaskDTO `json:"task,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// TasksResponse represents a list of tasks response
type TasksResponse struct {
	Success bool       `json:"success"`
	Tasks   []*TaskDTO `json:"tasks,omitempty"`
	Error   string     `json:"error,omitempty"`
}
