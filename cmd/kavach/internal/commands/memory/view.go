// Package memory provides memory bank commands.
// view.go: View memory files with Rust CLI tools (bat, eza).
// DACE: Single responsibility - memory viewing only.
package memory

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/claude/shared/pkg/rust"
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var viewCategory string
var viewKey string
var viewTreeFlag bool

var viewCmd = &cobra.Command{
	Use:   "view [file]",
	Short: "View memory files with syntax highlighting",
	Long: `[VIEW]
desc: Display memory bank content with Rust CLI enhancements
tools: bat (syntax highlight), eza (tree view), fd (find)

[FLAGS]
--category, -c:  Category to view (decisions, kanban, etc.)
--key, -k:       Specific key/file within category
--tree:          Show directory tree (uses eza)

[EXAMPLES]
kavach memory view                     # Interactive tree view
kavach memory view --tree              # Full memory bank tree
kavach memory view -c kanban           # List kanban files
kavach memory view -c kanban -k status # View specific file

[RUST_CLI]
bat:  Syntax highlighting for TOON files
eza:  Directory tree with icons
fd:   Fast file search`,
	Run: runViewCmd,
}

func init() {
	viewCmd.Flags().StringVarP(&viewCategory, "category", "c", "", "Category to view")
	viewCmd.Flags().StringVarP(&viewKey, "key", "k", "", "Key/file to view")
	viewCmd.Flags().BoolVar(&viewTreeFlag, "tree", false, "Show directory tree")
}

func runViewCmd(cmd *cobra.Command, args []string) {
	tools := rust.Detect()
	memDir := util.MemoryDir()
	project := util.DetectProject()

	// Show available tools
	fmt.Println("[RUST_CLI]")
	fmt.Printf("bat: %s\n", toolStatus(tools.HasBat()))
	fmt.Printf("eza: %s\n", toolStatus(tools.HasEza()))
	fmt.Printf("fd: %s\n", toolStatus(tools.HasFd()))
	fmt.Printf("rg: %s\n", toolStatus(tools.HasRg()))
	fmt.Println()

	// Tree view
	if viewTreeFlag {
		showMemoryTree(tools, memDir)
		return
	}

	// Category + Key specified
	if viewCategory != "" && viewKey != "" {
		viewFile(tools, memDir, project, viewCategory, viewKey)
		return
	}

	// Category only - list files
	if viewCategory != "" {
		listCategoryFiles(tools, memDir, project, viewCategory)
		return
	}

	// Default: show tree if eza available, else list categories
	if tools.HasEza() {
		showMemoryTree(tools, memDir)
	} else {
		listCategories(memDir)
	}
}

func toolStatus(available bool) string {
	if available {
		return "available"
	}
	return "not installed"
}

func showMemoryTree(tools *rust.Tools, memDir string) {
	fmt.Println("[MEMORY_BANK_TREE]")
	if tools.HasEza() {
		tree, err := rust.LsTree(memDir, 3)
		if err == nil && tree != "" {
			fmt.Println(tree)
			return
		}
	}
	// Fallback: simple list
	listCategories(memDir)
}

func listCategories(memDir string) {
	fmt.Println("[CATEGORIES]")
	categories := []string{"decisions", "graph", "kanban", "patterns", "proposals", "research", "roadmaps", "STM"}
	for _, cat := range categories {
		path := filepath.Join(memDir, cat)
		if util.DirExists(path) {
			fmt.Printf("  %s/\n", cat)
		}
	}
	fmt.Println()
	fmt.Println("[ROOT_FILES]")
	roots := []string{"GOVERNANCE.toon", "index.toon", "volatile.toon"}
	for _, f := range roots {
		path := filepath.Join(memDir, f)
		if util.FileExists(path) {
			fmt.Printf("  %s\n", f)
		}
	}
}

func listCategoryFiles(tools *rust.Tools, memDir, project, category string) {
	projectPath := filepath.Join(memDir, category, project)

	fmt.Printf("[%s/%s]\n", category, project)

	// Use fd if available
	if tools.HasFd() {
		files, err := rust.FindTOON(projectPath)
		if err == nil {
			for _, f := range files {
				fmt.Printf("  %s\n", filepath.Base(f))
			}
			return
		}
	}

	// Fallback: os.ReadDir
	entries, err := os.ReadDir(projectPath)
	if err != nil {
		fmt.Printf("  (empty or not found)\n")
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			fmt.Printf("  %s\n", e.Name())
		}
	}
}

func viewFile(tools *rust.Tools, memDir, project, category, key string) {
	// Build file path
	fileName := key
	if filepath.Ext(key) == "" {
		fileName = key + ".toon"
	}

	filePath := filepath.Join(memDir, category, project, fileName)

	// Try project path first, then global
	if !util.FileExists(filePath) {
		filePath = filepath.Join(memDir, category, "global", fileName)
	}

	if !util.FileExists(filePath) {
		fmt.Printf("File not found: %s/%s/%s\n", category, project, fileName)
		return
	}

	fmt.Printf("[FILE: %s]\n", filePath)
	fmt.Println()

	// Use bat for syntax highlighting
	if tools.HasBat() {
		content, err := rust.CatTOON(filePath)
		if err == nil {
			fmt.Println(string(content))
			return
		}
	}

	// Fallback: plain read
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	fmt.Println(string(content))
}
