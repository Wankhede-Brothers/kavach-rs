// Package rust provides Rust CLI tool integration.
// cat.go: File reading with bat (syntax highlighting).
// DACE: Single responsibility - cat/bat operations.
package rust

import (
	"os"
	"os/exec"
)

// Cat reads a file, using bat if available for syntax highlighting.
// Falls back to os.ReadFile if bat unavailable.
func Cat(path string) ([]byte, error) {
	tools := Detect()
	if tools.HasBat() {
		return batPlain(tools.Bat, path)
	}
	return os.ReadFile(path)
}

// CatHighlight reads a file with syntax highlighting (bat only).
// Returns empty if bat unavailable.
func CatHighlight(path string) ([]byte, error) {
	tools := Detect()
	if !tools.HasBat() {
		return os.ReadFile(path)
	}
	return batColored(tools.Bat, path)
}

// batPlain runs bat --plain (no line numbers, no decoration).
func batPlain(batPath, filePath string) ([]byte, error) {
	cmd := exec.Command(batPath, "--plain", "--paging=never", filePath)
	return cmd.Output()
}

// batColored runs bat with colors and line numbers.
func batColored(batPath, filePath string) ([]byte, error) {
	cmd := exec.Command(batPath, "--color=always", "--paging=never", filePath)
	return cmd.Output()
}

// CatTOON reads a TOON file with appropriate highlighting.
func CatTOON(path string) ([]byte, error) {
	tools := Detect()
	if tools.HasBat() {
		// Use YAML highlighting for TOON (similar syntax)
		cmd := exec.Command(tools.Bat, "--plain", "--paging=never", "-l", "yaml", path)
		return cmd.Output()
	}
	return os.ReadFile(path)
}
