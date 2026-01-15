// Package memory provides memory bank commands.
// sync_scratchpad.go: Scratchpad sync operations.
// DACE: Micro-modular split from sync.go
package memory

import (
	"fmt"
	"os"

	"github.com/claude/shared/pkg/util"
)

func updateScratchpad(project, today string, inProgress, completed []TodoItem) {
	// P0 FIX: Use util.EnsureScratchpadDir for consistent path handling
	scratchpadDir, err := util.EnsureScratchpadDir(project)
	if err != nil {
		return
	}

	scratchpadPath := util.ScratchpadPath(project)
	_ = scratchpadDir // suppress unused warning

	f, err := os.Create(scratchpadPath)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintln(f, "# Scratchpad - SP/3.0")
	fmt.Fprintln(f, "# Auto-synced from TodoWrite")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "[SCRATCHPAD:%s]\n", project)
	fmt.Fprintf(f, "updated: %s\n", today)
	fmt.Fprintln(f)

	fmt.Fprintln(f, "[CURRENT_TASK]")
	if len(inProgress) > 0 {
		fmt.Fprintf(f, "task: %s\n", inProgress[0].Content)
		fmt.Fprintf(f, "status: in_progress\n")
		fmt.Fprintf(f, "active: %s\n", inProgress[0].ActiveForm)
	} else {
		fmt.Fprintln(f, "task: null")
		fmt.Fprintln(f, "status: idle")
	}
	fmt.Fprintln(f)

	fmt.Fprintln(f, "[COMPLETED_THIS_SESSION]")
	if len(completed) == 0 {
		fmt.Fprintln(f, "tasks[]:")
	} else {
		fmt.Fprintln(f, "tasks[]:")
		for _, t := range completed {
			fmt.Fprintf(f, "  - %s\n", t.Content)
		}
	}
	fmt.Fprintln(f)

	fmt.Fprintln(f, "[PENDING]")
	// Will be filled by pending tasks
	fmt.Fprintln(f, "count: calculated_from_kanban")
}

func updateScratchpadManual(project, today, task, status string) {
	// P0 FIX: Use util.EnsureScratchpadDir for consistent path handling
	_, err := util.EnsureScratchpadDir(project)
	if err != nil {
		return
	}

	scratchpadPath := util.ScratchpadPath(project)

	f, err := os.Create(scratchpadPath)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintln(f, "# Scratchpad - SP/3.0")
	fmt.Fprintf(f, "[SCRATCHPAD:%s]\n", project)
	fmt.Fprintf(f, "updated: %s\n", today)
	fmt.Fprintln(f)
	fmt.Fprintln(f, "[CURRENT_TASK]")
	fmt.Fprintf(f, "task: %s\n", task)
	fmt.Fprintf(f, "status: %s\n", status)
}
