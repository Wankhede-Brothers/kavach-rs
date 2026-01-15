# Memory Bank

Persistent memory across Claude Code sessions using TOON format with project isolation.

## Architecture

```
~/.local/shared/shared-ai/memory/       # Linux
~/Library/Application Support/.../       # macOS
%APPDATA%\shared-ai\memory\             # Windows

memory/
├── GOVERNANCE.toon        # Rules and policies
├── index.toon             # Project index
├── volatile.toon          # Session state (cleared on restart)
│
├── decisions/             # Architecture decisions
│   ├── global/            # Shared decisions
│   │   └── decisions.toon
│   └── {project}/         # Project-specific
│       └── decisions.toon
│
├── patterns/              # Code patterns
│   ├── global/
│   │   └── pending-tasks.toon  # Universal TODO rules
│   └── {project}/
│
├── research/              # Research with TTL (TABULA_RASA)
│   └── {project}/
│
├── kanban/                # 5-stage pipeline
│   └── {project}/         # Backlog→InProgress→Testing→Verified→Done
│       └── kanban.toon
│
├── proposals/             # Feature proposals
├── roadmaps/              # Project roadmaps
├── graph/                 # Knowledge graphs
└── STM/                   # Short-term memory
    └── scratchpad.toon    # Session scratchpad
```

---

## Project Isolation

Memory Bank queries are scoped to the active project:

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROJECT DETECTION                             │
│                    (Priority Order)                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. KAVACH_PROJECT env var                                      │
│     └── Explicit override                                        │
│                                                                  │
│  2. .git root detection                                          │
│     └── Walk up tree, find .git, use directory name             │
│                                                                  │
│  3. .claude/project.json marker                                  │
│     └── Read "name" or "project" field                          │
│                                                                  │
│  4. Memory Bank project match                                    │
│     └── Match working dir against known projects in kanban/     │
│                                                                  │
│  5. "global" fallback                                            │
│     └── Shared across all projects                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Query Scoping

```bash
kavach memory bank           # Active project + global ONLY
kavach memory bank --all     # ALL projects (explicit)
kavach memory bank --project # Project files only
```

---

## TOON Format

Text-Only Object Notation - 40% savings vs JSON.

### Example: decisions.toon

```toon
# Architecture Decisions - {project}

DECISION:auth-strategy
  date: 2026-01-16
  status: approved
  title: JWT Authentication
  context: Need stateless auth for API
  options[3]
    1. Session-based
    2. JWT tokens
    3. OAuth2 only
  decision: JWT tokens
  rationale: Stateless, scalable, standard
  consequences[2]
    + Horizontal scaling
    - Token revocation complexity

DECISION:database-choice
  date: 2026-01-15
  status: approved
  title: PostgreSQL with SQLx
  context: Need type-safe database access
  decision: PostgreSQL + SQLx
  rationale: Compile-time query checking
```

### Example: kanban.toon

```toon
# Kanban Board - {project}

KANBAN:active
  project: my-project
  updated: 2026-01-16

[BACKLOG]
TASK:impl-auth
  title: Implement authentication
  priority: high
  created: 2026-01-14

[IN_PROGRESS]
TASK:api-endpoints
  title: Design REST endpoints
  assignee: backend-engineer
  started: 2026-01-16

[TESTING]
TASK:user-model
  title: User model with SQLx
  tests: unit,integration

[VERIFIED]
TASK:db-schema
  title: Database schema
  verified_by: aegis-guardian
  date: 2026-01-15

[DONE]
TASK:project-setup
  title: Initialize project
  completed: 2026-01-14
```

---

## Kanban Pipeline

5-stage production pipeline with Aegis-Guard verification:

```
Backlog → In-Progress → Testing → Verified → Done
                           ↓          ↓
                      lint/bugs    algorithm
                      warnings     dead code
                      unit tests   suppressed
```

### Commands

```bash
kavach memory kanban              # View kanban board
kavach memory kanban --add        # Add task to backlog
kavach memory kanban --move       # Move task between stages
kavach memory kanban --project    # Project-specific board
```

---

## Fact TTL (Time-To-Live)

Research facts automatically expire based on category:

```toon
TTL:CATEGORIES
  syntax: 7d      # Code patterns
  config: 5d      # Configuration
  behavior: 3d    # Runtime behavior
  migration: 30d  # Breaking changes
  security: 1d    # CVEs, advisories

RULE: IF today > verified + ttl → STALE
ACTION: Re-research required
```

### Example: research.toon

```toon
# Research Cache - {project}

FACT:axum-0.8-router
  category: syntax
  verified: 2026-01-16
  ttl: 7d
  expires: 2026-01-23
  src: https://docs.rs/axum/0.8/
  val: Router::new().route("/", get(handler))

FACT:sqlx-query-macro
  category: syntax
  verified: 2026-01-15
  ttl: 7d
  expires: 2026-01-22
  src: https://docs.rs/sqlx/
  val: sqlx::query!("SELECT * FROM users")
```

---

## Event Bus Integration

Memory operations publish events for telemetry:

```go
// Memory write event
eventBus.Publish(events.EventMemoryWrite, "kavach", map[string]interface{}{
    "category": "decisions",
    "key":      "auth-strategy",
    "project":  "my-project",
})
```

### Event Types

| Event | Trigger | Data |
|-------|---------|------|
| `EventMemoryWrite` | Memory Bank write | category, key, project |
| `EventSessionStart` | Session init | session_id, project |
| `EventAgentInvoke` | Agent spawned | agent, task |
| `EventSkillInvoke` | Skill executed | skill, args |

---

## File Locking

Concurrent Memory Bank writes are protected:

```go
lockMgr := lock.GetLockManager()

// Acquire with timeout
err := lockMgr.AcquireWithTimeout(path, 5*time.Second)
if err != nil {
    return fmt.Errorf("lock acquisition failed: %w", err)
}
defer lockMgr.Release(path)

// Safe to write
bank.SaveFile(path, doc)
```

This prevents race conditions when multiple Claude Code sessions write to the same file.

---

## Commands

```bash
# Query Memory Bank
kavach memory bank               # Project-scoped summary
kavach memory bank --all         # All projects
kavach memory bank --category X  # Specific category

# Write to Memory Bank
kavach memory write --category decisions --key auth-strategy --value "JWT"

# Kanban operations
kavach memory kanban             # View board
kavach memory kanban --status    # Pipeline status

# Short-term memory
kavach memory stm                # Session scratchpad
kavach memory stm --set key val  # Set scratchpad value
```

---

## Category Reference

| Category | Purpose | TTL | Scope |
|----------|---------|-----|-------|
| `decisions` | Architecture decisions | Permanent | Project |
| `patterns` | Code patterns | 7d | Project/Global |
| `research` | Framework research | 1-30d | Project |
| `kanban` | Task board | Permanent | Project |
| `proposals` | Feature proposals | Permanent | Project |
| `roadmaps` | Project roadmaps | Permanent | Project |
| `graph` | Knowledge graph | Permanent | Project |
| `STM` | Session scratchpad | Session | Global |

---

## GOVERNANCE.toon

Root governance rules for the Memory Bank:

```toon
# Memory Bank Governance

GOVERNANCE:rules
  version: 4.0
  protocol: SP/3.0

[ENFORCEMENT]
TABULA_RASA: WebSearch BEFORE code
NO_AMNESIA: Query Memory Bank for context
PROJECT_ISOLATION: Scope queries to active project
MICRO_MODULAR: Max 100 lines per file

[TTL:DEFAULTS]
syntax: 7d
config: 5d
behavior: 3d
migration: 30d
security: 1d

[LOCK]
timeout: 5s
retry: 3
backoff: exponential
```

---

## Migration from JSON

If migrating from older JSON-based Memory Bank:

```bash
# Convert JSON to TOON
kavach migrate --from json --to toon

# Verify migration
kavach memory bank --validate
```

TOON format provides:
- 40% file size reduction
- Human-readable structure
- Git-friendly diffs
- Native block notation
