# Compact Rules - SP/3.0
# Loaded during PreCompact hook
# Updated: Lost-in-Middle mitigation

DACE:COMPACT
  target: <10000 tokens post-compact

PRE_COMPACT:ACTIONS
  1. kavach session compact (saves to TOON)
  2. Scratchpad saved to STM/projects/{project}/
  3. Session state saved to STM/session-state.toon
  4. post_compact: true (triggers recovery)

POST_COMPACT:CRITICAL_RECOVERY
  # These rules are LOST during compaction - MUST reinject
  LOAD_MODULES:
    - amnesia.md (Memory Bank awareness)
    - tabula-rasa.md (Stale weights prevention)
    - date.md (Date injection)

  REINJECT_CONTEXT:
    date: ${HOOK_INJECTED_DATE}
    cutoff: 2025-01
    memory: ~/.local/shared/shared-ai/memory/
    binary: kavach (BINARY_FIRST)

POST_COMPACT:RULES
  1. DO_NOT re-read files already discussed
  2. DO_NOT summarize conversation history
  3. DO query memory bank: kavach memory bank
  4. DO lazy-load: Grep before Read
  5. DO continue task from scratchpad

POST_COMPACT:BEHAVIOR
  mode: SUGGEST not BLOCK
  Read: ALLOWED with lazy-load hints
  Binary: PREFERRED for context queries
  Files: Load on-demand, not eagerly

LOST_IN_MIDDLE:MITIGATION
  problem: Critical rules in middle of context get lost
  solution_1: Sandwich pattern in CLAUDE.md (TOP + BOTTOM)
  solution_2: Post-compact recovery injection
  solution_3: Periodic reinforcement (every 15 turns)

RESUME:COMMAND
  run: kavach session init
  action: Auto-detects post_compact, restores context
  injects: date, cutoff, memory bank path, critical rules
