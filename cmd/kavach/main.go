// Package main provides the entry point for Kavach - Brahmastra Stack CLI.
// Universal enforcement binary for AI coding assistants (Claude Code, OpenCode).
package main

import (
	"fmt"
	"os"

	"github.com/claude/cmd/kavach/internal/commands"
)

// Version information set at build time.
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	if err := commands.Execute(Version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
