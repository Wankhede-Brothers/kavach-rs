// Package memory provides memory bank commands.
// sync_scratchpad.go: Scratchpad sync operations.
// DACE: Micro-modular split from sync.go
package memory

import (
	"fmt"
	"os"

	"github.com/claude/shared/pkg/util"
)

func updateScratchpadManual(project, today, task, status string) {
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

	fmt.Fprintln(f, "# Scratchpad - SP/1.0")
	fmt.Fprintf(f, "[SCRATCHPAD:%s]\n", project)
	fmt.Fprintf(f, "updated: %s\n", today)
	fmt.Fprintln(f)
	fmt.Fprintln(f, "[CURRENT_TASK]")
	fmt.Fprintf(f, "task: %s\n", task)
	fmt.Fprintf(f, "status: %s\n", status)
}
