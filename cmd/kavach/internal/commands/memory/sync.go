// Package memory provides memory bank commands.
// sync.go: Memory sync command entry point.
// DACE: Micro-modular - types, scratchpad, kanban, helpers in separate files.
package memory

import (
	"fmt"
	"strings"
	"time"

	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/stmlog"
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var syncHookMode bool
var syncTask string
var syncStatus string

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync task completion to Memory Bank",
	Long: `[SYNC]
desc: Sync tool events to Memory Bank (STM session log + scratchpad)
hook: PostToolUse:TaskCreate|TaskUpdate|Write|Edit|Bash|Task
purpose: Keep scratchpad and session log in sync with task progress

[TRIGGERS]
- PostToolUse:TodoWrite - When todo list is updated
- Manual: kavach memory sync --task "task name" --status completed

[UPDATES]
- kanban.toon: Move task between columns
- scratchpad.toon: Update current task state
- roadmap.toon: Update phase progress

[HOOK_MODE]
Reads tool_result from stdin (Claude Code hook format)
Parses TodoWrite output and syncs to Memory Bank

[USAGE]
kavach memory sync --hook < tool_result.json
kavach memory sync --task "Implement API" --status completed`,
	Run: runSyncCmd,
}

func init() {
	syncCmd.Flags().BoolVar(&syncHookMode, "hook", false, "Hook mode (reads from stdin)")
	syncCmd.Flags().StringVar(&syncTask, "task", "", "Task name to sync")
	syncCmd.Flags().StringVar(&syncStatus, "status", "", "Task status (pending, in_progress, completed)")
}

func runSyncCmd(cmd *cobra.Command, args []string) {
	// BUG FIX: Use exact matching for writes to prevent updating wrong project
	project := util.DetectProjectForWrite()
	today := time.Now().Format("2006-01-02")

	if syncHookMode {
		runSyncHookMode(project, today)
		return
	}

	// Manual mode
	if syncTask != "" && syncStatus != "" {
		updateScratchpadManual(project, today, syncTask, syncStatus)
		fmt.Printf("[SYNC] %s -> %s\n", syncTask, syncStatus)
		return
	}

	cmd.Help()
}

func runSyncHookMode(project, today string) {
	input := hook.MustReadHookInput()
	toolName := input.ToolName

	switch toolName {
	case "TaskCreate":
		syncTaskCreate(input, project, today)
	case "TaskUpdate":
		syncTaskUpdate(input, project, today)
	case "Write", "Edit":
		syncFileChange(input, project, today)
	case "Bash":
		syncBashResult(input, project, today)
	case "Task":
		syncAgentResult(input, project, today)
	default:
		// Unknown tool — silent pass
	}

	hook.ExitSilent()
}

func syncTaskCreate(input *hook.Input, project, today string) {
	subject := input.GetString("subject")
	description := input.GetString("description")
	if subject == "" {
		return
	}
	stmlog.AppendEvent(project,"task_created", subject, description)
}

func syncTaskUpdate(input *hook.Input, project, today string) {
	taskID := input.GetString("taskId")
	status := input.GetString("status")
	if taskID == "" || status == "" {
		return
	}
	subject := input.GetString("subject")
	stmlog.AppendEvent(project,"task_"+status, subject, taskID)

	// Update scratchpad and kanban on completion
	if status == "completed" && subject != "" {
		updateScratchpadManual(project, today, subject, "completed")
		UpdateKanbanTimestamp(project, today)
	}
	if status == "in_progress" && subject != "" {
		updateScratchpadManual(project, today, subject, "in_progress")
	}
}

func syncFileChange(input *hook.Input, project, today string) {
	filePath := input.GetString("file_path")
	if filePath == "" {
		return
	}
	// Track modified file in session state for task scoping
	session := enforce.GetOrCreateSession()
	session.AddFileModified(filePath)

	stmlog.AppendEvent(project,"file_"+input.ToolName, filePath, "")
}

func syncBashResult(input *hook.Input, project, today string) {
	command := input.GetString("command")
	if command == "" {
		return
	}
	// Phase 7b: Only log significant commands (builds, tests, deploys, git)
	for _, sig := range []string{"build", "test", "deploy", "cargo", "go ", "bun ", "npm ", "git commit", "git push", "git merge"} {
		if containsStr(command, sig) {
			stmlog.AppendEvent(project,"bash_"+sig, command, "")
			return
		}
	}
	// Skip sync for trivial bash commands (ls, cat, etc.)
}

func syncAgentResult(input *hook.Input, project, today string) {
	desc := input.GetString("description")
	agentType := input.GetString("subagent_type")
	if desc == "" {
		return
	}
	stmlog.AppendEvent(project,"agent_"+agentType, desc, "")
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && strings.Contains(strings.ToLower(s), sub)
}
