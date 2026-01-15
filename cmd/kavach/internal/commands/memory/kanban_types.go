// Package memory provides memory bank commands.
// kanban_types.go: Type definitions for Kanban board.
// DACE: Micro-modular split from kanban.go
package memory

// Kanban Column Constants - Production Pipeline
const (
	ColBacklog    = "backlog"     // Tasks waiting to be started
	ColInProgress = "in_progress" // Currently being worked on
	ColTesting    = "testing"     // Aegis-Guard: Lints, Warnings, Core Bugs
	ColVerified   = "verified"    // Aegis-Guard: Algorithm, Dead Code, Suppressed Elements
	ColDone       = "done"        // Production Ready
)

// Verification Status Constants
const (
	VerifyPending = "pending"
	VerifyPassed  = "passed"
	VerifyFailed  = "failed"
	VerifyBlocked = "blocked"
)

// KanbanCard represents a task card with verification state
type KanbanCard struct {
	ID           string
	Column       string
	Title        string
	Priority     string
	Type         string
	Assignee     string
	AegisStatus  string   // pending, passed, failed, blocked
	LintIssues   int      // Count of lint issues
	Warnings     int      // Count of warnings
	CoreBugs     int      // Count of core bugs
	DeadCode     bool     // Dead code detected
	Suppressed   bool     // Suppressed elements detected
	AlgoVerified bool     // Algorithm verified for production
	FailReasons  []string // Reasons for verification failure
}

// KanbanBoard represents the full board with pipeline stages
type KanbanBoard struct {
	Project    string
	WorkDir    string
	Updated    string
	Phases     map[int][]KanbanCard
	LoopCount  int  // Number of CEO loops
	Production bool // Ready for production
}

// AegisReport represents verification report for CEO
type AegisReport struct {
	TaskID      string
	Stage       string // "testing" or "verified"
	Status      string // "passed" or "failed"
	LintIssues  int
	Warnings    int
	CoreBugs    int
	DeadCode    bool
	Suppressed  bool
	AlgoOK      bool
	FailReasons []string
	Timestamp   string
}
