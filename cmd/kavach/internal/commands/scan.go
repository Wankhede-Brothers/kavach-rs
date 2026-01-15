// Package commands provides CLI commands for kavach.
// scan.go: DACE-compliant tree scanner for micro-modular discovery.
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "DACE tree scanner for micro-modular file discovery",
	Long: `[SCAN]
desc: Scan directory tree for DACE-compliant file discovery
principle: Find smallest relevant files without reading content

[DACE:BENEFITS]
- Shows structure without loading content
- Identifies files by size (lines)
- Flags violations (>100 lines)
- Helps LLM navigate micro-modular codebase

[USAGE]
kavach scan                    # Scan current directory
kavach scan ./pkg              # Scan specific path
kavach scan --depth 3          # Limit depth
kavach scan --type go          # Filter by extension`,
	Run: runScanCmd,
}

var (
	scanDepth   int
	scanType    string
	scanMinSize int
)

func init() {
	scanCmd.Flags().IntVar(&scanDepth, "depth", 5, "Max depth to scan")
	scanCmd.Flags().StringVar(&scanType, "type", "", "Filter by file extension (go, rs, ts)")
	scanCmd.Flags().IntVar(&scanMinSize, "min", 0, "Min lines to show")
	rootCmd.AddCommand(scanCmd)
}

// FileInfo holds scanned file metadata.
type FileInfo struct {
	Path   string
	Lines  int
	Depth  int
	Status string // ✓ ok, ⚠ warn, ✗ violation
}

func runScanCmd(cmd *cobra.Command, args []string) {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("[SCAN:DACE]")
	fmt.Printf("path: %s\n", absPath)
	fmt.Printf("depth: %d\n", scanDepth)
	if scanType != "" {
		fmt.Printf("filter: .%s\n", scanType)
	}
	fmt.Println()

	files := scanDirectory(absPath, 0)

	// Sort by path
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	// Group by directory
	printTree(files, absPath)

	// Summary
	var ok, warn, violation int
	for _, f := range files {
		switch f.Status {
		case "✓":
			ok++
		case "⚠":
			warn++
		case "✗":
			violation++
		}
	}

	fmt.Println()
	fmt.Println("[SUMMARY]")
	fmt.Printf("total: %d files\n", len(files))
	fmt.Printf("✓ ok (<50): %d\n", ok)
	fmt.Printf("⚠ warn (50-100): %d\n", warn)
	fmt.Printf("✗ violation (>100): %d\n", violation)
}

func scanDirectory(root string, depth int) []FileInfo {
	var files []FileInfo

	if depth > scanDepth {
		return files
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return files
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden and common ignore patterns
		if strings.HasPrefix(name, ".") ||
			name == "node_modules" ||
			name == "vendor" ||
			name == "target" ||
			name == "__pycache__" {
			continue
		}

		fullPath := filepath.Join(root, name)

		if entry.IsDir() {
			files = append(files, scanDirectory(fullPath, depth+1)...)
			continue
		}

		// Filter by extension if specified
		ext := strings.TrimPrefix(filepath.Ext(name), ".")
		if scanType != "" && ext != scanType {
			continue
		}

		// Only scan code files
		if !isCodeExtension(ext) {
			continue
		}

		lines := countFileLines(fullPath)
		if lines < scanMinSize {
			continue
		}

		status := "✓"
		if lines > 100 {
			status = "✗"
		} else if lines > 50 {
			status = "⚠"
		}

		files = append(files, FileInfo{
			Path:   fullPath,
			Lines:  lines,
			Depth:  depth,
			Status: status,
		})
	}

	return files
}

func printTree(files []FileInfo, root string) {
	fmt.Println("[TREE]")

	for _, f := range files {
		relPath, _ := filepath.Rel(root, f.Path)
		indent := strings.Count(relPath, string(os.PathSeparator))
		prefix := strings.Repeat("  ", indent)

		fmt.Printf("%s%s %s (%d lines)\n", prefix, f.Status, filepath.Base(f.Path), f.Lines)
	}
}

func countFileLines(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}

	lines := 1
	for _, b := range data {
		if b == '\n' {
			lines++
		}
	}
	return lines
}

func isCodeExtension(ext string) bool {
	codeExts := map[string]bool{
		"go": true, "rs": true, "ts": true, "tsx": true,
		"js": true, "jsx": true, "py": true, "rb": true,
		"java": true, "kt": true, "swift": true, "c": true,
		"cpp": true, "h": true, "hpp": true, "zig": true,
		"toml": true, "yaml": true, "yml": true, "json": true,
		"md": true, "toon": true,
	}
	return codeExts[ext]
}
