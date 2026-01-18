# Project Structure

**Universal enforcement layer for AI coding assistants**

Part of the **Brahmastra Stack**: Kavach CLI + Sutra Protocol (SP/1.0) + TOON Format + DACE

---

## Architecture Principle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         UNIFIED INFRASTRUCTURE                               │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                        kavach binary                                  │   │
│  │                     (Single Installation)                             │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                     Memory Bank (TOON format)                         │   │
│  │                     Project-Isolated Storage                          │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│         ┌──────────────────────────┼──────────────────────────┐             │
│         ▼                          ▼                          ▼             │
│  ┌─────────────┐           ┌─────────────┐           ┌─────────────┐       │
│  │  Event Bus  │           │  File Lock  │           │   Gates     │       │
│  │ (telemetry) │           │(concurrency)│           │ (enforce)   │       │
│  └─────────────┘           └─────────────┘           └─────────────┘       │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
                                     │
        ┌────────────────────────────┼────────────────────────────┐
        │                            │                            │
        ▼                            ▼                            ▼
┌───────────────┐          ┌───────────────┐          ┌───────────────┐
│  Claude Code  │          │   OpenCode    │          │  Other CLI    │
│  (CLAUDE.md)  │          │  (AGENT.md)   │          │  (hooks)      │
└───────────────┘          └───────────────┘          └───────────────┘
```

---

## Directory Layout

```
kavach/
│
├── cmd/kavach/                    # SOURCE: Main Binary
│   ├── main.go                    # Entry point (~80 lines)
│   ├── go.mod                     # Module definition
│   └── internal/
│       └── commands/              # Cobra command tree
│           ├── root.go            # Root + dispatch
│           ├── gates/             # Gate commands (4 core)
│           │   ├── register.go
│           │   ├── enforcer.go    # Full pipeline (TABULA_RASA)
│           │   ├── ceo.go         # Task orchestration
│           │   ├── bash.go        # Command sanitization
│           │   └── read.go        # File access control
│           ├── memory/            # Memory commands
│           │   ├── register.go
│           │   ├── bank.go        # Project-scoped queries
│           │   ├── write.go       # With file locking + events
│           │   ├── kanban.go      # 5-stage pipeline
│           │   └── stm.go         # Short-term memory
│           ├── session/           # Session commands
│           │   ├── register.go
│           │   ├── init.go        # With event publishing
│           │   ├── end.go         # Persist state
│           │   ├── compact.go     # Pre-compact save
│           │   └── validate.go    # State validation
│           ├── agents/            # Agent management (micro-modular)
│           │   ├── types.go       # Agent struct definitions
│           │   ├── builtin.go     # 10 built-in agents
│           │   ├── discover.go    # Find agents in paths
│           │   ├── loader.go      # Load agent definitions
│           │   ├── inject.go      # Memory Bank injection
│           │   └── output.go      # TOON/Sutra output
│           ├── skills/            # Skill management (micro-modular)
│           │   ├── types.go       # Skill struct definitions
│           │   ├── builtin.go     # 15 built-in skills
│           │   ├── discover.go    # Find skills in paths
│           │   ├── loader.go      # Load skill definitions
│           │   ├── inject.go      # Memory Bank injection
│           │   └── output.go      # TOON/Sutra output
│           └── status.go          # System health
│
├── shared/                        # SOURCE: Shared Packages
│   ├── go.mod
│   ├── pkg/
│   │   ├── hook/                  # Hook I/O utilities
│   │   │   ├── input.go           # ReadHookInput()
│   │   │   └── output.go          # Approve(), Block()
│   │   ├── toon/                  # TOON format parser
│   │   │   ├── parser.go          # Parse TOON files
│   │   │   └── writer.go          # Write TOON files
│   │   ├── patterns/              # Dynamic pattern matching
│   │   │   ├── patterns.go        # IsCodeFile(), IsBlocked()
│   │   │   └── forbidden.go       # Forbidden phrase detection
│   │   ├── enforce/               # Enforcement logic
│   │   │   ├── session.go         # Session state management
│   │   │   └── rules.go           # Enforcement rules
│   │   ├── agentic/               # Research gate (TABULA_RASA)
│   │   │   ├── research.go        # Framework detection
│   │   │   └── gate.go            # Research requirement check
│   │   ├── dsa/                   # Data structures
│   │   │   └── lru.go             # LRU cache
│   │   └── util/                  # Common utilities
│   │       ├── platform.go        # Platform detection
│   │       ├── paths.go           # Path utilities
│   │       ├── project.go         # Project detection
│   │       └── time.go            # Date utilities
│   ├── events/                    # Event bus (telemetry)
│   │   ├── bus.go                 # Event bus singleton
│   │   └── types.go               # Event type definitions
│   └── lock/                      # File locking (concurrency)
│       └── filelock.go            # Lock manager
│
├── examples/                      # EXAMPLES: Reference implementations
│   ├── agents/                    # 6 agent definitions
│   │   ├── ceo/
│   │   │   └── AGENT.md           # CEO orchestrator
│   │   ├── aegis-guardian/
│   │   │   └── AGENT.md           # Final verification
│   │   ├── backend-engineer/
│   │   │   └── AGENT.md           # Rust/API engineer
│   │   ├── frontend-engineer/
│   │   │   └── AGENT.md           # TypeScript/React engineer
│   │   ├── nlu-intent-analyzer/
│   │   │   └── AGENT.md           # Fast intent parser
│   │   └── research-director/
│   │       └── AGENT.md           # Evidence-based research
│   └── skills/                    # 15 skill definitions
│       ├── api-design/
│       │   ├── SKILL.md           # Main skill (KAVACH:DYNAMIC)
│       │   └── references.toon    # Dynamic WebSearch queries
│       ├── arch/
│       ├── cloud-infrastructure-mastery/
│       ├── create-claude-components/
│       ├── debug-like-expert/
│       ├── dsa/
│       ├── frontend/
│       ├── heal/
│       ├── high-performance-data-processing/
│       ├── resume/
│       ├── rust/
│       ├── security/
│       ├── sql/
│       ├── sutra-protocol/
│       └── testing/
│           ├── SKILL.md
│           └── references.toon
│
├── configs/                       # SETUP: Platform Configs
│   ├── linux/
│   │   └── settings.json          # Claude Code hooks
│   ├── darwin/
│   │   └── settings.json
│   └── windows/
│       └── settings.json
│
├── templates/memory/              # SETUP: Memory Bank Templates
│   ├── GOVERNANCE.toon            # Root governance rules
│   ├── index.toon                 # Project index
│   ├── volatile.toon              # Session state template
│   └── {category}/
│       └── TEMPLATE.toon          # Category template
│
├── install/                       # SETUP: Installation
│   ├── install.sh                 # Linux/macOS
│   └── install.ps1                # Windows PowerShell
│
├── docs/                          # DOCS: Documentation
│   ├── ARCHITECTURE.md            # System design
│   ├── HOOKS.md                   # Hook reference
│   ├── MEMORY_BANK.md             # Memory Bank schema
│   ├── STRUCTURE.md               # This file
│   ├── API.md                     # Command reference
│   └── INSTALLATION.md            # Install guide
│
├── go.work                        # Go workspace
├── Makefile                       # Build system
├── README.md                      # Project overview
├── CLAUDE.md                      # Project instructions
└── LICENSE
```

---

## Directory Categories

| Category | Directories | Purpose |
|----------|-------------|---------|
| **SOURCE** | `cmd/`, `shared/` | Go source code |
| **EXAMPLES** | `examples/agents/`, `examples/skills/` | Reference implementations |
| **SETUP** | `configs/`, `install/`, `templates/` | Installation & configuration |
| **DOCS** | `docs/` | Documentation |

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

---

## Micro-Modular Architecture

Each command group follows micro-modular principles (max 100 lines per file):

### Agents Package
```
agents/
├── types.go       # Agent struct, Level enum
├── builtin.go     # 10 built-in agent definitions
├── discover.go    # FindAgents(), search paths
├── loader.go      # LoadAgent(), parse AGENT.md
├── inject.go      # InjectMemoryBank(), context injection
└── output.go      # OutputTOON(), OutputSutra()
```

### Skills Package
```
skills/
├── types.go       # Skill struct, Trigger enum
├── builtin.go     # 15 built-in skill definitions
├── discover.go    # FindSkills(), search paths
├── loader.go      # LoadSkill(), parse SKILL.md
├── inject.go      # InjectMemoryBank(), context injection
└── output.go      # OutputTOON(), OutputSutra()
```

---

## Platform Paths

| Platform | Binary | Memory Bank |
|----------|--------|-------------|
| Linux | `~/.local/bin/kavach` | `~/.local/shared/shared-ai/memory/` |
| macOS | `~/.local/bin/kavach` | `~/Library/Application Support/shared-ai/memory/` |
| Windows | `%USERPROFILE%\.local\bin\kavach.exe` | `%APPDATA%\shared-ai\memory\` |

---

## Build Commands

```bash
# Build
make build              # Build kavach binary
make install            # Build + install to ~/.local/bin
make release            # Cross-platform release

# Development
make workspace          # Initialize go.work
make test               # Run tests
make lint               # Run linter

# Cleanup
make clean              # Remove build artifacts
```

---

## Installation

### Linux/macOS
```bash
curl -fsSL https://raw.githubusercontent.com/Wankhede-Brothers/kavach-go/main/install/install.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/Wankhede-Brothers/kavach-go/main/install/install.ps1 | iex
```

### From Source
```bash
git clone https://github.com/Wankhede-Brothers/kavach-go.git
cd kavach-go
make install
```

---

## DACE Architecture

Skills follow Dynamic Agentic Context Engineering:

```toon
SKILL:example
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand

KAVACH:DYNAMIC
  # Binary commands for dynamic context
  context: kavach skills --get example --inject
  references: ~/.claude/skills/example/references.toon
  research: kavach memory bank | grep -i keyword
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded content - WebSearch for current patterns
```

### DACE Principles
- **lazy_load**: Load context on-demand
- **skill_first**: Use kavach binary before spawning agents
- **on_demand**: Inject research only when needed
- **no_hardcode**: WebSearch for current patterns

---

## File Size Targets

| Component | Max Lines | Purpose |
|-----------|-----------|---------|
| main.go | 80 | Entry point |
| root.go | 100 | Dispatcher |
| Each gate | 120 | Gate logic |
| Each memory | 100 | Memory ops |
| Each session | 100 | Session ops |
| types.go | 60 | Type definitions |
| register.go | 40 | Command registration |
