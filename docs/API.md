# API Reference

Complete command reference for `kavach`.

---

## Command Structure

```
kavach <group> <command> [flags]
```

**Groups:**
- `gates` - Enforcement gates
- `memory` - Memory bank operations
- `session` - Session management

**Top-level:**
- `status` - System status
- `agents` - List agents
- `skills` - List skills
- `--version` - Show version

---

## Gates

Enforcement gates for tool validation.

### enforcer

Full pipeline enforcement (recommended for Write/Edit).

```bash
# Hook mode
echo '{"tool_name":"Write","tool_input":{"file_path":"main.go"}}' | kavach gates enforcer --hook
```

**Pipeline:** intent -> ceo -> quality -> aegis

**Output:**
```json
{"decision": "approve", "reason": "Gate passed"}
```

---

### ceo

Task orchestration gate. Validates agent delegation.

```bash
echo '{"tool_name":"Task","tool_input":{"subagent_type":"backend"}}' | kavach gates ceo --hook
```

**Validates:**
- Agent type exists
- Proper delegation chain
- NLU intent alignment

---

### bash

Command sanitization gate.

```bash
echo '{"tool_name":"Bash","tool_input":{"command":"rm -rf /"}}' | kavach gates bash --hook
```

**Blocks:**
- Destructive commands (`rm -rf`, `dd`, etc.)
- System modification (`chmod 777`, etc.)
- Credential access (`cat ~/.ssh/id_rsa`)

---

### read

File access control gate.

```bash
echo '{"tool_name":"Read","tool_input":{"file_path":"/etc/passwd"}}' | kavach gates read --hook
```

**Blocks:**
- System files (`/etc/passwd`, `/etc/shadow`)
- Credential files (`.env`, `*.pem`, `*.key`)
- SSH keys

---

### intent

Intent classification (UserPromptSubmit).

```bash
echo '{"tool_name":"prompt","tool_input":{"prompt":"fix the bug"}}' | kavach gates intent --hook
```

**Output:**
```json
{
  "hookEventName": "UserPromptSubmit",
  "additionalContext": "[INTENT]\nclassified: fix\ndate: 2026-01-16\nstatus: ok\n"
}
```

**Classifications:**
| Intent | Trigger Words |
|--------|---------------|
| implement | implement, create, build, add, write |
| fix | fix, bug, error, broken, issue |
| refactor | refactor, improve, optimize, clean |
| research | research, find, search, explore |
| question | what, how, why, where, when, ? |
| status | status, progress, state |
| review | review, check, verify, validate |
| general | (fallback) |

---

### ast

AST syntax validation for Edit operations.

```bash
echo '{"tool_name":"Edit","tool_input":{"file_path":"main.go"}}' | kavach gates ast --hook
```

**Validates:**
- Go: `go/parser`
- Python: AST parse
- JavaScript/TypeScript: ESLint parse

---

### research

TABULA_RASA enforcement.

```bash
echo '{"tool_name":"Write","tool_input":{"file_path":"new.go"}}' | kavach gates research --hook
```

**Rule:** WebSearch must precede code generation.

**Blocks:** Code creation without prior research evidence.

---

### content

Content validation gate.

```bash
echo '{"tool_name":"Write","tool_input":{"content":"..."}}' | kavach gates content --hook
```

**Validates:**
- No hardcoded secrets
- No TODO/FIXME in production code
- Proper error handling

---

### quality

Code quality chain (AST + lint).

```bash
echo '{"tool_name":"Edit","tool_input":{"file_path":"main.go"}}' | kavach gates quality --hook
```

**Pipeline:** ast -> lint -> verify

---

### skill

Skill validation gate.

```bash
echo '{"tool_name":"Task","tool_input":{"skill":"commit"}}' | kavach gates skill --hook
```

**Validates:** Skill exists in skill registry.

---

### lint

Lint checks only.

```bash
echo '{"tool_name":"Edit","tool_input":{"file_path":"main.go"}}' | kavach gates lint --hook
```

---

## Memory

Memory bank operations.

### bank

Query project memory.

```bash
# Full memory summary
kavach memory bank

# Health check with stats
kavach memory bank --status

# Project-specific
kavach memory bank --project myproject

# Reindex
kavach memory bank --scan
```

**Output (TOON format):**
```
[MEMORY_BANK]
project: my-project
date: 2026-01-16
status: healthy

[DECISIONS]
count: 12
latest: API restructure

[PATTERNS]
count: 8
latest: Error handling pattern

[RESEARCH]
count: 5
fresh: 3
stale: 2
```

---

### kanban

Sprint board management.

```bash
# Full board
kavach memory kanban

# Status summary
kavach memory kanban --status
```

**Output:**
```
[KANBAN]
sprint: 2026-W03
status: active

[TODO]
- Implement auth endpoints
- Add unit tests

[IN_PROGRESS]
- API documentation

[DONE]
- Database schema
- User model
```

---

### write

Write to memory bank.

```bash
kavach memory write --category decisions --key auth --value "JWT with refresh tokens"
```

**Categories:**
- `decisions` - Architecture decisions
- `patterns` - Code patterns
- `research` - Research findings (with TTL)
- `kanban` - Sprint tasks
- `proposals` - Feature proposals
- `roadmaps` - Project roadmaps

---

### stm

Short-term memory (scratchpad).

```bash
# Read
kavach memory stm

# Write
kavach memory stm --set "Current task: implement auth"
```

---

### inject

RPC context injection.

```bash
kavach memory inject --context "API context loaded"
```

---

### spec

Spec injection for Task agents.

```bash
kavach memory spec --agent backend --task "Implement REST endpoints"
```

---

## Session

Session lifecycle management.

### init

Initialize session.

```bash
kavach session init
```

**Output (TOON format):**
```
[SESSION]
id: abc123
date: 2026-01-16
project: my-project
cutoff: 2025-01

[ENFORCE]
TABULA_RASA: WebSearch BEFORE code
NO_AMNESIA: Query memory bank
DATE_INJECTION: 2026-01-16
```

---

### end

End session, persist state.

```bash
kavach session end
```

---

### compact

Pre-compact save (~700 tokens).

```bash
kavach session compact
```

**Use:** Before context compaction.

---

### resume

Post-compact resume (~300 tokens).

```bash
kavach session resume
```

**Use:** After context compaction to restore state.

---

### validate

Validate session state.

```bash
kavach session validate
```

---

### save

Save session state explicitly.

```bash
kavach session save
```

---

## System Commands

### status

System health check.

```bash
kavach status
```

**Output:**
```
[STATUS]
binary: kavach v0.1.0
memory: healthy (12 entries)
session: active (abc123)
platform: linux
```

---

### agents

List available agents.

```bash
# List all
kavach agents

# Get specific
kavach agents --get ceo
```

**Agents:**
| Level | Agent | Model | Purpose |
|-------|-------|-------|---------|
| L-1 | nlu-intent-analyzer | haiku | Fast intent parsing |
| L0 | ceo | opus | Orchestration |
| L0 | research-director | opus | Research coordination |
| L1 | backend-engineer | sonnet | Backend implementation |
| L1 | frontend-engineer | sonnet | Frontend implementation |
| L1 | devops-engineer | sonnet | Infrastructure |
| L1 | security-engineer | sonnet | Security review |
| L1 | qa-lead | sonnet | Test strategy |
| L1.5 | code-reviewer | sonnet | Code review |
| L2 | aegis-guardian | opus | Final verification |

---

### skills

List available skills.

```bash
# List all
kavach skills

# Get specific
kavach skills --get commit
```

---

### --version

Show version.

```bash
kavach --version
```

---

## Hook I/O Format

### Input (stdin)

All hooks receive JSON on stdin:

```json
{
  "tool_name": "Write",
  "tool_input": {
    "file_path": "/path/to/file.go",
    "content": "package main..."
  }
}
```

### Output (stdout)

**PreToolUse - Approve:**
```json
{
  "decision": "approve",
  "reason": "Gate passed"
}
```

**PreToolUse - Block:**
```json
{
  "decision": "block",
  "reason": "TABULA_RASA: WebSearch required"
}
```

**PreToolUse - Modify:**
```json
{
  "decision": "modify",
  "reason": "Context injected",
  "additionalContext": "[CONTEXT]\ndate: 2026-01-16\n"
}
```

**UserPromptSubmit:**
```json
{
  "hookEventName": "UserPromptSubmit",
  "additionalContext": "[INTENT]\nclassified: implement\ndate: 2026-01-16\n"
}
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (with message) |

---

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `KAVACH_DEBUG` | Enable debug output |
| `KAVACH_MEMORY` | Override memory path |
| `KAVACH_CONFIG` | Override config path |
