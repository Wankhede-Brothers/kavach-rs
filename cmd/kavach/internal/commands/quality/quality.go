// Package quality provides standalone quality commands.
// quality.go: Code quality analysis with SP/1.0 TOON output.
// P3 FIX: Standalone quality command for direct CLI usage.
package quality

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
	qualityFormat  string
	qualityVerbose bool
)

// QualityCmd is the top-level quality command.
var QualityCmd = &cobra.Command{
	Use:   "quality [file|dir...]",
	Short: "Analyze code quality metrics",
	Long: `[QUALITY]
desc: Analyze code quality and complexity metrics
output: SP/1.0 TOON format

[METRICS]
lines:        Total lines of code
functions:    Function count and complexity
dace_score:   DACE micro-modular compliance (0-100)
ast_valid:    AST syntax validation
imports:      Import analysis

[USAGE]
kavach quality file.go          # Analyze single file
kavach quality src/             # Analyze directory
kavach quality --format=json    # JSON output
kavach quality --verbose        # Detailed output`,
	Run: runQuality,
}

func init() {
	QualityCmd.Flags().StringVar(&qualityFormat, "format", "toon", "Output format (toon, json)")
	QualityCmd.Flags().BoolVar(&qualityVerbose, "verbose", false, "Verbose output")
}

// QualityResult represents quality metrics for a file.
type QualityResult struct {
	File       string `json:"file"`
	Lines      int    `json:"lines"`
	Functions  int    `json:"functions"`
	Imports    int    `json:"imports"`
	DACEScore  int    `json:"dace_score"`
	ASTValid   bool   `json:"ast_valid"`
	ASTError   string `json:"ast_error,omitempty"`
	Complexity string `json:"complexity"`
}

func runQuality(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	var results []QualityResult

	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %v\n", arg, err)
			continue
		}

		if info.IsDir() {
			filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				ext := util.GetExtension(path)
				if isAnalyzable(ext) {
					result := analyzeFile(path)
					results = append(results, result)
				}
				return nil
			})
		} else {
			result := analyzeFile(arg)
			results = append(results, result)
		}
	}

	outputQualityResults(results)
}

func isAnalyzable(ext string) bool {
	analyzable := map[string]bool{
		".go": true, ".rs": true, ".ts": true, ".tsx": true,
		".js": true, ".jsx": true, ".py": true,
	}
	return analyzable[ext]
}

func analyzeFile(path string) QualityResult {
	result := QualityResult{
		File:     path,
		ASTValid: true,
	}

	content, err := os.ReadFile(path)
	if err != nil {
		result.ASTValid = false
		result.ASTError = err.Error()
		return result
	}

	lines := strings.Split(string(content), "\n")
	result.Lines = len(lines)

	ext := util.GetExtension(path)

	// Count functions (simple heuristic)
	result.Functions = countFunctions(string(content), ext)

	// Count imports
	result.Imports = countImports(string(content), ext)

	// AST validation
	switch ext {
	case ".go":
		if errMsg := validate.GoSyntax(string(content)); errMsg != "" {
			result.ASTValid = false
			result.ASTError = errMsg
		}
	case ".json":
		if errMsg := validate.JSONSyntax(string(content)); errMsg != "" {
			result.ASTValid = false
			result.ASTError = errMsg
		}
	}

	// Calculate DACE score (0-100)
	result.DACEScore = calculateDACEScore(result)

	// Determine complexity
	result.Complexity = determineComplexity(result)

	return result
}

func countFunctions(content, ext string) int {
	count := 0
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch ext {
		case ".go":
			if strings.HasPrefix(trimmed, "func ") {
				count++
			}
		case ".rs":
			if strings.HasPrefix(trimmed, "fn ") || strings.HasPrefix(trimmed, "pub fn ") {
				count++
			}
		case ".ts", ".tsx", ".js", ".jsx":
			if strings.Contains(trimmed, "function ") ||
				strings.Contains(trimmed, "=> {") ||
				strings.Contains(trimmed, "async ") {
				count++
			}
		case ".py":
			if strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ") {
				count++
			}
		}
	}

	return count
}

func countImports(content, ext string) int {
	count := 0
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch ext {
		case ".go":
			if strings.HasPrefix(trimmed, "import ") || trimmed == "import (" {
				count++
			}
			if strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"") {
				count++
			}
		case ".rs":
			if strings.HasPrefix(trimmed, "use ") {
				count++
			}
		case ".ts", ".tsx", ".js", ".jsx":
			if strings.HasPrefix(trimmed, "import ") || strings.HasPrefix(trimmed, "require(") {
				count++
			}
		case ".py":
			if strings.HasPrefix(trimmed, "import ") || strings.HasPrefix(trimmed, "from ") {
				count++
			}
		}
	}

	return count
}

func calculateDACEScore(r QualityResult) int {
	score := 100

	// Deduct for lines over 100 (DACE micro-modular)
	if r.Lines > 100 {
		deduct := (r.Lines - 100) / 10
		if deduct > 50 {
			deduct = 50
		}
		score -= deduct
	}

	// Deduct for too many functions (suggests file should be split)
	if r.Functions > 10 {
		deduct := (r.Functions - 10) * 2
		if deduct > 20 {
			deduct = 20
		}
		score -= deduct
	}

	// Deduct for AST errors
	if !r.ASTValid {
		score -= 30
	}

	if score < 0 {
		score = 0
	}

	return score
}

func determineComplexity(r QualityResult) string {
	if r.Lines <= 50 && r.Functions <= 5 {
		return "low"
	}
	if r.Lines <= 100 && r.Functions <= 10 {
		return "medium"
	}
	return "high"
}

func outputQualityResults(results []QualityResult) {
	if qualityFormat == "json" {
		fmt.Println("[")
		for i, r := range results {
			fmt.Printf("  {\"file\": %q, \"lines\": %d, \"dace_score\": %d, \"ast_valid\": %v}",
				r.File, r.Lines, r.DACEScore, r.ASTValid)
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
	totalLines := 0
	totalFunctions := 0
	avgDACE := 0
	for _, r := range results {
		totalLines += r.Lines
		totalFunctions += r.Functions
		avgDACE += r.DACEScore
	}
	if len(results) > 0 {
		avgDACE /= len(results)
	}

	fmt.Println("[QUALITY_SUMMARY]")
	fmt.Printf("files: %d\n", len(results))
	fmt.Printf("total_lines: %d\n", totalLines)
	fmt.Printf("total_functions: %d\n", totalFunctions)
	fmt.Printf("avg_dace_score: %d\n", avgDACE)
	fmt.Println()

	if qualityVerbose {
		for _, r := range results {
			fmt.Printf("[FILE:%s]\n", r.File)
			fmt.Printf("  lines: %d\n", r.Lines)
			fmt.Printf("  functions: %d\n", r.Functions)
			fmt.Printf("  imports: %d\n", r.Imports)
			fmt.Printf("  dace_score: %d\n", r.DACEScore)
			fmt.Printf("  complexity: %s\n", r.Complexity)
			if !r.ASTValid {
				fmt.Printf("  ast_error: %s\n", r.ASTError)
			}
			fmt.Println()
		}
	}
}
