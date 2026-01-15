// Package util provides common utility functions.
// detect_memory.go: Project detection from Memory Bank with fuzzy matching.
// DACE: Single responsibility - memory bank matching only.
package util

import "os"

// detectFromMemoryBank matches workdir against known projects.
// Uses fuzzy matching: "Nicole Carpenter" → "nicole-carpenter-freelance"
func detectFromMemoryBank(wd string) string {
	kanbanDir := KanbanPath()
	entries, err := os.ReadDir(kanbanDir)
	if err != nil {
		return ""
	}

	baseName := base(wd)
	normalizedBase := normalizeForMatch(baseName)

	// Exact match first
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "global" || e.Name() == "TEMPLATE.toon" {
			continue
		}
		if e.Name() == baseName {
			return e.Name()
		}
	}

	// Fuzzy match: "Nicole Carpenter" → "nicole-carpenter"
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "global" || e.Name() == "TEMPLATE.toon" {
			continue
		}
		normalizedProject := normalizeForMatch(e.Name())
		if len(normalizedBase) >= 3 && len(normalizedProject) >= 3 {
			if hasPrefix(normalizedProject, normalizedBase) ||
				hasPrefix(normalizedBase, normalizedProject) {
				return e.Name()
			}
		}
	}

	// Path component match
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "global" || e.Name() == "TEMPLATE.toon" {
			continue
		}
		if isPathComponent(wd, e.Name()) {
			return e.Name()
		}
	}

	return ""
}

// normalizeForMatch converts to lowercase, spaces→hyphens.
func normalizeForMatch(s string) string {
	result := ""
	for _, ch := range s {
		if ch >= 'A' && ch <= 'Z' {
			result += string(ch + 32) // lowercase
		} else if ch == ' ' || ch == '_' {
			result += "-"
		} else if ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '-' {
			result += string(ch)
		}
	}
	return result
}

// isPathComponent checks if name is a path component.
func isPathComponent(path, name string) bool {
	components := splitPath(path)
	for _, comp := range components {
		if comp == name {
			return true
		}
	}
	return false
}

// splitPath splits path into components.
func splitPath(path string) []string {
	var components []string
	for path != "" {
		dir, file := split(path)
		if file != "" {
			components = append([]string{file}, components...)
		}
		if dir == path {
			break
		}
		path = clean(dir)
	}
	return components
}
