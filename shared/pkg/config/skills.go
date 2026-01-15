// Package config provides dynamic configuration loading from TOON files.
// skills.go: Skill-specific config functions with priority support.
// P1 FIX #1: Dynamic skill priority loading (NO HARDCODING).
package config

import (
	"sort"
	"strconv"
	"strings"
)

// SkillConfig holds skill configuration with priority.
type SkillConfig struct {
	Name     string
	Priority int
	Keywords []string
}

// GetSkillsByPriority returns skills sorted by priority (lowest number = highest priority).
// P1 FIX #1: Loads from config/skill-patterns.toon instead of hardcoding.
func GetSkillsByPriority() []SkillConfig {
	patterns := GetSkillPatterns()
	var skills []SkillConfig

	// Extract skills with their priorities
	for section, values := range patterns {
		if !strings.HasPrefix(section, "SKILL:") {
			continue
		}

		skillName := strings.TrimPrefix(section, "SKILL:")
		skill := SkillConfig{
			Name:     skillName,
			Priority: 999, // Default low priority
			Keywords: []string{},
		}

		// Parse values for priority and keywords
		for _, v := range values {
			if strings.HasPrefix(v, "priority:") {
				priStr := strings.TrimSpace(strings.TrimPrefix(v, "priority:"))
				if pri, err := strconv.Atoi(priStr); err == nil {
					skill.Priority = pri
				}
			} else {
				// It's a keyword
				skill.Keywords = append(skill.Keywords, v)
			}
		}

		skills = append(skills, skill)
	}

	// Sort by priority (ascending - lower number = higher priority)
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Priority < skills[j].Priority
	})

	return skills
}

// GetSkillNames returns skill names sorted by priority.
// P1 FIX #1: Dynamic skill list from config.
func GetSkillNames() []string {
	skills := GetSkillsByPriority()
	names := make([]string, len(skills))
	for i, s := range skills {
		names[i] = s.Name
	}
	return names
}

// GetSkillKeywords returns keywords for a specific skill.
func GetSkillKeywords(skillName string) []string {
	patterns := GetSkillPatterns()
	section := "SKILL:" + skillName

	var keywords []string
	for _, v := range patterns[section] {
		if !strings.HasPrefix(v, "priority:") {
			keywords = append(keywords, v)
		}
	}
	return keywords
}
