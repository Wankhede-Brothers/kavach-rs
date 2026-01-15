// Package memory provides memory bank commands.
// sync_kanban.go: Kanban sync operations.
// DACE: Micro-modular split from sync.go
package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/util"
)

func updateKanban(project, today string, completed, inProgress, pending []TodoItem) {
	kanbanDir := util.MemoryBankPath("kanban")
	projectDir := filepath.Join(kanbanDir, project)
	os.MkdirAll(projectDir, 0755)

	kanbanPath := filepath.Join(projectDir, "kanban.toon")

	// Read existing kanban to preserve structure
	existingContent := ""
	if data, err := os.ReadFile(kanbanPath); err == nil {
		existingContent = string(data)
	}

	// If no existing kanban, create minimal one
	if existingContent == "" {
		f, err := os.Create(kanbanPath)
		if err != nil {
			return
		}
		defer f.Close()

		fmt.Fprintln(f, "# Kanban - SP/3.0")
		fmt.Fprintln(f, "# Auto-synced from TodoWrite")
		fmt.Fprintln(f)
		fmt.Fprintf(f, "KANBAN:%s\n", project)
		fmt.Fprintf(f, "updated: %s\n", today)
		fmt.Fprintf(f, "loop_count: 0\n")
		fmt.Fprintln(f)

		fmt.Fprintln(f, "PHASE_0_CARDS")

		// Add completed tasks to done
		for i, t := range completed {
			fmt.Fprintf(f, "P0-%d,done,%s,medium,task,passed,0,0,0\n", i+1, sanitizeTitle(t.Content))
		}

		// Add in_progress tasks
		for i, t := range inProgress {
			fmt.Fprintf(f, "P0-%d,in_progress,%s,high,task,pending,0,0,0\n", len(completed)+i+1, sanitizeTitle(t.Content))
		}

		// Add pending tasks to backlog
		for i, t := range pending {
			fmt.Fprintf(f, "P0-%d,backlog,%s,medium,task,pending,0,0,0\n", len(completed)+len(inProgress)+i+1, sanitizeTitle(t.Content))
		}
		return
	}

	// Update existing kanban - just update the timestamp for now
	// Full kanban update would require parsing and merging
	updatedContent := updateKanbanTimestamp(existingContent, today)
	os.WriteFile(kanbanPath, []byte(updatedContent), 0644)
}

func updateKanbanTimestamp(content, today string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "updated:") {
			lines[i] = "updated: " + today
			break
		}
	}
	return strings.Join(lines, "\n")
}
