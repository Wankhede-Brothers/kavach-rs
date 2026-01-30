// Package gates provides hook gates for Claude Code.
// ceo.go: CEO orchestration gate with DACE skill injection.
// Deprecated: Use umbrella gates (pre-tool routes Task to preToolCEO).
// Kept for direct CLI invocation only: kavach gates ceo --hook < input.json
package gates

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/claude/shared/pkg/config"
	"github.com/claude/shared/pkg/dag"
	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/telemetry"
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

	span := telemetry.StartSpan("ceo")
	defer span.End()

	input := hook.MustReadHookInput()
	span.SetTool("Task")

	if input.ToolName != "Task" {
		hook.ExitSilent()
	}

	subagentType := input.GetString("subagent_type")
	if subagentType == "" {
		hook.ExitBlockTOON("CEO", "Task_requires_subagent_type")
	}

	// Validate agent using dynamic config (not hardcoded)
	// SECURITY FIX: BLOCK unknown agents instead of just warning
	if !config.IsValidAgent(subagentType) {
		hook.ExitBlockTOON("CEO", "unknown_agent:"+subagentType)
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

		// Load session for DAG and task context
		session := enforce.GetOrCreateSession()

		// If a task is already active, include it in context
		if session.HasTask() {
			orchDirective["current_task"] = session.CurrentTask
		}
		breakdown := extractBreakdown(prompt)
		agents := resolveAgents(session, subagentType)

		// Auto-decompose: if no explicit breakdown but intent has multiple SubAgents,
		// generate research + implementation steps automatically
		if len(breakdown) <= 1 && len(agents) > 1 {
			breakdown = autoDecompose(prompt, agents)
		}

		if len(breakdown) > 1 {
			nodes := dag.Decompose(breakdown, agents)
			state, err := dag.Schedule(session.SessionID, prompt, nodes)
			if err == nil {
				if saveErr := dag.Save(state); saveErr != nil {
					fmt.Fprintf(os.Stderr, "[CEO_DAG] Save error: %v\n", saveErr)
					orchDirective["WARNING"] = "DAG state NOT persisted: " + saveErr.Error()
				}
				directive := dag.BuildDirective(state)
				hook.ExitModifyTOONWithModule("CEO_DAG_DISPATCH", orchDirective, directive)
			} else {
				fmt.Fprintf(os.Stderr, "[CEO_DAG] Schedule error: %v\n", err)
				orchDirective["WARNING"] = "DAG scheduling failed: " + err.Error()
			}
		}

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

// extractBreakdown splits a multi-step prompt into individual steps.
// Supports: numbered lists (1. 2.), bullet lists (- ), semicolons,
// and natural language conjunctions (then, after that, next, finally).
func extractBreakdown(prompt string) []string {
	// Try structured lists first (numbered, bulleted)
	lines := strings.Split(prompt, "\n")
	var steps []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		i := 0
		for i < len(trimmed) && trimmed[i] >= '0' && trimmed[i] <= '9' {
			i++
		}
		if i > 0 && i < len(trimmed) && (trimmed[i] == '.' || trimmed[i] == ')') {
			steps = append(steps, strings.TrimSpace(trimmed[i+1:]))
		} else if strings.HasPrefix(trimmed, "- ") {
			steps = append(steps, strings.TrimSpace(trimmed[2:]))
		}
	}
	if len(steps) > 1 {
		return steps
	}

	// Try semicolons ("research auth; implement login; add tests")
	if strings.Contains(prompt, ";") {
		steps = nil
		for _, p := range strings.Split(prompt, ";") {
			if t := strings.TrimSpace(p); t != "" {
				steps = append(steps, t)
			}
		}
		if len(steps) > 1 {
			return steps
		}
	}

	// Try natural language conjunctions
	conjSplits := splitByConjunctions(prompt)
	if len(conjSplits) > 1 {
		return conjSplits
	}

	return nil
}

// splitByConjunctions splits a prompt by sequential conjunctions.
func splitByConjunctions(prompt string) []string {
	lower := " " + strings.ToLower(prompt) + " "
	conjunctions := []string{" then ", " after that ", " and then ", " next ", " finally ", " also "}

	// Find earliest conjunction
	bestPos, bestLen := -1, 0
	for _, conj := range conjunctions {
		if pos := strings.Index(lower, conj); pos >= 0 && (bestPos < 0 || pos < bestPos) {
			bestPos = pos
			bestLen = len(conj)
		}
	}
	if bestPos < 0 {
		return nil
	}

	// Recursive split: left part + split(right part)
	// Adjust for the leading space we added
	left := strings.TrimSpace(prompt[:bestPos])
	right := strings.TrimSpace(prompt[bestPos+bestLen-1:]) // -1 for leading space offset
	var result []string
	if left != "" {
		result = append(result, left)
	}
	rightParts := splitByConjunctions(right)
	if len(rightParts) > 0 {
		result = append(result, rightParts...)
	} else if right != "" {
		result = append(result, right)
	}
	return result
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
