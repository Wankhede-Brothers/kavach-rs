// Package util provides common utility functions.
// detect_git.go: Project detection from .git and .claude markers.
// DACE: Single responsibility - git/marker detection only.
package util

import "os"

// detectGitProject walks up to find .git directory.
func detectGitProject(wd string) string {
	dir := wd
	for {
		gitPath := join(dir, ".git")
		if DirExists(gitPath) || FileExists(gitPath) {
			return base(dir)
		}
		p := parent(dir)
		if p == dir {
			break
		}
		dir = p
	}
	return ""
}

// detectClaudeMarker checks for Claude project markers.
func detectClaudeMarker(wd string) string {
	dir := wd
	for {
		// Check .claude/project.json
		projectJSON := join(dir, ".claude", "project.json")
		if FileExists(projectJSON) {
			if name := readProjectName(projectJSON); name != "" {
				return name
			}
			return base(dir)
		}

		// Check .claude-project
		markerFile := join(dir, ".claude-project")
		if FileExists(markerFile) {
			if data, err := os.ReadFile(markerFile); err == nil && len(data) > 0 {
				name := trimSpaces(string(data))
				if name != "" {
					return name
				}
			}
			return base(dir)
		}

		p := parent(dir)
		if p == dir {
			break
		}
		dir = p
	}
	return ""
}

// readProjectName reads name/project from JSON file.
func readProjectName(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	content := string(data)
	for _, key := range []string{`"name"`, `"project"`} {
		if idx := indexOf(content, key); idx >= 0 {
			rest := content[idx+len(key):]
			if colonIdx := indexOf(rest, ":"); colonIdx >= 0 {
				rest = rest[colonIdx+1:]
				if quoteStart := indexOf(rest, `"`); quoteStart >= 0 {
					rest = rest[quoteStart+1:]
					if quoteEnd := indexOf(rest, `"`); quoteEnd >= 0 {
						return rest[:quoteEnd]
					}
				}
			}
		}
	}
	return ""
}
