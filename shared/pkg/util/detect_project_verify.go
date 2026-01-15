// Package util provides common utility functions.
// detect_project_verify.go: Project existence and path verification.
// DACE: Single responsibility - verification functions only.
package util

import (
	"os"
	"path/filepath"
	"strings"
)

// projectExistsInMemoryBank checks if project folder exists.
func projectExistsInMemoryBank(project string) bool {
	kanbanPath := filepath.Join(KanbanPath(), project)
	return DirExists(kanbanPath)
}

// verifyProjectPath checks if project's registered path matches working dir.
func verifyProjectPath(project, wd string) bool {
	indexPath := IndexPath()
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return false
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, project+",") {
			parts := strings.Split(trimmed, ",")
			if len(parts) >= 2 {
				registeredPath := expandPath(strings.TrimSpace(parts[1]))
				return strings.HasPrefix(wd, registeredPath) || wd == registeredPath
			}
		}
	}
	return false
}

// createPathBasedProject creates a project from the working directory path.
func createPathBasedProject(wd string) string {
	project := filepath.Base(wd)
	project = strings.ToLower(project)
	project = strings.ReplaceAll(project, " ", "-")

	if err := registerNewProject(project, wd); err != nil {
		return "global" // Fallback
	}
	return project
}

// detectStack tries to detect the project's tech stack.
func detectStack(wd string) string {
	var stack []string

	if FileExists(filepath.Join(wd, "Cargo.toml")) {
		stack = append(stack, "rust")
	}
	if FileExists(filepath.Join(wd, "go.mod")) {
		stack = append(stack, "go")
	}
	if FileExists(filepath.Join(wd, "package.json")) {
		stack = append(stack, "node")
	}
	if FileExists(filepath.Join(wd, "pyproject.toml")) || FileExists(filepath.Join(wd, "requirements.txt")) {
		stack = append(stack, "python")
	}

	if len(stack) == 0 {
		return "*"
	}
	return strings.Join(stack, "|")
}
