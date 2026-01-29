// Package gates provides hook gates for Claude Code.
// pretool.go: Pre-tool umbrella gate (PreToolUse for non-write tools).
// Routes by tool name to exactly ONE L3 gate.
package gates

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/claude/shared/pkg/agentic"
	"github.com/claude/shared/pkg/config"
	"github.com/claude/shared/pkg/context"
	"github.com/claude/shared/pkg/dag"
	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/patterns"
	"github.com/claude/shared/pkg/telemetry"
	"github.com/spf13/cobra"
)

var preToolHookMode bool

var preToolCmd = &cobra.Command{
	Use:   "pre-tool",
	Short: "Pre-tool umbrella gate (bash|read|ceo|skill|content|task|context)",
	Run:   runPreToolGate,
}

func init() {
	preToolCmd.Flags().BoolVar(&preToolHookMode, "hook", false, "Hook mode")
}

func runPreToolGate(cmd *cobra.Command, args []string) {
	if !preToolHookMode {
		cmd.Help()
		return
	}

	span := telemetry.StartSpan("pre-tool")
	defer span.End()

	input := hook.MustReadHookInput()
	span.SetTool(input.ToolName)

	switch input.ToolName {
	case "Bash":
		preToolBash(input)
	case "Read", "Glob", "Grep":
		preToolRead(input)
	case "Task":
		preToolCEO(input)
	case "Skill":
		preToolSkill(input)
	case "WebFetch":
		preToolContent(input)
	case "TaskCreate", "TaskUpdate", "TaskGet", "TaskList", "TaskOutput":
		preToolTask(input)
	case "AskUserQuestion":
		// Context tracking only — silent pass
		hook.ExitSilent()
	default:
		hook.ExitSilent()
	}
}

// preToolBash: bash command sanitization.
func preToolBash(input *hook.Input) {
	command := input.GetString("command")
	if command == "" {
		hook.ExitBlockTOON("BASH", "empty_command")
	}
	if config.IsBlockedCommand(command) {
		hook.ExitBlockTOON("BASH", "blocked_command")
	}
	if patterns.IsBlocked(command) {
		hook.ExitBlockTOON("BASH", "blocked_command")
	}

	// Phase 11c: Check forbidden phrases in Bash description
	desc := input.GetString("description")
	if desc != "" {
		rg := agentic.NewResearchGate()
		if violations := rg.CheckForForbiddenPhrases(desc); len(violations) > 0 {
			hook.ExitBlockTOON("RESEARCH_GATE", "forbidden_phrase_in_bash:"+violations[0])
		}
	}

	// Legacy CLI detection
	if legacy, rust, reason := detectLegacyCommand(command); legacy != "" {
		hook.ExitBlockTOON("RUST_CLI", "LEGACY_BLOCKED:"+legacy+":USE:"+rust+":"+reason)
	}

	// Sudo warning
	if strings.HasPrefix(strings.TrimSpace(command), "sudo") {
		hook.ExitModifyTOON("BASH", map[string]string{"warn": "sudo_detected"})
	}

	// Risky command warnings
	cfg := config.LoadGatesConfig()
	cmdLower := strings.ToLower(command)
	for _, warn := range cfg.Bash.WarnCommands {
		if strings.Contains(cmdLower, strings.ToLower(warn)) {
			hook.ExitModifyTOON("BASH", map[string]string{"warn": warn + "_detected"})
		}
	}

	hook.ExitSilent()
}

// preToolRead: sensitive path blocking for Read/Glob/Grep.
func preToolRead(input *hook.Input) {
	var filePath string
	switch input.ToolName {
	case "Read":
		filePath = input.GetString("file_path")
	case "Glob", "Grep":
		filePath = input.GetString("path")
	}

	if filePath == "" && input.ToolName == "Read" {
		hook.ExitBlockTOON("READ", "no_file_path")
	}
	if filePath == "" {
		hook.ExitSilent()
	}

	if config.IsBlockedPath(filePath) {
		hook.ExitBlockTOON("READ", "blocked_path")
	}
	if config.IsBlockedExtension(filePath) {
		hook.ExitBlockTOON("READ", "blocked_extension")
	}
	if patterns.IsSensitive(filePath) {
		hook.ExitBlockTOON("READ", "sensitive_file")
	}

	if config.IsWarnPath(filePath) {
		hook.ExitModifyTOON("READ", map[string]string{"warn": "may_contain_secrets"})
	}
	if patterns.IsLargeFile(filePath) {
		hook.ExitModifyTOON("READ", map[string]string{"warn": "large_file"})
	}

	// Hot-context: hint if this file was already read recently
	if input.ToolName == "Read" && context.WasFileRecentlyRead(filePath) {
		hook.ExitModifyTOON("HOT_CONTEXT", map[string]string{
			"file":   filePath,
			"status": "recently_read",
			"hint":   "file_already_in_context",
		})
	}

	hook.ExitSilent()
}

// preToolCEO: CEO orchestration for Task tool.
func preToolCEO(input *hook.Input) {
	subagentType := input.GetString("subagent_type")
	if subagentType == "" {
		hook.ExitBlockTOON("CEO", "Task_requires_subagent_type")
	}
	if !config.IsValidAgent(subagentType) {
		hook.ExitBlockTOON("CEO", "unknown_agent:"+subagentType)
	}

	prompt := input.GetString("prompt")

	// Phase 11c: Check forbidden phrases in Task prompts
	if prompt != "" {
		rg := agentic.NewResearchGate()
		if violations := rg.CheckForForbiddenPhrases(prompt); len(violations) > 0 {
			hook.ExitBlockTOON("RESEARCH_GATE", "forbidden_phrase_in_task:"+violations[0])
		}
	}

	skill := detectSkillFromConfig(prompt)
	today := time.Now().Format("2006-01-02")

	if config.IsEngineer(subagentType) {
		orchDirective := map[string]string{
			"agent":      subagentType,
			"date":       today,
			"cutoff":     "2025-01",
			"CEO_FLOW":   "DELEGATE->VERIFY->AEGIS",
			"AFTER_TASK": "Verify result meets requirements",
			"IF_FAIL":    "Re-delegate with specific feedback",
			"IF_PASS":    "Run kavach orch aegis for final verification",
		}
		if skill != "" {
			orchDirective["skill"] = skill
			orchDirective["inject"] = "Invoke " + skill + " for domain expertise"
		}

		session := enforce.GetOrCreateSession()
		if session.HasTask() {
			orchDirective["current_task"] = session.CurrentTask
		}

		// DAG decomposition (reuse existing logic from ceo.go)
		breakdown := extractBreakdown(prompt)
		agents := resolveAgents(session, subagentType)
		if len(breakdown) <= 1 && len(agents) > 1 {
			breakdown = autoDecompose(prompt, agents)
		}

		if len(breakdown) > 1 {
			runDAGSchedule(session, prompt, breakdown, agents, orchDirective)
		}

		hook.ExitModifyTOON("CEO_ORCHESTRATION", orchDirective)
	}

	hook.ExitSilent()
}

// preToolSkill: skill name validation.
func preToolSkill(input *hook.Input) {
	if input.ToolName != "Skill" {
		hook.ExitSilent()
	}

	skillName := input.GetString("skill")
	if skillName == "" {
		hook.ExitBlockTOON("SKILL", "no_skill_name")
	}

	skillName = strings.ToLower(skillName)
	validSkills := config.GetValidSkills()
	if !validSkills[skillName] {
		hook.ExitModifyTOON("SKILL_WARN", map[string]string{
			"skill":  skillName,
			"status": "unrecognized_but_allowed",
		})
	}

	hook.ExitModifyTOON("SKILL", map[string]string{
		"skill":  skillName,
		"status": "routed",
	})
}

// preToolContent: URL safety for WebFetch.
func preToolContent(input *hook.Input) {
	// WebFetch content validation — delegate to existing content logic
	content := input.GetString("content")
	if content == "" {
		hook.ExitSilent()
	}

	contentLower := strings.ToLower(content)
	sensitivePatterns := []string{
		"password =", "secret =", "api_key =", "token =",
		"private_key", "BEGIN RSA PRIVATE", "BEGIN OPENSSH PRIVATE",
	}
	for _, pattern := range sensitivePatterns {
		if strings.Contains(contentLower, strings.ToLower(pattern)) {
			hook.ExitBlockTOON("CONTENT", "sensitive:"+pattern)
		}
	}

	hook.ExitSilent()
}

// preToolTask: task management validation.
func preToolTask(input *hook.Input) {
	session := enforce.GetOrCreateSession()

	switch input.ToolName {
	case "TaskCreate":
		if input.GetString("subject") == "" {
			hook.ExitBlockTOON("TASK_GATE", "TaskCreate:missing_subject")
		}
		if input.GetString("description") == "" {
			hook.ExitBlockTOON("TASK_GATE", "TaskCreate:missing_description")
		}
		taskListID := getTaskListID()
		today := time.Now().Format("2006-01-02")
		hook.ExitModifyTOON("TASK_CREATE", map[string]string{
			"task_list_id": taskListID,
			"created_date": today,
			"session_id":   session.SessionID,
		})
	case "TaskUpdate":
		if input.GetString("taskId") == "" {
			hook.ExitBlockTOON("TASK_GATE", "TaskUpdate:missing_taskId")
		}
		validStatuses := []string{"pending", "in_progress", "completed", "deleted", ""}
		status := input.GetString("status")
		if status != "" && !contains(validStatuses, status) {
			hook.ExitBlockTOON("TASK_GATE", "TaskUpdate:invalid_status:"+status)
		}
	case "TaskGet":
		if input.GetString("taskId") == "" {
			hook.ExitBlockTOON("TASK_GATE", "TaskGet:missing_taskId")
		}
	case "TaskList":
		taskListID := getTaskListID()
		if taskListID != "" {
			hook.ExitModifyTOON("TASK_LIST_CONTEXT", map[string]string{
				"task_list_id":  taskListID,
				"multi_session": "true",
			})
		}
	case "TaskOutput":
		if input.GetString("task_id") == "" {
			hook.ExitBlockTOON("TASK_GATE", "TaskOutput:missing_task_id")
		}
	}

	hook.ExitSilent()
}

// runDAGSchedule handles DAG decomposition and scheduling.
func runDAGSchedule(session *enforce.SessionState, prompt string, breakdown, agents []string, orchDirective map[string]string) {
	nodes := dag.Decompose(breakdown, agents)
	state, err := dag.Schedule(session.SessionID, prompt, nodes)
	if err == nil {
		if err := dag.Save(state); err != nil {
			fmt.Fprintf(os.Stderr, "[CEO_DAG] Save error: %v\n", err)
		}
		directive := dag.BuildDirective(state)
		hook.ExitModifyTOONWithModule("CEO_DAG_DISPATCH", orchDirective, directive)
	}
}
