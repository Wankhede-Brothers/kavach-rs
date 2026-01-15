package types

import "time"

// GateState represents persistent session state for gates.
type GateState struct {
	CEOInvoked   bool   `json:"ceo_invoked"`
	ResearchDone bool   `json:"research_done"`
	Timestamp    string `json:"timestamp"`
	SessionID    string `json:"session_id"`
}

// NewGateState creates a new gate state with current timestamp.
func NewGateState(sessionID string) *GateState {
	return &GateState{
		Timestamp: time.Now().Format(time.RFC3339),
		SessionID: sessionID,
	}
}

// GateDecision is an alias for HookResponse (backward compatibility).
type GateDecision = HookResponse

// DSARequest represents a DSA operation request.
type DSARequest struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path,omitempty"`
	Query   string `json:"query,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

// DSAResult represents a DSA operation result.
type DSAResult struct {
	Success     bool        `json:"success"`
	Count       int         `json:"result_count"`
	Results     []FileMatch `json:"results"`
	Algorithm   string      `json:"algorithm"`
	ExecutionMs int         `json:"execution_time_ms"`
	Note        string      `json:"note,omitempty"`
}

// FileMatch represents a file match from DSA operations.
type FileMatch struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Extension string `json:"extension,omitempty"`
}
