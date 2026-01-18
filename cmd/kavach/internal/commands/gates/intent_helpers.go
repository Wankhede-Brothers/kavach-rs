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

// isImplementationIntent returns true for intent types that require research first.
// TABULA_RASA: ALL intents except pure greetings require research.
func isImplementationIntent(intentType string) bool {
	// ALL non-trivial intents require research - training weights are STALE
	requiresResearch := []string{
		"implement", "debug", "refactor", "optimize", "fix",
		"audit", "docs", "unclassified", "research",
	}
	for _, t := range requiresResearch {
		if intentType == t {
			return true
		}
	}
	return false
}

// containsTechnicalTerms checks if prompt contains terms that need current docs.
// These terms indicate user needs FRESH information, not stale training weights.
func containsTechnicalTerms(prompt string) bool {
	terms := []string{
		"version", "install", "config", "setup", "deploy",
		"error", "fix", "issue", "bug", "fail",
		"how to", "how do", "what is", "which",
		"bun", "node", "npm", "yarn", "pnpm",
		"docker", "kubernetes", "terraform", "cloudflare",
		"react", "vue", "angular", "nextjs", "astro",
		"rust", "go", "python", "java",
		"api", "endpoint", "route", "handler",
		"database", "postgres", "mysql", "redis",
		"env", "environment", "variable",
	}
	lower := strings.ToLower(prompt)
	for _, term := range terms {
		if strings.Contains(lower, term) {
			return true
		}
	}
	return false
}
