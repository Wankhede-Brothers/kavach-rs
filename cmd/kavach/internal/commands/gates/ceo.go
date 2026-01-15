// Package gates provides hook gates for Claude Code.
// ceo.go: CEO orchestration gate with DACE skill injection.
// NO HARDCODING - All patterns loaded from config/*.toon at runtime.
package gates

import (
	"strings"
	"time"

	"github.com/claude/shared/pkg/config"
	"github.com/claude/shared/pkg/hook"
	"github.com/spf13/cobra"
)

var ceoHookMode bool

var ceoCmd = &cobra.Command{
	Use:   "ceo",
	Short: "CEO orchestration gate",
	Long: `[CEO_GATE]
desc: Orchestration gate with DACE skill injection
flow: Task -> Skill detection -> Agent validation -> Context injection
config: Patterns loaded from config/skill-patterns.toon (NO HARDCODING)`,
	Run: runCEOGate,
}

func init() {
	ceoCmd.Flags().BoolVar(&ceoHookMode, "hook", false, "Hook mode")
}

func runCEOGate(cmd *cobra.Command, args []string) {
	if !ceoHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()

	if input.ToolName != "Task" {
		hook.ExitSilent()
	}

	subagentType := input.GetString("subagent_type")
	if subagentType == "" {
		hook.ExitBlockTOON("CEO", "Task_requires_subagent_type")
	}

	// Validate agent using dynamic config (not hardcoded)
	if !config.IsValidAgent(subagentType) {
		hook.ExitModifyTOON("CEO", map[string]string{
			"status": "unknown_agent:" + subagentType,
		})
	}

	// DACE: Detect skill from task prompt using dynamic config
	prompt := input.GetString("prompt")
	skill := detectSkillFromConfig(prompt)
	today := time.Now().Format("2006-01-02")

	// P1 FIX: CEO orchestration flow directive
	if config.IsEngineer(subagentType) {
		orchDirective := map[string]string{
			"agent":  subagentType,
			"date":   today,
			"cutoff": "2025-01",
		}

		if skill != "" {
			orchDirective["skill"] = skill
			orchDirective["inject"] = "Invoke " + skill + " for domain expertise"
		}

		// P1 FIX: Add orchestration flow instructions
		orchDirective["CEO_FLOW"] = "DELEGATE->VERIFY->AEGIS"
		orchDirective["AFTER_TASK"] = "Verify result meets requirements"
		orchDirective["IF_FAIL"] = "Re-delegate with specific feedback"
		orchDirective["IF_PASS"] = "Run kavach orch aegis for final verification"

		hook.ExitModifyTOON("CEO_ORCHESTRATION", orchDirective)
	}

	// Non-engineer agents (ceo, research-director) - just validate
	hook.ExitSilent()
}

// detectSkillFromConfig detects skill using patterns from config file
// P1 FIX #1: NO HARDCODING - skills loaded dynamically with priority from config/skill-patterns.toon
func detectSkillFromConfig(prompt string) string {
	lower := strings.ToLower(prompt)

	// P1 FIX #1: Get skills sorted by priority from config (not hardcoded)
	skills := config.GetSkillsByPriority()

	for _, skill := range skills {
		if matchesKeywords(lower, skill.Keywords) {
			return skill.Name
		}
	}

	return ""
}

// matchesKeywords checks if text contains any keyword (word boundary)
func matchesKeywords(text string, keywords []string) bool {
	padded := " " + text + " "
	for _, kw := range keywords {
		patterns := []string{" " + kw + " ", " " + kw + ",", " " + kw + "."}
		for _, p := range patterns {
			if strings.Contains(padded, p) {
				return true
			}
		}
		if strings.Contains(kw, " ") && strings.Contains(text, kw) {
			return true
		}
	}
	return false
}
