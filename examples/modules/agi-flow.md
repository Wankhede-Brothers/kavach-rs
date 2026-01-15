# AGI Flow - SP/3.0
# Real-time learning system - NEVER trust training weights

PARADIGM:SHIFT
  OLD: Training weights → Generate answer (STALE)
  NEW: Intent → Memory Bank → WebSearch → Generate (FRESH)
  WHY: Training cutoff = 2025-01, world moves fast

AGI:FLOW
  ```
  User Prompt (vague/technical)
         ↓
  [Hook: UserPromptSubmit] → kavach gates intent --hook
         ↓
  NLU-Intent (Level -1, haiku)
    - Classify: debug/implement/research/optimize/refactor
    - Detect domain: frontend/backend/security/database/infra
    - Recommend: skills + agents
         ↓
  [BEFORE] Memory Bank Check
    - kavach memory bank (load context)
    - kavach memory kanban (current tasks)
         ↓
  CEO (Level 0, opus) - NEVER writes code
    - Analyze scope
    - Spawn research-director if needed
    - Delegate to engineers with skill context
         ↓
  [TABULA_RASA] WebSearch "{topic} {today's date}"
    - FORBIDDEN: Assuming from training weights
    - REQUIRED: Fresh patterns from 2026
         ↓
  Engineer (Level 1, sonnet) - Sub-agent
    - Invoke skill for domain expertise
    - Implement with latest patterns
         ↓
  Aegis (Level 2, opus) - Verification
    - Test coverage
    - Security audit
    - Quality gate
         ↓
  [AFTER] Memory Bank Sync
    - kavach memory sync (update kanban, scratchpad)
    - Persist learnings for future
         ↓
  DONE (PRODUCTION_READY)
  ```

HOOKS:AUTO_TRIGGER
  SessionStart: kavach session init (date injection)
  UserPromptSubmit: kavach gates intent --hook (NLU + context)
  PreToolUse:Task: kavach gates ceo --hook (skill injection)
  PostToolUse:Write: kavach orch aegis --hook (verification)
  PostToolUse:TodoWrite: kavach memory sync --hook (memory update)
  PreCompact: kavach session compact (save state)
  Stop: kavach session end (persist)

MEMORY:SYNC_PROTOCOL
  BEFORE_EXECUTION:
    1. kavach memory bank (load full context)
    2. kavach memory kanban (current task state)
    3. Check scratchpad for work-in-progress

  AFTER_EXECUTION:
    1. TodoWrite triggers → kavach memory sync
    2. Update kanban.toon (task status)
    3. Update scratchpad.toon (learnings)
    4. Update patterns.toon (if new pattern discovered)

TABULA_RASA:ENFORCEMENT
  cutoff: 2025-01
  today: ${HOOK_INJECTED_DATE}
  RULE: WebSearch BEFORE code
  FORMAT: "{topic} {year}" (e.g., "React patterns 2026")
  FORBIDDEN: "I_think", "I_believe", "Based_on_my_knowledge"
  REQUIRED: "According to [source]", "Latest docs show"

NO_AMNESIA:ENFORCEMENT
  memory_bank: ~/.local/shared/shared-ai/memory/
  query: kavach memory bank
  RULE: Check Memory Bank EVERY session start
  FORBIDDEN: "I_dont_have_memory_access"
  REQUIRED: Acknowledge Memory Bank exists

LOSS_IN_MIDDLE:MITIGATION
  1. Sandwich pattern in CLAUDE.md (rules at TOP + BOTTOM)
  2. Post-compact recovery (re-inject after /compact)
  3. Periodic reinforcement (every 15 turns)
  4. Critical rules in recency window

VAGUE_TO_TECHNICAL:MAPPING
  "make it faster" → optimize → /dsa, /arch
  "its broken" → debug → /debug-like-expert
  "add login" → implement + security → /security
  "looks ugly" → implement + frontend → /frontend
  "go live" → implement + infra → /cloud-infrastructure-mastery
  "save data" → implement + database → /sql

MICRO_MODULAR:RULES
  max_lines: 100 (hard limit)
  mode: lazy_load (load on demand)
  modules: Load only when skill/agent invoked
  DACE: Dynamic Agentic Context Engineering
