// Package agents - loader.go
// DACE: Single responsibility - file parsing
package agents

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ParseFile parses an agent .md file with YAML frontmatter
func ParseFile(path string) *Agent {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	agent := &Agent{
		Name: strings.TrimSuffix(filepath.Base(path), ".md"),
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
			parseFrontmatterLine(agent, line)
		}

		if frontmatterDone && agent.Description == "" {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				agent.Description = trimmed
			}
		}
	}

	return agent
}

func parseFrontmatterLine(agent *Agent, line string) {
	if strings.HasPrefix(line, "name:") {
		agent.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
	} else if strings.HasPrefix(line, "model:") {
		agent.Model = strings.TrimSpace(strings.TrimPrefix(line, "model:"))
	} else if strings.HasPrefix(line, "description:") {
		agent.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
	} else if strings.HasPrefix(line, "level:") {
		levelStr := strings.TrimSpace(strings.TrimPrefix(line, "level:"))
		agent.Level = ParseLevel(levelStr)
	}
}
