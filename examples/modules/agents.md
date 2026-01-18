# Agent Hierarchy - SP/1.0
# Lazy-loaded when Task tool is used

AGENTS:HIERARCHY
  Level -1: nlu-intent-analyzer (haiku)
    role: Fast intent parsing, routes to CEO
    spawned_by: System (UserPromptSubmit hook)

  Level 0: Decision Makers (opus)
    ceo: Orchestration, delegation, NEVER writes code
    research-director: Evidence-based findings, deep research

  Level 1: Engineers (sonnet) - SUB-AGENTS
    backend-engineer: Rust, Axum, API, Database
    frontend-engineer: TypeScript, React, UI
    devops-engineer: Docker, CI/CD, Infrastructure
    security-engineer: Auth, OWASP, Encryption
    qa-lead: Testing, coverage, fixtures

  Level 2: aegis-guardian (opus) - VERIFICATION LAYER
    role: Final verification, quality gate
    NOT_A_SUBAGENT: Runs via hook or explicit invocation

SPAWN:RULES
  CEO_SPAWNS_SUBAGENT:
    - Code implementation → backend-engineer
    - UI/Component work → frontend-engineer
    - Complex research → research-director
    - Codebase exploration → Explore agent

  CEO_EXECUTES_DIRECTLY:
    - Status query → kavach status (no sub-agent)
    - Memory query → kavach memory bank (no sub-agent)
    - Simple file read → Read tool (no sub-agent)
    - Skill invocation → Skill tool (skills are context)

  AEGIS_TRIGGERED_BY:
    - PostToolUse:Write hook (automatic)
    - kavach orch aegis (explicit)
    - CEO invocation (before DONE)

PIPELINE:FLOW
  TASK arrives
    ↓
  CEO evaluates:
    ├─ Simple query? → Execute directly (no sub-agent)
    ├─ Code task? → Spawn Engineer (sub-agent)
    └─ Research? → Spawn Research-Director (sub-agent)
    ↓
  Sub-agent completes
    ↓
  Aegis verifies (hook-triggered, not spawned by engineer)
    ↓
  DONE

DELEGATION:FORMAT
  [DELEGATE]
  from: ceo
  to: {agent}
  skill: {skill}
  task: {task_description}
  model: {haiku|sonnet|opus}

SUTRA:PROTOCOL
  CORE: META,TASK,FACT,UNKNOWN,CONTEXT,VERIFY,ERROR,RESULT
  COMM: DELEGATE,RESULT,HANDOFF,PROMISE
  STATUS: PASS,FAIL,PENDING,RUNNING,BLOCKED,WARNING

CRITICAL:RULES
  1. CEO never writes code - always delegates to engineers
  2. Engineers are sub-agents spawned via Task tool
  3. Aegis is NOT a sub-agent - it's a verification layer
  4. Simple queries don't need sub-agents
  5. Skills inject context, they are NOT agents
