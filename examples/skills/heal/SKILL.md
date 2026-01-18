---
name: heal
description: 5-Layer Code Analysis - Static, semantic, security
license: MIT
compatibility: claude-code
metadata:
  category: quality
  triggers: [heal, lint, fix, code quality]
  protocol: SP/1.0
---

```toon
# Heal Skill - SP/1.0 + DACE

SKILL:heal
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand
  triggers[4]: heal,lint,fix,lint_and_fix
  goal: 5-layer code analysis
  success: Issues categorized, fixed, reported
  fail: Missing layers, wrong severity

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get heal --inject
  references: ~/.claude/skills/heal/references.toon
  research: kavach memory bank | grep -i heal
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded lint rules - WebSearch for current tools
  topics[5]
    STATIC: WebSearch "static analysis {LANGUAGE} {YEAR}"
    SEMANTIC: WebSearch "type checking {LANGUAGE} {YEAR}"
    SECURITY: WebSearch "SAST tools {YEAR}"
    DEAD_CODE: WebSearch "dead code detection {YEAR}"
    SUPPRESSED: WebSearch "code suppression audit"

COMMANDS[4]{cmd,action}
  !heal,Full 5-layer
  !heal --static,Layer 1 only
  !heal --security,Layer 3 only
  !heal --rust|--ts|--py|--go,Language-specific

LAYERS[5]{layer,focus,tools}
  1. STATIC,Syntax + lint,clippy | eslint | ruff
  2. SEMANTIC,Type safety + null handling,Type checker
  3. SECURITY,Vulns + secrets + OWASP,audit tools
  4. PERFORMANCE,O(n) + allocations,profiler
  5. RELIABILITY,Error handling + edge cases,tests

SEVERITY[4]{level,examples,action}
  CRITICAL,SQL injection + secrets,HALT
  HIGH,.unwrap() prod + missing error,FLAG
  MEDIUM,Suboptimal patterns,SUGGEST
  LOW,Style issues,AUTO-FIX

TOOLS[4]{lang,static,security}
  Rust,clippy,cargo audit
  TS,eslint,npm audit
  Python,ruff,safety
  Go,golangci-lint,govulncheck

SECRETS_PATTERNS[5]
  [A-Za-z0-9_-]{20,}
  AKIA[0-9A-Z]{16}
  sk_live_[0-9a-zA-Z]{24}
  ghp_[0-9A-Za-z]{36}
  -----BEGIN.*PRIVATE KEY-----

REPORT
  format: |
    LAYER 1: STATIC      [OK/FAIL]
    LAYER 2: SEMANTIC    [OK/FAIL]
    LAYER 3: SECURITY    [OK/FAIL]
    LAYER 4: PERFORMANCE [OK/FAIL]
    LAYER 5: RELIABILITY [OK/FAIL]
    VERDICT: ✅ READY | ⚠️ WARNINGS | ❌ BLOCKED

RULES
  do[4]
    All 5 layers
    Categorize by severity
    Block on security
    Auto-fix low severity
  dont[3]
    Skip layers
    Treat all equal
    Allow "temp" secrets

FOOTER
  protocol: SP/1.0
  research_gate: enforced
```
