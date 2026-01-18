---
name: resume
description: Project-Aware Context Recovery - Memory Bank Loading
license: MIT
compatibility: claude-code
metadata:
  category: context
  triggers: [resume, /resume, context recovery, session start]
  protocol: SP/1.0
---

```toon
# Resume Skill - SP/1.0 + DACE

SKILL:resume
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand
  triggers[4]: /resume,session recovery,post-compact,next day
  goal: Load project context from Memory Bank
  success: Full project awareness without re-reading files
  fail: Re-reading already-known files, losing task state

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get resume --inject
  references: ~/.claude/skills/resume/references.toon
  memory: kavach memory bank
  session: kavach session init
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: Use kavach commands for session recovery - NO manual file reads
  topics[5]
    MEMORY: kavach memory bank
    KANBAN: kavach memory kanban
    SESSION: kavach session init
    STATUS: kavach status
    VALIDATE: kavach session validate

FLOW:RESUME[4]{step,action}
  1. DETECT,Get working directory → extract project
  2. LOAD,Memory Bank filtered by project
  3. INJECT,Date + task state + file summaries
  4. WIRE,Ready to continue exactly where left off

LOAD:SOURCES[4]
  scratchpad.toon → current task
  hot-context.json → files read
  decisions/{project}/decisions.toon → project decisions
  CLAUDE.md → already loaded by hook

MEMORY:SOURCES[4]{source,content,tokens_saved}
  scratchpad.toon,Current task + focus,1k
  hot-context.json,Files read (summaries),30-50k
  decisions/*.toon,Project decisions,5k
  CLAUDE.md,Already loaded by hook,5k

RULES
  do[4]
    Filter by current project
    Inject today's date
    Load summaries not full files
    Continue from last task
  dont[4]
    Load other projects' context
    Re-read files in hot-context
    Forget task state
    Miss date injection

FOOTER
  protocol: SP/1.0
  binary: kavach session resume
```
