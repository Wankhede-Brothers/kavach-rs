// Package hook provides hook input/output utilities for the umbrella CLI.
package hook

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/claude/shared/pkg/types"
)

// Input is an alias to types.HookInput for convenience.
type Input = types.HookInput

// ReadHookInput reads and parses JSON hook input from stdin.
func ReadHookInput() (*types.HookInput, error) {
	return ReadHookInputFrom(os.Stdin)
}

// ReadHookInputFrom reads and parses JSON hook input from a reader.
func ReadHookInputFrom(r io.Reader) (*types.HookInput, error) {
	reader := bufio.NewReader(r)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var input types.HookInput
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}

	return &input, nil
}

// MustReadHookInput reads hook input or exits with error JSON.
func MustReadHookInput() *types.HookInput {
	input, err := ReadHookInput()
	if err != nil {
		OutputError("failed to read hook input: " + err.Error())
		os.Exit(1)
	}
	return input
}

// GetStringFromInput extracts a string value from tool input by key.
func GetStringFromInput(input *types.HookInput, key string) string {
	if input == nil || input.ToolInput == nil {
		return ""
	}
	if val, ok := input.ToolInput[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetBoolFromInput extracts a boolean value from tool input by key.
func GetBoolFromInput(input *types.HookInput, key string) bool {
	if input == nil || input.ToolInput == nil {
		return false
	}
	if val, ok := input.ToolInput[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// GetIntFromInput extracts an integer value from tool input by key.
func GetIntFromInput(input *types.HookInput, key string) int {
	if input == nil || input.ToolInput == nil {
		return 0
	}
	if val, ok := input.ToolInput[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	return 0
}
