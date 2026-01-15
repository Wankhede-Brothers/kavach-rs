# Example Hooks

This directory contains example hooks showing how to create custom enforcement hooks for Claude Code.

## Hook Lifecycle

Hooks execute at specific points in the Claude Code workflow:

| Lifecycle | When | Use For |
|-----------|------|---------|
| **SessionStart** | New session begins | Initialize context, load preferences |
| **UserPromptSubmit** | User sends a message | Validate input, check rates |
| **PreToolUse** | Before tool executes | Gate tool usage, validate parameters |
| **PostToolUse** | After tool completes | Update memory, log events |
| **Stop** | Session ends | Save state, cleanup |

## Hook Input/Output

### Input (JSON via stdin)
```json
{
  "lifecycle": "PostToolUse",
  "tool_name": "Bash",
  "tool_result": {
    "output": "command output",
    "error": null
  },
  "agent": "backend-engineer",
  "model": "opus"
}
```

### Output (JSON via stdout)
```json
{
  "decision": "allow",
  "reason": "No issues detected",
  "modifications": {}
}
```

### Exit Codes
- `0` - Continue (decision was processed)
- `1` - Error (hook failed)
- `2` - Block (explicitly block operation)

## Example: Custom Linter Hook

```bash
#!/bin/bash
# hooks/my-linter.sh

# Read input
input=$(cat)
tool_name=$(echo "$input" | jq -r '.tool_name')

# Only check after Write/Edit tools
if [[ "$tool_name" == "Write" ]] || [[ "$tool_name" == "Edit" ]]; then
    # Get file path from tool_use_id
    # Run custom linter
    # Block if issues found
fi

# Output decision
echo '{"decision":"allow","reason":"No lint issues"}'
exit 0
```

## Example: Go Hook Template

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type HookInput struct {
    Lifecycle  string          `json:"lifecycle"`
    ToolName   string          `json:"tool_name,omitempty"`
    ToolResult *ToolResult     `json:"tool_result,omitempty"`
    Agent      string          `json:"agent"`
    Model      string          `json:"model"`
}

type HookOutput struct {
    Decision     string                 `json:"decision"`
    Reason       string                 `json:"reason"`
    Modifications map[string]interface{} `json:"modifications,omitempty"`
}

func main() {
    var input HookInput
    decoder := json.NewDecoder(os.Stdin)
    if err := decoder.Decode(&input); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Process input...

    output := HookOutput{
        Decision: "allow",
        Reason:   "All checks passed",
    }

    encoder := json.NewEncoder(os.Stdout)
    encoder.Encode(output)
}
```

## Registering Hooks

Add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "type": "command",
        "command": "/path/to/hook/my-hook"
      }
    ]
  }
}
```

## Best Practices

1. **Fast Execution** - Hooks should complete in < 1 second
2. **Idempotent** - Same input → same output
3. **Error Handling** - Return errors, don't crash
4. **Logging** - Use stderr for debug output
5. **JSON Output** - Always return valid JSON
