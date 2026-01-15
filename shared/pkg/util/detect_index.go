// Package util provides common utility functions.
// detect_index.go: Project detection from index.toon.
// DACE: Single responsibility - index.toon parsing only.
package util

import "os"

// detectFromIndex reads index.toon PROJECTS section.
// Format: PROJECTS[N]{id,path,stack,aliases}
// BUG FIX: Returns MOST SPECIFIC match (longest path) not first match.
func detectFromIndex(wd string) string {
	indexPath := IndexPath()
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return ""
	}

	content := string(data)
	lines := splitLines(content)

	var bestMatch string
	var bestPathLen int

	inProjects := false
	for _, line := range lines {
		trimmed := trimSpaces(line)

		if hasPrefix(trimmed, "PROJECTS[") {
			inProjects = true
			continue
		}

		if inProjects && len(trimmed) > 0 && trimmed[0] >= 'A' && trimmed[0] <= 'Z' && !hasPrefix(line, " ") {
			break
		}

		if inProjects && (hasPrefix(line, "  ") || hasPrefix(line, "\t")) {
			parts := splitByComma(trimmed)
			if len(parts) >= 2 {
				projectID := parts[0]
				projectPath := expandPath(parts[1])

				if hasPrefix(wd, projectPath) || wd == projectPath {
					// Keep the LONGEST (most specific) match
					if len(projectPath) > bestPathLen {
						bestMatch = projectID
						bestPathLen = len(projectPath)
					}
				}
			}
		}
	}

	return bestMatch
}

// expandPath expands ~ to home directory.
func expandPath(path string) string {
	if hasPrefix(path, "~/") {
		return HomeDir() + path[1:]
	}
	if path == "~" {
		return HomeDir()
	}
	return path
}
