// Package orch provides orchestration subcommands.
// verify.go: Post-completion verification before marking task done.
// DACE: Blocks task completion until Aegis passes.
package orch

import (
	"fmt"
	"os"
	"time"

	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Post-completion verification",
	Long: `[VERIFY]
desc: Run full verification pipeline before marking task done
purpose: Ensure Aegis passes before task moves to DONE column

[PIPELINE]
1. Run Aegis Stage 1 (TESTING)
2. Run Aegis Stage 2 (VERIFIED)
3. If PASS: Allow task to move to DONE
4. If FAIL: Block and report to CEO

[USAGE]
kavach orch verify              # Verify current project
kavach orch verify --task ID    # Verify specific task`,
	Run: runVerifyCmd,
}

var verifyTaskID string

func init() {
	verifyCmd.Flags().StringVar(&verifyTaskID, "task", "", "Task ID to verify")
}

func runVerifyCmd(cmd *cobra.Command, args []string) {
	project := util.DetectProject()
	today := time.Now().Format("2006-01-02")

	fmt.Println("[VERIFY:PIPELINE]")
	fmt.Printf("project: %s\n", project)
	fmt.Printf("date: %s\n", today)
	if verifyTaskID != "" {
		fmt.Printf("task: %s\n", verifyTaskID)
	}
	fmt.Println()

	// Run full verification
	result := runVerification(project)

	// Output result
	outputAegisResult(os.Stdout, project, today, result)

	// Update kanban if verification passed
	if result.Status == "passed" {
		fmt.Println("\n[ACTION]")
		fmt.Println("kanban: Task can move to DONE")
		fmt.Println("promise: <promise>PRODUCTION_READY</promise>")
	} else {
		fmt.Println("\n[ACTION]")
		fmt.Println("kanban: Task stays in current column")
		fmt.Println("report: CEO notified of failures")
		fmt.Println("loop: CONTINUES until verification passes")
	}
}
