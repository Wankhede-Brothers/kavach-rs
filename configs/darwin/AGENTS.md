# Brahmastra Stack - SP/1.0 (macOS/OpenCode)

META:SYSTEM
  protocol: SP/1.0
  platform: darwin
  binary: /usr/local/bin/kavach
  memory: ~/Library/Application Support/shared-ai/memory/
  stack: Brahmastra

CRITICAL:STALE_WEIGHTS
  cutoff: 2025-01
  today: ${HOOK_INJECTED_DATE}
  RULE: WebSearch BEFORE code

PATHS:MACOS
  binary: /usr/local/bin/kavach
  memory: ~/Library/Application Support/shared-ai/memory/
  settings: ~/Library/Application Support/opencode/settings.json

DACE:CORE
  mode: lazy_load,skill_first,on_demand
  max_lines: 100

AGENTS
  L-1: nlu-intent-analyzer
  L0: ceo, research-director
  L1: backend, frontend, devops, security
  L2: aegis-guardian

HOOKS
  SessionStart: kavach session init
  PreToolUse: kavach gates enforcer --hook
  Stop: kavach session end

MEMORY_BANK
  path: ~/Library/Application Support/shared-ai/memory/
  format: TOON only
  QUERY: kavach memory bank

META:END
  protocol: SP/1.0
  stack: Brahmastra
