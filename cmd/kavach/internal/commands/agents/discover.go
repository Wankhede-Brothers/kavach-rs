// Package agents - discover.go
// DACE: Single responsibility - agent discovery from directories
package agents

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/claude/shared/pkg/util"
)

// Discover loads agents from all sources with proper precedence
func Discover() []*Agent {
	agents := make([]*Agent, 0)

	// 1. Load built-in agents
	agents = append(agents, GetBuiltin()...)

	// 2. Discover from ~/.claude/agents/ (global)
	globalDir := filepath.Join(util.HomeDir(), ".claude", "agents")
	agents = merge(agents, fromDir(globalDir))

	// 3. Discover from .claude/agents/ (project - highest precedence)
	projectDir := filepath.Join(util.WorkingDir(), ".claude", "agents")
	agents = merge(agents, fromDir(projectDir))

	// Sort by level
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Level < agents[j].Level
	})

	return agents
}

// fromDir discovers agent .md files from a directory
func fromDir(dir string) []*Agent {
	agents := make([]*Agent, 0)
	if !util.DirExists(dir) {
		return agents
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return agents
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		agent := ParseFile(filepath.Join(dir, e.Name()))
		if agent != nil {
			agents = append(agents, agent)
		}
	}

	return agents
}

// merge combines agents, discovered takes precedence over builtin
func merge(builtin, discovered []*Agent) []*Agent {
	result := make([]*Agent, 0, len(builtin)+len(discovered))
	seen := make(map[string]bool)

	// Add discovered first (takes precedence)
	for _, a := range discovered {
		result = append(result, a)
		seen[a.Name] = true
	}

	// Add builtin if not overridden
	for _, a := range builtin {
		if !seen[a.Name] {
			result = append(result, a)
		}
	}

	return result
}

// Find locates an agent by name
func Find(agents []*Agent, name string) *Agent {
	for _, a := range agents {
		if a.Name == name {
			return a
		}
	}
	return nil
}
