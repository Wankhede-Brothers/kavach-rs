// Package gates provides hook gates for Claude Code.
// intent_helpers.go: Helper functions for intent classification.
// DACE: Micro-modular split from intent.go
package gates

import "strings"

func isSimpleQuery(prompt string) bool {
	simple := []string{"hello", "hi", "hey", "thanks", "thank you", "bye", "yes", "no", "ok", "okay"}
	trimmed := strings.TrimSpace(prompt)
	for _, s := range simple {
		if trimmed == s {
			return true
		}
	}
	return false
}

func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

func isStatusQuery(prompt string) bool {
	triggers := []string{"status", "project status", "what is the status", "show status", "check status"}
	for _, t := range triggers {
		if strings.Contains(prompt, t) {
			return true
		}
	}
	return false
}

// isImplementationIntent returns true for intent types that require code changes.
// P1 FIX: These intents reset research state to enforce task-scoped research.
func isImplementationIntent(intentType string) bool {
	implementationTypes := []string{
		"implement", "debug", "refactor", "optimize", "fix",
		"audit", "docs", "unclassified",
	}
	for _, t := range implementationTypes {
		if intentType == t {
			return true
		}
	}
	return false
}
