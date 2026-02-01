# Installation Guide

Complete installation guide for Kavach across all platforms.

---

## Quick Install

### Linux/macOS

```bash
curl -fsSL https://raw.githubusercontent.com/Wankhede-Brothers/kavach-rs/main/install/install.sh | bash
```

### Windows (PowerShell Admin)

```powershell
irm https://raw.githubusercontent.com/Wankhede-Brothers/kavach-rs/main/install/install.ps1 | iex
```

### With Options

```bash
# Full install: kavach + Rust CLI tools
curl -fsSL .../install.sh | bash -s -- --full

# For OpenCode instead of Claude Code
curl -fsSL .../install.sh | bash -s -- --cli opencode

# Just Rust CLI tools
curl -fsSL .../install.sh | bash -s -- --rust-cli
```

---

## Build from Source

### Prerequisites

| Dependency | Version | Purpose |
|------------|---------|---------|
| Rust | stable | Compilation |
| Git | 2.0+ | Source download |
| Make | 4.0+ | Build system (optional) |

### Steps

```bash
# 1. Clone repository
git clone https://github.com/Wankhede-Brothers/kavach-rs.git
cd kavach-rs

# 2. Build and install
just build
just install

# 3. Verify
kavach status
```

---

## Platform-Specific Setup

### Linux

**Binary Location:** `~/.local/bin/kavach` or `~/.claude/bin/kavach`

**Memory Bank:** `~/.local/share/shared-ai/memory/`

**Config:** `~/.claude/` (for Claude Code)

```bash
# Manual install
mkdir -p ~/.local/bin
cp kavach ~/.local/bin/

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"

# Initialize memory bank
mkdir -p ~/.local/share/shared-ai/memory/{decisions,patterns,research,kanban,proposals,roadmaps,graph,STM}

# Copy config (for Claude Code)
cp configs/linux/CLAUDE.md ~/.claude/
```

### macOS

**Binary Location:** `~/.local/bin/kavach`

**Memory Bank:** `~/Library/Application Support/shared-ai/memory/`

**Config:** `~/.claude/`

```bash
# Manual install
mkdir -p ~/.local/bin
cp kavach ~/.local/bin/

# Add to PATH (add to ~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"

# Initialize memory bank
mkdir -p ~/Library/Application\ Support/shared-ai/memory/{decisions,patterns,research,kanban,proposals,roadmaps,graph,STM}

# Copy config
cp configs/darwin/CLAUDE.md ~/.claude/
```

### Windows

**Binary Location:** `%USERPROFILE%\.local\bin\kavach.exe`

**Memory Bank:** `%APPDATA%\shared-ai\memory\`

**Config:** `%USERPROFILE%\.claude\`

```powershell
# Manual install (PowerShell Admin)
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.local\bin"
Copy-Item kavach.exe "$env:USERPROFILE\.local\bin\"

# Add to PATH
$env:PATH += ";$env:USERPROFILE\.local\bin"
[Environment]::SetEnvironmentVariable("PATH", $env:PATH, "User")

# Initialize memory bank
$categories = @("decisions", "patterns", "research", "kanban", "proposals", "roadmaps", "graph", "STM")
foreach ($cat in $categories) {
    New-Item -ItemType Directory -Force -Path "$env:APPDATA\shared-ai\memory\$cat"
}

# Copy config
Copy-Item configs\windows\CLAUDE.md "$env:USERPROFILE\.claude\"
```

---

## Claude Code Configuration

After installing the binary, configure Claude Code hooks.

### Full Configuration

Create/update `~/.claude/settings.json`:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [{"type": "command", "command": "kavach session init"}]
      }
    ],
    "UserPromptSubmit": [
      {
        "hooks": [{"type": "command", "command": "kavach session init"}]
      }
    ],
    "PreToolUse": [
      {
        "matcher": "Task",
        "hooks": [{"type": "command", "command": "kavach gates ceo --hook"}]
      },
      {
        "matcher": "Bash",
        "hooks": [{"type": "command", "command": "kavach gates bash --hook"}]
      },
      {
        "matcher": "Read",
        "hooks": [{"type": "command", "command": "kavach gates read --hook"}]
      },
      {
        "matcher": "Write",
        "hooks": [{"type": "command", "command": "kavach gates enforcer --hook"}]
      },
      {
        "matcher": "Edit",
        "hooks": [{"type": "command", "command": "kavach gates enforcer --hook"}]
      }
    ],
    "Stop": [
      {
        "hooks": [{"type": "command", "command": "kavach session end"}]
      }
    ],
    "PreCompact": [
      {
        "hooks": [{"type": "command", "command": "kavach session compact"}]
      }
    ]
  }
}
```

### Minimal Configuration

For basic enforcement only:

```json
{
  "hooks": {
    "SessionStart": [
      {"hooks": [{"type": "command", "command": "kavach session init"}]}
    ],
    "PreToolUse": [
      {"matcher": "Write", "hooks": [{"type": "command", "command": "kavach gates enforcer --hook"}]},
      {"matcher": "Edit", "hooks": [{"type": "command", "command": "kavach gates enforcer --hook"}]}
    ]
  }
}
```

---

## OpenCode Configuration

### 1. Copy AGENTS.md

```bash
# Linux
cp configs/linux/AGENTS.md ~/.config/opencode/

# macOS
cp configs/darwin/AGENTS.md ~/Library/Application\ Support/opencode/

# Windows
Copy-Item configs\windows\AGENTS.md "$env:APPDATA\opencode\"
```

### 2. Configure Hooks

See `examples/hooks/` for OpenCode hook configuration.

---

## Verification

### Test Binary

```bash
kavach status                    # System health check
kavach memory bank               # Memory Bank summary
kavach agents                    # List agents
kavach skills                    # List skills
```

### Test Hooks

```bash
# Test session init
kavach session init

# Expected output:
# [SESSION]
# id: sess_abc123
# date: 2026-01-18
# project: kavach
# ...

# Test gate (Write/Edit)
echo '{"tool_name":"Write","tool_input":{"file_path":"test.go"}}' | kavach gates enforcer --hook

# Expected output:
# {"decision":"approve","reason":"Gate passed"}
```

### Test Memory Bank

```bash
kavach memory bank               # Project-scoped
kavach memory bank --all         # All projects
kavach memory kanban             # Kanban dashboard
```

---

## Troubleshooting

### "command not found"

Binary not in PATH:

```bash
# Linux/macOS
export PATH="$HOME/.local/bin:$PATH"
# Add to ~/.bashrc or ~/.zshrc for persistence
```

```powershell
# Windows
$env:PATH += ";$env:USERPROFILE\.local\bin"
```

### "unknown command"

Old binary version:

```bash
# Update
git pull origin main
just build
just install
```

### Hook Errors

1. Check binary path: `which kavach`
2. Verify binary permissions: `chmod +x ~/.local/bin/kavach`
3. Check settings.json syntax: `cat ~/.claude/settings.json | jq .`

### JSON Validation Failed

Wrong output schema:
- PreToolUse requires `decision` and `reason`
- SessionStart requires `hookEventName` and `additionalContext`

---

## Uninstall

### Linux/macOS

```bash
rm -f ~/.local/bin/kavach
rm -rf ~/.local/share/shared-ai/    # Memory Bank (optional)
```

### Windows

```powershell
Remove-Item -Force "$env:USERPROFILE\.local\bin\kavach.exe"
Remove-Item -Recurse -Force "$env:APPDATA\shared-ai"  # Memory Bank (optional)
```

---

## Upgrade

```bash
# Pull latest
cd kavach-rs
git pull origin main

# Rebuild
just build
just install

# Verify
kavach status
```

---

## Multiple CLI Tools

Kavach works with multiple AI coding assistants simultaneously:

| Tool | Config File | Hook Config |
|------|-------------|-------------|
| Claude Code | CLAUDE.md | settings.json |
| OpenCode | AGENTS.md | hooks.yaml |
| Other | Custom | Custom |

All tools share:
- Same binary (`kavach`)
- Same Memory Bank location
- Same enforcement rules

Only tool-specific:
- Config file format
- Hook wiring syntax

---

## Support

- Issues: [GitHub Issues](https://github.com/Wankhede-Brothers/kavach-rs/issues)
- Documentation: [docs/](../docs/)
- Examples: [examples/](../examples/)
