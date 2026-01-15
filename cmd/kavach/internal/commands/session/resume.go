package session

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/toon"
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume session after compact (DACE-optimized)",
	Long: `[RESUME]
desc: Restore state from Memory Bank after compaction

[DACE:COMPLIANT]
- ~80 tokens output
- Loads from TOON files
- Pointers to commands, not data

[USAGE]
kavach session resume`,
	Run: runResumeCmd,
}

func runResumeCmd(cmd *cobra.Command, args []string) {
	ctx := enforce.NewContext()
	session := enforce.GetOrCreateSession()
	bank := toon.NewMemoryBank()

	wasPostCompact := session.IsPostCompact()
	if wasPostCompact {
		session.ClearPostCompact()
	}

	// DACE: Ultra-minimal output (~80 tokens)
	fmt.Println("[RESUME:DACE]")
	fmt.Printf("date: %s\n", ctx.Today)
	fmt.Printf("session: %s\n", session.ID)
	fmt.Printf("project: %s\n", session.Project)
	if wasPostCompact {
		fmt.Printf("compact_recovered: true\n")
	}
	fmt.Println()

	fmt.Println("[STATE]")
	fmt.Printf("research_done: %s | memory: %s | ceo: %s\n\n",
		boolStr(session.ResearchDone), boolStr(session.MemoryQueried), boolStr(session.CEOInvoked))

	fmt.Println("[ENFORCE]")
	fmt.Printf("TABULA_RASA: cutoff=%s | WebSearch BEFORE code\n", session.TrainingCutoff)
	fmt.Println("NO_AMNESIA: query memory bank")
	fmt.Println()

	// Check for task to continue
	taskFound := false
	if session.CurrentTask != "" {
		fmt.Printf("[TASK] %s | status: %s\n\n", session.CurrentTask, session.TaskStatus)
		taskFound = true
	} else {
		scratchpad := loadProjectScratchpad(session.Project)
		if scratchpad != nil {
			if taskBlock := scratchpad.Get("TASK"); taskBlock != nil {
				if intent := taskBlock.Get("intent"); intent != "" && intent != "null" {
					fmt.Printf("[TASK] %s | status: %s\n\n", intent, taskBlock.Get("status"))
					taskFound = true
				}
			}
		}
	}

	if !taskFound {
		fmt.Println("[TASK] none | Ask user for next task")
		fmt.Println()
	}

	// DACE: Pointer only
	total := 0
	for _, count := range bank.GetCategoryStats() {
		total += count
	}
	fmt.Printf("[MEMORY] %d docs | query: kavach memory bank\n", total)
}

// loadProjectScratchpad loads the project-specific scratchpad.toon.
func loadProjectScratchpad(project string) *toon.Document {
	if project == "" {
		return nil
	}

	scratchpadPath := filepath.Join(util.STMPath(), "projects", project, "scratchpad.toon")
	if !util.FileExists(scratchpadPath) {
		return nil
	}

	f, err := os.Open(scratchpadPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	parser := toon.NewParser()
	doc, _ := parser.Parse(f)
	return doc
}
