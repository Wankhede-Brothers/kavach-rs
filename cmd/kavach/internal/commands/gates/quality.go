// Package gates provides hook gates for Claude Code.
// quality.go: Code quality gate.
// DACE: Uses shared packages for validation and utilities.
package gates

import (
	"path/filepath"

	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/util"
	"github.com/claude/shared/pkg/validate"
	"github.com/spf13/cobra"
)

var qualityHookMode bool

var qualityCmd = &cobra.Command{
	Use:   "quality",
	Short: "Code quality gate",
	Run:   runQualityGate,
}

func init() {
	qualityCmd.Flags().BoolVar(&qualityHookMode, "hook", false, "Hook mode")
}

func runQualityGate(cmd *cobra.Command, args []string) {
	if !qualityHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()

	if input.ToolName != "Write" && input.ToolName != "Edit" {
		hook.ExitSilent()
	}

	filePath := input.GetString("file_path")
	content := input.GetString("content")
	if input.ToolName == "Edit" {
		content = input.GetString("new_string")
	}

	if content == "" {
		hook.ExitSilent()
	}

	// P2 FIX: Only validate .go files in project (kavach uses TOON, not JSON)
	ext := util.GetExtension(filePath)
	if ext != ".go" {
		hook.ExitSilent()
	}

	// P2 FIX: Skip files outside project directory
	if !isInProjectDir(filePath) {
		hook.ExitSilent()
	}

	// Go syntax validation (with fixed string/comment handling)
	if err := validate.GoSyntax(content); err != "" {
		hook.ExitBlockTOON("QUALITY", "ast_go:"+err)
	}

	// DACE: Line count check for Go files
	lineCount := util.CountLines(content)
	if lineCount > 100 {
		hook.ExitBlockTOON("DACE", "exceeds_100_lines:"+util.Itoa(lineCount))
	}

	hook.ExitSilent()
}

// isInProjectDir checks if file is within the current working directory.
// Uses filepath.Rel for safe path comparison (Go best practice).
func isInProjectDir(filePath string) bool {
	wd := util.WorkingDir()
	if wd == "" {
		return false
	}
	// filepath.Rel returns error if paths are on different drives or unrelated
	rel, err := filepath.Rel(wd, filePath)
	if err != nil {
		return false
	}
	// If relative path starts with "..", file is outside project
	return len(rel) < 2 || rel[:2] != ".."
}
