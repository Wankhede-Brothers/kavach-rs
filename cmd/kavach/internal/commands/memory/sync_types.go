// Package memory provides memory bank commands.
// sync_types.go: Type definitions for sync operations.
// DACE: Micro-modular split from sync.go
package memory

// TodoItem represents a single todo from TodoWrite
type TodoItem struct {
	Content    string `json:"content"`
	Status     string `json:"status"`
	ActiveForm string `json:"activeForm"`
}

// TodoWriteResult represents the PostToolUse result for TodoWrite
type TodoWriteResult struct {
	Tool   string `json:"tool"`
	Result struct {
		Todos []TodoItem `json:"todos"`
	} `json:"tool_result"`
}
