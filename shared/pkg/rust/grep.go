// Package rust provides Rust CLI tool integration.
// grep.go: Content searching with rg (ripgrep).
// DACE: Single responsibility - grep/rg operations.
package rust

import (
	"os/exec"
	"strings"
)

// Grep searches for pattern in files.
// Uses rg if available, returns matching lines.
func Grep(pattern, path string) ([]string, error) {
	tools := Detect()
	if tools.HasRg() {
		return rgGrep(tools.Rg, pattern, path)
	}
	return nil, nil // No fallback - rg strongly preferred
}

// GrepFiles returns files containing pattern.
func GrepFiles(pattern, dir string) ([]string, error) {
	tools := Detect()
	if tools.HasRg() {
		return rgFiles(tools.Rg, pattern, dir)
	}
	return nil, nil
}

// GrepTOON searches TOON files for pattern.
func GrepTOON(pattern, dir string) ([]string, error) {
	tools := Detect()
	if tools.HasRg() {
		return rgType(tools.Rg, pattern, dir, "toon")
	}
	return nil, nil
}

// rgGrep uses ripgrep to search file content.
func rgGrep(rgPath, pattern, path string) ([]string, error) {
	cmd := exec.Command(rgPath, "--no-heading", pattern, path)
	out, err := cmd.Output()
	if err != nil {
		// rg returns exit 1 for no matches
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}
	return splitNonEmpty(string(out)), nil
}

// rgFiles uses ripgrep to find files with matches.
func rgFiles(rgPath, pattern, dir string) ([]string, error) {
	cmd := exec.Command(rgPath, "--files-with-matches", pattern, dir)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}
	return splitNonEmpty(string(out)), nil
}

// rgType uses ripgrep with type filter.
func rgType(rgPath, pattern, dir, fileType string) ([]string, error) {
	cmd := exec.Command(rgPath, "--type-add", fileType+":*."+fileType, "-t", fileType, "--files-with-matches", pattern, dir)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}
	return splitNonEmpty(string(out)), nil
}

// splitNonEmpty splits by newline, filtering empty strings.
func splitNonEmpty(s string) []string {
	var result []string
	for _, line := range strings.Split(s, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
