package memory

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/claude/shared/events"
	"github.com/claude/shared/lock"
	"github.com/claude/shared/pkg/toon"
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var writeCategory string
var writeKey string

var writeProject string

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write to memory bank",
	Long: `[WRITE]
desc: Write TOON-formatted content to memory bank (project-isolated)
hook: PostToolUse (persist learnings)
purpose: Store decisions, patterns, research for future reference

[FLAGS]
-c, --category: Target category (required)
-k, --key:      Entry filename without .toon (required)
-p, --project:  Project name (default: auto-detect from working dir)

[CATEGORIES]
decisions, kanban, patterns, proposals, research, roadmaps

[PROJECT_ISOLATION]
default: Uses DetectProject() from working directory
override: --project flag for explicit project
path: ~/.local/shared/shared-ai/memory/{category}/{project}/{key}.toon

[INPUT]
stdin: TOON-formatted content

[USAGE]
echo '[DECISION]...' | kavach memory write -c decisions -k D001
echo '[KANBAN]...' | kavach memory write -c kanban -k sprint-1 -p my-project

[OUTPUT]
Success: Path where file was saved (project-scoped)
Error:   Parse error or missing flags`,
	Run: runWriteCmd,
}

func init() {
	writeCmd.Flags().StringVarP(&writeCategory, "category", "c", "", "Memory category (decisions, patterns, etc.)")
	writeCmd.Flags().StringVarP(&writeKey, "key", "k", "", "Entry key/filename")
	writeCmd.Flags().StringVarP(&writeProject, "project", "p", "", "Project name (default: auto-detect)")
}

func runWriteCmd(cmd *cobra.Command, args []string) {
	if writeCategory == "" {
		fmt.Fprintln(os.Stderr, "Error: --category required")
		os.Exit(1)
	}

	if writeKey == "" {
		fmt.Fprintln(os.Stderr, "Error: --key required")
		os.Exit(1)
	}

	// Read content from stdin
	scanner := bufio.NewScanner(os.Stdin)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if content == "" {
		fmt.Fprintln(os.Stderr, "Error: no content provided on stdin")
		os.Exit(1)
	}

	// Parse as TOON
	parser := toon.NewParser()
	doc, err := parser.ParseString(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing TOON: %v\n", err)
		os.Exit(1)
	}

	// Determine project (flag override or auto-detect from working directory)
	// BUG FIX: Use exact matching for writes to prevent updating wrong project
	project := writeProject
	if project == "" {
		project = util.DetectProjectForWrite()
	}

	// Ensure project directory exists
	projectDir := util.MemoryBankPath(writeCategory) + "/" + project
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating project directory: %v\n", err)
		os.Exit(1)
	}

	// Save to memory bank with project isolation
	path := projectDir + "/" + writeKey + ".toon"
	bank := toon.NewMemoryBank()

	// DACE: Acquire file lock before write (prevents concurrent Memory Bank writes)
	lockMgr := lock.GetLockManager()
	if err := lockMgr.AcquireWithTimeout(path, 5*time.Second); err != nil {
		fmt.Fprintf(os.Stderr, "Error acquiring lock: %v\n", err)
		os.Exit(1)
	}
	defer lockMgr.Release(path)

	if err := bank.SaveFile(path, doc); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving: %v\n", err)
		os.Exit(1)
	}

	// DACE: Publish EventMemoryWrite for subscribers (session tracking, hooks)
	eventBus := events.GetEventBus()
	eventBus.Publish(events.EventMemoryWrite, "kavach", map[string]interface{}{
		"category": writeCategory,
		"key":      writeKey,
		"path":     path,
		"project":  util.DetectProject(),
	})

	fmt.Printf("Saved to %s\n", path)
}
