// Package memory provides memory bank commands.
// bank.go: Memory bank command entry point.
// DACE: Micro-modular - helpers and display in separate files.
package memory

import "github.com/spf13/cobra"

var bankStatusFlag bool
var bankScanFlag bool
var bankAllFlag bool

var bankCmd = &cobra.Command{
	Use:   "bank",
	Short: "Memory bank operations",
	Long: `[BANK]
desc: Query and manage TOON-based memory bank
path: ~/.local/shared/shared-ai/memory/
hook: SessionStart (load context)

[CATEGORIES]
decisions:  Past decisions with rationale
graph:      Knowledge graph connections
kanban:     Task/project tracking
patterns:   Code patterns and solutions
proposals:  Feature/design proposals
research:   Research findings (TABULA_RASA)
roadmaps:   Project roadmaps
STM:        Short-term memory (scratchpad)

[PROJECT_ISOLATION]
default: Active project + global ONLY (prevents context pollution)
detection: .git root > .claude/project.json > memory bank match > global
override: KAVACH_PROJECT env var for explicit project

[FLAGS]
--status:  Health check (file counts, existence)
--scan:    Reindex all categories
--all:     Show ALL projects (default: active project only)

[USAGE]
kavach memory bank           # Active project + global only
kavach memory bank --all     # All projects
kavach memory bank --status  # Health check
kavach memory bank --scan    # Reindex`,
	Run: runBankCmd,
}

func init() {
	bankCmd.Flags().BoolVar(&bankStatusFlag, "status", false, "Show memory bank health status")
	bankCmd.Flags().BoolVar(&bankScanFlag, "scan", false, "Reindex memory bank")
	bankCmd.Flags().BoolVar(&bankAllFlag, "all", false, "Show all projects (default: active project only)")
}

func runBankCmd(cmd *cobra.Command, args []string) {
	if bankStatusFlag {
		showMemoryStatus()
		return
	}

	if bankScanFlag {
		scanMemoryBank()
		return
	}

	if bankAllFlag {
		showAllProjectsMemory()
		return
	}

	// DEFAULT: Project-scoped memory (active project + global only)
	showProjectScopedMemory()
}
