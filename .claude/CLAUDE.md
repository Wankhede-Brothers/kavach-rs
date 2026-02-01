# kavach-go — Project-Specific Rules
# Behavioral rules (sparring, zero-stale, memory bank, etc.) loaded from ~/.claude/CLAUDE.md

HOOKS:AUTO
  # Lifecycle
  SessionStart → kavach session init
  Stop → kavach session end
  PreCompact → kavach session compact
  # Intent
  UserPromptSubmit → kavach gates intent --hook
  # Umbrella Gates (hierarchical — 4 umbrellas replace 20 individual gates)
  PreToolUse:Write|Edit|NotebookEdit → kavach gates pre-write --hook
    └── security.chain → security.content → guard.code-guard → research
  PostToolUse:Write|Edit|NotebookEdit → kavach gates post-write --hook
    └── antiprod(P0→P3) → quality → lint → context → memory
  PreToolUse:* → kavach gates pre-tool --hook
    └── bash | read | ceo | skill | content | task | context
  PostToolUse:* → kavach gates post-tool --hook
    └── memory | context | research | task

CRITICAL:ANTI_PATTERNS
  NEVER: console.log in production code (use structured logger)
  NEVER: TODO/FIXME/HACK comments (implement or create ticket)
  NEVER: http://localhost in non-config files
  NEVER: .catch(() => {}) empty error handlers
  NEVER: Non-null assertions (!.) — use optional chaining (?.)
  NEVER: fetch() without error handling
  NEVER: as any — use proper type narrowing
  NEVER: .unwrap() in Rust handlers — use ? operator
  NEVER: process.env without fallback value
  NEVER: placeholder/fake values outside test files
  GATE: kavach gates post-write (umbrella — antiprod P0→P3)

DACE:RULES
  max_lines: 100
  mode: lazy_load,skill_first
  modules: ~/.claude/modules/

DACE:PROJECT_PROFILE
  lang: rust
  test: just test
  build: just build
  lint: just lint
  fmt: just fmt
