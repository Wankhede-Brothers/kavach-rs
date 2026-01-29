// Package gates provides hook gates for Claude Code.
// prewrite.go: Pre-write umbrella gate (PreToolUse:Write|Edit|NotebookEdit).
// Hierarchy: SECURITY(chain,content) → GUARD(code-guard) → RESEARCH
package gates

import (
	"github.com/claude/shared/pkg/agentic"
	"github.com/claude/shared/pkg/chain"
	"github.com/claude/shared/pkg/config"
	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/patterns"
	"github.com/claude/shared/pkg/types"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var preWriteHookMode bool

var preWriteCmd = &cobra.Command{
	Use:   "pre-write",
	Short: "Pre-write umbrella gate (security → guard → research)",
	Run:   runPreWriteGate,
}

func init() {
	preWriteCmd.Flags().BoolVar(&preWriteHookMode, "hook", false, "Hook mode")
}

func runPreWriteGate(cmd *cobra.Command, args []string) {
	if !preWriteHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()
	session := enforce.GetOrCreateSession()

	// L2: SECURITY — chain verification (Intent → CEO → Aegis → Research)
	if blocked, reason, context := runSecurityChain(input, session); blocked {
		hook.Output(&types.HookResponse{
			HookSpecificOutput: &types.HookSpecificOutput{
				HookEventName:            "PreToolUse",
				PermissionDecision:       "deny",
				PermissionDecisionReason: reason,
				AdditionalContext:        context,
			},
		})
		os.Exit(0)
	}

	// L2: SECURITY — content (secrets/credentials detection)
	if blocked, reason := runContentCheck(input); blocked {
		hook.ExitBlockTOON("CONTENT", reason)
	}

	// L2: GUARD — code-guard (prevent premature code removal)
	if input.ToolName == "Edit" {
		runCodeGuardCheck(input)
	}

	// L2: RESEARCH — TABULA_RASA enforcement
	runResearchCheck(input, session)

	// Check write blocked paths
	filePath := input.GetString("file_path")
	if filePath != "" && config.IsBlockedWritePath(filePath) {
		hook.ExitBlockTOON("ENFORCER", "Write:blocked_path:"+filePath)
	}

	hook.ExitSilent()
}

// runSecurityChain runs the multi-agent verification chain.
// Returns (blocked, reason, context).
func runSecurityChain(input *hook.Input, session *enforce.SessionState) (bool, string, string) {
	prompt := getPromptFromInput(input)
	runner := chain.NewRunner(session.ID)
	state := runner.RunFull(prompt, input.ToolName, input.ToolInput, session.ResearchDone)

	if state.IsBlocked() {
		return true, state.GetBlockReason(), runner.ToTOON()
	}
	return false, "", ""
}

// runContentCheck checks for secrets and credentials in content.
func runContentCheck(input *hook.Input) (bool, string) {
	content := input.GetString("content")
	if input.ToolName == "Edit" {
		content = input.GetString("new_string")
	}
	if content == "" {
		return false, ""
	}

	sensitivePatterns := []string{
		"password =", "secret =", "api_key =", "token =",
		"private_key", "BEGIN RSA PRIVATE", "BEGIN OPENSSH PRIVATE",
	}
	contentLower := strings.ToLower(content)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(contentLower, strings.ToLower(pattern)) {
			return true, "sensitive:" + pattern
		}
	}
	return false, ""
}

// runCodeGuardCheck checks Edit operations for premature code removal.
func runCodeGuardCheck(input *hook.Input) {
	oldString := input.GetString("old_string")
	newString := input.GetString("new_string")
	filePath := input.GetString("file_path")

	if !patterns.IsCodeFile(filePath) {
		return
	}

	// Function removal check
	removedFunctions := detectFunctionRemoval(oldString, newString)
	if len(removedFunctions) > 0 {
		if containsStubPatterns(oldString) {
			hook.ExitBlockTOON("CODE_GUARD",
				"BLOCK_REMOVAL:unimplemented_code:functions:"+strings.Join(removedFunctions, ",")+
					":REASON:Never remove TODO/stub functions without implementing them first")
		}
		if len(newString) < len(oldString)/2 {
			hook.ExitBlockTOON("CODE_GUARD",
				"BLOCK_REMOVAL:significant_code_reduction:functions:"+strings.Join(removedFunctions, ",")+
					":REASON:Verify use case before removing functions.")
		}
	}

	// Stub removal without implementation
	if containsStubPatterns(oldString) && !containsStubPatterns(newString) {
		if len(newString) <= len(oldString) {
			hook.ExitBlockTOON("CODE_GUARD",
				"BLOCK_REMOVAL:stub_removed_without_implementation:file:"+filePath+
					":REASON:TODO/FIXME removed but code not expanded.")
		}
	}

	// Complete deletion
	if strings.TrimSpace(newString) == "" && len(oldString) > 50 {
		hook.ExitBlockTOON("CODE_GUARD",
			"BLOCK_REMOVAL:complete_deletion:file:"+filePath+
				":REASON:Cannot delete significant code block.")
	}

	// Rust impl block removal
	if strings.Contains(oldString, "impl ") && !strings.Contains(newString, "impl ") {
		hook.ExitBlockTOON("CODE_GUARD",
			"BLOCK_REMOVAL:impl_block:file:"+filePath+
				":REASON:Cannot remove impl block without understanding trait implementation.")
	}
}

// runResearchCheck enforces TABULA_RASA (research before code).
func runResearchCheck(input *hook.Input, session *enforce.SessionState) {
	filePath := input.GetString("file_path")
	if filePath == "" {
		return
	}
	if !patterns.IsCodeFile(filePath) || session.ResearchDone {
		return
	}

	rg := agentic.NewResearchGate()
	query := rg.BuildSearchQuery("implementation patterns")
	hook.ExitBlockTOON("TABULA_RASA",
		"WebSearch_required_before_code:suggest:"+query)
}
