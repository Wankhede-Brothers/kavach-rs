# DACE Routing - SP/1.0
# Maps tasks to Skills + Agents for full pipeline

DACE:PIPELINE
  flow: TASK → SKILL → AGENT → AEGIS → DONE
  principle: skill_first (domain expertise before delegation)

ROUTING:BACKEND{task_pattern,skill,agent}
  "API endpoint",/api-design,backend-engineer
  "REST API",/api-design,backend-engineer
  "gRPC service",/api-design,backend-engineer
  "GraphQL",/api-design,backend-engineer
  "Database schema",/sql,backend-engineer
  "SQL query",/sql,backend-engineer
  "GIN index",/sql,backend-engineer
  "Full-text search",/sql,backend-engineer
  "PostgreSQL",/sql,backend-engineer
  "Rust implementation",/rust,backend-engineer
  "Axum handler",/rust,backend-engineer
  "Rate limiting",/security,backend-engineer
  "Session management",/security,backend-engineer
  "Webhook handler",/api-design,backend-engineer

ROUTING:FRONTEND{task_pattern,skill,agent}
  "React component",/frontend,frontend-engineer
  "UI component",/frontend,frontend-engineer
  "Form validation",/frontend,frontend-engineer
  "Loading states",/frontend,frontend-engineer
  "Skeletons",/frontend,frontend-engineer
  "Toast notifications",/frontend,frontend-engineer
  "Navigation",/frontend,frontend-engineer
  "Responsive design",/frontend,frontend-engineer
  "Dashboard",/frontend,frontend-engineer
  "TypeScript",/frontend,frontend-engineer

ROUTING:INFRASTRUCTURE{task_pattern,skill,agent}
  "Deploy",/cloud-infrastructure-mastery,backend-engineer
  "DigitalOcean",/cloud-infrastructure-mastery,backend-engineer
  "Kubernetes",/cloud-infrastructure-mastery,backend-engineer
  "Health endpoints",/api-design,backend-engineer
  "Cloudflare",/cloud-infrastructure-mastery,backend-engineer
  "Docker",/cloud-infrastructure-mastery,backend-engineer
  "CI/CD",/cloud-infrastructure-mastery,backend-engineer

ROUTING:SECURITY{task_pattern,skill,agent}
  "Authentication",/security,backend-engineer
  "Authorization",/security,backend-engineer
  "MFA",/security,backend-engineer
  "RSA signature",/security,backend-engineer
  "Encryption",/security,backend-engineer
  "OWASP",/security,backend-engineer
  "Escrow",/security,backend-engineer

ROUTING:QUALITY{task_pattern,skill,agent}
  "Debug",/debug-like-expert,research-director
  "Bug fix",/debug-like-expert,backend-engineer
  "Test",/testing,backend-engineer
  "Unit test",/testing,backend-engineer
  "Integration test",/testing,backend-engineer
  "Code analysis",/heal,aegis-guardian
  "Refactor",/heal,backend-engineer

ROUTING:ARCHITECTURE{task_pattern,skill,agent}
  "System design",/arch,research-director
  "Architecture",/arch,research-director
  "Scalability",/arch,research-director
  "Algorithm",/dsa,backend-engineer
  "Data structure",/dsa,backend-engineer
  "Performance",/dsa,backend-engineer

AGENT:HIERARCHY
  Level -1: nlu-intent-analyzer (haiku) → Parse intent
  Level 0:  ceo (opus) → Orchestrate, delegate
  Level 0:  research-director (opus) → Research, evidence
  Level 1:  backend-engineer (sonnet) → Rust, API, DB
  Level 1:  frontend-engineer (sonnet) → React, UI
  Level 2:  aegis-guardian (opus) → Verify, approve

DELEGATION:FORMAT
  [DELEGATE]
  from: ceo
  to: {agent}
  skill: {skill}
  task: {task_description}

  [SKILL:INJECT]
  invoke: /{skill}
  context: kavach skills --get {skill} --inject

  [RESEARCH:REQUIRED]
  cutoff: 2025-01
  today: ${date}
  action: WebSearch "{topic} {year}"

OUTPUT:EXAMPLE
  ┌──────────┬────────────────────────┬─────────────┬─────────────────────┐
  │ Priority │ Task                   │ Skill       │ Agent               │
  ├──────────┼────────────────────────┼─────────────┼─────────────────────┤
  │ 🔴       │ GIN Index for search   │ /sql        │ backend-engineer    │
  │ 🔴       │ Deploy Gatus           │ /cloud-inf  │ backend-engineer    │
  │ 🔴       │ Health endpoints       │ /api-design │ backend-engineer    │
  │ 🟠       │ Redis sessions         │ /security   │ backend-engineer    │
  │ 🟠       │ Form validation UI     │ /frontend   │ frontend-engineer   │
  └──────────┴────────────────────────┴─────────────┴─────────────────────┘
