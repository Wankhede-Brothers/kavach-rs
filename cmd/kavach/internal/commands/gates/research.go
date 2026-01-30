// Package gates provides hook gates for Claude Code.
// research.go: Research enforcement gate (TABULA RASA).
// DACE: Uses shared/pkg/patterns for dynamic code detection.
package gates

import (
	"strings"

	"github.com/claude/shared/pkg/agentic"
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
		if (patterns.IsCodeFile(filePath) || patterns.IsInfraFile(filePath)) && !session.ResearchDone {
			hook.ExitBlockTOON("TABULA_RASA", "WebSearch_required,cutoff:"+session.TrainingCutoff+",today:"+ctx.Today)
		}

		// Soft topic mismatch warning (same as prewrite)
		if session.ResearchDone && len(session.ResearchTopics) > 0 {
			frameworks := agentic.ExtractFrameworkFromTask(filePath)
			topicsJoined := strings.ToLower(strings.Join(session.ResearchTopics, " "))
			for _, fw := range frameworks {
				if !strings.Contains(topicsJoined, fw) {
					hook.ExitModifyTOON("RESEARCH_TOPIC_WARN", map[string]string{
						"warning": "File references '" + fw + "' but no matching research topic found",
						"suggest": "WebSearch " + fw + " before writing to " + filePath,
					})
				}
			}
		}
	}

	hook.ExitSilent()
}
