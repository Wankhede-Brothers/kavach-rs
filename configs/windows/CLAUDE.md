# Brahmastra Stack - SP/1.0 (Windows)

META:SYSTEM
  protocol: SP/1.0
  platform: windows
  binary: %LOCALAPPDATA%\bin\kavach.exe
  memory: %LOCALAPPDATA%\shared-ai\memory\
  stack: Brahmastra

CRITICAL:STALE_WEIGHTS
  cutoff: 2025-01
  today: ${HOOK_INJECTED_DATE}
  RULE: WebSearch BEFORE code
  FORBIDDEN: "I think","I believe"

PATHS:WINDOWS
  binary: %LOCALAPPDATA%\bin\kavach.exe
  memory: %LOCALAPPDATA%\shared-ai\memory\
  settings: %USERPROFILE%\.claude\settings.json

DACE:CORE
  mode: lazy_load,skill_first,on_demand
  max_lines: 100 (hard block)
  warn_lines: 50 (suggest split)

CRITICAL:RUST_CLI_FIRST
  RULE: Use Rust CLI tools INSTEAD of legacy commands
  CHECK: kavach memory view (shows tool status)
  INSTALL: scoop install bat eza fd ripgrep sd procs dust bottom delta

RUST_CLI:COMMANDS{legacy,rust}
  type,bat --plain
  dir,eza --icons
  where,fd
  findstr,rg
  ps,procs
  diff,delta

ZIG_CLI:STACK
  bun: JS/TS runtime
  ghostty: GPU terminal (WSL)

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
  path: %LOCALAPPDATA%\shared-ai\memory\
  format: TOON only
  QUERY: kavach memory bank

ANTI:PATTERNS{bad,good}
  type file.txt,bat file.txt
  dir /s,eza --tree
  where /r,fd pattern
  findstr pattern,rg pattern
  "I think",WebSearch
  Read(memory/),kavach memory bank

META:END
  protocol: SP/1.0
  stack: Brahmastra
