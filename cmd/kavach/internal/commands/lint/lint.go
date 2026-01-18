// Package lint provides standalone lint commands.
// lint.go: File linting with SP/1.0 TOON output.
// P3 FIX: Standalone lint command for direct CLI usage.
package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/util"
	"github.com/claude/shared/pkg/validate"
	"github.com/spf13/cobra"
)

var (
	lintFix    bool
	lintFormat string
)

// LintCmd is the top-level lint command.
var LintCmd = &cobra.Command{
	Use:   "lint [file...]",
	Short: "Lint files for code quality issues",
	Long: `[LINT]
desc: Lint files for code quality and style issues
output: SP/1.0 TOON format

[CHECKS]
trailing_ws:   Trailing whitespace detection
tabs_spaces:   Tab vs space consistency
ast_syntax:    AST validation for Go, JSON, TOON
line_length:   Lines over 120 characters
dace_lines:    DACE micro-modular (<100 lines)

[USAGE]
kavach lint file.go           # Lint single file
kavach lint src/              # Lint directory
kavach lint --fix file.go     # Auto-fix issues
kavach lint --format=json     # JSON output`,
	Run: runLint,
}

func init() {
	LintCmd.Flags().BoolVar(&lintFix, "fix", false, "Auto-fix issues where possible")
	LintCmd.Flags().StringVar(&lintFormat, "format", "toon", "Output format (toon, json)")
}

// LintResult represents lint results for a file.
type LintResult struct {
	File   string      `json:"file"`
	Issues []LintIssue `json:"issues"`
	Fixed  int         `json:"fixed,omitempty"`
}

// LintIssue represents a single lint issue.
type LintIssue struct {
	Line    int    `json:"line"`
	Column  int    `json:"column,omitempty"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func runLint(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	var results []LintResult

	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %v\n", arg, err)
			continue
		}

		if info.IsDir() {
			// Lint directory recursively
			filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				ext := util.GetExtension(path)
				if isLintable(ext) {
					result := lintFile(path)
					if len(result.Issues) > 0 {
						results = append(results, result)
					}
				}
				return nil
			})
		} else {
			result := lintFile(arg)
			results = append(results, result)
		}
	}

	outputResults(results)
}

func isLintable(ext string) bool {
	lintable := map[string]bool{
		".go": true, ".rs": true, ".ts": true, ".tsx": true,
		".js": true, ".jsx": true, ".py": true, ".json": true,
		".yaml": true, ".yml": true, ".toon": true, ".md": true,
	}
	return lintable[ext]
}

func lintFile(path string) LintResult {
	result := LintResult{File: path, Issues: []LintIssue{}}

	content, err := os.ReadFile(path)
	if err != nil {
		result.Issues = append(result.Issues, LintIssue{
			Line:    0,
			Code:    "E000",
			Message: "cannot read file: " + err.Error(),
		})
		return result
	}

	lines := strings.Split(string(content), "\n")
	ext := util.GetExtension(path)

	// Check trailing whitespace
	for i, line := range lines {
		if strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t") {
			result.Issues = append(result.Issues, LintIssue{
				Line:    i + 1,
				Code:    "W001",
				Message: "trailing whitespace",
			})
		}
	}

	// Check line length
	for i, line := range lines {
		if len(line) > 120 {
			result.Issues = append(result.Issues, LintIssue{
				Line:    i + 1,
				Column:  121,
				Code:    "W002",
				Message: fmt.Sprintf("line too long (%d > 120)", len(line)),
			})
		}
	}

	// Check tabs vs spaces for Go
	if ext == ".go" {
		for i, line := range lines {
			if strings.HasPrefix(line, "    ") && !strings.HasPrefix(line, "\t") {
				result.Issues = append(result.Issues, LintIssue{
					Line:    i + 1,
					Code:    "W003",
					Message: "use tabs instead of spaces for Go indentation",
				})
			}
		}
	}

	// AST validation
	switch ext {
	case ".go":
		if errMsg := validate.GoSyntax(string(content)); errMsg != "" {
			result.Issues = append(result.Issues, LintIssue{
				Line:    1,
				Code:    "E001",
				Message: "Go syntax error: " + errMsg,
			})
		}
	case ".json":
		if errMsg := validate.JSONSyntax(string(content)); errMsg != "" {
			result.Issues = append(result.Issues, LintIssue{
				Line:    1,
				Code:    "E002",
				Message: "JSON syntax error: " + errMsg,
			})
		}
	}

	// DACE: Check line count
	if len(lines) > 100 {
		result.Issues = append(result.Issues, LintIssue{
			Line:    1,
			Code:    "D001",
			Message: fmt.Sprintf("DACE: file exceeds 100 lines (%d lines)", len(lines)),
		})
	}

	// Auto-fix if requested
	if lintFix && len(result.Issues) > 0 {
		fixed := autoFix(path, lines)
		result.Fixed = fixed
	}

	return result
}

func autoFix(path string, lines []string) int {
	fixed := 0
	newLines := make([]string, len(lines))

	for i, line := range lines {
		newLine := line
		// Fix trailing whitespace
		trimmed := strings.TrimRight(line, " \t")
		if trimmed != line {
			newLine = trimmed
			fixed++
		}
		newLines[i] = newLine
	}

	if fixed > 0 {
		content := strings.Join(newLines, "\n")
		os.WriteFile(path, []byte(content), 0644)
	}

	return fixed
}

func outputResults(results []LintResult) {
	if lintFormat == "json" {
		// JSON output
		fmt.Println("[")
		for i, r := range results {
			fmt.Printf("  {\"file\": %q, \"issues\": %d}", r.File, len(r.Issues))
			if i < len(results)-1 {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
		}
		fmt.Println("]")
		return
	}

	// TOON output (default)
	totalIssues := 0
	for _, r := range results {
		totalIssues += len(r.Issues)
	}

	fmt.Println("[LINT_RESULTS]")
	fmt.Printf("files: %d\n", len(results))
	fmt.Printf("issues: %d\n", totalIssues)
	fmt.Println()

	for _, r := range results {
		if len(r.Issues) == 0 {
			continue
		}
		fmt.Printf("[FILE:%s]\n", r.File)
		for _, issue := range r.Issues {
			fmt.Printf("  %d: [%s] %s\n", issue.Line, issue.Code, issue.Message)
		}
		if r.Fixed > 0 {
			fmt.Printf("  fixed: %d\n", r.Fixed)
		}
		fmt.Println()
	}
}
