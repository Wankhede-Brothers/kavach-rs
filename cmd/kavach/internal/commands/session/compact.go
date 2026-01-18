package session

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/toon"
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var compactCmd = &cobra.Command{
	Use:   "compact",
	Short: "Pre-compact session save",
	Long: `[COMPACT]
desc: Save state before context compaction
hook: PreCompact

[DACE:COMPLIANT]
- Saves to TOON files, outputs ~50 tokens
- Pointers to Memory Bank, not data
- Task state preserved in scratchpad.toon

[SAVES_TO]
- STM/session-state.toon
- STM/projects/{project}/scratchpad.toon

[USAGE]
kavach session compact`,
	Run: runCompactCmd,
}

func runCompactCmd(cmd *cobra.Command, args []string) {
	ctx := enforce.NewContext()
	session := enforce.GetOrCreateSession()
	bank := toon.NewMemoryBank()

	// Mark post-compact mode and save scratchpad
	session.MarkPostCompact()
	saveScratchpad(session)

	// DACE: Ultra-minimal output (~50 tokens)
	// All state is in Memory Bank files - just confirm and point
	fmt.Println("[COMPACT:DACE]")
	fmt.Printf("date: %s\n", ctx.Today)
	fmt.Printf("session: %s\n", session.ID)
	fmt.Printf("project: %s\n", session.Project)
	fmt.Printf("compact_count: %d\n", session.CompactCount)

	// Show task if exists
	if session.CurrentTask != "" {
		fmt.Printf("task_saved: %s\n", session.CurrentTask)
	}

	// DACE: Pointer to memory bank total
	total := 0
	for _, count := range bank.GetCategoryStats() {
		total += count
	}
	fmt.Printf("memory_docs: %d\n\n", total)

	fmt.Println("[POST_COMPACT]")
	fmt.Println("run: kavach session init")
	fmt.Println("This auto-restores context from Memory Bank")
}

// saveScratchpad saves current task state to project scratchpad.
func saveScratchpad(session *enforce.SessionState) {
	if session.Project == "" {
		return
	}

	projectDir := filepath.Join(util.STMPath(), "projects", session.Project)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return
	}

	scratchpadPath := filepath.Join(projectDir, "scratchpad.toon")
	f, err := os.Create(scratchpadPath)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintln(f, "# Project Scratchpad - SP/1.0")
	fmt.Fprintln(f, "# Preserved during compact")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "[SCRATCHPAD:%s]\n", session.Project)
	fmt.Fprintf(f, "workdir: %s\n", session.WorkDir)
	fmt.Fprintf(f, "updated: %s\n", session.Today)
	fmt.Fprintln(f)

	fmt.Fprintln(f, "[TASK]")
	if session.CurrentTask != "" {
		fmt.Fprintf(f, "intent: %s\n", session.CurrentTask)
		fmt.Fprintf(f, "status: %s\n", session.TaskStatus)
	} else {
		fmt.Fprintln(f, "intent: null")
		fmt.Fprintln(f, "status: idle")
	}
	fmt.Fprintln(f)

	fmt.Fprintln(f, "[CONTEXT]")
	if len(session.FilesModified) == 0 {
		fmt.Fprintln(f, "files[]:")
	} else {
		fmt.Fprintln(f, "files[]:")
		for _, file := range session.FilesModified {
			fmt.Fprintf(f, "  - %s\n", file)
		}
	}
}
