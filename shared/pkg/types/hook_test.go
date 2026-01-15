// Package types provides consolidated type definitions.
// hook_test.go: Tests for hook types.
package types

import "testing"

func TestHookInput_GetString(t *testing.T) {
	tests := []struct {
		name  string
		input *HookInput
		key   string
		want  string
	}{
		{
			name: "get file_path from ToolInput",
			input: &HookInput{
				ToolInput: map[string]interface{}{
					"file_path": "/test/file.go",
				},
			},
			key:  "file_path",
			want: "/test/file.go",
		},
		{
			name: "get prompt from Prompt field",
			input: &HookInput{
				Prompt: "fix the bug",
			},
			key:  "prompt",
			want: "fix the bug",
		},
		{
			name: "get nonexistent key",
			input: &HookInput{
				ToolInput: map[string]interface{}{},
			},
			key:  "nonexistent",
			want: "",
		},
		{
			name:  "nil ToolInput",
			input: &HookInput{},
			key:   "file_path",
			want:  "",
		},
		{
			name: "non-string value",
			input: &HookInput{
				ToolInput: map[string]interface{}{
					"count": 42,
				},
			},
			key:  "count",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.GetString(tt.key)
			if got != tt.want {
				t.Errorf("GetString(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestHookInput_GetToolName(t *testing.T) {
	input := &HookInput{ToolName: "Read"}
	if got := input.GetToolName(); got != "Read" {
		t.Errorf("GetToolName() = %v, want Read", got)
	}
}

func TestHookInput_GetPrompt(t *testing.T) {
	input := &HookInput{Prompt: "test prompt"}
	if got := input.GetPrompt(); got != "test prompt" {
		t.Errorf("GetPrompt() = %v, want test prompt", got)
	}
}

func TestNewApprove(t *testing.T) {
	resp := NewApprove("test reason")
	if resp.Decision != "approve" {
		t.Errorf("Decision = %v, want approve", resp.Decision)
	}
	if resp.Reason != "test reason" {
		t.Errorf("Reason = %v, want test reason", resp.Reason)
	}
}

func TestNewBlock(t *testing.T) {
	resp := NewBlock("blocked reason")
	if resp.Decision != "block" {
		t.Errorf("Decision = %v, want block", resp.Decision)
	}
	if resp.Reason != "blocked reason" {
		t.Errorf("Reason = %v, want blocked reason", resp.Reason)
	}
}

func TestNewModify(t *testing.T) {
	resp := NewModify("modify reason", "additional context")
	if resp.Decision != "approve" {
		t.Errorf("Decision = %v, want approve", resp.Decision)
	}
	if resp.Reason != "modify reason" {
		t.Errorf("Reason = %v, want modify reason", resp.Reason)
	}
	if resp.AdditionalContext != "additional context" {
		t.Errorf("AdditionalContext = %v, want additional context", resp.AdditionalContext)
	}
}
