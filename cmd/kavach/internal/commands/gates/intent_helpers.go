// Package gates provides hook gates for Claude Code.
// intent_helpers.go: Helper functions for intent classification.
// DACE: Micro-modular split from intent.go
package gates

import (
	"strings"

	"github.com/claude/shared/pkg/config"
)

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

// isImplementationIntent returns true for intent types that require research.
// Called by: intent_output.go to decide whether to inject BLOCK:RESEARCH.
func isImplementationIntent(intentType string) bool {
	for _, t := range []string{"implement", "debug", "refactor", "optimize", "fix", "audit", "docs", "unclassified"} {
		if intentType == t {
			return true
		}
	}
	return false
}

// isDelegationRequired returns true for intent types that must go through CEO delegation.
// These are intents where the CEO should orchestrate sub-agents rather than Claude writing code directly.
func isDelegationRequired(intentType string) bool {
	for _, t := range []string{"implement", "debug", "refactor", "optimize", "audit"} {
		if intentType == t {
			return true
		}
	}
	return false
}

// containsTechnicalTerms checks if prompt contains terms needing current docs.
// Called by: intent NLU to boost research requirement confidence.
// Loaded dynamically from frameworks.toon — no hardcoded lists.
func containsTechnicalTerms(prompt string) bool {
	lower := strings.ToLower(prompt)
	sections := config.GetFrameworkPatterns()
	for _, items := range sections {
		for _, term := range items {
			if strings.Contains(lower, strings.ToLower(term)) {
				return true
			}
		}
	}
	return false
}
