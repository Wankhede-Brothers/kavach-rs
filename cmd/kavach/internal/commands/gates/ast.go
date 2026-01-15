// Package gates provides hook gates for Claude Code.
// ast.go: AST validation gate.
// DACE: Uses shared/pkg/validate for syntax validation.
package gates

import (
	"strings"

	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/util"
	"github.com/claude/shared/pkg/validate"
	"github.com/spf13/cobra"
)

var astHookMode bool

var astCmd = &cobra.Command{
	Use:   "ast",
	Short: "AST validation gate",
	Run:   runASTGate,
}

func init() {
	astCmd.Flags().BoolVar(&astHookMode, "hook", false, "Hook mode")
}

func runASTGate(cmd *cobra.Command, args []string) {
	if !astHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()

	if input.ToolName != "Write" && input.ToolName != "Edit" {
		hook.ExitSilent()
	}

	filePath := input.GetString("file_path")
	if filePath == "" {
		hook.ExitSilent()
	}

	content := input.GetString("content")
	if input.ToolName == "Edit" {
		content = input.GetString("new_string")
	}

	if strings.TrimSpace(content) == "" {
		hook.ExitSilent()
	}

	// Validate syntax based on extension
	ext := strings.ToLower(util.GetExtension(filePath))

	switch ext {
	case ".go":
		if err := validate.GoSyntax(content); err != "" {
			hook.ExitBlockTOON("AST", "Go:"+err)
		}
	case ".json":
		if err := validate.JSONSyntax(content); err != "" {
			hook.ExitBlockTOON("AST", "JSON:"+err)
		}
	}

	hook.ExitSilent()
}
