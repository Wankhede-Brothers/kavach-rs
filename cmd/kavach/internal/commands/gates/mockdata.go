// Package gates provides hook gates for Claude Code.
// mockdata.go: PostToolUse:Write/Edit gate — blocks mock/hardcoded data.
package gates

import (
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/patterns"
	"github.com/spf13/cobra"
)

var mockdataHookMode bool

var mockdataCmd = &cobra.Command{
	Use:   "mockdata",
	Short: "Block mock/hardcoded data in Write/Edit",
	Run:   runMockDataGate,
}

func init() {
	mockdataCmd.Flags().BoolVar(&mockdataHookMode, "hook", false, "Hook mode")
}

func runMockDataGate(cmd *cobra.Command, args []string) {
	if !mockdataHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()

	// Only check Write and Edit tools
	if input.ToolName != "Write" && input.ToolName != "Edit" {
		hook.ExitSilent()
	}

	filePath := input.GetString("file_path")

	// Get content: Write uses "content", Edit uses "new_string"
	content := input.GetString("content")
	if input.ToolName == "Edit" {
		content = input.GetString("new_string")
	}

	if content == "" || filePath == "" {
		hook.ExitSilent()
	}

	detected, reason := patterns.DetectMockData(filePath, content)
	if detected {
		hook.ExitBlockTOON("MOCK_DATA", reason+
			"\nINSTRUCTION: Replace hardcoded data with real API fetch."+
			"\nUse: useState + useEffect + fetch('/api/...') pattern."+
			"\nFor backend: Use sqlx::query_as from database tables."+
			"\nNEVER use const mock/dummy/fake/sample arrays.")
	}

	hook.ExitSilent()
}
