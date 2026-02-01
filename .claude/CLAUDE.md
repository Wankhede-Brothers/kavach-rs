# kavach-go — Project-Specific Rules (Rust-only CLI)
# Behavioral rules (sparring, zero-stale, memory bank, etc.) loaded from ~/.claude/CLAUDE.md

PROJECT:META
  lang: rust
  binary: kavach (53 commands — 23 gates + 30 non-gate)
  crate: crates/kavach-cli/
  build: just build
  test: just test
  lint: just lint
  fmt: just fmt
  install: just install

HOOKS:AUTO
  # Lifecycle
  SessionStart → kavach session init
  Stop → kavach session end
  PreCompact → kavach session compact
  # Intent
  UserPromptSubmit → kavach gates intent --hook
  # Subagent lifecycle
  SubagentStart → kavach gates subagent --hook
  SubagentStop → kavach gates subagent --hook
  # Umbrella Gates (hierarchical — 4 umbrellas replace 20 individual gates)
  PreToolUse:Write|Edit|NotebookEdit → kavach gates pre-write --hook
    └── security.chain → security.content → guard.code-guard → research
  PostToolUse:Write|Edit|NotebookEdit → kavach gates post-write --hook
    └── antiprod(P0→P3) → quality → lint → context → memory
  PreToolUse:Bash|Read|Glob|Grep|Task|Skill|WebFetch|WebSearch|AskUserQuestion|TaskCreate|TaskUpdate|TaskGet|TaskList|TaskOutput → kavach gates pre-tool --hook
    └── bash | read | ceo | skill | content | task | context
  PostToolUse:Bash|Read|Glob|Grep|Task|WebSearch|WebFetch|TaskCreate|TaskUpdate|TaskOutput → kavach gates post-tool --hook
    └── memory | context | research | task

CRITICAL:ANTI_PATTERNS
  NEVER: .unwrap() in Rust handlers — use ? operator
  NEVER: TODO/FIXME/HACK comments (implement or create ticket)
  NEVER: placeholder/fake values outside test files
  NEVER: dead code — remove it, don't #[allow(dead_code)]
  GATE: kavach gates post-write (umbrella — antiprod P0→P3)

DACE:RULES
  max_lines: 100
  mode: lazy_load,skill_first
  modules: ~/.claude/modules/
