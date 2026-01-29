// Package gates provides hook gates for Claude Code.
// research.go: Research enforcement gate (TABULA RASA).
// DACE: Uses shared/pkg/patterns for dynamic code detection.
package gates

import (
	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/patterns"
	"github.com/spf13/cobra"
)

var researchHookMode bool

var researchCmd = &cobra.Command{
	Use:   "research",
	Short: "Research enforcement gate (TABULA RASA)",
	Run:   runResearchGate,
}

func init() {
	researchCmd.Flags().BoolVar(&researchHookMode, "hook", false, "Hook mode")
}

func runResearchGate(cmd *cobra.Command, args []string) {
	if !researchHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()
	ctx := enforce.NewContext()
	session := enforce.GetOrCreateSession()

	// Both WebSearch and WebFetch count as research
	if input.ToolName == "WebSearch" || input.ToolName == "WebFetch" {
		session.MarkResearchDone()
		hook.ExitSilent()
	}

	if input.ToolName == "Write" || input.ToolName == "Edit" {
		filePath := input.GetString("file_path")

		// Use shared patterns for code detection
		if patterns.IsCodeFile(filePath) && !session.ResearchDone {
			hook.ExitBlockTOON("TABULA_RASA", "WebSearch_required,cutoff:"+session.TrainingCutoff+",today:"+ctx.Today)
		}
	}

	hook.ExitSilent()
}
