// Package gates provides hook gates for Claude Code.
// ceo_decompose.go: Auto-decomposition of complex intents into parallel DAG nodes.
// DACE: Micro-modular split from ceo.go
package gates

import (
	"strings"

	"github.com/claude/shared/pkg/enforce"
)

// resolveAgents returns the agent list from intent bridge (session state),
// falling back to the single subagentType from the Task call.
func resolveAgents(session *enforce.SessionState, fallback string) []string {
	if len(session.IntentSubAgents) > 0 {
		return session.IntentSubAgents
	}
	return []string{fallback}
}

// autoDecompose generates breakdown steps from a natural language prompt
// when intent has multiple SubAgents but no explicit numbered/bullet list.
// Research agents get a "Research {topic}" step; others get "Implement {topic}".
func autoDecompose(prompt string, agents []string) []string {
	topic := extractTopic(prompt)
	var steps []string
	for _, agent := range agents {
		if isResearchAgent(agent) {
			steps = append(steps, "Research "+topic)
		} else {
			steps = append(steps, "Implement "+topic)
		}
	}
	// Deduplicate identical steps
	return dedup(steps)
}

func isResearchAgent(agent string) bool {
	return strings.Contains(agent, "research")
}

// extractTopic pulls the core topic from a prompt (first 80 chars, trimmed).
func extractTopic(prompt string) string {
	// Strip common prefixes
	for _, prefix := range []string{"please ", "can you ", "i need to ", "help me "} {
		prompt = strings.TrimPrefix(strings.ToLower(prompt), prefix)
	}
	if len(prompt) > 80 {
		prompt = prompt[:80]
	}
	return strings.TrimSpace(prompt)
}

func dedup(items []string) []string {
	seen := make(map[string]bool, len(items))
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
