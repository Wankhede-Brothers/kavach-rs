# Modules - SP/1.0 Lazy-Load System

Modules are lazy-loaded context files that extend CLAUDE.md without bloating the initial context window.

## Lost-in-Middle Mitigation

These modules implement the **Attention Sink + Recency Anchor** pattern to combat the "Lost in the Middle" problem in LLMs.

```
┌─────────────────────────────────────────────────────────────────┐
│                    MODULE ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  CLAUDE.md (Always loaded - Sandwich Pattern)                  │
│  ├── TOP: Critical rules (Attention Sink)                      │
│  ├── MIDDLE: Module references (Lazy-load)                     │
│  └── BOTTOM: Critical rules repeated (Recency Anchor)          │
│                                                                 │
│  ~/.claude/modules/ (Lazy-loaded on demand)                    │
│  ├── amnesia.md      → Post-compact: Memory Bank awareness     │
│  ├── tabula-rasa.md  → Post-compact: Stale weights prevention  │
│  ├── date.md         → Post-compact: Date injection            │
│  ├── compact.md      → PreCompact hook reference               │
│  ├── agents.md       → Task tool: Agent hierarchy              │
│  ├── dace.md         → Write/Edit: DACE principles             │
│  ├── hooks.md        → Reference: Hook documentation           │
│  ├── reuse.md        → Write/Edit: Reusability patterns        │
│  ├── structure.md    → Write/Edit: File structure patterns     │
│  └── tools.md        → Bash tool: Modern CLI tools             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Critical Modules (Lost-in-Middle Recovery)

These modules are injected after `/compact` or during periodic reinforcement:

| Module | Purpose | Trigger |
|--------|---------|---------|
| `amnesia.md` | Prevents "I have no memory" claims | Post-compact, Reinforcement |
| `tabula-rasa.md` | Prevents stale knowledge usage | Post-compact, Reinforcement |
| `date.md` | Ensures correct date awareness | Post-compact, Reinforcement |

## Lazy-Load Modules

These modules are loaded on-demand when specific tools are used:

| Module | Purpose | Trigger |
|--------|---------|---------|
| `agents.md` | Agent hierarchy and delegation | Task tool |
| `compact.md` | Compact behavior rules | PreCompact hook |
| `dace.md` | Micro-modular file principles | Write/Edit tools |
| `hooks.md` | Hook reference documentation | Reference only |
| `reuse.md` | Code reusability patterns | Write/Edit tools |
| `structure.md` | File structure patterns | Write/Edit tools |
| `tools.md` | Modern Rust/Zig CLI tools | Bash tool |

## Installation

Copy modules to your Claude config:

```bash
mkdir -p ~/.claude/modules
cp examples/modules/*.md ~/.claude/modules/
```

## How It Works

### 1. Post-Compact Recovery

When `/compact` runs, the intent hook detects `post_compact: true` and injects:

```
[CONTEXT:RECOVERY]
trigger: post_compact detected

[NO_AMNESIA]
memory_bank: ~/.local/shared/shared-ai/memory/
RULE: Memory Bank EXISTS and is queryable

[TABULA_RASA]
cutoff: 2025-01
today: 2026-01-16
RULE: WebSearch BEFORE code

[DATE_INJECTION]
today: 2026-01-16
RULE: Use injected date, NEVER guess
```

### 2. Periodic Reinforcement

Every 15 turns, the intent hook injects a reinforcement block to combat attention decay:

```
[CONTEXT:REINFORCE]
trigger: periodic (attention decay mitigation)
turn: 15

CRITICAL:BINARY_FIRST
  kavach BEFORE Read/Explore/Task

CRITICAL:TABULA_RASA
  cutoff: 2025-01, today: 2026-01-16
  WebSearch BEFORE code assumptions

CRITICAL:NO_AMNESIA
  Memory Bank EXISTS
  Query: kavach memory bank
```

## Customization

### Adjust Reinforcement Frequency

Edit `session/types.go`:

```go
ReinforceEveryN int // Default: 15 turns
```

### Add Custom Modules

1. Create `~/.claude/modules/your-module.md`
2. Follow SP/1.0 TOON format
3. Add load trigger to intent hook if needed

## References

- [Lost in the Middle - Stanford Research](https://arxiv.org/abs/2307.03172)
- [Attention Sinks - MIT StreamingLLM](https://arxiv.org/abs/2309.17453)
- [SP/1.0 Sutra Protocol](../docs/ARCHITECTURE.md)
