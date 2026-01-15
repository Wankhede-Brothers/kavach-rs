// Package rust provides Rust CLI tool integration.
// find.go: File finding with fd (faster than find).
// DACE: Single responsibility - find/fd operations.
package rust

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// Find searches for files matching pattern in directory.
// Uses fd if available, falls back to filepath.Glob.
func Find(dir, pattern string) ([]string, error) {
	tools := Detect()
	if tools.HasFd() {
		return fdFind(tools.Fd, dir, pattern)
	}
	return globFind(dir, pattern)
}

// FindTOON finds all .toon files in directory.
func FindTOON(dir string) ([]string, error) {
	tools := Detect()
	if tools.HasFd() {
		return fdFindExt(tools.Fd, dir, "toon")
	}
	return globFind(dir, "*.toon")
}

// fdFind uses fd to search for files.
func fdFind(fdPath, dir, pattern string) ([]string, error) {
	cmd := exec.Command(fdPath, "--type", "f", pattern, dir)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return splitLines(string(out)), nil
}

// fdFindExt uses fd to search by extension.
func fdFindExt(fdPath, dir, ext string) ([]string, error) {
	cmd := exec.Command(fdPath, "--type", "f", "--extension", ext, dir)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return splitLines(string(out)), nil
}

// globFind uses filepath.Glob as fallback.
func globFind(dir, pattern string) ([]string, error) {
	fullPattern := filepath.Join(dir, "**", pattern)
	return filepath.Glob(fullPattern)
}

// splitLines splits output by newlines, filtering empty.
func splitLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
