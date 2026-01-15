// Package util provides common utility functions for the umbrella CLI.
// paths.go: Core path utilities (max 100 lines).
// DACE: Single responsibility - path functions only.
package util

import (
	"os"
	"path/filepath"
)

// HomeDir returns the user's home directory.
func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// ClaudeDir returns the .claude directory path.
func ClaudeDir() string {
	return filepath.Join(HomeDir(), ".claude")
}

// SharedAIDir returns the shared-ai directory path.
// DEPRECATED: Use GetPaths() for platform-aware paths.
func SharedAIDir() string {
	return filepath.Join(HomeDir(), ".local", "shared", "shared-ai")
}

// MemoryDir returns the memory bank directory path.
func MemoryDir() string {
	legacyPath := filepath.Join(HomeDir(), ".local", "shared", "shared-ai", "memory")
	if DirExists(legacyPath) {
		return legacyPath
	}
	paths := GetPaths(DetectCLI())
	return paths.Memory
}

// DetectCLI determines which CLI is being used.
func DetectCLI() CLIType {
	if DirExists(filepath.Join(HomeDir(), ".claude")) {
		return CLIClaudeCode
	}
	if os.Getenv("OPENCODE_HOME") != "" {
		return CLIOpenCode
	}
	return CLIClaudeCode
}

// MemoryBankPath returns the path to a memory bank category.
func MemoryBankPath(category string) string {
	return filepath.Join(MemoryDir(), category)
}

// MemoryFile returns path to a top-level memory file.
func MemoryFile(name string) string {
	return filepath.Join(MemoryDir(), name)
}

// GovernancePath returns the path to GOVERNANCE.toon.
func GovernancePath() string {
	return MemoryFile("GOVERNANCE.toon")
}

// IndexPath returns the path to index.toon.
func IndexPath() string {
	return MemoryFile("index.toon")
}

// VolatilePath returns the path to volatile.toon.
func VolatilePath() string {
	return MemoryFile("volatile.toon")
}

// BinDir returns the claude bin directory path.
func BinDir() string {
	return filepath.Join(ClaudeDir(), "bin")
}

// ProjectsDir returns the projects directory path.
func ProjectsDir() string {
	return filepath.Join(ClaudeDir(), "projects")
}

// SettingsPath returns the path to settings.json.
func SettingsPath() string {
	return filepath.Join(ClaudeDir(), "settings.json")
}

// STMPath returns the path to Short-Term Memory directory.
func STMPath() string {
	return MemoryBankPath("STM")
}

// GraphPath returns the path to the graph directory.
func GraphPath() string {
	return MemoryBankPath("graph")
}

// KanbanPath returns the path to the kanban directory.
func KanbanPath() string {
	return MemoryBankPath("kanban")
}

// ProjectMemoryPath returns project-specific memory path.
func ProjectMemoryPath(workDir, category string) string {
	projectID := sanitizePath(workDir)
	return filepath.Join(MemoryBankPath(category), projectID)
}

// WorkingDir returns the current working directory.
func WorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

// sanitizePath converts a path to a safe directory name.
func sanitizePath(path string) string {
	result := ""
	for _, ch := range path {
		if ch == '/' || ch == '\\' {
			result += "_"
		} else if ch == '.' || ch == '-' || ch == '_' ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') {
			result += string(ch)
		}
	}
	return result
}
