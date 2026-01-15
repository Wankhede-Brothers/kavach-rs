package session

import (
	"fmt"
	"os"

	"github.com/claude/shared/pkg/enforce"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate session state",
	Long: `[VALIDATE]
desc: Validate current session state and enforcement compliance
purpose: Health check for session integrity

[CHECKS]
date:        Session date matches today
memory_bank: Memory bank accessible
GOVERNANCE:  Governance file loaded
research:    Research state tracked

[USAGE]
kavach session validate

[OUTPUT]
[VALIDATE] date, session ID
[CHECKS]   Each check with status
[RESULT]   PASS or FAIL with error count

[EXIT_CODE]
0: All checks passed
1: One or more checks failed`,
	Run: runValidateCmd,
}

func runValidateCmd(cmd *cobra.Command, args []string) {
	ctx := enforce.NewContext()
	session := enforce.GetOrCreateSession()

	fmt.Println("[VALIDATE]")
	fmt.Println("date: " + ctx.Today)
	fmt.Println("session: " + session.ID)
	fmt.Println()

	errors := 0

	fmt.Println("[CHECKS]")

	// Check date matches
	if session.Today != ctx.Today {
		fmt.Println("date: EXPIRED")
		errors++
	} else {
		fmt.Println("date: valid")
	}

	// Check memory bank accessible
	if ctx.MemoryBank != nil {
		fmt.Println("memory_bank: accessible")
	} else {
		fmt.Println("memory_bank: ERROR")
		errors++
	}

	// Check governance loaded
	if ctx.Governance != nil {
		fmt.Println("GOVERNANCE: loaded")
	} else {
		fmt.Println("GOVERNANCE: missing")
	}

	// Check research state
	if session.ResearchDone {
		fmt.Println("research_done: true")
	} else {
		fmt.Println("research_done: false (WebSearch required before code)")
	}

	fmt.Println()
	fmt.Println("[RESULT]")
	if errors > 0 {
		fmt.Printf("status: FAIL\nerrors: %d\n", errors)
		os.Exit(1)
	}
	fmt.Println("status: PASS")
}
