// Package gates provides hook gates for Claude Code.
// postwrite.go: Post-write umbrella gate (PostToolUse:Write|Edit|NotebookEdit).
// Hierarchy: ANTIPROD(P0→P3) → QUALITY → LINT → CONTEXT → MEMORY
package gates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/context"
	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/patterns"
	"github.com/claude/shared/pkg/stmlog"
	"github.com/claude/shared/pkg/telemetry"
	"github.com/claude/shared/pkg/util"
	"github.com/claude/shared/pkg/validate"
	"github.com/spf13/cobra"
)

var postWriteHookMode bool

var postWriteCmd = &cobra.Command{
	Use:   "post-write",
	Short: "Post-write umbrella gate (antiprod → quality → lint → context → memory)",
	Run:   runPostWriteGate,
}

func init() {
	postWriteCmd.Flags().BoolVar(&postWriteHookMode, "hook", false, "Hook mode")
}

func runPostWriteGate(cmd *cobra.Command, args []string) {
	if !postWriteHookMode {
		cmd.Help()
		return
	}

	span := telemetry.StartSpan("post-write")
	defer span.End()

	input := hook.MustReadHookInput()
	span.SetTool(input.ToolName)
	session := enforce.GetOrCreateSession()
	span.SetSessionLoaded(true)

	filePath := input.GetString("file_path")
	content := input.GetString("content")
	if input.ToolName == "Edit" {
		content = input.GetString("new_string")
	}

	// L2: ANTIPROD — P0→P3 hierarchy (stops at first P-level hit)
	if content != "" && filePath != "" {
		if reason := runAntiProdCheck(filePath, content); reason != "" {
			hook.ExitBlockTOON("ANTIPROD", reason)
		}
	}

	// L2: QUALITY — AST, folder depth, line count (.go files only)
	if content != "" && filePath != "" {
		runQualityCheck(filePath, content, session)
	}

	// L2: LINT — whitespace, tabs
	if content != "" && filePath != "" {
		runLintCheck(filePath, content)
	}

	// L2: CONTEXT — hot-context tracking
	if filePath != "" {
		if input.ToolName == "Write" {
			context.TrackFileWrite(filePath)
		} else if input.ToolName == "Edit" {
			context.TrackFileEdit(filePath)
		}
	}

	// L2: MEMORY — inline STM sync (previously separate kavach memory sync --hook)
	if filePath != "" {
		session.AddFileModified(filePath)
		stmlog.AppendEvent("", "file_"+input.ToolName, filePath, "")
	}

	hook.ExitSilent()
}

// runAntiProdCheck runs P0→P3 anti-production checks.
// Returns block reason on first hit, empty string if clean.
func runAntiProdCheck(filePath, content string) string {
	results := patterns.DetectAntiProd(filePath, content)
	if len(results) == 0 {
		return ""
	}

	// Return the highest priority (lowest P-level) result
	r := results[0]
	return fmt.Sprintf("%s:%s\nINSTRUCTION: %s", r.Code, r.Match, r.Message)
}

// codeExtensions lists file types subject to DACE quality checks.
var codeExtensions = map[string]bool{
	".go": true, ".rs": true, ".ts": true, ".tsx": true,
	".js": true, ".jsx": true, ".py": true, ".astro": true,
}

// runQualityCheck runs code quality validation for all code files.
// DACE structural checks (line count, folder depth) apply universally.
// AST validation is language-specific (Go only for now).
func runQualityCheck(filePath, content string, session *enforce.SessionState) {
	ext := util.GetExtension(filePath)
	if !codeExtensions[ext] {
		return
	}

	wd := util.WorkingDir()
	if wd == "" {
		return
	}
	rel, err := filepath.Rel(wd, filePath)
	if err != nil || (len(rel) >= 2 && rel[:2] == "..") {
		return
	}

	// Folder depth (universal)
	depth := 0
	for _, c := range rel {
		if c == filepath.Separator {
			depth++
		}
	}
	if depth > 7 {
		hook.ExitBlockTOON("DACE", "folder_depth_exceeds_7:"+util.Itoa(depth))
	}

	// Go-specific AST validation
	if ext == ".go" {
		if errMsg := validate.GoSyntax(content); errMsg != "" {
			hook.ExitBlockTOON("QUALITY", "ast_go:"+errMsg)
		}
	}

	// Line count (universal)
	lineCount := util.CountLines(content)
	if lineCount > 100 {
		hook.ExitBlockTOON("DACE", "exceeds_100_lines:"+util.Itoa(lineCount))
	}

	session.AddFileModified(filePath)
}

// runLintCheck runs basic lint validation.
func runLintCheck(filePath, content string) {
	lines := strings.Split(content, "\n")
	var issues []string

	for i, line := range lines {
		if strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t") {
			issues = append(issues, "trailing_ws:"+util.Itoa(i+1))
		}
	}

	ext := util.GetExtension(filePath)
	if ext == ".go" {
		for i, line := range lines {
			if strings.HasPrefix(line, "    ") && !strings.HasPrefix(line, "\t") {
				issues = append(issues, "spaces:"+util.Itoa(i+1))
			}
		}
	}

	if len(issues) > 0 {
		max := 3
		if len(issues) < max {
			max = len(issues)
		}
		// Lint issues are warnings, not blocks — use stderr
		fmt.Fprintf(os.Stderr, "[LINT] %s\n", strings.Join(issues[:max], ","))
	}
}
