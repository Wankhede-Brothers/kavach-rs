# Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              KAVACH SYSTEM                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         AI CODING ASSISTANTS                         │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │   │
│  │  │ Claude Code │    │  OpenCode   │    │  Other CLI  │             │   │
│  │  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘             │   │
│  │         └──────────────────┼──────────────────┘                     │   │
│  │                            │                                        │   │
│  │                     HOOK SYSTEM                                     │   │
│  │    SessionStart → UserPromptSubmit → PreToolUse → Stop             │   │
│  └────────────────────────────┼────────────────────────────────────────┘   │
│                               │                                             │
│                               ▼                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        KAVACH BINARY                                 │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │                                                                      │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│  │  │    GATES     │  │    MEMORY    │  │   SESSION    │              │   │
│  │  │  enforcer    │  │     bank     │  │     init     │              │   │
│  │  │  ceo         │  │    kanban    │  │   validate   │              │   │
│  │  │  bash        │  │    write     │  │     end      │              │   │
│  │  │  read        │  │     stm      │  │   compact    │              │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘              │   │
│  │                                                                      │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│  │  │   AGENTS     │  │    SKILLS    │  │    STATUS    │              │   │
│  │  │    list      │  │     list     │  │    health    │              │   │
│  │  │    get       │  │     get      │  │   enforce    │              │   │
│  │  │   inject     │  │    inject    │  │    state     │              │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘              │   │
│  │                                                                      │   │
│  └────────────────────────────┼────────────────────────────────────────┘   │
│                               │                                             │
│         ┌─────────────────────┼─────────────────────┐                      │
│         ▼                     ▼                     ▼                      │
│  ┌─────────────┐       ┌─────────────┐       ┌─────────────┐              │
│  │ MEMORY BANK │       │  EVENT BUS  │       │  FILE LOCK  │              │
│  │ (TOON files)│       │ (telemetry) │       │(concurrency)│              │
│  │             │       │             │       │             │              │
│  │ decisions/  │       │ SessionStart│       │ Acquire()   │              │
│  │ patterns/   │       │ MemoryWrite │       │ Release()   │              │
│  │ research/   │       │ AgentInvoke │       │ Timeout()   │              │
│  │ kanban/     │       │ SkillInvoke │       │             │              │
│  └─────────────┘       └─────────────┘       └─────────────┘              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Hook Flow

```
USER INPUT
    │
    ▼
┌─────────────────┐
│ SessionStart    │──► kavach session init
└────────┬────────┘    • Inject date (TABULA_RASA)
         │             • Load session state
         │             • Publish EventSessionStart
         ▼
┌─────────────────┐
│UserPromptSubmit │──► kavach session init
└────────┬────────┘    • Re-inject date on each prompt
         │
         ▼
┌─────────────────┐
│ PreToolUse      │──► kavach gates <gate> --hook
└────────┬────────┘
         │
    ┌────┴────┬─────────┬─────────┐
    ▼         ▼         ▼         ▼
┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐
│ Task  │ │ Bash  │ │ Read  │ │Write/ │
│       │ │       │ │       │ │ Edit  │
│ceo    │ │bash   │ │read   │ │enforcer│
│gate   │ │gate   │ │gate   │ │gate   │
└───────┘ └───────┘ └───────┘ └───────┘
    │
    ▼
┌─────────────────┐
│ TOOL EXECUTION  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ PostToolUse     │──► Memory write (with lock + events)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Stop            │──► kavach session end
└─────────────────┘    • Persist state
```

---

## Agent Hierarchy (SP/3.0)

```
                    ┌─────────────────┐
                    │      USER       │
                    └────────┬────────┘
                             │
                             ▼
              ┌──────────────────────────┐
              │  L-1: NLU Intent Analyzer │  (haiku - 20x faster)
              │  Parses ALL user requests │
              │  Outputs Sutra Protocol   │
              └────────────┬─────────────┘
                           │
                           ▼
         ┌─────────────────────────────────────┐
         │           L0: DECISION MAKERS        │  (opus)
         │  ┌─────────────┐  ┌───────────────┐ │
         │  │     CEO     │  │   Research    │ │
         │  │ Orchestrator│  │   Director    │ │
         │  │ Never writes│  │ Evidence-based│ │
         │  │    code     │  │   findings    │ │
         │  └──────┬──────┘  └───────┬───────┘ │
         └─────────┼─────────────────┼─────────┘
                   │                 │
                   ▼                 ▼
    ┌──────────────────────────────────────────────┐
    │              L1: ENGINEERS                    │  (sonnet)
    │  ┌──────────┐ ┌──────────┐                   │
    │  │ Backend  │ │ Frontend │                   │
    │  │ Rust,API │ │ TS,React │                   │
    │  └──────────┘ └──────────┘                   │
    └──────────────────┬───────────────────────────┘
                       │
                       ▼
         ┌─────────────────────────────┐
         │   L1.5: CODE REVIEWER       │  (sonnet)
         │   Post-implementation       │
         └─────────────┬───────────────┘
                       │
                       ▼
         ┌─────────────────────────────┐
         │   L2: AEGIS GUARDIAN        │  (opus)
         │   Final verification        │
         │   Quality gate              │
         └─────────────┬───────────────┘
                       │
                       ▼
              ┌────────────────┐
              │    <promise>   │
              │PRODUCTION_READY│
              │   </promise>   │
              └────────────────┘
```

---

## Gate Pipeline

```
TOOL REQUEST (stdin JSON)
      │
      ▼
┌─────────────────────────────────────────────────────────────────┐
│                       ENFORCER GATE                              │
│                                                                  │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐    │
│  │  RESEARCH GATE │  │  PATTERN CHECK │  │ FORBIDDEN CHECK│    │
│  │  (pkg/agentic) │  │  (pkg/patterns)│  │ (phrases)      │    │
│  │                │  │                │  │                │    │
│  │ • Framework    │  │ • IsCodeFile() │  │ • No guessing  │    │
│  │   detection    │  │ • IsBlocked()  │  │ • No assuming  │    │
│  │ • Research     │  │ • IsSensitive()│  │ • Verify first │    │
│  │   required?    │  │ • ValidAgent() │  │                │    │
│  └────────────────┘  └────────────────┘  └────────────────┘    │
│                                                                  │
│  Decision: APPROVE (silent) | BLOCK (with reason)               │
└─────────────────────────────────────────────────────────────────┘
      │
      ├── Task Tool ──► CEO Gate (validate agent hierarchy)
      │
      ├── Bash Tool ──► Bash Gate (command sanitization)
      │
      ├── Read Tool ──► Read Gate (sensitive file blocking)
      │
      └── Write/Edit ──► Enforcer Gate (TABULA_RASA check)
```

---

## Project Isolation

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

Memory Bank Query Scoping:
┌─────────────────────────────────────────────────────────────────┐
│  kavach memory bank           → Active project + global ONLY    │
│  kavach memory bank --all     → ALL projects (explicit)         │
│  kavach memory bank --project → Project files only              │
└─────────────────────────────────────────────────────────────────┘
```

---

## Memory Bank Structure

```
~/.local/shared/shared-ai/memory/       # Linux
~/Library/Application Support/.../      # macOS
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

## Event Bus

```go
// Event Types
const (
    EventMemoryWrite   = "memory_write"
    EventSessionStart  = "session_start"
    EventAgentInvoke   = "agent_invoke"
    EventSkillInvoke   = "skill_invoke"
    EventError         = "error"
)

// Event Structure
type Event struct {
    Type      EventType
    Source    string      // "kavach"
    SessionID string
    Timestamp time.Time
    Payload   interface{}
}

// Usage
eventBus := events.GetEventBus()
eventBus.Publish(events.EventMemoryWrite, "kavach", map[string]interface{}{
    "category": "decisions",
    "key":      "auth-strategy",
    "project":  "my-project",
})
```

---

## File Locking

```go
// Prevent concurrent Memory Bank writes
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

---

## DACE (Dynamic Agentic Context Engineering)

```
DACE:CORE
  mode: lazy_load,skill_first,on_demand
  output_tokens: 2048
  max_lines: 100 (hard block)
  warn_lines: 50 (suggest split)

DACE:COMMAND_PRIORITY
  1. kavach [command]     # FIRST: Binary commands
  2. Memory Bank TOON     # SECOND: If binary insufficient
  3. Read specific file   # LAST: Only when explicitly needed
  NEVER: Explore agents for known patterns

DACE:SKILL_STRUCTURE
  SKILL.md          # Main skill with KAVACH:DYNAMIC block
  references.toon   # Dynamic WebSearch queries (NO hardcoded content)
```

---

## Enforcement Principles

### TABULA_RASA
```
Training Cutoff: 2025-01
Today: $(kavach status | grep today)

RULE: LLM weights are STALE
ACTION: WebSearch BEFORE code
BLOCK: Code generation without research
FORBIDDEN: assumptions without verification
```

### NO_AMNESIA
```
RULE: Query Memory Bank for context
ACTION: kavach memory bank
BLOCK: Claims of no memory access
```

### PROJECT_ISOLATION
```
RULE: Scope queries to active project
ACTION: DetectProject() with priority-based detection
BENEFIT: No context pollution across parallel sessions
```

### MICRO_MODULAR
```
MAX_LINES: 100 per file (150 soft limit)
RULE: Single responsibility per file
PATTERN: types.go, loader.go, output.go, etc.
```

---

## Shared Packages

| Package | Purpose | Used By |
|---------|---------|---------|
| `pkg/hook` | Hook I/O (stdin/stdout) | All gates |
| `pkg/toon` | TOON parser/writer | Memory operations |
| `pkg/patterns` | Dynamic pattern matching | Enforcer gate |
| `pkg/enforce` | Session state management | Session commands |
| `pkg/agentic` | Research gate (TABULA_RASA) | Enforcer gate |
| `pkg/dsa` | LRU cache, data structures | Memory bank |
| `pkg/util` | Path utilities, detection | All commands |
| `events/` | Event bus (telemetry) | Session, memory, agents |
| `lock/` | File locking (concurrency) | Memory write |
