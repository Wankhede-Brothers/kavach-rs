---
name: arch
description: System Design - Numbers not adjectives
license: MIT
compatibility: claude-code
metadata:
  category: design
  triggers: [system design, architecture, scalability, distributed]
  protocol: SP/1.0
---

```toon
# Architecture Skill - SP/1.0 + DACE

SKILL:arch
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand
  triggers[4]: system design,architecture,scalability,distributed
  goal: Scalable with documented trade-offs
  success: Capacity calculated, CAP analyzed
  fail: "Use microservices" without numbers

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get arch --inject
  references: ~/.claude/skills/arch/references.toon
  research: kavach memory bank | grep -i arch
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded numbers - WebSearch for current benchmarks
  topics[6]
    SYSTEM_DESIGN: WebSearch "system design {YEAR}"
    MICROSERVICES: WebSearch "microservices patterns {YEAR}"
    DATABASE_DESIGN: WebSearch "database scaling {YEAR}"
    CACHING: WebSearch "caching strategies {YEAR}"
    MESSAGING: WebSearch "message queue {YEAR}"
    LOAD_BALANCING: WebSearch "load balancing {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[4]
    TODAY=$(date +%Y-%m-%d)
    Calculate capacity estimates
    Analyze CAP trade-offs
    Document complexity
  forbidden[2]
    ✗ NEVER recommend without numbers
    ✓ ALWAYS document complexity

DESIGN:STEPS[8]{step,action}
  1,Requirements (functional + non-functional)
  2,Estimates (QPS + Storage + Bandwidth)
  3,API Design
  4,Data Model (SQL vs NoSQL)
  5,High-Level Architecture
  6,Deep Dives (bottlenecks)
  7,Failure Modes
  8,Monitoring

CAPACITY
  formulas[4]
    QPS = DAU × requests / 86400
    Peak = QPS × 3
    Storage = records × size × retention
    Bandwidth = QPS × payload_size
  quick[2]
    1M req/day ≈ 12 QPS
    1KB × 1M users = 1GB

DATABASE[5]{need,choice}
  ACID,PostgreSQL
  High writes,Cassandra
  Flexible schema,MongoDB
  Graph,Neo4j
  Cache,Redis

CAP
  rule: Pick CP or AP (P is mandatory)

SCALING[4]{type,description}
  horizontal,Add machines (no SPOF)
  vertical,Bigger machine (SPOF risk)
  caching,Cache-aside + write-through
  sharding,Hash + range + geographic

RULES
  do[4]
    Show calculations
    Analyze trade-offs
    Consider failures
    Start simple
  dont[4]
    "Use microservices" blindly
    Skip capacity math
    Assume availability
    Over-engineer

VALIDATE[3]{check,status}
  Capacity calculated,[ ]
  Trade-offs documented,[ ]
  Failures analyzed,[ ]

FOOTER
  protocol: SP/1.0
  research_gate: enforced
```
