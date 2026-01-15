// Package memory provides memory bank commands.
// kanban_loader.go: Kanban board loading and parsing.
// DACE: Micro-modular split from kanban.go
package memory

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/claude/shared/pkg/util"
)

// detectKanbanProject finds the active project for kanban
func detectKanbanProject(kanbanDir string) string {
	project := util.DetectProject()

	projectPath := filepath.Join(kanbanDir, project)
	if util.FileExists(filepath.Join(projectPath, "kanban.toon")) {
		return project
	}

	entries, _ := os.ReadDir(kanbanDir)
	for _, e := range entries {
		if e.IsDir() && e.Name() != "global" && e.Name() != "TEMPLATE.toon" {
			return e.Name()
		}
	}

	return "global"
}

// loadKanbanTOON loads kanban board from TOON file
func loadKanbanTOON(kanbanDir, project string) *KanbanBoard {
	board := &KanbanBoard{
		Project: project,
		Phases:  make(map[int][]KanbanCard),
		Updated: time.Now().Format("2006-01-02"),
	}

	toonPath := filepath.Join(kanbanDir, project, "kanban.toon")
	if !util.FileExists(toonPath) {
		return board
	}

	file, err := os.Open(toonPath)
	if err != nil {
		return board
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentPhase := -1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "KANBAN:") {
			board.Project = strings.TrimPrefix(line, "KANBAN:")
			continue
		}

		if strings.HasPrefix(line, "workdir:") {
			board.WorkDir = strings.TrimSpace(strings.TrimPrefix(line, "workdir:"))
			continue
		}

		if strings.HasPrefix(line, "updated:") {
			board.Updated = strings.TrimSpace(strings.TrimPrefix(line, "updated:"))
			continue
		}

		if strings.HasPrefix(line, "loop_count:") {
			fmt.Sscanf(line, "loop_count:%d", &board.LoopCount)
			continue
		}

		if strings.HasPrefix(line, "PHASE_") && strings.Contains(line, "_CARDS") {
			parts := strings.Split(line, "_")
			if len(parts) >= 2 {
				fmt.Sscanf(parts[1], "%d", &currentPhase)
			}
			continue
		}

		if currentPhase >= 0 && (strings.HasPrefix(line, "p") || strings.HasPrefix(line, "P")) {
			card := parseKanbanCardLine(line)
			if card.ID != "" {
				board.Phases[currentPhase] = append(board.Phases[currentPhase], card)
			}
		}
	}

	return board
}

// parseKanbanCardLine parses a card line from TOON format
func parseKanbanCardLine(line string) KanbanCard {
	// Format: id,column,title,priority,type,aegis_status,lint,warn,bugs
	parts := strings.SplitN(line, ",", 10)
	if len(parts) < 4 {
		return KanbanCard{}
	}

	card := KanbanCard{
		ID:          strings.TrimSpace(parts[0]),
		Column:      strings.TrimSpace(parts[1]),
		Title:       strings.TrimSpace(parts[2]),
		Priority:    strings.TrimSpace(parts[3]),
		AegisStatus: VerifyPending,
	}

	if len(parts) >= 5 {
		card.Type = strings.TrimSpace(parts[4])
	}
	if len(parts) >= 6 {
		card.AegisStatus = strings.TrimSpace(parts[5])
	}
	if len(parts) >= 7 {
		fmt.Sscanf(parts[6], "%d", &card.LintIssues)
	}
	if len(parts) >= 8 {
		fmt.Sscanf(parts[7], "%d", &card.Warnings)
	}
	if len(parts) >= 9 {
		fmt.Sscanf(parts[8], "%d", &card.CoreBugs)
	}

	return card
}
