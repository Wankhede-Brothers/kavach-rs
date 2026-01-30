// Package patterns provides dynamic pattern loading from TOON config.
// check.go: Pattern checking utilities.
// DACE: Reusable check functions for all gates.
package patterns

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SanitizePath validates and sanitizes a path against traversal attacks.
// Returns cleaned path if valid, error if path traversal detected or outside allowed bases.
func SanitizePath(path string, allowedBases []string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	// Clean the path to remove . and .. components
	cleaned := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path traversal detected: %s", path)
	}

	// If no allowed bases specified, just return cleaned path
	if len(allowedBases) == 0 {
		return cleaned, nil
	}

	// Resolve to absolute path for comparison
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if path is within any allowed base
	for _, base := range allowedBases {
		absBase, err := filepath.Abs(base)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, absBase+string(filepath.Separator)) || absPath == absBase {
			return cleaned, nil
		}
	}

	return "", fmt.Errorf("path outside allowed directories: %s", path)
}

// ValidateIdentifier checks if a string is a safe identifier (alphanumeric + underscore/hyphen).
func ValidateIdentifier(name string) error {
	if name == "" {
		return fmt.Errorf("empty identifier")
	}
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-') {
			return fmt.Errorf("invalid character in identifier: %c", r)
		}
	}
	return nil
}

// IsSensitive checks if path matches any sensitive pattern.
func IsSensitive(path string) bool {
	cfg := Load()
	pathLower := strings.ToLower(path)
	for _, p := range cfg.Sensitive {
		if strings.Contains(pathLower, p) {
			return true
		}
	}
	return false
}

// IsBlocked checks if command matches any blocked pattern.
// Security: Case-insensitive matching for robust blocking.
func IsBlocked(cmd string) bool {
	cfg := Load()
	cmdLower := strings.ToLower(cmd)
	for _, p := range cfg.Blocked {
		if strings.Contains(cmdLower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// IsCodeFile checks if path is a code file.
func IsCodeFile(path string) bool {
	cfg := Load()
	pathLower := strings.ToLower(path)
	for _, ext := range cfg.CodeExts {
		if strings.HasSuffix(pathLower, ext) {
			return true
		}
	}
	return false
}

// IsInfraFile checks if path is an infrastructure/config file that needs research.
// Matches: Dockerfile, docker-compose*, *.yml, *.yaml, *.tf, *.tfvars,
// *.toml (non-Cargo.toml), Makefile, Jenkinsfile, .github/*, Caddyfile, nginx.conf, etc.
func IsInfraFile(path string) bool {
	pathLower := strings.ToLower(path)
	base := strings.ToLower(filepath.Base(path))

	// Exact filename matches
	infraFiles := []string{
		"dockerfile", "makefile", "jenkinsfile", "caddyfile",
		"nginx.conf", "docker-compose.yml", "docker-compose.yaml",
	}
	for _, f := range infraFiles {
		if base == f || strings.HasPrefix(base, "docker-compose") {
			return true
		}
	}

	// Extension matches
	infraExts := []string{".yml", ".yaml", ".tf", ".tfvars", ".hcl"}
	for _, ext := range infraExts {
		if strings.HasSuffix(pathLower, ext) {
			return true
		}
	}

	// .toml but not Cargo.toml (Cargo.toml is a code project file)
	if strings.HasSuffix(pathLower, ".toml") && base != "cargo.toml" {
		return true
	}

	// Path-based matches
	if strings.Contains(pathLower, ".github/") || strings.Contains(pathLower, ".gitlab-ci") {
		return true
	}

	return false
}

// IsLargeFile checks if path is a potentially large file.
func IsLargeFile(path string) bool {
	cfg := Load()
	pathLower := strings.ToLower(path)
	for _, ext := range cfg.LargeExts {
		if strings.HasSuffix(pathLower, ext) {
			return true
		}
	}
	return false
}

// IsValidAgent checks if agent is in valid agents list.
func IsValidAgent(agent string) bool {
	cfg := Load()
	for _, agents := range cfg.ValidAgents {
		for _, a := range agents {
			if a == agent {
				return true
			}
		}
	}
	// Also check built-in agents
	builtins := []string{"Explore", "Plan", "Bash"}
	for _, b := range builtins {
		if b == agent {
			return true
		}
	}
	return false
}

// ClassifyIntent classifies prompt into intent category.
func ClassifyIntent(prompt string) string {
	cfg := Load()
	promptLower := strings.ToLower(prompt)
	for intent, words := range cfg.IntentWords {
		for _, word := range words {
			if strings.Contains(promptLower, word) {
				return intent
			}
		}
	}
	return "general"
}
