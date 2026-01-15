// Package hook provides hook input/output utilities.
// input_test.go: Tests for hook input parsing.
package hook

import (
	"strings"
	"testing"

	"github.com/claude/shared/pkg/types"
)

func TestReadHookInputFrom(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantTool string
		wantErr  bool
	}{
		{
			name:     "valid Read tool input",
			json:     `{"tool_name":"Read","tool_input":{"file_path":"/test/file.go"}}`,
			wantTool: "Read",
			wantErr:  false,
		},
		{
			name:     "valid Bash tool input",
			json:     `{"tool_name":"Bash","tool_input":{"command":"ls -la"}}`,
			wantTool: "Bash",
			wantErr:  false,
		},
		{
			name:     "valid UserPromptSubmit",
			json:     `{"prompt":"fix the bug"}`,
			wantTool: "",
			wantErr:  false,
		},
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantErr: true,
		},
		{
			name:     "empty input",
			json:     `{}`,
			wantTool: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.json)
			input, err := ReadHookInputFrom(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadHookInputFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && input.ToolName != tt.wantTool {
				t.Errorf("ToolName = %v, want %v", input.ToolName, tt.wantTool)
			}
		})
	}
}

func TestGetStringFromInput(t *testing.T) {
	input := &types.HookInput{
		ToolName: "Read",
		ToolInput: map[string]interface{}{
			"file_path": "/test/file.go",
			"command":   "ls",
		},
	}

	tests := []struct {
		key  string
		want string
	}{
		{"file_path", "/test/file.go"},
		{"command", "ls"},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		got := GetStringFromInput(input, tt.key)
		if got != tt.want {
			t.Errorf("GetStringFromInput(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}

	// Test nil input
	if got := GetStringFromInput(nil, "key"); got != "" {
		t.Errorf("GetStringFromInput(nil, key) = %v, want empty", got)
	}

	// Test nil ToolInput
	nilInput := &types.HookInput{ToolName: "Test"}
	if got := GetStringFromInput(nilInput, "key"); got != "" {
		t.Errorf("GetStringFromInput with nil ToolInput = %v, want empty", got)
	}
}

func TestGetBoolFromInput(t *testing.T) {
	input := &types.HookInput{
		ToolInput: map[string]interface{}{
			"verbose": true,
			"quiet":   false,
			"string":  "not a bool",
		},
	}

	tests := []struct {
		key  string
		want bool
	}{
		{"verbose", true},
		{"quiet", false},
		{"string", false},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		got := GetBoolFromInput(input, tt.key)
		if got != tt.want {
			t.Errorf("GetBoolFromInput(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}

	// Test nil input
	if got := GetBoolFromInput(nil, "key"); got != false {
		t.Errorf("GetBoolFromInput(nil, key) = %v, want false", got)
	}
}

func TestGetIntFromInput(t *testing.T) {
	input := &types.HookInput{
		ToolInput: map[string]interface{}{
			"count":    5,
			"float":    10.5,
			"string":   "not an int",
			"negative": -3,
		},
	}

	tests := []struct {
		key  string
		want int
	}{
		{"count", 5},
		{"float", 10},
		{"string", 0},
		{"negative", -3},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
		got := GetIntFromInput(input, tt.key)
		if got != tt.want {
			t.Errorf("GetIntFromInput(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}

	// Test nil input
	if got := GetIntFromInput(nil, "key"); got != 0 {
		t.Errorf("GetIntFromInput(nil, key) = %v, want 0", got)
	}
}
