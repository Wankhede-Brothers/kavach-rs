// Package memory provides memory bank commands.
// kanban_output.go: Kanban output formatters (TOON, Sutra, Visual).
// DACE: Micro-modular split from kanban.go
package memory

import (
	"fmt"
	"strings"
	"time"
)

// outputSutraProtocol outputs in SP/3.0 format for agent communication
func outputSutraProtocol(board *KanbanBoard) {
	fmt.Println("[META]")
	fmt.Println("protocol: SP/3.0")
	fmt.Println("from: kavach/kanban")
	fmt.Println("to: CEO")
	fmt.Printf("date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("project: %s\n", board.Project)
	fmt.Println()

	counts := countByColumn(board)

	fmt.Println("[KANBAN_STATE]")
	fmt.Printf("backlog: %d\n", counts[ColBacklog])
	fmt.Printf("in_progress: %d\n", counts[ColInProgress])
	fmt.Printf("testing: %d\n", counts[ColTesting])
	fmt.Printf("verified: %d\n", counts[ColVerified])
	fmt.Printf("done: %d\n", counts[ColDone])
	fmt.Printf("loop_count: %d\n", board.LoopCount)
	fmt.Println()

	fmt.Println("[AEGIS_QUEUE]")
	testingTasks := getTasksByColumn(board, ColTesting)
	if len(testingTasks) > 0 {
		fmt.Println("stage: TESTING")
		fmt.Println("checks: lint,warnings,core_bugs,unit_tests")
		for _, t := range testingTasks {
			fmt.Printf("task: %s,%s,%s\n", t.ID, t.Title, t.AegisStatus)
		}
	}
	verifyTasks := getTasksByColumn(board, ColVerified)
	if len(verifyTasks) > 0 {
		fmt.Println("stage: VERIFIED")
		fmt.Println("checks: algorithm,dead_code,suppressed,hidden_bugs")
		for _, t := range verifyTasks {
			fmt.Printf("task: %s,%s,%s\n", t.ID, t.Title, t.AegisStatus)
		}
	}
	fmt.Println()

	failedTasks := getFailedTasks(board)
	if len(failedTasks) > 0 {
		fmt.Println("[AEGIS_FAILURES]")
		fmt.Println("action: REPORT_TO_CEO")
		fmt.Println("result: LOOP_CONTINUES")
		for _, t := range failedTasks {
			fmt.Printf("failed: %s,%s,lint:%d,warn:%d,bugs:%d\n",
				t.ID, t.Title, t.LintIssues, t.Warnings, t.CoreBugs)
		}
		fmt.Println()
	}

	total := counts[ColBacklog] + counts[ColInProgress] + counts[ColTesting] + counts[ColVerified] + counts[ColDone]
	progress := 0
	if total > 0 {
		progress = (counts[ColDone] * 100) / total
	}

	fmt.Println("[PROMISE]")
	if progress == 100 && len(failedTasks) == 0 {
		fmt.Println("status: PRODUCTION_READY")
		fmt.Println("signal: <promise>PRODUCTION_READY</promise>")
	} else {
		fmt.Println("status: IN_PROGRESS")
		fmt.Printf("progress: %d%%\n", progress)
		fmt.Println("signal: LOOP_CONTINUES")
	}
}

// outputTOONStatus outputs compact TOON format
func outputTOONStatus(board *KanbanBoard) {
	counts := countByColumn(board)
	priorities := countByPriority(board)

	total := counts[ColBacklog] + counts[ColInProgress] + counts[ColTesting] + counts[ColVerified] + counts[ColDone]
	progress := 0
	if total > 0 {
		progress = (counts[ColDone] * 100) / total
	}

	fmt.Println("[KANBAN]")
	fmt.Printf("project: %s\n", board.Project)
	fmt.Printf("updated: %s\n", board.Updated)
	fmt.Printf("loop_count: %d\n", board.LoopCount)
	fmt.Println()

	fmt.Println("[PIPELINE]")
	fmt.Printf("backlog: %d\n", counts[ColBacklog])
	fmt.Printf("in_progress: %d\n", counts[ColInProgress])
	fmt.Printf("testing: %d\n", counts[ColTesting])
	fmt.Printf("verified: %d\n", counts[ColVerified])
	fmt.Printf("done: %d\n", counts[ColDone])
	fmt.Println()

	fmt.Println("[PRIORITY]")
	fmt.Printf("critical: %d\n", priorities["critical"])
	fmt.Printf("high: %d\n", priorities["high"])
	fmt.Printf("medium: %d\n", priorities["medium"])
	fmt.Printf("low: %d\n", priorities["low"])
	fmt.Println()

	fmt.Println("[AEGIS]")
	failedCount := len(getFailedTasks(board))
	fmt.Printf("failed: %d\n", failedCount)
	fmt.Printf("action: %s\n", getAegisAction(failedCount))
	fmt.Println()

	fmt.Println("[PROGRESS]")
	bar := renderProgressBar(progress, 20)
	fmt.Printf("%s %d%% (%d/%d)\n", bar, progress, counts[ColDone], total)
}

// outputVisualDashboard outputs human-readable visual dashboard
func outputVisualDashboard(board *KanbanBoard) {
	counts := countByColumn(board)
	total := counts[ColBacklog] + counts[ColInProgress] + counts[ColTesting] + counts[ColVerified] + counts[ColDone]
	progress := 0
	if total > 0 {
		progress = (counts[ColDone] * 100) / total
	}

	// Header
	fmt.Println("+" + strings.Repeat("=", 79) + "+")
	fmt.Printf("|                    KANBAN DASHBOARD: %-38s |\n", board.Project)
	fmt.Println("+" + strings.Repeat("=", 79) + "+")
	fmt.Printf("|  Updated: %-20s  Loop Count: %-5d  Progress: %3d%%        |\n",
		board.Updated, board.LoopCount, progress)
	fmt.Println("+" + strings.Repeat("=", 79) + "+")
	fmt.Println()

	// Pipeline
	fmt.Println("+" + strings.Repeat("-", 79) + "+")
	fmt.Println("|                           PRODUCTION PIPELINE                                 |")
	fmt.Println("+----------+-----------+-----------+-----------+-----------+-------------------+")
	fmt.Println("| BACKLOG  |IN_PROGRESS|  TESTING  | VERIFIED  |   DONE    |      STATUS       |")
	fmt.Println("+----------+-----------+-----------+-----------+-----------+-------------------+")
	fmt.Printf("|   %3d    |    %3d    |    %3d    |    %3d    |    %3d    | %s |\n",
		counts[ColBacklog], counts[ColInProgress], counts[ColTesting],
		counts[ColVerified], counts[ColDone], getStatusIcon(progress))
	fmt.Println("+----------+-----------+-----------+-----------+-----------+-------------------+")
	fmt.Println()

	// Progress bar
	fmt.Println("+" + strings.Repeat("-", 79) + "+")
	fmt.Printf("|  Progress: %s %3d%%                            |\n",
		renderProgressBar(progress, 40), progress)
	fmt.Println("+" + strings.Repeat("-", 79) + "+")
	fmt.Println()

	// Aegis status
	failedTasks := getFailedTasks(board)
	testingTasks := getTasksByColumn(board, ColTesting)
	verifyTasks := getTasksByColumn(board, ColVerified)

	fmt.Println("+" + strings.Repeat("-", 79) + "+")
	fmt.Println("|                           AEGIS-GUARD STATUS                                  |")
	fmt.Println("+" + strings.Repeat("-", 79) + "+")

	if len(testingTasks) > 0 {
		fmt.Println("|  TESTING STAGE (Lint, Warnings, Core Bugs):                                  |")
		for _, t := range testingTasks {
			icon := getAegisIcon(t.AegisStatus)
			fmt.Printf("|    %s %-50s [%s]        |\n",
				icon, truncate(t.Title, 50), t.AegisStatus)
		}
	}

	if len(verifyTasks) > 0 {
		fmt.Println("|  VERIFIED STAGE (Algorithm, Dead Code, Suppressed):                          |")
		for _, t := range verifyTasks {
			icon := getAegisIcon(t.AegisStatus)
			fmt.Printf("|    %s %-50s [%s]        |\n",
				icon, truncate(t.Title, 50), t.AegisStatus)
		}
	}

	if len(failedTasks) > 0 {
		fmt.Println("+" + strings.Repeat("-", 79) + "+")
		fmt.Println("|  FAILURES REPORTED TO CEO - LOOP CONTINUES                                   |")
		for _, t := range failedTasks {
			fmt.Printf("|    X %-60s          |\n", truncate(t.Title, 60))
		}
	}
	fmt.Println("+" + strings.Repeat("-", 79) + "+")
	fmt.Println()

	// Production Promise
	fmt.Println("+" + strings.Repeat("-", 79) + "+")
	if progress == 100 && len(failedTasks) == 0 {
		fmt.Println("|                    [OK] PROMISE: PRODUCTION_READY                            |")
	} else {
		fmt.Println("|                    [..] PROMISE: IN_PROGRESS (LOOP CONTINUES)                |")
	}
	fmt.Println("+" + strings.Repeat("-", 79) + "+")
}
