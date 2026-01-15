package toon

import (
	"os"
	"path/filepath"

	"github.com/claude/shared/pkg/util"
)

// MemoryBank provides TOON-based memory bank operations.
// Base path: ~/.local/shared/shared-ai/memory/
type MemoryBank struct {
	parser  *Parser
	workDir string
}

// NewMemoryBank creates a new memory bank accessor.
func NewMemoryBank() *MemoryBank {
	return &MemoryBank{
		parser:  NewParser(),
		workDir: util.WorkingDir(),
	}
}

// NewMemoryBankForProject creates a memory bank for a specific project.
func NewMemoryBankForProject(workDir string) *MemoryBank {
	return &MemoryBank{
		parser:  NewParser(),
		workDir: workDir,
	}
}

// LoadCategory loads all TOON files from a memory category.
// Categories: decisions, graph, kanban, patterns, proposals, research, roadmaps, STM
func (m *MemoryBank) LoadCategory(category string) ([]*Document, error) {
	path := util.MemoryBankPath(category)
	return m.loadFromDir(path)
}

// LoadProjectCategory loads TOON files for current project within a category.
func (m *MemoryBank) LoadProjectCategory(category string) ([]*Document, error) {
	if m.workDir == "" {
		return nil, nil
	}
	path := util.ProjectMemoryPath(m.workDir, category)
	return m.loadFromDir(path)
}

// loadFromDir loads all TOON files from a directory (recursive).
// DACE: Memory Bank uses TOON format only, not markdown.
// P2 FIX: Track skipped files instead of silent ignore.
func (m *MemoryBank) loadFromDir(path string) ([]*Document, error) {
	if !util.DirExists(path) {
		return nil, nil
	}

	var docs []*Document
	var skippedCount int

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			skippedCount++ // P2 FIX: Track instead of silent skip
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(filePath)
		if ext != ".toon" {
			return nil // TOON only - not an error
		}

		doc, loadErr := m.LoadFile(filePath)
		if loadErr != nil {
			skippedCount++ // P2 FIX: Track invalid files
			return nil
		}

		docs = append(docs, doc)
		return nil
	})

	// Note: skippedCount is tracked but not returned to avoid API change
	// Future: Add LoadResult struct with stats
	return docs, err
}

// LoadFile loads a single TOON file.
func (m *MemoryBank) LoadFile(path string) (*Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return m.parser.Parse(f)
}

// LoadGovernance loads the top-level GOVERNANCE.toon file.
func (m *MemoryBank) LoadGovernance() (*Document, error) {
	return m.LoadFile(util.GovernancePath())
}

// LoadIndex loads the top-level index.toon file.
func (m *MemoryBank) LoadIndex() (*Document, error) {
	return m.LoadFile(util.IndexPath())
}

// LoadVolatile loads the top-level volatile.toon file.
func (m *MemoryBank) LoadVolatile() (*Document, error) {
	return m.LoadFile(util.VolatilePath())
}

// LoadDecisions loads the decisions category.
func (m *MemoryBank) LoadDecisions() ([]*Document, error) {
	return m.LoadCategory("decisions")
}

// LoadPatterns loads the patterns category.
func (m *MemoryBank) LoadPatterns() ([]*Document, error) {
	return m.LoadCategory("patterns")
}

// LoadResearch loads the research category.
func (m *MemoryBank) LoadResearch() ([]*Document, error) {
	return m.LoadCategory("research")
}

// LoadRoadmaps loads the roadmaps category.
func (m *MemoryBank) LoadRoadmaps() ([]*Document, error) {
	return m.LoadCategory("roadmaps")
}

// LoadProposals loads the proposals category.
func (m *MemoryBank) LoadProposals() ([]*Document, error) {
	return m.LoadCategory("proposals")
}

// LoadGraph loads the graph category.
func (m *MemoryBank) LoadGraph() ([]*Document, error) {
	return m.LoadCategory("graph")
}

// LoadKanban loads the kanban category.
func (m *MemoryBank) LoadKanban() ([]*Document, error) {
	return m.LoadCategory("kanban")
}

// LoadSTM loads the STM (Short-Term Memory) category.
func (m *MemoryBank) LoadSTM() ([]*Document, error) {
	return m.LoadCategory("STM")
}

// SaveFile saves a TOON document to a file.
func (m *MemoryBank) SaveFile(path string, doc *Document) error {
	if err := util.EnsureParentDir(path); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := NewWriter(f)
	return writer.WriteDocument(doc)
}

// SaveToCategory saves a document to a category.
func (m *MemoryBank) SaveToCategory(category, filename string, doc *Document) error {
	path := filepath.Join(util.MemoryBankPath(category), filename)
	return m.SaveFile(path, doc)
}

// SaveToProjectCategory saves a document to project-specific category.
func (m *MemoryBank) SaveToProjectCategory(category, filename string, doc *Document) error {
	if m.workDir == "" {
		return m.SaveToCategory(category, filename, doc)
	}
	path := filepath.Join(util.ProjectMemoryPath(m.workDir, category), filename)
	return m.SaveFile(path, doc)
}

// Query searches for blocks matching criteria across categories.
func (m *MemoryBank) Query(category string, blockName string) ([]*Block, error) {
	docs, err := m.LoadCategory(category)
	if err != nil {
		return nil, err
	}

	var results []*Block
	for _, doc := range docs {
		if block := doc.Get(blockName); block != nil {
			results = append(results, block)
		}
	}

	return results, nil
}

// ListCategories returns all available memory categories.
func (m *MemoryBank) ListCategories() []string {
	return []string{
		"decisions",
		"graph",
		"kanban",
		"patterns",
		"proposals",
		"research",
		"roadmaps",
		"STM",
	}
}

// GetCategoryStats returns file count for each category.
func (m *MemoryBank) GetCategoryStats() map[string]int {
	stats := make(map[string]int)
	for _, cat := range m.ListCategories() {
		docs, _ := m.LoadCategory(cat)
		stats[cat] = len(docs)
	}
	return stats
}
