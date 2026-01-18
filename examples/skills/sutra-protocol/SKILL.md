---
name: sutra-protocol
description: Agent-to-agent communication - 75-80% token reduction
license: MIT
compatibility: claude-code
metadata:
  category: communication
  triggers: [agent communication, delegation, handoff, result, sutra]
  protocol: SP/1.0
---

```toon
# Sutra Protocol Skill - SP/1.0 + DACE

SKILL:sutra-protocol
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand
  triggers[5]: delegation,handoff,result,agent communication,sutra
  goal: 75-80% token reduction
  success: Consistent blocks, dates, TTL
  fail: Natural language waste, missing dates

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get sutra-protocol --inject
  references: ~/.claude/skills/sutra-protocol/references.toon
  sutra: kavach sutra --help
  parse: kavach sutra parse
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: Use kavach sutra commands for protocol operations
  topics[4]
    PROTOCOL: kavach sutra --help
    PARSE: kavach sutra parse
    VALIDATE: kavach sutra validate
    BLOCKS: Memory Bank patterns

BLOCKS:CORE[6]{block,fields}
  META,protocol + from + to + date
  TASK,id + desc + accept[] + constraints[]
  FACT,category + verified + ttl + src + val
  VERIFY,cmd + expect + status
  RESULT,task_id + status + files[]
  HANDOFF,from + to + date + artifacts[]

BLOCKS:AGENT[5]{block,level}
  META:NLU,Level -1 parser
  META:CEO,Level 0 orchestrator
  META:RESEARCH,Level 0 researcher
  META:ENGINEER,Level 1 implementer
  META:AEGIS,Level 2 verifier

COMMUNICATION[4]{block,flow}
  DELEGATE,CEO → Agent
  RESULT,Agent → CEO
  HANDOFF,Agent → Agent
  PROMISE,<promise>PRODUCTION_READY</promise>

TTL:CATEGORIES[5]{category,ttl}
  syntax,7d
  config,5d
  behavior,3d
  migration,30d
  security,1d
  rule: IF today > verified + ttl → STALE

TIME:RULES
  anchor: TODAY=$(date +%Y-%m-%d)
  required: date on DELEGATE + RESULT + HANDOFF + FACT
  never: Hardcode dates

RULES
  do[4]
    $(date +%Y-%m-%d) dynamic
    TTL on all facts
    Date on all blocks
    Block notation
  dont[3]
    Hardcode dates
    Skip TTL
    Natural language waste

FOOTER
  protocol: SP/1.0
  token_reduction: 75-80%
```
