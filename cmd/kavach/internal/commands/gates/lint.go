// Package gates provides hook gates for Claude Code.
// lint.go: Lint validation gate.
// DACE: Uses shared/pkg/util for utilities.
package gates

import (
	"strings"

	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var lintHookMode bool

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint validation gate",
	Run:   runLintGate,
}

func init() {
	lintCmd.Flags().BoolVar(&lintHookMode, "hook", false, "Hook mode")
}

func runLintGate(cmd *cobra.Command, args []string) {
	if !lintHookMode {
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

	issues := checkLintIssues(content, filePath)

	if len(issues) > 0 {
		hook.ExitModifyTOON("LINT", map[string]string{
			"issues": strings.Join(issues[:util.Min(3, len(issues))], ","),
		})
	}

	hook.ExitSilent()
}

func checkLintIssues(content, filePath string) []string {
	issues := []string{}
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t") {
			issues = append(issues, "trailing_ws:"+util.Itoa(i+1))
		}
	}

	ext := util.GetExtension(filePath)
	if ext == ".go" {
		for i, line := range lines {
			if strings.HasPrefix(line, "    ") && !strings.HasPrefix(line, "\t") {
				issues = append(issues, "spaces:"+util.Itoa(i+1))
			}
		}
	}

	return issues
}
