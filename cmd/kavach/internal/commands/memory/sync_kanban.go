// Package memory provides memory bank commands.
// sync_kanban.go: Kanban sync operations.
// DACE: Micro-modular split from sync.go
package memory

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/util"
)

// UpdateKanbanTimestamp updates the kanban file's timestamp for a project.
func UpdateKanbanTimestamp(project, today string) {
	kanbanDir := filepath.Join(util.MemoryBankPath("kanban"), project)
	kanbanPath := filepath.Join(kanbanDir, "kanban.toon")

	data, err := os.ReadFile(kanbanPath)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "updated:") {
			lines[i] = "updated: " + today
			break
		}
	}
	os.WriteFile(kanbanPath, []byte(strings.Join(lines, "\n")), 0644)
}
