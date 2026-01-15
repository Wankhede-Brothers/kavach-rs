// Package patterns provides dynamic pattern loading from TOON config.
// types.go: Pattern type definitions.
// DACE: Single responsibility - types only.
package patterns

// PatternSet holds a named set of patterns loaded from config.
type PatternSet struct {
	Name     string
	Patterns []string
	Source   string
}

// AgentSet holds valid agent definitions.
type AgentSet struct {
	Level  string
	Agents []string
}

// Config holds all dynamic patterns.
type Config struct {
	Sensitive   []string
	Blocked     []string
	CodeExts    []string
	LargeExts   []string
	ValidAgents map[string][]string
	IntentWords map[string][]string
	LoadedFrom  string
}
