// Package enforce provides enforcement context and session management.
// session.go: Facade for session subpackage (backward compatibility).
// DACE: Re-exports from micro-modular session package.
package enforce

import (
	"github.com/claude/shared/pkg/enforce/session"
)

// SessionState is an alias to session.SessionState for backward compatibility.
type SessionState = session.SessionState

// NewSessionState creates a new session with date injection.
func NewSessionState(workDir string) *SessionState {
	return session.NewSessionState(workDir)
}

// LoadSessionState loads existing session state from TOON file.
func LoadSessionState() (*SessionState, error) {
	return session.LoadSessionState()
}

// GetOrCreateSession returns existing or new session.
func GetOrCreateSession() *SessionState {
	return session.GetOrCreateSession()
}
