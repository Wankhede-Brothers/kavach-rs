// Package skills - inject.go
// DACE: Single responsibility - research/pattern injection from Memory Bank
package skills

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/util"
)

// InjectContext loads research and patterns from Memory Bank into skill
func InjectContext(skill *Skill) {
	memDir := util.MemoryDir()
	project := util.DetectProject()

	// Load research
	researchPath := filepath.Join(memDir, "research", project, "research.toon")
	if util.FileExists(researchPath) {
		skill.Research = loadEntries(researchPath, skill.Name)
	}

	// Load patterns
	patternsPath := filepath.Join(memDir, "patterns", project, "patterns.toon")
	if util.FileExists(patternsPath) {
		skill.Patterns = loadEntries(patternsPath, skill.Name)
	}
}

func loadEntries(path, skillName string) []string {
	entries := make([]string, 0)
	data, err := os.ReadFile(path)
	if err != nil {
		return entries
	}

	lines := strings.Split(string(data), "\n")
	lowSkill := strings.ToLower(skillName)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		lowLine := strings.ToLower(trimmed)
		if strings.Contains(lowLine, lowSkill) {
			entries = append(entries, trimmed)
		}
	}

	return entries
}
