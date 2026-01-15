// Package gates provides hook gates for Claude Code.
// intent_nlu.go: NLU classification logic using config patterns.
// DACE: Micro-modular split from intent.go
package gates

import (
	"strings"

	"github.com/claude/shared/pkg/config"
)

// classifyIntentFromConfig performs NLU using patterns from config file
// P1 FIX #8: Uses priority-based matching (first match wins for intent type)
// NO HARDCODING - patterns loaded from config/nlu-patterns.toon
func classifyIntentFromConfig(prompt string) *IntentClassification {
	if isSimpleQuery(prompt) {
		return nil
	}

	patterns := config.GetNLUPatterns()
	intent := &IntentClassification{
		ResearchReq: true,
		Confidence:  "medium",
	}

	// P1 FIX #8: Intent type classification with priority (first match wins)
	// Priority order: debug > performance > refactor > research > implement
	// Using else-if chain prevents later matches from overwriting earlier ones
	if matchesConfigPatterns(prompt, patterns["DEBUG:PATTERNS"]) {
		intent.Type = "debug"
		intent.Skills = []string{"/debug-like-expert"}
		intent.Agent = "ceo"
		intent.SubAgents = []string{"research-director", "backend-engineer"}
		intent.Confidence = "high"
	} else if matchesConfigPatterns(prompt, patterns["PERFORMANCE:PATTERNS"]) {
		intent.Type = "optimize"
		intent.Skills = []string{"/dsa", "/arch"}
		intent.Agent = "ceo"
		intent.SubAgents = []string{"research-director", "backend-engineer"}
		intent.Confidence = "high"
	} else if matchesConfigPatterns(prompt, patterns["REFACTOR:PATTERNS"]) {
		intent.Type = "refactor"
		intent.Skills = []string{"/heal"}
		intent.Agent = "ceo"
		intent.SubAgents = []string{"backend-engineer", "aegis-guardian"}
	} else if matchesConfigPatterns(prompt, patterns["RESEARCH:PATTERNS"]) {
		intent.Type = "research"
		intent.Agent = "research-director"
		intent.Confidence = "high"
	} else if matchesConfigPatterns(prompt, patterns["DOCS:PATTERNS"]) {
		// BUG-002 FIX: Handle documentation intent
		intent.Type = "docs"
		intent.Agent = "research-director"
		intent.SubAgents = []string{"backend-engineer"}
		intent.Confidence = "medium"
	} else if matchesConfigPatterns(prompt, patterns["AUDIT:PATTERNS"]) {
		// BUG-002 FIX: Handle audit/review intent
		intent.Type = "audit"
		intent.Skills = []string{"/security", "/heal"}
		intent.Agent = "ceo"
		intent.SubAgents = []string{"security-engineer", "aegis-guardian"}
		intent.Confidence = "high"
	} else if matchesConfigPatterns(prompt, patterns["IMPLEMENT:PATTERNS"]) {
		intent.Type = "implement"
		intent.Agent = "ceo"
	}

	// Domain classification is additive (can have multiple domains)
	// These don't overwrite each other - they add skills/subagents
	// P2 FIX: Route to specialized engineers (security, database, devops)
	if matchesConfigPatterns(prompt, patterns["SECURITY:DOMAIN"]) {
		intent.Domain = "security"
		intent.Skills = appendUnique(intent.Skills, "/security")
		intent.SubAgents = appendUnique(intent.SubAgents, "security-engineer")
	}

	if matchesConfigPatterns(prompt, patterns["FRONTEND:DOMAIN"]) {
		if intent.Domain == "" {
			intent.Domain = "frontend"
		}
		intent.Skills = appendUnique(intent.Skills, "/frontend")
		intent.SubAgents = appendUnique(intent.SubAgents, "frontend-engineer")
	}

	if matchesConfigPatterns(prompt, patterns["DATABASE:DOMAIN"]) {
		if intent.Domain == "" {
			intent.Domain = "database"
		}
		intent.Skills = appendUnique(intent.Skills, "/sql")
		intent.SubAgents = appendUnique(intent.SubAgents, "database-engineer")
	}

	if matchesConfigPatterns(prompt, patterns["INFRA:DOMAIN"]) {
		if intent.Domain == "" {
			intent.Domain = "infrastructure"
		}
		intent.Skills = appendUnique(intent.Skills, "/cloud-infrastructure-mastery")
		intent.SubAgents = appendUnique(intent.SubAgents, "devops-engineer")
	}

	if matchesConfigPatterns(prompt, patterns["DEVOPS:DOMAIN"]) {
		if intent.Domain == "" {
			intent.Domain = "devops"
		}
		intent.Skills = appendUnique(intent.Skills, "/cloud-infrastructure-mastery")
		intent.SubAgents = appendUnique(intent.SubAgents, "devops-engineer")
	}

	if matchesConfigPatterns(prompt, patterns["TESTING:DOMAIN"]) {
		if intent.Domain == "" {
			intent.Domain = "testing"
		}
		intent.Skills = appendUnique(intent.Skills, "/testing")
		intent.SubAgents = appendUnique(intent.SubAgents, "aegis-guardian")
	}

	if matchesConfigPatterns(prompt, patterns["API:DOMAIN"]) {
		if intent.Domain == "" {
			intent.Domain = "backend"
		}
		intent.Skills = appendUnique(intent.Skills, "/api-design")
		intent.SubAgents = appendUnique(intent.SubAgents, "backend-engineer")
	}

	// TABULA_RASA: Never silently bypass research requirements
	// If no pattern matches but query is non-trivial, still require research
	if intent.Type == "" && intent.Domain == "" {
		// Default to "unclassified" - still requires research
		intent.Type = "unclassified"
		intent.Confidence = "low"
	}

	if intent.Agent == "" {
		intent.Agent = "ceo"
	}

	return intent
}

func matchesConfigPatterns(text string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(text, p) {
			return true
		}
	}
	return false
}
