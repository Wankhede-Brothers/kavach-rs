# Examples

Configuration examples for `kavach` (Brahmastra Stack) across different AI coding assistants.

---

## Structure

```
examples/
├── settings/           # Claude Code configurations
│   ├── settings.json.example           # Full configuration
│   └── settings-minimal.json.example   # Minimal configuration
├── hooks/              # Hook configurations
│   ├── opencode.yaml.example           # OpenCode hooks
│   └── README.md                       # Hook documentation
├── modules/            # Lazy-load context modules (10 modules)
│   ├── amnesia.md           # CRITICAL: Memory Bank awareness
│   ├── tabula-rasa.md       # CRITICAL: Stale weights prevention
│   ├── date.md              # CRITICAL: Date injection
│   ├── compact.md           # PreCompact hook rules
│   ├── agents.md            # Agent hierarchy reference
│   ├── dace.md              # DACE micro-modular principles
│   ├── hooks.md             # Hook reference documentation
│   ├── reuse.md             # Code reusability patterns
│   ├── structure.md         # File structure patterns
│   └── tools.md             # Modern Rust/Zig CLI tools
├── agents/             # Agent definitions (6 agents)
│   ├── aegis-guardian.md    # L2 - Final verification
│   ├── backend-engineer.md  # L1 - Backend implementation
│   ├── ceo.md               # L0 - Orchestration
│   ├── frontend-engineer.md # L1 - Frontend implementation
│   ├── nlu-intent-analyzer.md # L-1 - Intent parsing
│   └── research-director.md # L0 - Research coordination
└── skills/             # Skill definitions (15 skills)
    ├── api-design/          # REST, gRPC, GraphQL
    ├── arch/                # System Design
    ├── cloud-infrastructure-mastery/  # K8s, Terraform, AWS
    ├── create-claude-components/      # Skills, Hooks, Agents
    ├── debug-like-expert/   # Systematic debugging
    ├── dsa/                 # Data Structures & Algorithms
    ├── frontend/            # React, Vue, Svelte
    ├── heal/                # 5-Layer code analysis
    ├── high-performance-data-processing/  # Parquet, Arrow, Polars
    ├── resume/              # Session recovery
    ├── rust/                # Rust Engineering
    ├── security/            # OWASP, Auth, Encryption
    ├── sql/                 # PostgreSQL & SQLx
    ├── sutra-protocol/      # Agent communication
    └── testing/             # Universal testing
```

---

## DACE Architecture

All examples follow **Dynamic Agentic Context Engineering (DACE)**:

```
DACE:CORE
  mode: lazy_load,skill_first,on_demand
  output_tokens: 2048
  max_lines: 100

DACE:PRINCIPLES
  lazy_load:    Load context on-demand, not upfront
  skill_first:  Use kavach binary before spawning agents
  on_demand:    Inject research/patterns only when needed
  no_hardcode:  WebSearch for current patterns
```

---

## Modules (Lost-in-Middle Mitigation)

Modules extend CLAUDE.md without bloating context. Critical modules are injected after `/compact` to prevent context loss.

### Critical Modules (Post-Compact Recovery)

| Module | Purpose | Injection |
|--------|---------|-----------|
| `amnesia.md` | Prevents "I have no memory" claims | Post-compact, Reinforcement |
| `tabula-rasa.md` | Prevents stale knowledge usage | Post-compact, Reinforcement |
| `date.md` | Ensures correct date awareness | Post-compact, Reinforcement |

### Lazy-Load Modules

| Module | Purpose | Trigger |
|--------|---------|---------|
| `compact.md` | Compact behavior rules | PreCompact hook |
| `agents.md` | Agent hierarchy reference | Task tool |
| `dace.md` | Micro-modular file principles | Write/Edit tools |
| `hooks.md` | Hook reference documentation | Reference only |
| `reuse.md` | Code reusability patterns | Write/Edit tools |
| `structure.md` | File structure patterns | Write/Edit tools |
| `tools.md` | Modern Rust/Zig CLI tools | Bash tool |

### Installation

```bash
# Copy modules to Claude config
mkdir -p ~/.claude/modules
cp examples/modules/*.md ~/.claude/modules/
```

See [modules/README.md](modules/README.md) for detailed documentation.

---

## Skill Structure (SP/3.0)

Each skill contains:

```
skill-name/
├── SKILL.md         # Main skill definition with KAVACH:DYNAMIC
└── references.toon  # Dynamic research topics (NO hardcoded content)
```

### SKILL.md Template

```markdown
---
name: skill-name
description: One-line description
license: MIT
compatibility: claude-code
metadata:
  category: category
  triggers: [keyword1, keyword2]
  protocol: SP/3.0
---

```toon
SKILL:skill-name
  protocol: SP/3.0
  dace: lazy_load,skill_first,on_demand
  triggers[N]: trigger1,trigger2
  goal: What it achieves
  success: Acceptance criteria
  fail: Anti-patterns

KAVACH:DYNAMIC
  context: kavach skills --get skill-name --inject
  references: ~/.claude/skills/skill-name/references.toon
  research: kavach memory bank | grep -i keyword
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded content - WebSearch for current patterns
  topics[N]
    TOPIC1: WebSearch "topic {YEAR}"
    TOPIC2: WebSearch "topic {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[N]
    kavach status (inject today's date)
    WebSearch topic BEFORE implementing
    WebFetch official docs
  forbidden: assumptions without verification
```
```

---

## Agent Hierarchy

```
[AGENTS:HIERARCHY]
Level -1: NLU Intent Analyzer (haiku)
  └── Parses ALL user requests → Sutra Protocol output

Level 0: Decision Makers (opus)
  ├── CEO - Orchestration, never writes code
  └── Research Director - Evidence-based findings

Level 1: Engineers (sonnet)
  ├── backend-engineer   - Rust, API
  └── frontend-engineer  - TypeScript, React

Level 1.5: Code Reviewer (sonnet)
  └── Post-implementation review

Level 2: Aegis Guardian (opus)
  └── Final verification, quality gate
```

---

## Skills Reference

| Skill | Category | Purpose | Research Gate |
|-------|----------|---------|---------------|
| api-design | API | REST, gRPC, GraphQL patterns | WebSearch "API design {YEAR}" |
| arch | Design | System design with numbers | WebSearch "system design {YEAR}" |
| cloud-infrastructure-mastery | Cloud | K8s, Terraform, AWS/GCP | WebSearch "kubernetes {YEAR}" |
| create-claude-components | Meta | Skills, Hooks, Agents creation | WebSearch "claude code {YEAR}" |
| debug-like-expert | Debug | Systematic investigation | WebSearch "debugging {YEAR}" |
| dsa | Algorithms | O(1) over O(n) | WebSearch "algorithms {YEAR}" |
| frontend | UI | React, Vue, Svelte, a11y | WebSearch "react {YEAR}" |
| heal | Quality | 5-layer code analysis | WebSearch "static analysis {YEAR}" |
| high-performance-data-processing | Data | Parquet, Arrow, Polars, Rayon | WebSearch "polars {YEAR}" |
| resume | Context | Session recovery from Memory Bank | kavach memory bank |
| rust | Language | Rust best practices | WebSearch "rust {YEAR}" |
| security | Security | OWASP, Auth, Encryption | WebSearch "OWASP {YEAR}" |
| sql | Database | PostgreSQL, Index Scan optimization | WebSearch "postgresql {YEAR}" |
| sutra-protocol | Protocol | 75-80% token reduction | kavach sutra --help |
| testing | QA | Test pyramid, coverage | WebSearch "testing {YEAR}" |

---

## Quick Start

### Claude Code

```bash
# 1. Install kavach
go install github.com/claude/cmd/kavach@latest
# OR build from source
cd kavach-go && go build -o kavach ./cmd/kavach
cp kavach ~/.local/bin/

# 2. Copy modules (Lost-in-Middle mitigation)
mkdir -p ~/.claude/modules
cp examples/modules/*.md ~/.claude/modules/

# 3. Copy agents
mkdir -p ~/.claude/agents
cp examples/agents/*.md ~/.claude/agents/

# 4. Copy skills
mkdir -p ~/.claude/skills
cp -r examples/skills/* ~/.claude/skills/

# 5. Copy settings
cp examples/settings/settings.json.example ~/.claude/settings.json

# 6. Update paths
sed -i "s|/home/USER|$HOME|g" ~/.claude/settings.json

# 7. Verify
kavach status
kavach agents
kavach skills
```

---

## Hook Reference

| Event | Gate | Purpose |
|-------|------|---------|
| SessionStart | `kavach session init` | Initialize session, inject date |
| UserPromptSubmit | `kavach session init` | Date injection on each prompt |
| PreToolUse:Task | `kavach gates ceo --hook` | Validate agent delegation |
| PreToolUse:Bash | `kavach gates bash --hook` | Sanitize commands |
| PreToolUse:Read | `kavach gates read --hook` | Block sensitive files |
| PreToolUse:Write | `kavach gates enforcer --hook` | Full pipeline validation |
| PreToolUse:Edit | `kavach gates enforcer --hook` | Full pipeline validation |
| Stop | `kavach session end` | Save session state |
| PreCompact | `kavach session compact` | Pre-compact save |

---

## Key Principles

### TABULA RASA
- Training cutoff: 2025-01
- Rule: WebSearch BEFORE writing code
- Forbidden: Assumptions without verification

### NO AMNESIA
- Memory Bank: `~/.local/shared/shared-ai/memory/`
- Query: `kavach memory bank`
- Forbidden: Claims of no memory access

### BINARY FIRST
- Rule: Use kavach commands BEFORE reading files
- Query: `kavach status && kavach memory bank`
- Forbidden: Spawning agents when binary command exists

---

## License

MIT
