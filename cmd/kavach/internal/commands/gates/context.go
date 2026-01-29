// Package gates provides hook gates for Claude Code.
// context.go: Hot-context tracking gate.
// P3 FIX #16: Tracks files read for DACE token optimization.
// P1 FIX: Extended to track Read, Write, Edit, Task operations.
package gates

import (
	"github.com/claude/shared/pkg/context"
	"github.com/claude/shared/pkg/hook"
	"github.com/spf13/cobra"
)

var contextHookMode bool

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Hot-context tracking gate",
	Long: `[CONTEXT_GATE]
desc: Tracks file operations for DACE token optimization
hook: PostToolUse:Read,Write,Edit,Task
output: Updates hot-context.json with file metadata`,
	Run: runContextGate,
}

func init() {
	contextCmd.Flags().BoolVar(&contextHookMode, "hook", false, "Hook mode")
}

func runContextGate(cmd *cobra.Command, args []string) {
	if !contextHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()

	// Track all file operation tools (Read, Glob, Grep, Write, Edit, Task)
	switch input.ToolName {
	case "Read":
		filePath := input.GetString("file_path")
		if filePath != "" {
			context.TrackFileRead(filePath)
		}
	case "Glob":
		// Glob uses "path" field (optional, defaults to cwd)
		path := input.GetString("path")
		pattern := input.GetString("pattern")
		if path != "" {
			context.TrackFileRead(path)
		} else if pattern != "" {
			context.TrackFileRead("glob:" + pattern)
		}
	case "Grep":
		path := input.GetString("path")
		pattern := input.GetString("pattern")
		if path != "" {
			context.TrackFileRead(path)
		} else if pattern != "" {
			context.TrackFileRead("grep:" + pattern)
		}
	case "Write":
		filePath := input.GetString("file_path")
		if filePath != "" {
			context.TrackFileWrite(filePath)
		}
	case "Edit":
		filePath := input.GetString("file_path")
		if filePath != "" {
			context.TrackFileEdit(filePath)
		}
	case "Task":
		agentType := input.GetString("subagent_type")
		if agentType != "" {
			context.TrackAgentCompletion(agentType)
		}
	default:
		// Unknown tool, exit silently
	}

	hook.ExitSilent()
}
