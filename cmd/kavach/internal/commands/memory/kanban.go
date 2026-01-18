// Package memory provides memory bank commands.
// kanban.go: Sprint/Kanban dashboard command entry point.
// DACE: Micro-modular - types, loader, output, helpers in separate files.
package memory

import (
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var kanbanStatusFlag bool
var kanbanProjectFlag string
var kanbanVisualFlag bool
var kanbanSutraFlag bool

var kanbanCmd = &cobra.Command{
	Use:   "kanban",
	Short: "Sprint/Kanban visual dashboard with Aegis-Guard verification",
	Long: `[KANBAN]
desc: Production pipeline dashboard with automated verification gates
path: ~/.local/shared/shared-ai/memory/kanban/<project>/kanban.toon
protocol: SP/1.0 (Sutra Protocol)

[PIPELINE_STAGES]
1. BACKLOG      → Tasks waiting to be started
2. IN_PROGRESS  → Currently being worked on by engineers
3. TESTING      → Aegis-Guard Stage 1: Lints, Warnings, Core Bugs
4. VERIFIED     → Aegis-Guard Stage 2: Algorithm, Dead Code, Suppressed
5. DONE         → Production Ready (PROMISE fulfilled)

[AEGIS_GUARD_TESTING]
checks: lint warnings, compiler warnings, core bugs, unit tests
action: If FAIL → Report to CEO → LOOP continues

[AEGIS_GUARD_VERIFIED]
checks: algorithm, hidden bugs, dead code, suppressed elements
action: If FAIL → Report to CEO → LOOP continues until PROMISE

[FLAGS]
--status:       Quick status overview (TOON format)
--visual:       Human-readable visual dashboard
--sutra:        Output in Sutra Protocol format for agents
--project, -p:  Specify project name

[USAGE]
kavach memory kanban                    # Default: Visual dashboard
kavach memory kanban --status           # Quick TOON status
kavach memory kanban --sutra            # Agent communication format
kavach memory kanban -p my-project      # Specific project`,
	Run: runKanbanCmd,
}

func init() {
	kanbanCmd.Flags().BoolVar(&kanbanStatusFlag, "status", false, "Show sprint status only (TOON)")
	kanbanCmd.Flags().BoolVar(&kanbanVisualFlag, "visual", false, "Human-readable visual dashboard")
	kanbanCmd.Flags().BoolVar(&kanbanSutraFlag, "sutra", false, "Sutra Protocol format for agents")
	kanbanCmd.Flags().StringVarP(&kanbanProjectFlag, "project", "p", "", "Project name")
}

func runKanbanCmd(cmd *cobra.Command, args []string) {
	kanbanDir := util.MemoryBankPath("kanban")

	// Determine project
	project := kanbanProjectFlag
	if project == "" {
		project = detectKanbanProject(kanbanDir)
	}

	// Load kanban data
	board := loadKanbanTOON(kanbanDir, project)

	// Output based on flags
	if kanbanSutraFlag {
		outputSutraProtocol(board)
	} else if kanbanStatusFlag {
		outputTOONStatus(board)
	} else {
		// Default: Visual Dashboard for humans
		outputVisualDashboard(board)
	}
}
