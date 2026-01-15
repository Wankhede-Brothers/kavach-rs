// Package gates provides hook gates for Claude Code.
// intent_types.go: Type definitions for intent classification.
// DACE: Micro-modular split from intent.go
package gates

// IntentClassification holds NLU classification results
type IntentClassification struct {
	Type        string
	Domain      string
	Skills      []string
	Agent       string
	SubAgents   []string
	ResearchReq bool
	Confidence  string
}
