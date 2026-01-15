// Package gates provides hook gates for Claude Code.
// intent.go: Intent classification gate entry point.
// DACE: Micro-modular - types, nlu, output, helpers in separate files.
package gates

import (
	"strings"
	"time"

	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/spf13/cobra"
)

var intentHookMode bool

var intentCmd = &cobra.Command{
	Use:   "intent",
	Short: "Intent classification gate with AGI-like NLU",
	Long: `[NLU_INTENT_GATE]
desc: AGI-like intent classification for vague natural language
config: Patterns loaded from config/nlu-patterns.toon (NO HARDCODING)
flow: User prompt -> NLU classification -> Agent/Skill recommendation`,
	Run: runIntentGate,
}

func init() {
	intentCmd.Flags().BoolVar(&intentHookMode, "hook", false, "Hook mode")
}

func runIntentGate(cmd *cobra.Command, args []string) {
	if !intentHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()
	prompt := strings.ToLower(input.GetString("prompt"))

	if prompt == "" {
		hook.ExitUserPromptSubmitSilent()
	}

	session := enforce.GetOrCreateSession()
	session.IncrementTurn()
	today := time.Now().Format("2006-01-02")

	var contextBlocks []string

	// 1. POST-COMPACT RECOVERY
	if session.IsPostCompact() {
		contextBlocks = append(contextBlocks, postCompactRecovery(session))
		session.ClearPostCompact()
		session.MarkReinforcementDone()
	}

	// 2. PERIODIC REINFORCEMENT
	if session.NeedsReinforcement() && !session.IsPostCompact() {
		contextBlocks = append(contextBlocks, periodicReinforcement(session))
		session.MarkReinforcementDone()
	}

	// 3. STATUS QUERIES
	if isStatusQuery(prompt) {
		contextBlocks = append(contextBlocks, statusDirective())
		hook.ExitUserPromptSubmitWithContext(strings.Join(contextBlocks, "\n\n"))
	}

	// 4. AGI NLU: Classify intent using dynamic config
	intent := classifyIntentFromConfig(prompt)
	if intent != nil {
		// P1 FIX: Reset research state for implementation intents
		// This makes research TASK-SCOPED, not session-scoped
		if isImplementationIntent(intent.Type) && intent.ResearchReq {
			session.ResetTaskResearch()
		}
		contextBlocks = append(contextBlocks, formatIntentDirective(intent, today))
	}

	if len(contextBlocks) > 0 {
		hook.ExitUserPromptSubmitWithContext(strings.Join(contextBlocks, "\n\n"))
	}

	hook.ExitUserPromptSubmitSilent()
}
