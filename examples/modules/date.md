# DATE_INJECTION - Time Awareness - SP/3.0
# CRITICAL: Loaded during post-compact recovery
# Purpose: Ensure LLM knows current date

CRITICAL:DATE_INJECTION
  source: SessionStart hook
  command: kavach session init
  format: YYYY-MM-DD (ISO 8601)
  variable: ${HOOK_INJECTED_DATE}

  INJECTION_POINTS:
    - SessionStart hook output
    - CLAUDE.md ${HOOK_INJECTED_DATE} substitution
    - Post-compact recovery block
    - Periodic reinforcement block

DATE:AWARENESS
  today: ${HOOK_INJECTED_DATE}
  cutoff: 2025-01
  rule: ALWAYS use injected date, NEVER guess

DATE:USAGE_RULES
  WebSearch: Include year in queries (e.g., "React 2026")
  FACT blocks: verified: ${today}
  File headers: date: ${today}
  Commits: Use current date context

DATE:FORBIDDEN
  - Hardcoding year (e.g., "2024", "2025")
  - Guessing current date
  - Using training data dates
  - Assuming date from context

DATE:VERIFICATION
  command: kavach status
  output: today: YYYY-MM-DD
  fallback: date +%Y-%m-%d

RECOVERY:AFTER_COMPACT
  CRITICAL: Date context often lost in summary
  ACTION: kavach status (re-injects today's date)
  VERIFY: Check date before any time-sensitive operation

REINFORCE
  Today: ${HOOK_INJECTED_DATE}
  Source: kavach session init
  Rule: NEVER hardcode dates
