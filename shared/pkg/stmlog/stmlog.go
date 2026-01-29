// Package stmlog provides STM session log appending.
// Used by gate hooks to persist tool events without importing the memory command package.
package stmlog

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/claude/shared/lock"
	"github.com/claude/shared/pkg/util"
)

// dedupCache tracks recent event fingerprints to prevent double-logging.
// Key: sha256(eventType+subject)[:16], Value: timestamp of last write.
var (
	dedupMu    sync.Mutex
	dedupCache = make(map[string]time.Time)
	dedupTTL   = 2 * time.Second
)

// isDuplicate returns true if the same event was logged within the TTL window.
func isDuplicate(eventType, subject string) bool {
	h := sha256.Sum256([]byte(eventType + "|" + subject))
	key := fmt.Sprintf("%x", h[:8])

	dedupMu.Lock()
	defer dedupMu.Unlock()

	now := time.Now()

	// Prune expired entries (cheap — cache is small)
	for k, t := range dedupCache {
		if now.Sub(t) > dedupTTL {
			delete(dedupCache, k)
		}
	}

	if last, ok := dedupCache[key]; ok && now.Sub(last) < dedupTTL {
		return true
	}
	dedupCache[key] = now
	return false
}

// AppendEvent writes a single event line to the STM session_log.toon.
func AppendEvent(project, eventType, subject, detail string) {
	if project == "" {
		project = util.DetectProjectForWrite()
	}

	// Bug 1: Dedup guard — skip if same event fired within 2s
	if isDuplicate(eventType, subject) {
		return
	}

	stmDir, err := util.EnsureScratchpadDir(project)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[STMLOG] scratchpad dir error: %v\n", err)
		return
	}
	logPath := stmDir + "/session_log.toon"

	// Bug 4: File locking to prevent race conditions on concurrent appends
	lm := lock.GetLockManager()
	if err := lm.AcquireWithTimeout(logPath, lock.DefaultTimeout); err != nil {
		fmt.Fprintf(os.Stderr, "[STMLOG] lock error: %v\n", err)
		return
	}
	defer lm.Release(logPath)

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[STMLOG] log open error: %v\n", err)
		return
	}
	defer f.Close()

	// Cap session log at ~200 events (~50KB)
	info, _ := f.Stat()
	if info != nil && info.Size() > 50*1024 {
		return
	}

	ts := time.Now().Format("15:04:05")
	line := fmt.Sprintf("[%s] %s: %s", ts, eventType, sanitize(subject))
	if detail != "" {
		line += " | " + sanitize(detail)
	}
	fmt.Fprintln(f, line)
}

// IsBashSignificant returns true if the command is worth logging.
func IsBashSignificant(command string) bool {
	lower := strings.ToLower(command)
	for _, sig := range []string{"build", "test", "deploy", "cargo", "go ", "bun ", "npm ", "git commit", "git push", "git merge"} {
		if strings.Contains(lower, sig) {
			return true
		}
	}
	return false
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, ",", ";")
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > 60 {
		s = s[:57] + "..."
	}
	return s
}
