// Package util provides common utility functions.
// detect_project_write.go: Exact project matching for write operations.
// DACE: Single responsibility - write-safe project detection entry point.
// BUG FIX: Never uses fuzzy matching, auto-creates new projects.
package util

import "os"

// DetectProjectForWrite returns project for write operations.
// CRITICAL: Uses EXACT working directory matching only.
// If no match found, creates new project entry in index.toon.
// Never uses fuzzy matching to prevent updating wrong project.
func DetectProjectForWrite() string {
	// Priority 1: Explicit environment variable
	if envProject := os.Getenv("KAVACH_PROJECT"); envProject != "" {
		return envProject
	}

	wd := WorkingDir()
	if wd == "" {
		return "global"
	}

	// Priority 2: Exact match from index.toon (by working directory path)
	if indexProject := detectFromIndex(wd); indexProject != "" {
		return indexProject
	}

	// Priority 3: Git root detection - get project name from git
	gitProject := detectGitProject(wd)
	if gitProject != "" {
		return resolveOrCreateProject(gitProject, wd)
	}

	// Priority 4: Claude project marker
	if markerProject := detectClaudeMarker(wd); markerProject != "" {
		return resolveOrCreateProject(markerProject, wd)
	}

	// Priority 5: Create new project from directory name
	return createPathBasedProject(wd)
}

// resolveOrCreateProject checks if project exists and path matches, or creates new.
func resolveOrCreateProject(project, wd string) string {
	if projectExistsInMemoryBank(project) {
		if verifyProjectPath(project, wd) {
			return project
		}
		// Path mismatch - create unique name
		return createPathBasedProject(wd)
	}
	// Project doesn't exist - register it
	if err := registerNewProject(project, wd); err == nil {
		return project
	}
	return createPathBasedProject(wd)
}
