// Package telemetry provides structured hook span tracking.
// span.go: Hook execution span recording to hooks.jsonl.
package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Span tracks a single hook execution.
type Span struct {
	Hook          string            `json:"hook"`
	Tool          string            `json:"tool,omitempty"`
	Timestamp     string            `json:"ts"`
	DurationMs    int64             `json:"duration_ms"`
	SessionLoaded bool              `json:"session_loaded"`
	Tokens        int               `json:"tokens"`
	Result        string            `json:"result"`
	Tier          int               `json:"tier,omitempty"`
	Checks        map[string]string `json:"checks,omitempty"`
	start         time.Time
}

// StartSpan begins tracking a hook execution.
func StartSpan(hook string) *Span {
	return &Span{
		Hook:      hook,
		Timestamp: time.Now().Format(time.RFC3339),
		start:     time.Now(),
	}
}

// SetTool sets the tool name for the span.
func (s *Span) SetTool(tool string) { s.Tool = tool }

// SetSessionLoaded marks whether session was loaded.
func (s *Span) SetSessionLoaded(loaded bool) { s.SessionLoaded = loaded }

// SetTokens sets the token count injected.
func (s *Span) SetTokens(tokens int) { s.Tokens = tokens }

// SetTier sets the intent tier level.
func (s *Span) SetTier(tier int) { s.Tier = tier }

// SetResult sets the hook outcome.
func (s *Span) SetResult(result string) { s.Result = result }

// AddCheck records a sub-check result.
func (s *Span) AddCheck(name, result string) {
	if s.Checks == nil {
		s.Checks = make(map[string]string)
	}
	s.Checks[name] = result
}

// End finalizes the span, computes duration, and appends to hooks.jsonl.
func (s *Span) End() {
	s.DurationMs = time.Since(s.start).Milliseconds()
	if s.Result == "" {
		s.Result = "allow"
	}
	appendSpan(s)
}

// EndWithResult finalizes with a specific result.
func (s *Span) EndWithResult(result string) {
	s.Result = result
	s.End()
}

// Today returns today's date string.
func Today() string {
	return time.Now().Format("2006-01-02")
}

// telemetryDir returns the telemetry output directory.
func telemetryDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "shared", "shared-ai", "memory", "telemetry")
}

// appendSpan writes a span as a JSON line to hooks.jsonl.
func appendSpan(s *Span) {
	dir := telemetryDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "[TELEMETRY] mkdir error: %v\n", err)
		return
	}
	path := filepath.Join(dir, "hooks.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[TELEMETRY] open error: %v\n", err)
		return
	}
	defer f.Close()

	data, err := json.Marshal(s)
	if err != nil {
		return
	}
	f.Write(data)
	f.Write([]byte("\n"))
}
