// Package memory provides memory bank commands.
// kanban_helpers.go: Helper functions for Kanban operations.
// DACE: Micro-modular split from kanban.go
package memory

import "strings"

// countByColumn counts cards by column
func countByColumn(board *KanbanBoard) map[string]int {
	counts := map[string]int{
		ColBacklog:    0,
		ColInProgress: 0,
		ColTesting:    0,
		ColVerified:   0,
		ColDone:       0,
	}

	for _, cards := range board.Phases {
		for _, c := range cards {
			counts[c.Column]++
		}
	}
	return counts
}

// countByPriority counts cards by priority
func countByPriority(board *KanbanBoard) map[string]int {
	counts := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
	}

	for _, cards := range board.Phases {
		for _, c := range cards {
			counts[c.Priority]++
		}
	}
	return counts
}

// getTasksByColumn returns tasks in a specific column
func getTasksByColumn(board *KanbanBoard, column string) []KanbanCard {
	var tasks []KanbanCard
	for _, cards := range board.Phases {
		for _, c := range cards {
			if c.Column == column {
				tasks = append(tasks, c)
			}
		}
	}
	return tasks
}

// getFailedTasks returns all failed verification tasks
func getFailedTasks(board *KanbanBoard) []KanbanCard {
	var failed []KanbanCard
	for _, cards := range board.Phases {
		for _, c := range cards {
			if c.AegisStatus == VerifyFailed {
				failed = append(failed, c)
			}
		}
	}
	return failed
}

// filterByColumn filters cards by column
func filterByColumn(cards []KanbanCard, column string) []KanbanCard {
	var result []KanbanCard
	for _, c := range cards {
		if c.Column == column {
			result = append(result, c)
		}
	}
	return result
}

// getStatusIcon returns status icon based on progress
func getStatusIcon(progress int) string {
	if progress == 100 {
		return "[OK] PRODUCTION "
	} else if progress >= 75 {
		return "[..] ALMOST READY"
	} else if progress >= 50 {
		return "[..] IN PROGRESS "
	}
	return "[  ] EARLY STAGE "
}

// getAegisIcon returns icon for aegis status
func getAegisIcon(status string) string {
	switch status {
	case VerifyPassed:
		return "[OK]"
	case VerifyFailed:
		return "[X]"
	case VerifyBlocked:
		return "[!]"
	default:
		return "[ ]"
	}
}

// getAegisAction returns action based on failures
func getAegisAction(failedCount int) string {
	if failedCount > 0 {
		return "REPORT_TO_CEO"
	}
	return "CONTINUE"
}

// renderProgressBar renders a text progress bar
func renderProgressBar(percent, width int) string {
	filled := (percent * width) / 100
	empty := width - filled
	return "[" + strings.Repeat("#", filled) + strings.Repeat("-", empty) + "]"
}

// truncate truncates string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
