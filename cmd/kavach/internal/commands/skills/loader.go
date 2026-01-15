// Package skills - loader.go
// DACE: Single responsibility - SKILL.md file parsing
package skills

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ParseFile parses a SKILL.md file with YAML frontmatter
func ParseFile(path string) *Skill {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Skill name from parent directory
	skill := &Skill{
		Name: filepath.Base(filepath.Dir(path)),
		Path: path,
	}

	scanner := bufio.NewScanner(file)
	inFrontmatter := false
	frontmatterDone := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				frontmatterDone = true
				inFrontmatter = false
				continue
			}
		}

		if inFrontmatter {
			parseFrontmatterLine(skill, line)
		}

		if frontmatterDone && skill.Description == "" {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				skill.Description = trimmed
			}
		}
	}

	return skill
}

func parseFrontmatterLine(skill *Skill, line string) {
	if strings.HasPrefix(line, "name:") {
		skill.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
	} else if strings.HasPrefix(line, "category:") {
		skill.Category = strings.TrimSpace(strings.TrimPrefix(line, "category:"))
	} else if strings.HasPrefix(line, "description:") {
		skill.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
	}
}
