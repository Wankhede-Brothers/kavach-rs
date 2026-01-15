# Brahmastra Stack - SP/3.0 (Linux)

META:SYSTEM
  protocol: SP/3.0
  platform: linux
  binary: ~/.local/bin/kavach
  memory: ~/.local/share/shared-ai/memory/
  stack: Brahmastra

CRITICAL:STALE_WEIGHTS
  cutoff: 2025-01
  today: ${HOOK_INJECTED_DATE}
  RULE: WebSearch BEFORE code
  FORBIDDEN: "I think","I believe"

PATHS:LINUX
  binary: ~/.local/bin/kavach
  memory: ~/.local/share/shared-ai/memory/
  settings: ~/.claude/settings.json

DACE:CORE
  mode: lazy_load,skill_first,on_demand
  max_lines: 100 (hard block)
  warn_lines: 50 (suggest split)

CRITICAL:RUST_CLI_FIRST
  RULE: Use Rust CLI tools INSTEAD of legacy coreutils
  CHECK: kavach memory view (shows tool status)

RUST_CLI:COMMANDS{legacy,rust}
  cat,bat --plain
  ls,eza --icons
  find,fd
  grep,rg
  sed,sd
  ps,procs
  du,dust
  top,btm
  diff,delta

ZIG_CLI:STACK
  bun: JS/TS runtime
  ghostty: GPU terminal

AGENTS
  L-1: nlu-intent-analyzer (haiku)
  L0: ceo, research-director (opus)
  L1: backend, frontend, devops, security (sonnet)
  L2: aegis-guardian (opus)

HOOKS
  SessionStart: kavach session init
  UserPromptSubmit: kavach gates intent --hook
  PreToolUse:Task: kavach gates ceo --hook
  PreToolUse:Bash: kavach gates bash --hook
  PreToolUse:Read: kavach gates read --hook
  PreToolUse:Write: kavach gates enforcer --hook
  PreToolUse:Edit: kavach gates enforcer --hook
  Stop: kavach session end

MEMORY_BANK
  path: ~/.local/share/shared-ai/memory/
  format: TOON only
  QUERY: kavach memory bank

ANTI:PATTERNS{bad,good}
  cat file.txt,bat file.txt
  ls -la,eza -la --icons
  find . -name,fd pattern
  grep pattern,rg pattern
  sed 's/x/y/',sd x y
  "I think",WebSearch
  Read(memory/),kavach memory bank

META:END
  protocol: SP/3.0
  stack: Brahmastra
