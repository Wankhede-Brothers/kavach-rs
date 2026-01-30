// Package session provides session state management.
// subsets.go: Lightweight session subset structs for lazy loading.
// Gates load only the subset they need, avoiding full session I/O.
package session

// SessionIdentity holds minimal identity fields (5 fields).
// Used by: preToolCEO, preToolTask
type SessionIdentity struct {
	ID        string
	Today     string
	Project   string
	WorkDir   string
	SessionID string
}

// SessionFlags extends identity with enforcement flags.
// Used by: preWrite, postTool WebSearch/WebFetch
type SessionFlags struct {
	SessionIdentity
	ResearchDone  bool
	MemoryQueried bool
	NLUParsed     bool
}

// SessionTracking extends flags with turn/task counters.
// Used by: postWrite quality, postTool TaskCreate/TaskUpdate
type SessionTracking struct {
	SessionFlags
	TurnCount      int
	PostCompact    bool
	CurrentTask    string
	TasksCreated   int
	TasksCompleted int
}

// IdentityFromState extracts SessionIdentity from full state.
func IdentityFromState(s *SessionState) *SessionIdentity {
	return &SessionIdentity{
		ID:        s.ID,
		Today:     s.Today,
		Project:   s.Project,
		WorkDir:   s.WorkDir,
		SessionID: s.SessionID,
	}
}

// FlagsFromState extracts SessionFlags from full state.
func FlagsFromState(s *SessionState) *SessionFlags {
	return &SessionFlags{
		SessionIdentity: *IdentityFromState(s),
		ResearchDone:    s.ResearchDone,
		MemoryQueried:   s.MemoryQueried,
		NLUParsed:       s.NLUParsed,
	}
}

// TrackingFromState extracts SessionTracking from full state.
func TrackingFromState(s *SessionState) *SessionTracking {
	return &SessionTracking{
		SessionFlags:   *FlagsFromState(s),
		TurnCount:      s.TurnCount,
		PostCompact:    s.PostCompact,
		CurrentTask:    s.CurrentTask,
		TasksCreated:   s.TasksCreated,
		TasksCompleted: s.TasksCompleted,
	}
}
