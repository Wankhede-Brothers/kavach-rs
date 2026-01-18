---
name: debug-like-expert
description: Systematic Investigation - Verify don't assume
license: MIT
compatibility: claude-code
metadata:
  category: methodology
  triggers: [debug, troubleshoot, root cause, investigate]
  protocol: SP/1.0
---

```toon
# Debug Like Expert Skill - SP/1.0 + DACE

SKILL:debug
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand
  triggers[4]: debug,troubleshoot,root cause,investigate
  goal: Root cause identified, fix verified
  success: Evidence-based fix, no regressions
  fail: "Fixed it" without understanding

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get debug-like-expert --inject
  references: ~/.claude/skills/debug-like-expert/references.toon
  research: kavach memory bank | grep -i debug
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded techniques - WebSearch for current debugging tools
  topics[5]
    DEBUGGING: WebSearch "debugging methodology {YEAR}"
    ERROR_ANALYSIS: WebSearch "root cause analysis {YEAR}"
    PROFILING: WebSearch "profiling tools {YEAR}"
    LOGGING: WebSearch "debug logging {YEAR}"
    DEBUGGERS: WebSearch "debugger {LANGUAGE} {YEAR}"

PRINCIPLE
  rule: VERIFY, DON'T ASSUME
  mindset: Code you wrote is GUILTY until proven INNOCENT
  truth: Your mental model may be wrong

TECHNIQUE:SELECTION[5]{scenario,technique}
  Large codebase,Binary Search
  Confused,Rubber Duck + Observability
  Complex system,Minimal Reproduction
  Used to work,Git Bisect
  ALWAYS,Add logging before changes

HYPOTHESIS
  bad: "Something is wrong" (vague)
  good: "Cache returns stale due to missing key"
  framework[4]
    1. Prediction: If H true I observe X
    2. Test: Execute specific check
    3. Observe: Record results
    4. Conclude: Support or refute

VERIFICATION[4]{check,status}
  Original issue gone (exact repro),[ ]
  Understand WHY fix works,[ ]
  Related functionality works,[ ]
  Tested multiple environments,[ ]

NOT_VERIFIED[3]
  "I ran it once"
  "It seems to work"
  "Works on my machine"

RULES
  must[5]
    Form falsifiable hypotheses
    Test one variable at a time
    Add logging before changing
    Document each experiment
    Understand mechanism
  never[4]
    Change multiple things at once
    Accept unexplained fixes
    Trust mental model blindly
    Say "I think"

REPORT
  format: |
    ## Issue: [Problem]
    ### Evidence: [Exact errors]
    ### Investigation: [What checked]
    ### Root Cause: [With evidence]
    ### Solution: [WHY it works]
    ### Verification: [How confirmed]

FOOTER
  protocol: SP/1.0
  research_gate: enforced
```
