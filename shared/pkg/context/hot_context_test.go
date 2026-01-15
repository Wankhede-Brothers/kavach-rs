// Package context provides hot-context tracking for DACE.
// hot_context_test.go: Tests for hot context tracking.
package context

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHotContext_TrackFile(t *testing.T) {
	ctx := &HotContext{
		Files: make(map[string]*FileEntry),
	}

	// Create a temp file to track
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(testFile, []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}

	ctx.TrackFile(testFile)

	if len(ctx.Files) != 1 {
		t.Errorf("expected 1 file tracked, got %d", len(ctx.Files))
	}

	entry, ok := ctx.Files[testFile]
	if !ok {
		t.Fatal("file not found in tracked files")
	}

	if entry.Path != testFile {
		t.Errorf("expected path %s, got %s", testFile, entry.Path)
	}

	if entry.Extension != ".go" {
		t.Errorf("expected extension .go, got %s", entry.Extension)
	}

	if !entry.IsCode {
		t.Error("expected IsCode to be true for .go file")
	}
}

func TestHotContext_WasRecentlyRead(t *testing.T) {
	ctx := &HotContext{
		Files: make(map[string]*FileEntry),
	}

	testPath := "/test/file.go"
	ctx.Files[testPath] = &FileEntry{
		Path:   testPath,
		ReadAt: time.Now().Format(time.RFC3339),
	}

	// Should be recently read (within 1 hour)
	if !ctx.WasRecentlyRead(testPath, time.Hour) {
		t.Error("expected file to be recently read")
	}

	// Test with old timestamp
	ctx.Files[testPath].ReadAt = time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	if ctx.WasRecentlyRead(testPath, time.Hour) {
		t.Error("expected file to NOT be recently read")
	}

	// Test non-existent file
	if ctx.WasRecentlyRead("/nonexistent", time.Hour) {
		t.Error("expected non-existent file to return false")
	}
}

func TestHotContext_GetRecentFiles(t *testing.T) {
	ctx := &HotContext{
		Files: make(map[string]*FileEntry),
	}

	now := time.Now()
	ctx.Files["/recent.go"] = &FileEntry{
		Path:   "/recent.go",
		ReadAt: now.Format(time.RFC3339),
	}
	ctx.Files["/old.go"] = &FileEntry{
		Path:   "/old.go",
		ReadAt: now.Add(-2 * time.Hour).Format(time.RFC3339),
	}

	recent := ctx.GetRecentFiles(time.Hour)
	if len(recent) != 1 {
		t.Errorf("expected 1 recent file, got %d", len(recent))
	}

	if recent[0].Path != "/recent.go" {
		t.Errorf("expected /recent.go, got %s", recent[0].Path)
	}
}

func TestHotContext_Cleanup(t *testing.T) {
	ctx := &HotContext{
		Files: make(map[string]*FileEntry),
	}

	now := time.Now()
	ctx.Files["/recent.go"] = &FileEntry{
		Path:   "/recent.go",
		ReadAt: now.Format(time.RFC3339),
	}
	ctx.Files["/old.go"] = &FileEntry{
		Path:   "/old.go",
		ReadAt: now.Add(-2 * time.Hour).Format(time.RFC3339),
	}

	removed := ctx.Cleanup(time.Hour)
	if removed != 1 {
		t.Errorf("expected 1 file removed, got %d", removed)
	}

	if len(ctx.Files) != 1 {
		t.Errorf("expected 1 file remaining, got %d", len(ctx.Files))
	}
}

func TestHotContext_Count(t *testing.T) {
	ctx := &HotContext{
		Files: make(map[string]*FileEntry),
	}

	if ctx.Count() != 0 {
		t.Error("expected count 0 for empty context")
	}

	ctx.Files["/test.go"] = &FileEntry{Path: "/test.go"}
	if ctx.Count() != 1 {
		t.Error("expected count 1")
	}
}

func TestIsCodeExtension(t *testing.T) {
	tests := []struct {
		ext    string
		isCode bool
	}{
		{".go", true},
		{".rs", true},
		{".ts", true},
		{".tsx", true},
		{".js", true},
		{".py", true},
		{".zig", true},
		{".md", false},
		{".txt", false},
		{".json", false},
		{".yaml", false},
		{"", false},
	}

	for _, tt := range tests {
		got := isCodeExtension(tt.ext)
		if got != tt.isCode {
			t.Errorf("isCodeExtension(%q) = %v, want %v", tt.ext, got, tt.isCode)
		}
	}
}

func TestTrackAgentCompletion(t *testing.T) {
	// Reset global state
	agentMu.Lock()
	agentCompletions = make(map[string]*AgentCompletion)
	agentMu.Unlock()

	TrackAgentCompletion("backend-engineer")
	TrackAgentCompletion("backend-engineer")
	TrackAgentCompletion("frontend-engineer")

	agentMu.Lock()
	defer agentMu.Unlock()

	if comp, ok := agentCompletions["backend-engineer"]; !ok {
		t.Error("expected backend-engineer completion tracked")
	} else if comp.Count != 2 {
		t.Errorf("expected count 2, got %d", comp.Count)
	}

	if comp, ok := agentCompletions["frontend-engineer"]; !ok {
		t.Error("expected frontend-engineer completion tracked")
	} else if comp.Count != 1 {
		t.Errorf("expected count 1, got %d", comp.Count)
	}
}
