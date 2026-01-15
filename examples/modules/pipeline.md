# Pipeline Architecture - SP/3.0
# Lazy-loaded when Task tool is used

PIPELINE:FLOW
  USER → NLU → CEO → ENGINEERS → AEGIS → PRODUCTION_READY
    │      │      │         │         │
  L-1    L0     L0        L1        L2

PIPELINE:LEVELS
  Level -1: nlu-intent-analyzer (haiku)
    Parses ALL user requests → Sutra Protocol

  Level 0: Decision Makers (opus)
    ceo: Orchestration, delegation
    research-director: Evidence-based findings

  Level 1: Engineers (sonnet, parallel)
    backend-engineer, frontend-engineer
    database-engineer, devops-engineer
    security-engineer, qa-lead

  Level 1.5: code-reviewer (sonnet)
    Post-implementation review

  Level 2: aegis-guardian (opus)
    Final verification, quality gate

PIPELINE:ENFORCEMENT
  PreToolUse Hook → enforcer.go → BLOCKS until pipeline active
  NLU MUST be invoked FIRST
  CEO MUST receive NLU output
  Engineers MUST receive CEO delegation
  Aegis MUST verify before completion

PIPELINE:COMPLETION
  signal: <promise>PRODUCTION_READY</promise>
  conditions:
    - ALL acceptance criteria met
    - aegis-guardian verification PASSED
    - All builds/tests GREEN
