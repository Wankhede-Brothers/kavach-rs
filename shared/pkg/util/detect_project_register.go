// Package util provides common utility functions.
// detect_project_register.go: Project registration and folder creation.
// DACE: Single responsibility - registration functions only.
package util

import (
	"os"
	"path/filepath"
)

// registerNewProject adds a new project to index.toon and creates folders.
func registerNewProject(project, wd string) error {
	// Create Memory Bank folders for the project
	categories := []string{"kanban", "decisions", "patterns", "research", "roadmaps", "proposals"}
	for _, cat := range categories {
		dir := filepath.Join(MemoryBankPath(cat), project)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create STM project folder
	stmDir := filepath.Join(STMPath(), "projects", project)
	if err := os.MkdirAll(stmDir, 0755); err != nil {
		return err
	}

	// Add to index.toon
	return appendToIndex(project, wd)
}

// itoa converts int to string without fmt import.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
