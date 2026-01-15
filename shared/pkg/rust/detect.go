// Package rust provides Rust CLI tool detection and integration.
// detect.go: Detect available Rust CLI tools.
// DACE: Single responsibility - tool detection only.
package rust

import (
	"os/exec"
	"sync"
)

// Tools holds paths to detected Rust CLI tools.
type Tools struct {
	Bat    string // cat replacement (syntax highlighting)
	Eza    string // ls replacement (icons, git)
	Fd     string // find replacement (faster)
	Rg     string // grep replacement (ripgrep)
	Sd     string // sed replacement (simpler)
	Procs  string // ps replacement
	Dust   string // du replacement
	Bottom string // top replacement
	Delta  string // diff replacement
}

var (
	tools     *Tools
	toolsOnce sync.Once
)

// Detect finds all available Rust CLI tools.
// Caches result - safe to call multiple times.
func Detect() *Tools {
	toolsOnce.Do(func() {
		tools = &Tools{
			Bat:    which("bat"),
			Eza:    which("eza"),
			Fd:     which("fd"),
			Rg:     which("rg"),
			Sd:     which("sd"),
			Procs:  which("procs"),
			Dust:   which("dust"),
			Bottom: which("btm"),
			Delta:  which("delta"),
		}
	})
	return tools
}

// which returns the path to a command or empty string.
func which(cmd string) string {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return ""
	}
	return path
}

// HasBat returns true if bat is available.
func (t *Tools) HasBat() bool { return t.Bat != "" }

// HasEza returns true if eza is available.
func (t *Tools) HasEza() bool { return t.Eza != "" }

// HasFd returns true if fd is available.
func (t *Tools) HasFd() bool { return t.Fd != "" }

// HasRg returns true if rg is available.
func (t *Tools) HasRg() bool { return t.Rg != "" }

// HasSd returns true if sd is available.
func (t *Tools) HasSd() bool { return t.Sd != "" }
