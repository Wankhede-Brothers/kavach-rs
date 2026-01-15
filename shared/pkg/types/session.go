package types

import "time"

// SessionInfo represents session metadata.
type SessionInfo struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	Project   string    `json:"project"`
	WorkDir   string    `json:"work_dir"`
}

// ProjectContext represents project-specific context.
type ProjectContext struct {
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	Language   string   `json:"language,omitempty"`
	Framework  string   `json:"framework,omitempty"`
	FilesRead  []string `json:"files_read,omitempty"`
	FilesCount int      `json:"files_count,omitempty"`
}

// TaskDefinition represents a task in the scratchpad.
type TaskDefinition struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// STMContext represents short-term memory context.
type STMContext struct {
	CurrentTask  *TaskDefinition `json:"current_task,omitempty"`
	Focus        string          `json:"focus,omitempty"`
	LastActivity string          `json:"last_activity,omitempty"`
	Project      string          `json:"project,omitempty"`
}

// NewSessionInfo creates a new session info with generated ID.
func NewSessionInfo(project, workDir string) *SessionInfo {
	return &SessionInfo{
		ID:        "sess_" + time.Now().Format("20060102150405"),
		StartTime: time.Now(),
		Project:   project,
		WorkDir:   workDir,
	}
}
