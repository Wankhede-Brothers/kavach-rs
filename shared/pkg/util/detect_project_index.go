// Package util provides common utility functions.
// detect_project_index.go: Index.toon manipulation for project registration.
// DACE: Single responsibility - index file operations only.
package util

import (
	"os"
	"strings"
	"time"
)

// appendToIndex adds a project entry to index.toon.
func appendToIndex(project, wd string) error {
	indexPath := IndexPath()
	data, err := os.ReadFile(indexPath)
	if err != nil {
		data = []byte("# Memory Bank Index - SP/3.0\n\nPROJECTS[1]{id,path,stack,aliases}\n")
	}

	content := string(data)
	if strings.Contains(content, project+",") {
		return nil // Already registered
	}

	lines := strings.Split(content, "\n")
	var newLines []string
	projectsFound := false
	inserted := false

	for i, line := range lines {
		if strings.HasPrefix(line, "PROJECTS[") {
			projectsFound = true
			newLines = append(newLines, incrementProjectCount(line))
			continue
		}

		if projectsFound && !inserted && !strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t") && line != "" {
			stack := detectStack(wd)
			newEntry := "  " + project + "," + wd + "," + stack + ","
			newLines = append(newLines, newEntry)
			inserted = true
		}

		newLines = append(newLines, lines[i])
	}

	if projectsFound && !inserted {
		stack := detectStack(wd)
		newEntry := "  " + project + "," + wd + "," + stack + ","
		newLines = append(newLines, newEntry)
	}

	newContent := updateIndexDate(strings.Join(newLines, "\n"))
	return os.WriteFile(indexPath, []byte(newContent), 0644)
}

// incrementProjectCount increments the count in PROJECTS[N] line.
func incrementProjectCount(line string) string {
	start := strings.Index(line, "[") + 1
	end := strings.Index(line, "]")
	if start > 0 && end > start {
		countStr := line[start:end]
		count := 0
		for _, ch := range countStr {
			if ch >= '0' && ch <= '9' {
				count = count*10 + int(ch-'0')
			}
		}
		return strings.Replace(line, "["+countStr+"]", "["+itoa(count+1)+"]", 1)
	}
	return line
}

// updateIndexDate updates the "updated" date in index content.
func updateIndexDate(content string) string {
	today := time.Now().Format("2006-01-02")
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, "updated:") {
			lines[i] = "  updated: " + today
			break
		}
	}
	return strings.Join(lines, "\n")
}
