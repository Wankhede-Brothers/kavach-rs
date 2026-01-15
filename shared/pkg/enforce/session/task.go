// Package session provides session state management.
// task.go: Task state management functions.
// DACE: Single responsibility - task tracking only.
package session

// SetTask updates the current task being worked on.
func (s *SessionState) SetTask(task, status string) {
	s.CurrentTask = task
	s.TaskStatus = status
	s.Save()
}

// AddFileModified tracks a file that was modified in this session.
func (s *SessionState) AddFileModified(filePath string) {
	for _, f := range s.FilesModified {
		if f == filePath {
			return // Already tracked
		}
	}
	s.FilesModified = append(s.FilesModified, filePath)
	s.Save()
}

// ClearTask clears the current task state.
func (s *SessionState) ClearTask() {
	s.CurrentTask = ""
	s.TaskStatus = ""
	s.FilesModified = []string{}
	s.Save()
}

// HasTask returns true if a task is currently active.
func (s *SessionState) HasTask() bool {
	return s.CurrentTask != ""
}
