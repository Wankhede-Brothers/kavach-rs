// Package gates provides hook gates for Claude Code.
// dag.go: DAG cycle detection gate entry point.
// DACE: Micro-modular - validator logic in dag_validator.go
package gates

import (
	"fmt"
	"strings"

	"github.com/claude/shared/pkg/hook"
	"github.com/spf13/cobra"
)

var dagHookMode bool
var dagCheckPath string

var dagCmd = &cobra.Command{
	Use:   "dag",
	Short: "DAG cycle detection gate",
	Long: `[DAG_GATE]
desc: Detects cycles in agent delegation chains
hook: PreToolUse:Task
purpose: Prevent infinite loops in agent orchestration

[USAGE]
kavach gates dag --hook              # PreToolUse hook mode
kavach gates dag --check A,B,C,A     # Manual cycle check`,
	Run: runDAGGate,
}

func init() {
	dagCmd.Flags().BoolVar(&dagHookMode, "hook", false, "Hook mode")
	dagCmd.Flags().StringVar(&dagCheckPath, "check", "", "Check path for cycles")
}

func runDAGGate(cmd *cobra.Command, args []string) {
	if dagCheckPath != "" {
		runManualCheck()
		return
	}

	if !dagHookMode {
		cmd.Help()
		return
	}

	runHookMode()
}

func runManualCheck() {
	path := strings.Split(dagCheckPath, ",")
	for i := range path {
		path[i] = strings.TrimSpace(path[i])
	}

	valid, cycle := ValidatePath(path)
	if !valid {
		fmt.Printf("[DAG_CYCLE_DETECTED]\ncycle: %s\n", strings.Join(cycle, " -> "))
		return
	}
	fmt.Println("[DAG_VALID]\nstatus: no_cycle")
}

func runHookMode() {
	input := hook.MustReadHookInput()

	if input.ToolName != "Task" {
		hook.ExitApproveTOON("DAG")
	}

	subagentType := input.GetString("subagent_type")
	if subagentType == "" {
		hook.ExitApproveTOON("DAG")
	}

	chainStr := input.GetString("_delegation_chain")
	var chain []string
	if chainStr != "" {
		chain = strings.Split(chainStr, ",")
	}
	chain = append(chain, subagentType)

	valid, cycle := ValidatePath(chain)
	if !valid {
		hook.ExitBlockTOON("DAG_CYCLE", "cycle:"+strings.Join(cycle, "->"))
	}

	hook.ExitModifyTOON("DAG", map[string]string{
		"_delegation_chain": strings.Join(chain, ","),
		"status":            "validated",
	})
}
