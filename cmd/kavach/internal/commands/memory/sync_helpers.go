// Package memory provides memory bank commands.
// sync_helpers.go: Helper functions for sync operations.
// DACE: Micro-modular split from sync.go
package memory

import "strings"

func categorizeTodos(todos []TodoItem) (completed, inProgress, pending []TodoItem) {
	for _, t := range todos {
		switch t.Status {
		case "completed":
			completed = append(completed, t)
		case "in_progress":
			inProgress = append(inProgress, t)
		default:
			pending = append(pending, t)
		}
	}
	return
}

func sanitizeTitle(title string) string {
	// Remove commas and newlines for CSV-like format
	title = strings.ReplaceAll(title, ",", ";")
	title = strings.ReplaceAll(title, "\n", " ")
	if len(title) > 60 {
		title = title[:57] + "..."
	}
	return title
}
