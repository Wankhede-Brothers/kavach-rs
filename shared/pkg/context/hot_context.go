// Package context provides hot-context tracking for DACE.
// hot_context.go: Tracks recently read files to avoid re-reading.
// P3 FIX #16: Implements hot-context.json tracking.
package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/claude/shared/pkg/util"
)

// HotContext tracks recently read files for token optimization.
type HotContext struct {
	Files     map[string]*FileEntry `json:"files"`
	UpdatedAt string                `json:"updated_at"`
	mu        sync.RWMutex
}

// FileEntry represents a tracked file.
type FileEntry struct {
	Path      string `json:"path"`
	ReadAt    string `json:"read_at"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
	IsCode    bool   `json:"is_code"`
}

// hotContextPath returns the path to hot-context.json.
func hotContextPath() string {
	return filepath.Join(util.STMPath(), "hot-context.json")
}

// LoadHotContext loads the hot-context from disk.
func LoadHotContext() *HotContext {
	ctx := &HotContext{
		Files: make(map[string]*FileEntry),
	}

	path := hotContextPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return ctx
	}

	json.Unmarshal(data, ctx)
	if ctx.Files == nil {
		ctx.Files = make(map[string]*FileEntry)
	}
	return ctx
}

// Save persists the hot-context to disk.
func (h *HotContext) Save() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.UpdatedAt = time.Now().Format(time.RFC3339)

	path := hotContextPath()
	if err := util.EnsureParentDir(path); err != nil {
		return err
	}

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// TrackFile records a file read.
func (h *HotContext) TrackFile(filePath string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	info, err := os.Stat(filePath)
	var size int64
	if err == nil {
		size = info.Size()
	}

	ext := filepath.Ext(filePath)
	isCode := isCodeExtension(ext)

	h.Files[filePath] = &FileEntry{
		Path:      filePath,
		ReadAt:    time.Now().Format(time.RFC3339),
		Size:      size,
		Extension: ext,
		IsCode:    isCode,
	}
}

// WasRecentlyRead checks if a file was read recently (within TTL).
func (h *HotContext) WasRecentlyRead(filePath string, ttl time.Duration) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	entry, ok := h.Files[filePath]
	if !ok {
		return false
	}

	readAt, err := time.Parse(time.RFC3339, entry.ReadAt)
	if err != nil {
		return false
	}

	return time.Since(readAt) < ttl
}

// GetRecentFiles returns files read within the TTL.
func (h *HotContext) GetRecentFiles(ttl time.Duration) []*FileEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var recent []*FileEntry
	now := time.Now()

	for _, entry := range h.Files {
		readAt, err := time.Parse(time.RFC3339, entry.ReadAt)
		if err != nil {
			continue
		}
		if now.Sub(readAt) < ttl {
			recent = append(recent, entry)
		}
	}

	return recent
}

// Cleanup removes entries older than TTL.
func (h *HotContext) Cleanup(ttl time.Duration) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	count := 0
	now := time.Now()

	for path, entry := range h.Files {
		readAt, err := time.Parse(time.RFC3339, entry.ReadAt)
		if err != nil || now.Sub(readAt) > ttl {
			delete(h.Files, path)
			count++
		}
	}

	return count
}

// Count returns the number of tracked files.
func (h *HotContext) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Files)
}

// isCodeExtension checks if extension is a code file.
func isCodeExtension(ext string) bool {
	codeExts := map[string]bool{
		".go": true, ".rs": true, ".ts": true, ".tsx": true,
		".js": true, ".jsx": true, ".py": true, ".rb": true,
		".java": true, ".c": true, ".cpp": true, ".h": true,
		".cs": true, ".swift": true, ".kt": true, ".scala": true,
		".zig": true, ".vue": true, ".svelte": true,
	}
	return codeExts[ext]
}

// Global hot context instance
var (
	globalHotContext *HotContext
	hotContextOnce   sync.Once
)

// GetHotContext returns the global hot context instance.
func GetHotContext() *HotContext {
	hotContextOnce.Do(func() {
		globalHotContext = LoadHotContext()
	})
	return globalHotContext
}

// TrackFileRead is a convenience function to track a file read.
func TrackFileRead(filePath string) {
	ctx := GetHotContext()
	ctx.TrackFile(filePath)
	ctx.Save()
}

// WasFileRecentlyRead checks if file was read in last hour.
func WasFileRecentlyRead(filePath string) bool {
	ctx := GetHotContext()
	return ctx.WasRecentlyRead(filePath, time.Hour)
}

// P1 FIX: Additional tracking functions for Write, Edit, and Agent operations

// TrackFileWrite tracks a file write operation.
func TrackFileWrite(filePath string) {
	ctx := GetHotContext()
	ctx.TrackFile(filePath)
	ctx.Save()
}

// TrackFileEdit tracks a file edit operation.
func TrackFileEdit(filePath string) {
	ctx := GetHotContext()
	ctx.TrackFile(filePath)
	ctx.Save()
}

// AgentCompletion tracks agent task completions.
type AgentCompletion struct {
	AgentType   string `json:"agent_type"`
	CompletedAt string `json:"completed_at"`
	Count       int    `json:"count"`
}

// agentCompletions stores agent completion stats.
var agentCompletions = make(map[string]*AgentCompletion)
var agentMu sync.Mutex

// TrackAgentCompletion records an agent task completion.
func TrackAgentCompletion(agentType string) {
	agentMu.Lock()
	defer agentMu.Unlock()

	if comp, ok := agentCompletions[agentType]; ok {
		comp.Count++
		comp.CompletedAt = time.Now().Format(time.RFC3339)
	} else {
		agentCompletions[agentType] = &AgentCompletion{
			AgentType:   agentType,
			CompletedAt: time.Now().Format(time.RFC3339),
			Count:       1,
		}
	}
}
