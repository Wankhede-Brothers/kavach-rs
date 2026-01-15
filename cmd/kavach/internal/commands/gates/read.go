// Package gates provides hook gates for Claude Code.
// read.go: Read file blocker gate.
// DACE: Uses shared/pkg/patterns for dynamic patterns.
package gates

import (
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/patterns"
	"github.com/spf13/cobra"
)

var readHookMode bool

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read file blocker gate",
	Run:   runReadGate,
}

func init() {
	readCmd.Flags().BoolVar(&readHookMode, "hook", false, "Hook mode")
}

func runReadGate(cmd *cobra.Command, args []string) {
	if !readHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()

	if input.ToolName != "Read" {
		hook.ExitSilent()
	}

	filePath := input.GetString("file_path")
	if filePath == "" {
		hook.ExitBlockTOON("READ", "no_file_path")
	}

	// Check sensitive files using shared patterns
	if patterns.IsSensitive(filePath) {
		hook.ExitBlockTOON("READ", "sensitive_file")
	}

	// Warn for large files
	if patterns.IsLargeFile(filePath) {
		hook.ExitModifyTOON("READ", map[string]string{
			"warn": "large_file",
		})
	}

	hook.ExitSilent()
}
