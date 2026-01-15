// Package util provides common utility functions.
// project.go: Project detection using priority-based system.
// DACE: Single responsibility - project detection only.
package util

import "os"

// DetectProject returns the current project name.
// Priority: env var → index.toon → .git → .claude marker → memory bank → global
func DetectProject() string {
	// Priority 1: Explicit environment variable
	if envProject := os.Getenv("KAVACH_PROJECT"); envProject != "" {
		return envProject
	}

	wd := WorkingDir()
	if wd == "" {
		return "global"
	}

	// Priority 2: index.toon PROJECTS section
	if indexProject := detectFromIndex(wd); indexProject != "" {
		return indexProject
	}

	// Priority 3: Git root detection
	if gitProject := detectGitProject(wd); gitProject != "" {
		return gitProject
	}

	// Priority 4: Claude project marker
	if markerProject := detectClaudeMarker(wd); markerProject != "" {
		return markerProject
	}

	// Priority 5: Memory bank fuzzy match
	if memProject := detectFromMemoryBank(wd); memProject != "" {
		return memProject
	}

	// Priority 6: Fallback
	return "global"
}

// GetProjectDir returns the root directory of the current project.
func GetProjectDir() string {
	wd := WorkingDir()
	if wd == "" {
		return ""
	}

	// Check for git root
	if dir := findGitRoot(wd); dir != "" {
		return dir
	}

	// Check for Claude marker
	if dir := findClaudeMarker(wd); dir != "" {
		return dir
	}

	return ""
}

// findGitRoot walks up to find .git directory.
func findGitRoot(wd string) string {
	dir := wd
	for {
		if DirExists(join(dir, ".git")) || FileExists(join(dir, ".git")) {
			return dir
		}
		p := parent(dir)
		if p == dir {
			break
		}
		dir = p
	}
	return ""
}

// findClaudeMarker walks up to find .claude marker.
func findClaudeMarker(wd string) string {
	dir := wd
	for {
		if FileExists(join(dir, ".claude", "project.json")) ||
			FileExists(join(dir, ".claude-project")) {
			return dir
		}
		p := parent(dir)
		if p == dir {
			break
		}
		dir = p
	}
	return ""
}
