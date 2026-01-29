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
// Extracts distinct concerns per agent rather than trivial string concat.
func autoDecompose(prompt string, agents []string) []string {
	topic := extractTopic(prompt)
	concerns := extractConcerns(prompt)

	// If we found distinct concerns, assign them to agents
	if len(concerns) >= len(agents) {
		var steps []string
		for i, agent := range agents {
			verb := "Implement"
			if isResearchAgent(agent) {
				verb = "Research"
			}
			steps = append(steps, verb+" "+concerns[i%len(concerns)])
		}
		return dedup(steps)
	}

	// Fallback: research phase then implement phase (2 steps minimum)
	var steps []string
	hasResearch := false
	for _, agent := range agents {
		if isResearchAgent(agent) {
			hasResearch = true
			steps = append(steps, "Research "+topic+" patterns and dependencies")
		} else {
			steps = append(steps, "Implement "+topic)
		}
	}
	if !hasResearch && len(steps) > 0 {
		// Prepend a research step so DAG has parallel value
		steps = append([]string{"Research " + topic + " patterns and dependencies"}, steps...)
	}
	return dedup(steps)
}

// extractConcerns pulls distinct technical concerns from a prompt.
// Looks for domain keywords that indicate separate work streams.
func extractConcerns(prompt string) []string {
	lower := strings.ToLower(prompt)
	var found []string
	seen := make(map[string]bool)
	domains := []struct {
		keyword, concern string
	}{
		{"auth", "authentication"},
		{"login", "authentication"},
		{"payment", "payment integration"},
		{"database", "database schema"},
		{"api", "API endpoints"},
		{"frontend", "frontend UI"},
		{"backend", "backend logic"},
		{"test", "test coverage"},
		{"deploy", "deployment"},
		{"security", "security hardening"},
		{"feed", "content feed"},
		{"search", "search functionality"},
	}
	for _, d := range domains {
		if strings.Contains(lower, d.keyword) && !seen[d.concern] {
			seen[d.concern] = true
			found = append(found, d.concern)
		}
	}
	return found
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
