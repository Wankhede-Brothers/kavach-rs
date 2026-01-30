// Package gates provides hook gates for Claude Code.
// intent.go: Intent classification gate with multi-tier cascade.
// Tier 0: Trivial → silent exit (0 tokens)
// Tier 1: Status query → ~50 tokens
// Tier 2: Session-aware (post-compact, reinforcement) → ~80 tokens
// Tier 3: Full NLU classification → ~200 tokens
package gates

import (
	"strings"
	"time"

	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/telemetry"
	"github.com/spf13/cobra"
)

var intentHookMode bool

var intentCmd = &cobra.Command{
	Use:   "intent",
	Short: "Intent classification gate with multi-tier cascade",
	Run:   runIntentGate,
}

func init() {
	intentCmd.Flags().BoolVar(&intentHookMode, "hook", false, "Hook mode")
}

func runIntentGate(cmd *cobra.Command, args []string) {
	if !intentHookMode {
		cmd.Help()
		return
	}

	span := telemetry.StartSpan("intent")
	defer span.End()

	input := hook.MustReadHookInput()
	prompt := strings.ToLower(strings.TrimSpace(input.GetString("prompt")))

	if prompt == "" {
		span.SetTier(0)
		span.SetTokens(0)
		hook.ExitUserPromptSubmitSilent()
	}

	// TIER 0: Trivial prompts — zero I/O, zero tokens
	if isSimpleQuery(prompt) {
		span.SetTier(0)
		span.SetTokens(0)
		span.SetResult("silent")
		hook.ExitUserPromptSubmitSilent()
	}

	// TIER 1: Status queries — SessionIdentity only, ~50 tokens
	if isStatusQuery(prompt) {
		span.SetTier(1)
		span.SetTokens(50)
		hook.ExitUserPromptSubmitWithContext(statusDirective())
	}

	// TIER 2+: Load session for post-compact and reinforcement checks
	session := enforce.GetOrCreateSession()
	session.IncrementTurn()
	session.ResetResearchForNewPrompt()
	span.SetSessionLoaded(true)
	today := time.Now().Format("2006-01-02")

	var contextBlocks []string

	// TIER 2: Post-compact recovery — SessionTracking, ~80 tokens
	if session.IsPostCompact() {
		span.SetTier(2)
		contextBlocks = append(contextBlocks, postCompactRecovery(session))
		session.ClearPostCompact()
		session.MarkReinforcementDone()
	}

	// TIER 2: Periodic reinforcement
	if session.NeedsReinforcement() && !session.IsPostCompact() {
		span.SetTier(2)
		contextBlocks = append(contextBlocks, periodicReinforcement(session))
		session.MarkReinforcementDone()
	}

	// TIER 3: Full NLU classification
	intent := classifyIntentFromConfig(prompt)
	if intent != nil && intent.Type != "unclassified" {
		span.SetTier(3)
		session.MarkNLUParsed()
		session.StoreIntent(intent.Type, intent.Domain, intent.SubAgents, intent.Skills)

		if isImplementationIntent(intent.Type) {
			intent.ResearchReq = true
			if containsTechnicalTerms(prompt) {
				intent.Confidence = "high"
			}
		}

		contextBlocks = append(contextBlocks, formatIntentDirective(intent, today))
		span.SetTokens(estimateTokens(contextBlocks))
	} else if intent != nil {
		// Unclassified but non-trivial — still set tier 3 with minimal output
		span.SetTier(3)
		span.SetTokens(0)
	}

	if len(contextBlocks) > 0 {
		hook.ExitUserPromptSubmitWithContext(strings.Join(contextBlocks, "\n\n"))
	}

	hook.ExitUserPromptSubmitSilent()
}

// estimateTokens gives a rough token count (4 chars per token).
func estimateTokens(blocks []string) int {
	total := 0
	for _, b := range blocks {
		total += len(b)
	}
	return total / 4
}
