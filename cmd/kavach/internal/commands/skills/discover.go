// Package skills - discover.go
// DACE: Single responsibility - skill discovery from directories
package skills

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/util"
)

// Discover loads skills from all sources with proper precedence
func Discover() []*Skill {
	skills := make([]*Skill, 0)

	// 1. Load built-in skills
	skills = append(skills, GetBuiltin()...)

	// 2. Discover from ~/.claude/skills/ (global)
	globalDir := filepath.Join(util.HomeDir(), ".claude", "skills")
	skills = merge(skills, fromDir(globalDir))

	// 3. Discover from .claude/skills/ (project - highest precedence)
	projectDir := filepath.Join(util.WorkingDir(), ".claude", "skills")
	skills = merge(skills, fromDir(projectDir))

	return skills
}

// fromDir discovers SKILL.md files from skill directories
func fromDir(dir string) []*Skill {
	skills := make([]*Skill, 0)
	if !util.DirExists(dir) {
		return skills
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return skills
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		// Look for SKILL.md (case-sensitive for Claude Code)
		skillPath := filepath.Join(dir, e.Name(), "SKILL.md")
		if !util.FileExists(skillPath) {
			continue
		}

		skill := ParseFile(skillPath)
		if skill != nil {
			skills = append(skills, skill)
		}
	}

	return skills
}

// merge combines skills, discovered takes precedence
func merge(builtin, discovered []*Skill) []*Skill {
	result := make([]*Skill, 0, len(builtin)+len(discovered))
	seen := make(map[string]bool)

	for _, s := range discovered {
		result = append(result, s)
		seen[s.Name] = true
	}

	for _, s := range builtin {
		if !seen[s.Name] {
			result = append(result, s)
		}
	}

	return result
}

// Find locates a skill by name
func Find(skills []*Skill, name string) *Skill {
	name = strings.ToLower(name)
	for _, s := range skills {
		if strings.ToLower(s.Name) == name {
			return s
		}
	}
	return nil
}

// ByCategory groups skills by category
func ByCategory(skills []*Skill) map[string][]*Skill {
	result := make(map[string][]*Skill)
	for _, s := range skills {
		cat := s.Category
		if cat == "" {
			cat = "general"
		}
		result[cat] = append(result[cat], s)
	}
	return result
}
