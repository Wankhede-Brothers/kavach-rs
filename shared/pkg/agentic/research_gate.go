// Package agentic provides research gate enforcement.
package agentic

import (
	"regexp"
	"strings"
	"time"

	"github.com/claude/shared/pkg/config"
)

// ResearchRequirement defines what research is needed.
type ResearchRequirement struct {
	Topic     string   // What to research
	Keywords  []string // Search keywords
	URLs      []string // Specific URLs to fetch
	Mandatory bool     // Must complete before proceeding
	Reason    string   // Why this research is needed
}

// ResearchGate enforces research-before-implementation pattern.
// CORE PRINCIPLE: NEVER trust training weights.
type ResearchGate struct {
	todayDate    string
	currentYear  string
	requirements []ResearchRequirement
}

// NewResearchGate creates a gate with current date injection.
func NewResearchGate() *ResearchGate {
	now := time.Now()
	return &ResearchGate{
		todayDate:   now.Format("2006-01-02"),
		currentYear: now.Format("2006"),
	}
}

// Today returns the current date (injected, not hardcoded).
func (g *ResearchGate) Today() string {
	return g.todayDate
}

// Year returns the current year.
func (g *ResearchGate) Year() string {
	return g.currentYear
}

// RequireResearch checks if a task needs research first.
func (g *ResearchGate) RequireResearch(task string) *ResearchRequirement {
	taskLower := strings.ToLower(task)

	// Framework/library tasks ALWAYS need research (loaded from frameworks.toon)
	frameworkPatterns := loadAllFrameworkPatterns()

	for _, fw := range frameworkPatterns {
		if strings.Contains(taskLower, fw) {
			return &ResearchRequirement{
				Topic:     fw,
				Keywords:  []string{fw + " " + g.currentYear + " best practices"},
				Mandatory: true,
				Reason:    "Framework patterns change frequently - research current version",
			}
		}
	}

	// API/syntax tasks need research
	if containsAny(taskLower, []string{"api", "syntax", "pattern", "implement"}) {
		return &ResearchRequirement{
			Topic:     "implementation",
			Keywords:  []string{"current best practices " + g.currentYear},
			Mandatory: true,
			Reason:    "Implementation patterns evolve - verify current approach",
		}
	}

	return nil
}

// BuildSearchQuery creates a research query with date injection.
func (g *ResearchGate) BuildSearchQuery(topic string) string {
	// Inject current year into search
	return topic + " " + g.currentYear + " documentation latest"
}

// ValidateResearchDone checks if research was performed.
func (g *ResearchGate) ValidateResearchDone(context string) bool {
	// Check for WebSearch/WebFetch patterns in context
	searchPatterns := []string{
		"WebSearch",
		"WebFetch",
		"researched",
		"documentation shows",
		"according to",
		"latest docs",
	}

	for _, pattern := range searchPatterns {
		if strings.Contains(context, pattern) {
			return true
		}
	}
	return false
}

// ForbiddenPhrases returns phrases that indicate skipped research.
var ForbiddenPhrases = []string{
	"Based on my knowledge",
	"I think",
	"I believe",
	"In my experience",
	"As I understand",
	"I recall",
	"From what I know",
}

// CheckForForbiddenPhrases scans for anti-patterns.
// P2 FIX #13: Now case-insensitive - "i think" matches "I think".
func (g *ResearchGate) CheckForForbiddenPhrases(text string) []string {
	var violations []string
	lowerText := strings.ToLower(text)
	for _, phrase := range ForbiddenPhrases {
		lowerPhrase := strings.ToLower(phrase)
		if strings.Contains(lowerText, lowerPhrase) {
			violations = append(violations, phrase)
		}
	}
	return violations
}

// ExtractFrameworkFromTask identifies frameworks mentioned in a task.
func ExtractFrameworkFromTask(task string) []string {
	fwPatterns := loadAllFrameworkPatterns()
	if len(fwPatterns) == 0 {
		return nil
	}
	escaped := make([]string, len(fwPatterns))
	for i, p := range fwPatterns {
		escaped[i] = regexp.QuoteMeta(p)
	}
	patterns := regexp.MustCompile(`(?i)(` + strings.Join(escaped, "|") + `)`)
	matches := patterns.FindAllString(task, -1)

	// Deduplicate
	seen := make(map[string]bool)
	var unique []string
	for _, m := range matches {
		lower := strings.ToLower(m)
		if !seen[lower] {
			seen[lower] = true
			unique = append(unique, lower)
		}
	}
	return unique
}

// defaultFrameworkPatterns is the fallback when frameworks.toon is not available.
var defaultFrameworkPatterns = []string{
	"axum", "tonic", "tokio", "react", "vue", "angular",
	"dioxus", "leptos", "yew", "astro", "tauri",
	"postgres", "sqlx", "diesel", "prisma",
	"terraform", "kubernetes", "docker",
}

// loadAllFrameworkPatterns flattens all sections from frameworks.toon into a single slice.
// Falls back to built-in defaults if config is empty.
func loadAllFrameworkPatterns() []string {
	sections := config.GetFrameworkPatterns()
	var all []string
	for _, items := range sections {
		all = append(all, items...)
	}
	if len(all) == 0 {
		return defaultFrameworkPatterns
	}
	return all
}

// Helper: check if string contains any of the patterns
func containsAny(s string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}
