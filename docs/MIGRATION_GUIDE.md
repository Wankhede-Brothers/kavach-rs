# Migration Guide

## Overview

This guide helps you migrate to Kavach from previous setups or configure existing Claude Code installations.

---

## Quick Install

```bash
# Linux/macOS - One-line install
curl -fsSL https://raw.githubusercontent.com/Wankhede-Brothers/kavach-rs/main/install/install.sh | bash

# Or build from source
git clone https://github.com/Wankhede-Brothers/kavach-rs.git
cd kavach-rs
cd crates/kavach-cli && cargo build --release
cp target/release/kavach ~/.local/bin/
```

---

## Migration Scenarios

### From Bash Hooks to Kavach

If you have custom bash hooks, kavach replaces them with a single binary:

**Before (multiple bash scripts):**
```json
{
  "hooks": {
    "PreToolUse": [
      {"matcher": "Write", "hooks": [{"type": "command", "command": "~/.claude/hooks/enforcer.sh"}]},
      {"matcher": "Bash", "hooks": [{"type": "command", "command": "~/.claude/hooks/bash-gate.sh"}]}
    ]
  }
}
```

**After (single kavach binary):**
```json
{
  "hooks": {
    "SessionStart": [
      {"hooks": [{"type": "command", "command": "kavach session init"}]}
    ],
    "PreToolUse": [
      {"matcher": "Write", "hooks": [{"type": "command", "command": "kavach gates enforcer --hook"}]},
      {"matcher": "Edit", "hooks": [{"type": "command", "command": "kavach gates enforcer --hook"}]},
      {"matcher": "Bash", "hooks": [{"type": "command", "command": "kavach gates bash --hook"}]},
      {"matcher": "Task", "hooks": [{"type": "command", "command": "kavach gates ceo --hook"}]}
    ],
    "Stop": [
      {"hooks": [{"type": "command", "command": "kavach session end"}]}
    ]
  }
}
```

### From claude-gate to kavach

If you were using the old `claude-gate` binary:

1. **Uninstall claude-gate:**
   ```bash
   rm -f ~/.claude/bin/claude-gate
   ```

2. **Install kavach:**
   ```bash
   curl -fsSL https://raw.githubusercontent.com/Wankhede-Brothers/kavach-rs/main/install/install.sh | bash
   ```

3. **Update settings.json:** Replace all `claude-gate` references with `kavach`
   ```bash
   sed -i 's/claude-gate/kavach/g' ~/.claude/settings.json
   ```

### Fresh Installation

See [INSTALLATION.md](INSTALLATION.md) for complete installation guide.

---

## Hook Configuration

### Full Configuration (Recommended)

Update `~/.claude/settings.json`:

```json
{
  "hooks": {
    "SessionStart": [
      {"hooks": [{"type": "command", "command": "kavach session init"}]}
    ],
    "UserPromptSubmit": [
      {"hooks": [{"type": "command", "command": "kavach session init"}]}
    ],
    "PreToolUse": [
      {"matcher": "Task", "hooks": [{"type": "command", "command": "kavach gates ceo --hook"}]},
      {"matcher": "Bash", "hooks": [{"type": "command", "command": "kavach gates bash --hook"}]},
      {"matcher": "Read", "hooks": [{"type": "command", "command": "kavach gates read --hook"}]},
      {"matcher": "Write", "hooks": [{"type": "command", "command": "kavach gates enforcer --hook"}]},
      {"matcher": "Edit", "hooks": [{"type": "command", "command": "kavach gates enforcer --hook"}]}
    ],
    "Stop": [
      {"hooks": [{"type": "command", "command": "kavach session end"}]}
    ],
    "PreCompact": [
      {"hooks": [{"type": "command", "command": "kavach session compact"}]}
    ]
  }
}
```

### Minimal Configuration

For basic TABULA_RASA enforcement only:

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

## Memory Bank Migration

### From RPC Server to TOON Files

Old Memory Bank used an RPC server. New kavach uses TOON files directly.

**Memory Bank location:**
- Linux: `~/.local/share/shared-ai/memory/`
- macOS: `~/Library/Application Support/shared-ai/memory/`
- Windows: `%APPDATA%\shared-ai\memory\`

**Initialize new Memory Bank:**
```bash
mkdir -p ~/.local/share/shared-ai/memory/{decisions,patterns,research,kanban,proposals,roadmaps,graph,STM}
kavach memory bank  # Verify
```

### Project Isolation

Memory Bank is now project-scoped:

```bash
kavach memory bank              # Current project only
kavach memory bank --all        # All projects
kavach memory bank --project    # Project files only
```

---

## Verification

After migration:

```bash
# Check binary
which kavach
kavach status

# Check hooks
kavach session init

# Check Memory Bank
kavach memory bank

# Check agents
kavach agents

# Check skills
kavach skills
```

---

## Troubleshooting

### Hook Not Executing

```bash
# Test binary directly
kavach session init

# Check PATH
echo $PATH | grep -o '[^:]*bin[^:]*'
export PATH="$HOME/.local/bin:$PATH"
```

### Permission Denied

```bash
chmod +x ~/.local/bin/kavach
```

### Memory Bank Empty

```bash
# Initialize structure
mkdir -p ~/.local/share/shared-ai/memory/{decisions,patterns,research,kanban,proposals,roadmaps,graph,STM}
```

### JSON Validation Failed

Check output schema matches Claude Code requirements:
- PreToolUse: `{"decision":"approve|block","reason":"..."}`
- SessionStart: `{"hookEventName":"...","additionalContext":"..."}`

---

## Rollback

If you need to rollback:

```bash
# Remove kavach
rm -f ~/.local/bin/kavach

# Restore old settings.json
cp ~/.claude/settings.json.backup ~/.claude/settings.json
```

---

## Support

- Issues: [GitHub Issues](https://github.com/Wankhede-Brothers/kavach-rs/issues)
- Documentation: [docs/](../docs/)
