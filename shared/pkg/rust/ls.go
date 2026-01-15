// Package rust provides Rust CLI tool integration.
// ls.go: Directory listing with eza (icons, git status).
// DACE: Single responsibility - ls/eza operations.
package rust

import (
	"os"
	"os/exec"
	"strings"
)

// Ls lists directory contents.
// Uses eza if available for icons and git status.
func Ls(dir string) ([]string, error) {
	tools := Detect()
	if tools.HasEza() {
		return ezaList(tools.Eza, dir)
	}
	return osReadDir(dir)
}

// LsTree shows directory tree.
// Uses eza --tree if available.
func LsTree(dir string, depth int) (string, error) {
	tools := Detect()
	if tools.HasEza() {
		return ezaTree(tools.Eza, dir, depth)
	}
	return "", nil // No tree fallback
}

// LsLong shows detailed listing.
func LsLong(dir string) (string, error) {
	tools := Detect()
	if tools.HasEza() {
		return ezaLong(tools.Eza, dir)
	}
	return "", nil
}

// ezaList uses eza to list directory.
func ezaList(ezaPath, dir string) ([]string, error) {
	cmd := exec.Command(ezaPath, "--oneline", dir)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return splitTrimmed(string(out)), nil
}

// ezaTree uses eza for tree view.
func ezaTree(ezaPath, dir string, depth int) (string, error) {
	args := []string{"--tree", "--icons"}
	if depth > 0 {
		args = append(args, "--level", itoa(depth))
	}
	args = append(args, dir)
	cmd := exec.Command(ezaPath, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ezaLong uses eza for long listing.
func ezaLong(ezaPath, dir string) (string, error) {
	cmd := exec.Command(ezaPath, "-la", "--icons", "--git", dir)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// osReadDir fallback using os.ReadDir.
func osReadDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	return names, nil
}

// splitTrimmed splits by newline and trims whitespace.
func splitTrimmed(s string) []string {
	var result []string
	for _, line := range strings.Split(s, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// itoa converts int to string (avoids strconv import).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
