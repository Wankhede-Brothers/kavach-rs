// Package gates provides hook gates for Claude Code.
// posttool.go: Post-tool umbrella gate (PostToolUse for non-write tools).
// Routes by tool name to 1-2 L3 gates.
package gates

import (
	"fmt"
	"os"

	"github.com/claude/shared/pkg/context"
	"github.com/claude/shared/pkg/dag"
	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/spf13/cobra"
)

var postToolHookMode bool

var postToolCmd = &cobra.Command{
	Use:   "post-tool",
	Short: "Post-tool umbrella gate (memory|context|research|task)",
	Run:   runPostToolGate,
}

func init() {
	postToolCmd.Flags().BoolVar(&postToolHookMode, "hook", false, "Hook mode")
}

func runPostToolGate(cmd *cobra.Command, args []string) {
	if !postToolHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()
	session := enforce.GetOrCreateSession()

	switch input.ToolName {
	case "Bash":
		// Memory sync only (handled externally)
		hook.ExitSilent()

	case "Read":
		filePath := input.GetString("file_path")
		if filePath != "" {
			context.TrackFileRead(filePath)
		}
		hook.ExitSilent()

	case "Glob":
		path := input.GetString("path")
		pattern := input.GetString("pattern")
		if path != "" {
			context.TrackFileRead(path)
		} else if pattern != "" {
			context.TrackFileRead("glob:" + pattern)
		}
		hook.ExitSilent()

	case "Grep":
		path := input.GetString("path")
		pattern := input.GetString("pattern")
		if path != "" {
			context.TrackFileRead(path)
		} else if pattern != "" {
			context.TrackFileRead("grep:" + pattern)
		}
		hook.ExitSilent()

	case "Task":
		agentType := input.GetString("subagent_type")
		if agentType != "" {
			context.TrackAgentCompletion(agentType)
		}
		hook.ExitSilent()

	case "WebSearch", "WebFetch":
		// Mark research done
		session.MarkResearchDone()
		hook.ExitSilent()

	case "TaskCreate":
		postToolTaskCreate(input, session)

	case "TaskUpdate":
		postToolTaskUpdate(input, session)

	case "TaskOutput":
		postToolTaskOutput(input)

	default:
		hook.ExitSilent()
	}
}

// postToolTaskCreate handles post-creation tracking.
func postToolTaskCreate(input *hook.Input, session *enforce.SessionState) {
	subject := input.GetString("subject")
	session.TasksCreated++
	session.SetCurrentTask(subject)
	session.Save()

	// DAG tracking
	if state, err := dag.Load(session.SessionID); err == nil {
		_, _, directive := dag.HandleTaskEvent(state, "TaskCreate", input.ToolInput)
		if err := dag.Save(state); err != nil {
			fmt.Fprintf(os.Stderr, "[TASK_DAG] Save error: %v\n", err)
		}
		if directive != "" {
			hook.ExitModifyTOONWithModule("TASK_CREATE_POST", map[string]string{
				"dag_directive": "active",
			}, directive)
		}
	}

	hook.ExitSilent()
}

// postToolTaskUpdate handles post-update tracking.
func postToolTaskUpdate(input *hook.Input, session *enforce.SessionState) {
	status := input.GetString("status")
	subject := input.GetString("subject")

	if status == "completed" || status == "deleted" {
		session.TasksCompleted++
		session.ClearTask()
	} else if status == "in_progress" && subject != "" {
		session.SetTask(subject, status)
	}
	session.Save()

	// DAG advancement
	if state, err := dag.Load(session.SessionID); err == nil {
		complete, needsAegis, directive := dag.HandleTaskEvent(state, "TaskUpdate", input.ToolInput)
		if err := dag.Save(state); err != nil {
			fmt.Fprintf(os.Stderr, "[TASK_DAG] Save error: %v\n", err)
		}
		if complete && needsAegis {
			hook.ExitModifyTOON("TASK_UPDATE_DAG_COMPLETE", map[string]string{
				"dag_status": "complete",
				"action":     "Run kavach orch aegis for final verification",
			})
		}
		if directive != "" {
			hook.ExitModifyTOONWithModule("TASK_UPDATE_DAG_ADVANCE", map[string]string{
				"dag_status": "advancing",
			}, directive)
		}
	}

	hook.ExitSilent()
}

// postToolTaskOutput handles task output retrieval tracking.
func postToolTaskOutput(input *hook.Input) {
	taskID := input.GetString("task_id")
	if taskID == "" {
		hook.ExitSilent()
	}

	health := GetTaskHealth()
	if issue := health.RecordTaskOutput(taskID, ""); issue != nil {
		hook.ExitModifyTOON("TASK_OUTPUT", map[string]string{
			"task_id":      taskID,
			"health_event": issue.IssueType,
			"health_desc":  issue.Description,
		})
	}

	hook.ExitSilent()
}
