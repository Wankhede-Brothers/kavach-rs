package types

// MemoryEntry represents a memory bank entry.
type MemoryEntry struct {
	Category string `json:"category"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	Verified string `json:"verified,omitempty"`
	TTL      string `json:"ttl,omitempty"`
	Source   string `json:"source,omitempty"`
}

// CompactSummary represents pre-compact session summary.
type CompactSummary struct {
	SessionID    string   `json:"session_id"`
	FilesRead    []string `json:"files_read"`
	FilesWritten []string `json:"files_written"`
	TasksTotal   int      `json:"tasks_total"`
	TasksDone    int      `json:"tasks_done"`
	Timestamp    string   `json:"timestamp"`
}

// HotContext represents hot context tracking data.
type HotContext struct {
	Project     string          `json:"project"`
	FilesRead   []FileReadEntry `json:"files_read"`
	TotalReads  int             `json:"total_reads"`
	LastUpdated string          `json:"last_updated"`
}

// FileReadEntry represents a file read tracking entry.
type FileReadEntry struct {
	Path      string `json:"path"`
	ReadCount int    `json:"read_count"`
	LastRead  string `json:"last_read"`
	Summary   string `json:"summary,omitempty"`
}

// Pattern represents a LTM pattern entry.
type Pattern struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Content  string `json:"content"`
	Verified string `json:"verified"`
	TTL      string `json:"ttl"`
	Source   string `json:"source,omitempty"`
}
