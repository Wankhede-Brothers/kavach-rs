// Package skills - types.go
// DACE: Single responsibility - type definitions only
package skills

// Skill represents a loaded skill definition
type Skill struct {
	Name        string
	Category    string
	Description string
	Triggers    []string
	Commands    []string
	Path        string
	Research    []string
	Patterns    []string
}

// Categories for skill organization
const (
	CatGit      = "git"
	CatSession  = "session"
	CatResearch = "research"
	CatMemory   = "memory"
	CatBuild    = "build"
	CatTest     = "test"
)
