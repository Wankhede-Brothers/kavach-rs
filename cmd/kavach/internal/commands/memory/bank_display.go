// Package memory provides memory bank commands.
// bank_display.go: Display functions for memory bank.
// DACE: Micro-modular split from bank.go
package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/toon"
	"github.com/claude/shared/pkg/util"
)

// showMemoryStatus displays memory bank health status
func showMemoryStatus() {
	memDir := util.MemoryDir()
	project := detectCurrentProject()

	fmt.Println("[MEMORY_BANK]")
	fmt.Println("path: ~/.local/shared/shared-ai/memory/")
	fmt.Printf("project: %s\n", project)
	fmt.Println()

	fmt.Println("[CATEGORIES]")
	categories := []string{"decisions", "graph", "kanban", "patterns", "proposals", "research", "roadmaps", "STM"}
	for _, cat := range categories {
		path := filepath.Join(memDir, cat)
		count := countFilesRecursive(path)
		fmt.Printf("%s: %d\n", cat, count)
	}
	fmt.Println()

	fmt.Println("[PROJECT_FILES]")
	for _, cat := range []string{"decisions", "kanban", "patterns", "proposals", "roadmaps"} {
		projectFile := filepath.Join(memDir, cat, project, cat+".toon")
		if util.FileExists(projectFile) {
			fmt.Printf("%s: OK\n", cat)
		} else {
			fmt.Printf("%s: -\n", cat)
		}
	}
	fmt.Println()

	fmt.Println("[ROOT_FILES]")
	topFiles := []string{"GOVERNANCE.toon", "index.toon", "volatile.toon"}
	for _, f := range topFiles {
		path := filepath.Join(memDir, f)
		status := "OK"
		if !util.FileExists(path) {
			status = "-"
		}
		fmt.Printf("%s: %s\n", f, status)
	}
}

// showProjectScopedMemory shows memory for ACTIVE project + global only
func showProjectScopedMemory() {
	memDir := util.MemoryDir()
	project := detectCurrentProject()
	projectDir := util.GetProjectDir()

	fmt.Println("[MEMORY]")
	fmt.Println("path: ~/.local/shared/shared-ai/memory/")
	fmt.Printf("project: %s\n", project)
	fmt.Printf("workdir: %s\n", util.WorkingDir())
	if projectDir != "" {
		fmt.Printf("project_root: %s\n", projectDir)
	}
	fmt.Println("scope: PROJECT_ISOLATED (active + global only)")
	fmt.Println()

	fmt.Println("[DETECTION]")
	if os.Getenv("KAVACH_PROJECT") != "" {
		fmt.Println("method: KAVACH_PROJECT env var")
	} else if projectDir != "" {
		fmt.Println("method: .git root detection")
	} else if project != "global" {
		fmt.Println("method: memory bank match")
	} else {
		fmt.Println("method: fallback (global)")
	}
	fmt.Println()

	categories := []string{"decisions", "kanban", "patterns", "proposals", "research", "roadmaps"}
	fmt.Println("[PROJECT_DOCS]")
	projectTotal := 0
	for _, cat := range categories {
		projectPath := filepath.Join(memDir, cat, project)
		count := countFilesInDir(projectPath)
		projectTotal += count
		if count > 0 {
			fmt.Printf("%s: %d\n", cat, count)
		}
	}
	fmt.Printf("project_total: %d\n", projectTotal)
	fmt.Println()

	fmt.Println("[GLOBAL_DOCS]")
	globalTotal := 0
	for _, cat := range categories {
		globalPath := filepath.Join(memDir, cat, "global")
		count := countFilesInDir(globalPath)
		globalTotal += count
		if count > 0 {
			fmt.Printf("%s: %d\n", cat, count)
		}
	}
	fmt.Printf("global_total: %d\n", globalTotal)
	fmt.Println()

	stmPath := filepath.Join(memDir, "STM")
	stmCount := countFilesInDir(stmPath)
	fmt.Println("[STM]")
	fmt.Printf("files: %d\n", stmCount)
	fmt.Println()

	fmt.Println("[ROOT]")
	rootFiles := []string{"GOVERNANCE.toon", "index.toon", "volatile.toon"}
	for _, f := range rootFiles {
		path := filepath.Join(memDir, f)
		if util.FileExists(path) {
			fmt.Printf("%s: OK\n", strings.TrimSuffix(f, ".toon"))
		}
	}
	fmt.Println()

	fmt.Printf("[TOTAL] %d (project: %d, global: %d, stm: %d)\n",
		projectTotal+globalTotal+stmCount, projectTotal, globalTotal, stmCount)
}

// showAllProjectsMemory shows memory for ALL projects
func showAllProjectsMemory() {
	bank := toon.NewMemoryBank()
	stats := bank.GetCategoryStats()
	project := detectCurrentProject()
	memDir := util.MemoryDir()

	fmt.Println("[MEMORY]")
	fmt.Println("path: ~/.local/shared/shared-ai/memory/")
	fmt.Printf("active_project: %s\n", project)
	fmt.Println("scope: ALL_PROJECTS (--all flag)")
	fmt.Println()

	fmt.Println("[DOCS]")
	total := 0
	for _, cat := range bank.ListCategories() {
		count := stats[cat]
		total += count
		fmt.Printf("%s: %d\n", cat, count)
	}
	fmt.Printf("total: %d\n", total)
	fmt.Println()

	fmt.Println("[PROJECTS]")
	entries, _ := os.ReadDir(filepath.Join(memDir, "kanban"))
	for _, e := range entries {
		if e.IsDir() && e.Name() != "TEMPLATE.toon" {
			indicator := ""
			if e.Name() == project {
				indicator = " <- ACTIVE"
			} else if e.Name() == "global" {
				indicator = " (shared)"
			}
			fmt.Printf("- %s%s\n", e.Name(), indicator)
		}
	}
	fmt.Println()

	fmt.Println("[ROOT]")
	if gov, _ := bank.LoadGovernance(); gov != nil {
		fmt.Printf("GOVERNANCE: %d blocks\n", len(gov.Blocks))
	} else {
		fmt.Println("GOVERNANCE: -")
	}
	if idx, _ := bank.LoadIndex(); idx != nil {
		fmt.Printf("index: %d blocks\n", len(idx.Blocks))
	} else {
		fmt.Println("index: -")
	}
}

// scanMemoryBank reindexes the memory bank
func scanMemoryBank() {
	bank := toon.NewMemoryBank()
	stats := bank.GetCategoryStats()
	project := detectCurrentProject()

	fmt.Println("[SCAN]")
	fmt.Println("path: ~/.local/shared/shared-ai/memory/")
	fmt.Printf("project: %s\n", project)
	fmt.Println()

	fmt.Println("[INDEX]")
	total := 0
	for _, cat := range bank.ListCategories() {
		count := stats[cat]
		total += count
		fmt.Printf("%s: %d\n", cat, count)
	}
	fmt.Printf("total: %d\n", total)
	fmt.Println()

	memDir := util.MemoryDir()
	fmt.Println("[PROJECTS]")
	entries, _ := os.ReadDir(filepath.Join(memDir, "kanban"))
	for _, e := range entries {
		if e.IsDir() && e.Name() != "global" {
			indicator := ""
			if e.Name() == project {
				indicator = " (current)"
			}
			fmt.Printf("- %s%s\n", e.Name(), indicator)
		}
	}
}
