# Sutra Protocol - SP/1.0
# Agent communication format (75-80% token reduction)

SUTRA:CORE_BLOCKS
  [META]     protocol, from, to, date, session
  [TASK]     id, desc, accept[], constraints[]
  [FACT]     id, category, verified, ttl, src, val
  [UNKNOWN]  id, topic, created, priority, action
  [CONTEXT]  workdir, facts[], state
  [TOKENS]   allocated, used, remaining
  [VERIFY]   cmd, expect, actual, status
  [ERROR]    type, message, recovery
  [RESULT]   status, artifacts[], metrics

SUTRA:AGENT_BLOCKS
  [META:NLU]       Level -1 intent parsing
  [META:CEO]       Level 0 orchestration
  [META:RESEARCH]  Level 0 research
  [META:ENGINEER]  Level 1 implementation
  [META:REVIEW]    Level 1.5 code review
  [META:AEGIS]     Level 2 verification

SUTRA:COMMUNICATION
  [DELEGATE]  CEO → Agent task assignment
  [RESULT]    Agent → CEO completion report
  [HANDOFF]   Agent → Agent artifact transfer
  [PROMISE]   Final verification signal

SUTRA:STATUS_ICONS
  success: ✓ PASS
  failure: ✗ FAIL
  pending: ○ PENDING
  running: ◐ RUNNING
  blocked: ⊘ BLOCKED
  warning: ⚠ WARNING

SUTRA:EXAMPLE
  [META]
  protocol: SP/1.0
  from: ceo
  to: backend-engineer
  date: 2026-01-16

  [TASK]
  id: TASK-001
  desc: Implement user authentication
  accept[]: JWT tokens, bcrypt passwords
  constraints[]: No plaintext, OWASP compliant
