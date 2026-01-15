// Package util provides common utility functions for the umbrella CLI.
// file_stm.go: STM (Short-Term Memory) directory management.
// DACE: Single responsibility - STM directory functions only.
// P0 FIX #2,#3: Ensures STM directory structure exists and consistent scratchpad paths.
package util

import "path/filepath"

// EnsureMemoryBankDirs creates all required memory bank directories.
// Called during session init to ensure STM and other directories exist.
// Uses os.MkdirAll with 0755 permissions per Go best practices.
// Reference: https://pkg.go.dev/os#MkdirAll
func EnsureMemoryBankDirs(project string) error {
	memDir := MemoryDir()

	// Core memory bank categories
	categories := []string{
		"decisions",
		"graph",
		"kanban",
		"patterns",
		"proposals",
		"research",
		"roadmaps",
		"STM",
	}

	// Create category directories
	for _, cat := range categories {
		if err := EnsureDir(filepath.Join(memDir, cat)); err != nil {
			return err
		}
	}

	// Create project-specific directories
	if project != "" && project != "global" {
		projectCategories := []string{"decisions", "kanban", "patterns", "proposals", "roadmaps"}
		for _, cat := range projectCategories {
			if err := EnsureDir(filepath.Join(memDir, cat, project)); err != nil {
				return err
			}
		}
	}

	// Create global directory in each category
	for _, cat := range []string{"decisions", "kanban", "patterns", "proposals", "roadmaps"} {
		if err := EnsureDir(filepath.Join(memDir, cat, "global")); err != nil {
			return err
		}
	}

	// Create STM/projects directory for scratchpads
	if err := EnsureDir(filepath.Join(memDir, "STM", "projects")); err != nil {
		return err
	}

	// Create project-specific STM directory
	if project != "" {
		if err := EnsureDir(filepath.Join(memDir, "STM", "projects", project)); err != nil {
			return err
		}
	}

	return nil
}

// EnsureScratchpadDir creates the scratchpad directory for a project.
func EnsureScratchpadDir(project string) (string, error) {
	scratchpadDir := filepath.Join(STMPath(), "projects", project)
	if err := EnsureDir(scratchpadDir); err != nil {
		return "", err
	}
	return scratchpadDir, nil
}

// ScratchpadPath returns the path to a project's scratchpad.toon file.
func ScratchpadPath(project string) string {
	return filepath.Join(STMPath(), "projects", project, "scratchpad.toon")
}
