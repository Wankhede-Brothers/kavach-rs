// Package telemetry provides structured hook span tracking.
// report.go: Aggregation and reporting for kavach telemetry report.
package telemetry

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Report holds aggregated telemetry data.
type Report struct {
	TotalSpans    int             `json:"total_spans"`
	TotalDuration int64           `json:"total_duration_ms"`
	ByHook        map[string]*Agg `json:"by_hook"`
	ByResult      map[string]int  `json:"by_result"`
	ByTool        map[string]*Agg `json:"by_tool"`
	Slowest       []Span          `json:"slowest"`
}

// Agg holds per-hook or per-tool aggregated metrics.
type Agg struct {
	Count      int   `json:"count"`
	TotalMs    int64 `json:"total_ms"`
	MaxMs      int64 `json:"max_ms"`
	TotalToken int   `json:"total_tokens"`
}

// GenerateReport reads hooks.jsonl and produces an aggregated report.
func GenerateReport() (*Report, error) {
	path := filepath.Join(telemetryDir(), "hooks.jsonl")
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open telemetry: %w", err)
	}
	defer f.Close()

	r := &Report{
		ByHook:   make(map[string]*Agg),
		ByResult: make(map[string]int),
		ByTool:   make(map[string]*Agg),
	}

	scanner := bufio.NewScanner(f)
	var spans []Span
	for scanner.Scan() {
		var s Span
		if json.Unmarshal(scanner.Bytes(), &s) != nil {
			continue
		}
		spans = append(spans, s)
		r.TotalSpans++
		r.TotalDuration += s.DurationMs
		r.ByResult[s.Result]++

		addToAgg(r.ByHook, s.Hook, &s)
		if s.Tool != "" {
			addToAgg(r.ByTool, s.Tool, &s)
		}
	}

	// Find top 5 slowest
	sort.Slice(spans, func(i, j int) bool {
		return spans[i].DurationMs > spans[j].DurationMs
	})
	limit := 5
	if len(spans) < limit {
		limit = len(spans)
	}
	r.Slowest = spans[:limit]

	return r, scanner.Err()
}

func addToAgg(m map[string]*Agg, key string, s *Span) {
	a, ok := m[key]
	if !ok {
		a = &Agg{}
		m[key] = a
	}
	a.Count++
	a.TotalMs += s.DurationMs
	a.TotalToken += s.Tokens
	if s.DurationMs > a.MaxMs {
		a.MaxMs = s.DurationMs
	}
}
