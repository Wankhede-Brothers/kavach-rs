// Package session provides session state management.
// task.go: Task state management functions.
// DACE: Single responsibility - task tracking only.
// Called by: gates/task.go (TaskCreate/TaskUpdate) and memory/sync.go (Write/Edit).
package session

// SetTask updates the current task being worked on.
// Called by: task gate on TaskUpdate with status change.
func (s *SessionState) SetTask(task, status string) {
	s.CurrentTask = task
	s.TaskStatus = status
	s.Save()
}

// AddFileModified tracks a file that was modified in this session.
// Returns true if the file was newly added, false if already tracked.
// Called by: memory sync on PostToolUse:Write/Edit.
func (s *SessionState) AddFileModified(filePath string) bool {
	for _, f := range s.FilesModified {
		if f == filePath {
			return false
		}
	}
	s.FilesModified = append(s.FilesModified, filePath)
	s.Save()
	return true
}

// ClearTask clears the current task state.
// Called by: task gate on task completion/deletion.
func (s *SessionState) ClearTask() {
	s.CurrentTask = ""
	s.TaskStatus = ""
	s.FilesModified = []string{}
	s.Save()
}

// HasTask returns true if a task is currently active.
// Called by: CEO gate to check if a task context exists.
func (s *SessionState) HasTask() bool {
	return s.CurrentTask != ""
}
