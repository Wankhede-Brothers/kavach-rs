package gates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/claude/shared/pkg/enforce"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/util"
)

func setupTestSpecs(t *testing.T) string {
	t.Helper()
	specsDir := util.MemoryBankPath("specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatalf("mkdir specs: %v", err)
	}
	// Write test spec files
	files := map[string]string{
		"security.toon":            "[SPEC:security]\nconstraint: use constant-time comparison",
		"implement.toon":           "[SPEC:implement]\nconstraint: read before modify",
		"security-implement.toon":  "[SPEC:security-implement]\nconstraint: OWASP review required",
		"default.toon":             "[SPEC:default]\nconstraint: follow conventions",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(specsDir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return specsDir
}

func newTestSession(domain, intentType string) *enforce.SessionState {
	s := &enforce.SessionState{
		ID:           "test_session",
		IntentDomain: domain,
		IntentType:   intentType,
	}
	return s
}

func TestSpecsDrivenGate_SecurityImplement(t *testing.T) {
	setupTestSpecs(t)
	session := newTestSession("security", "implement")
	input := &hook.Input{ToolName: "Write", ToolInput: map[string]interface{}{"file_path": "auth.go"}}

	result := specsDrivenGate(input, session)

	// Should include combined, domain, type, and default specs
	if !strings.Contains(result, "security-implement") {
		t.Error("expected security-implement spec in output")
	}
	if !strings.Contains(result, "OWASP review required") {
		t.Error("expected OWASP constraint from combined spec")
	}
	if !strings.Contains(result, "constant-time comparison") {
		t.Error("expected security domain spec")
	}
	if !strings.Contains(result, "read before modify") {
		t.Error("expected implement type spec")
	}
	if !strings.Contains(result, "follow conventions") {
		t.Error("expected default spec")
	}
}

func TestSpecsDrivenGate_DomainOnly(t *testing.T) {
	setupTestSpecs(t)
	session := newTestSession("security", "")
	input := &hook.Input{ToolName: "Write", ToolInput: map[string]interface{}{}}

	result := specsDrivenGate(input, session)

	if !strings.Contains(result, "constant-time comparison") {
		t.Error("expected security spec")
	}
	// No combined or type spec since intentType is empty
	if strings.Contains(result, "security-implement") {
		t.Error("should not include combined spec without intentType")
	}
}

func TestSpecsDrivenGate_NoIntent(t *testing.T) {
	setupTestSpecs(t)
	session := newTestSession("", "")
	input := &hook.Input{ToolName: "Write", ToolInput: map[string]interface{}{}}

	result := specsDrivenGate(input, session)

	// Only default.toon should match
	if !strings.Contains(result, "follow conventions") {
		t.Error("expected default spec when no intent")
	}
	if strings.Contains(result, "security") {
		t.Error("should not include domain spec without intent")
	}
}

func TestSpecsDrivenGate_DeduplicatesAcrossCalls(t *testing.T) {
	setupTestSpecs(t)
	session := newTestSession("security", "implement")
	input := &hook.Input{ToolName: "Write", ToolInput: map[string]interface{}{}}

	first := specsDrivenGate(input, session)
	if first == "" {
		t.Fatal("first call should return specs")
	}

	second := specsDrivenGate(input, session)
	if second != "" {
		t.Error("second call should return empty (already injected)")
	}
}

func TestSpecsDrivenGate_NonexistentDomain(t *testing.T) {
	setupTestSpecs(t)
	session := newTestSession("blockchain", "optimize")
	input := &hook.Input{ToolName: "Write", ToolInput: map[string]interface{}{}}

	result := specsDrivenGate(input, session)

	// Only default.toon exists for unknown domain/type
	if !strings.Contains(result, "follow conventions") {
		t.Error("expected default spec for unknown domain")
	}
	if strings.Contains(result, "security") {
		t.Error("should not include security spec for blockchain domain")
	}
}
