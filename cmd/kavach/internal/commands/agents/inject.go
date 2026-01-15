// Package agents - inject.go
// DACE: Single responsibility - research/pattern injection from Memory Bank
package agents

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/util"
)

// InjectContext loads research and patterns from Memory Bank into agent
func InjectContext(agent *Agent) {
	memDir := util.MemoryDir()
	project := util.DetectProject()

	// Load research
	researchPath := filepath.Join(memDir, "research", project, "research.toon")
	if util.FileExists(researchPath) {
		agent.Research = loadEntries(researchPath, agent.Name, []string{"verified:", "finding:", "fact:"})
	}

	// Load patterns
	patternsPath := filepath.Join(memDir, "patterns", project, "patterns.toon")
	if util.FileExists(patternsPath) {
		agent.Patterns = loadEntries(patternsPath, agent.Name, []string{"pattern:", "solution:", "template:"})
	}
}

func loadEntries(path, agentName string, markers []string) []string {
	entries := make([]string, 0)
	data, err := os.ReadFile(path)
	if err != nil {
		return entries
	}

	lines := strings.Split(string(data), "\n")
	lowAgent := strings.ToLower(agentName)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		lowLine := strings.ToLower(trimmed)

		// Check if relevant to agent
		if strings.Contains(lowLine, lowAgent) {
			entries = append(entries, trimmed)
			continue
		}

		// Check for marker keywords
		for _, marker := range markers {
			if strings.Contains(lowLine, marker) {
				entries = append(entries, trimmed)
				break
			}
		}
	}

	return entries
}
