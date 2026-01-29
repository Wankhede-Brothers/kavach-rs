// Package session provides session state management.
// helpers.go: Utility functions and path management.
// DACE: Single responsibility - helpers only.
package session

import (
	"os"
	"path/filepath"

	"github.com/claude/shared/pkg/util"
)

// StatePath returns path to session state file.
func StatePath() string {
	return filepath.Join(util.STMPath(), "session-state.toon")
}

// GetOrCreateSession returns existing or new session.
// IMPORTANT: Always uses DetectProject() for current project since workdir changes.
func GetOrCreateSession() *SessionState {
	state, err := LoadSessionState()
	if err != nil || state == nil {
		wd, _ := os.Getwd()
		state = NewSessionState(wd)
		state.Save()
	} else {
		// Always update project based on current working directory
		// (session may have been created from a different directory)
		wd, _ := os.Getwd()
		state.WorkDir = wd
		state.Project = util.DetectProject()
		state.Save() // Persist updated project/workdir
	}
	return state
}

// ToTOON converts session state to compact TOON format.
func (s *SessionState) ToTOON() string {
	research := "PENDING"
	if s.ResearchDone {
		research = "DONE"
	}
	memory := "PENDING"
	if s.MemoryQueried {
		memory = "DONE"
	}

	return "[SESSION]\n" +
		"id: " + s.ID + "\n" +
		"today: " + s.Today + "\n" +
		"project: " + s.Project + "\n" +
		"research: " + research + "\n" +
		"memory: " + memory + "\n" +
		"cutoff: " + s.TrainingCutoff + "\n"
}
