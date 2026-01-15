// Package memory provides memory bank commands.
// sync.go: Memory sync command entry point.
// DACE: Micro-modular - types, scratchpad, kanban, helpers in separate files.
package memory

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/claude/shared/pkg/hook"
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
desc: Sync task completion from TodoWrite to Memory Bank
hook: PostToolUse:TodoWrite
purpose: Keep kanban, scratchpad, roadmap in sync with task progress

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
	// Read from stdin (PostToolUse hook format)
	input := hook.MustReadHookInput()

	// Parse tool result
	toolResult := input.GetString("tool_result")
	if toolResult == "" {
		hook.ExitSilent()
		return
	}

	// Try to parse as TodoWrite result
	var result struct {
		Todos []TodoItem `json:"todos"`
	}

	// The tool_result might be a string or structured
	if err := json.Unmarshal([]byte(toolResult), &result); err != nil {
		// Try reading raw JSON from stdin
		scanner := bufio.NewScanner(os.Stdin)
		var rawInput string
		for scanner.Scan() {
			rawInput += scanner.Text()
		}
		if rawInput != "" {
			json.Unmarshal([]byte(rawInput), &result)
		}
	}

	if len(result.Todos) == 0 {
		hook.ExitSilent()
		return
	}

	// Sync todos to Memory Bank
	completed, inProgress, pending := categorizeTodos(result.Todos)

	// Update scratchpad with current task
	updateScratchpad(project, today, inProgress, completed)

	// Update kanban
	updateKanban(project, today, completed, inProgress, pending)

	// Output sync result
	fmt.Println("[SYNC:COMPLETE]")
	fmt.Printf("project: %s\n", project)
	fmt.Printf("completed: %d\n", len(completed))
	fmt.Printf("in_progress: %d\n", len(inProgress))
	fmt.Printf("pending: %d\n", len(pending))

	hook.ExitSilent()
}
