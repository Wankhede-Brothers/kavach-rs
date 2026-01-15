// Package agents - types.go
// DACE: Single responsibility - type definitions only
package agents

// Hierarchy levels (SP/3.0)
const (
	LevelNLU      = -1
	LevelCEO      = 0
	LevelEngineer = 1
	LevelReview   = 2
	LevelAegis    = 3
)

// Agent represents a loaded agent definition
type Agent struct {
	Name        string
	Level       int
	Model       string
	Description string
	Triggers    []string
	Tools       []string
	Path        string
	Research    []string
	Patterns    []string
}

// LevelName returns human-readable level name
func LevelName(level int) string {
	switch level {
	case LevelNLU:
		return "L-1 (NLU)"
	case LevelCEO:
		return "L0 (CEO)"
	case LevelEngineer:
		return "L1 (Engineers)"
	case LevelReview:
		return "L1.5 (Review)"
	case LevelAegis:
		return "L2 (Aegis)"
	default:
		return "Unknown"
	}
}

// ParseLevel converts string to level int
func ParseLevel(s string) int {
	switch s {
	case "-1", "nlu":
		return LevelNLU
	case "0", "ceo", "l0":
		return LevelCEO
	case "1", "engineer", "l1":
		return LevelEngineer
	case "1.5", "2", "review":
		return LevelReview
	case "3", "aegis", "l2":
		return LevelAegis
	default:
		return LevelEngineer
	}
}
