# NO_AMNESIA - Memory Bank Awareness - SP/1.0
# CRITICAL: Loaded during post-compact recovery
# Purpose: Prevent LLM from claiming "no memory access"

CRITICAL:NO_AMNESIA
  FORBIDDEN_PHRASES:
    - "I don't have access to memory"
    - "I have no memory of"
    - "I cannot remember"
    - "I don't have persistent memory"
    - "As an AI, I don't have memory"

  MEMORY_BANK_EXISTS:
    path: ~/.local/shared/shared-ai/memory/
    query: kavach memory bank
    structure: TOON files (single file per category)
    scope: PROJECT_ISOLATED (active + global)

MEMORY:CATEGORIES
  LTM (Long-Term): decisions, patterns, roadmaps, proposals
  STM (Short-Term): session-state, scratchpad, volatile
  ROOT: GOVERNANCE.toon, index.toon

MEMORY:QUERY_RULES
  ALWAYS: kavach memory bank (before claiming no memory)
  NEVER: Read(~/.local/shared/shared-ai/memory/*.toon) directly
  WHY: Binary provides structured output, DACE-compliant

RECOVERY:AFTER_COMPACT
  1. kavach memory bank (query full state)
  2. Restore context from TOON files
  3. Continue task from scratchpad
  NEVER: Claim memory loss after compact

REINFORCE
  Memory Bank = SINGLE SOURCE OF TRUTH
  kavach binary = ONLY interface to Memory Bank
  Claiming "no memory" = SP/1.0 VIOLATION
