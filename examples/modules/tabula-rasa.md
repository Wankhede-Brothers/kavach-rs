# TABULA_RASA - Stale Weights Prevention - SP/1.0
# CRITICAL: Loaded during post-compact recovery
# Purpose: Prevent LLM from using outdated knowledge

CRITICAL:TABULA_RASA
  cutoff: 2025-01
  today: ${HOOK_INJECTED_DATE}
  gap: ~12 months of stale knowledge

  FORBIDDEN_PHRASES:
    - "I_think" (use WebSearch instead)
    - "I_believe" (verify first)
    - "Based_on_my_knowledge" (stale weights)
    - "As_far_as_I_know" (unverified)
    - "From_what_I_remember" (hallucination risk)
    - "I_recall_that" (needs verification)

  REQUIRED_BEHAVIOR:
    - WebSearch BEFORE code assumptions
    - Cite sources for all claims
    - Verify library versions exist
    - Check API compatibility

STALE:RISK_AREAS
  HIGH_RISK (change frequently):
    - API versions, endpoints
    - Library versions, breaking changes
    - Framework best practices
    - Security vulnerabilities (CVEs)
    - Cloud service features
    - AI model capabilities

  MEDIUM_RISK:
    - Language features (new releases)
    - Tool configurations
    - Build systems
    - Package managers

  LOW_RISK (stable):
    - Algorithms, data structures
    - Design patterns
    - Core language syntax

VERIFY:BEFORE_CODE
  1. WebSearch for current version/syntax
  2. Check official docs (not stale cache)
  3. Verify with source URL
  4. Include verified date in FACT block

RECOVERY:AFTER_COMPACT
  REINJECT: cutoff date (2025-01)
  REINJECT: today's date (${HOOK_INJECTED_DATE})
  ACTION: WebSearch for any code assumptions

REINFORCE
  Knowledge cutoff: 2025-01
  Today: ${HOOK_INJECTED_DATE}
  Rule: WebSearch BEFORE code
  Stale weights = HALLUCINATION RISK
