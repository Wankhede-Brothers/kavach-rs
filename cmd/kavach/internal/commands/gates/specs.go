// Package gates provides hook gates for Claude Code.
// specs.go: Specs Driven Development (SDD) gate — injects matching spec files as TOON context.
// DACE: Single responsibility - spec resolution and injection only.
package gates

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/util"
)

// specsDrivenGate resolves and returns spec content for the current intent.
// Returns TOON-formatted spec content or empty string if no specs found.
func specsDrivenGate(input *hook.Input, session *enforce.SessionState) string {
	specsDir := util.MemoryBankPath("specs")
	if !util.DirExists(specsDir) {
		return ""
	}

	domain := strings.ToLower(session.IntentDomain)
	intentType := strings.ToLower(session.IntentType)

	// Resolution order: combined > domain > type > default (highest priority first)
	candidates := []string{}
	if domain != "" && intentType != "" {
		candidates = append(candidates, domain+"-"+intentType+".toon")
	}
	if domain != "" {
		candidates = append(candidates, domain+".toon")
	}
	if intentType != "" {
		candidates = append(candidates, intentType+".toon")
	}
	candidates = append(candidates, "default.toon")

	var parts []string
	for _, name := range candidates {
		if session.WasSpecInjected(name) {
			continue
		}
		path := filepath.Join(specsDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := strings.TrimSpace(string(data))
		if content == "" {
			continue
		}
		parts = append(parts, "[SPEC:"+name+"]\n"+content)
		session.MarkSpecInjected(name)
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "\n\n")
}
