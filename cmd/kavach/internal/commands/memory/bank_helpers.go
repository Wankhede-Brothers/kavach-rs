// Package memory provides memory bank commands.
// bank_helpers.go: Helper functions for memory bank operations.
// DACE: Micro-modular split from bank.go
package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/claude/shared/pkg/dsa"
	"github.com/claude/shared/pkg/util"
)

// DACE: LRU cache for memory bank file counts (avoids repeated fs scans)
// TTL: 5 minutes, Max: 100 entries
var fileCountCache = dsa.NewLRUCache[string, int](100, 5*time.Minute)

// detectCurrentProject returns the current project name
func detectCurrentProject() string {
	return util.DetectProject()
}

// countFilesInDir counts TOON files in a directory with LRU caching.
// DACE: Uses cache to avoid repeated filesystem scans (5 min TTL).
func countFilesInDir(dir string) int {
	// Check cache first
	if count, ok := fileCountCache.Get(dir); ok {
		return count
	}

	// Cache miss - scan directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".toon") {
			count++
		}
	}

	// Store in cache
	fileCountCache.Set(dir, count)
	return count
}

// countTOONEntries counts entries in a TOON file
func countTOONEntries(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}

	count := 0
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "ENTRIES[") {
			start := strings.Index(line, "[")
			end := strings.Index(line, "]")
			if start != -1 && end != -1 && end > start {
				fmt.Sscanf(line[start+1:end], "%d", &count)
			}
		}
	}
	return count
}

// countFilesRecursive counts TOON files recursively
func countFilesRecursive(dir string) int {
	count := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".toon" {
				count++
			}
		}
		return nil
	})
	return count
}
