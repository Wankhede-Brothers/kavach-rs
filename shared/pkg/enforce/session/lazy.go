// Package session provides session state management.
// lazy.go: Partial TOON parse for lightweight session loading.
// Stops parsing early once the required fields are found.
package session

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/claude/shared/pkg/util"
)

// LoadIdentity loads only identity fields from session state.
// Returns nil if no valid session exists.
func LoadIdentity() *SessionIdentity {
	path := StatePath()
	if !util.FileExists(path) {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	id := &SessionIdentity{}
	scanner := bufio.NewScanner(f)
	found := 0
	for scanner.Scan() && found < 5 {
		line := strings.TrimSpace(scanner.Text())
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			switch key {
			case "id":
				id.ID = val
				found++
			case "today":
				id.Today = val
				found++
			case "project":
				id.Project = val
				found++
			case "workdir":
				id.WorkDir = val
				found++
			case "session_id":
				id.SessionID = val
				found++
			}
		}
	}
	if id.Today != time.Now().Format("2006-01-02") {
		return nil
	}
	return id
}

// LoadFlags loads identity + enforcement flags.
func LoadFlags() *SessionFlags {
	path := StatePath()
	if !util.FileExists(path) {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	flags := &SessionFlags{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			switch key {
			case "id":
				flags.ID = val
			case "today":
				flags.Today = val
			case "project":
				flags.Project = val
			case "workdir":
				flags.WorkDir = val
			case "session_id":
				flags.SessionID = val
			case "research_done":
				flags.ResearchDone = val == "true"
			case "memory":
				flags.MemoryQueried = val == "true"
			case "nlu":
				flags.NLUParsed = val == "true"
			}
		}
	}
	if flags.Today != time.Now().Format("2006-01-02") {
		return nil
	}
	return flags
}

// LoadTracking loads identity + flags + tracking counters.
func LoadTracking() *SessionTracking {
	path := StatePath()
	if !util.FileExists(path) {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	t := &SessionTracking{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			switch key {
			case "id":
				t.ID = val
			case "today":
				t.Today = val
			case "project":
				t.Project = val
			case "workdir":
				t.WorkDir = val
			case "session_id":
				t.SessionID = val
			case "research_done":
				t.ResearchDone = val == "true"
			case "memory":
				t.MemoryQueried = val == "true"
			case "nlu":
				t.NLUParsed = val == "true"
			case "turn_count":
				t.TurnCount, _ = strconv.Atoi(val)
			case "post_compact":
				t.PostCompact = val == "true"
			case "task":
				t.CurrentTask = val
			case "tasks_created":
				t.TasksCreated, _ = strconv.Atoi(val)
			case "tasks_completed":
				t.TasksCompleted, _ = strconv.Atoi(val)
			}
		}
	}
	if t.Today != time.Now().Format("2006-01-02") {
		return nil
	}
	return t
}
