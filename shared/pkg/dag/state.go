// Package dag provides a parallel DAG scheduler for Kavach orchestration.
// state.go: JSON persistence for DAG state.
package dag

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// StatePath returns the file path for a session's DAG state.
func StatePath(sessionID string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "dag", sessionID+".json")
}

// Save persists DAG state to disk as JSON.
func Save(state *DAGState) error {
	path := StatePath(state.SessionID)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads DAG state from disk.
func Load(sessionID string) (*DAGState, error) {
	data, err := os.ReadFile(StatePath(sessionID))
	if err != nil {
		return nil, err
	}
	var state DAGState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &state, nil
}

// Delete removes DAG state for a session.
func Delete(sessionID string) error {
	return os.Remove(StatePath(sessionID))
}

// CleanupOld removes DAG state files older than maxAge.
// Called from session end to prevent accumulation.
func CleanupOld(maxAgeDays int) error {
	home, _ := os.UserHomeDir()
	dagDir := filepath.Join(home, ".claude", "dag")
	entries, err := os.ReadDir(dagDir)
	if err != nil {
		return nil // dir may not exist
	}
	cutoff := time.Now().AddDate(0, 0, -maxAgeDays)
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(dagDir, e.Name()))
		}
	}
	return nil
}
