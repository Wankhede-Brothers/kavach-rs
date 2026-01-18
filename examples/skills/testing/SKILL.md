---
name: testing
description: Universal Test Engineering - Research patterns first
license: MIT
compatibility: claude-code
metadata:
  category: testing
  triggers: [test, unit, integration, e2e, vitest, cargo test, coverage]
  protocol: SP/1.0
---

```toon
# Testing Skill - SP/1.0 + DACE

SKILL:testing
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand
  triggers[6]: test,unit,integration,e2e,coverage,mock
  goal: 80%+ coverage, proper pyramid
  success: Fast CI, no flaky tests
  fail: Shared state, implementation testing

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get testing --inject
  references: ~/.claude/skills/testing/references.toon
  research: kavach memory bank | grep -i test
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded test examples - WebSearch for current frameworks
  topics[6]
    UNIT: WebSearch "unit testing {LANGUAGE} {YEAR}"
    INTEGRATION: WebSearch "integration testing {YEAR}"
    E2E: WebSearch "e2e testing {YEAR}"
    PROPERTY: WebSearch "property testing {YEAR}"
    MOCKING: WebSearch "mocking {LANGUAGE} {YEAR}"
    COVERAGE: WebSearch "code coverage {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[4]
    TODAY=$(date +%Y-%m-%d)
    WebSearch "[framework] testing $YEAR patterns"
    WebFetch official framework docs
    Verify current API syntax
  forbidden[2]
    ✗ NEVER trust old test syntax
    ✓ ALWAYS verify current API

PYRAMID[3]{level,percentage,speed}
  Unit,60-70%,<10ms
  Integration,20-30%,<1s
  E2E,5-10%,<30s

COMMANDS[3]{rust,typescript}
  cargo test,bun test
  cargo test --lib,bun test:unit
  cargo llvm-cov,bun test --coverage

MOCKING
  mock_when[4]
    External HTTP APIs
    File system ops
    Random/Time ops
    Third-party services
  research: WebSearch "[lang] mocking $YEAR"

NAMING
  pattern: should_{expected}_when_{condition}
  example: should_return_error_when_input_empty
  avoid[3]: test1(),test_user(),it_works()

RULES
  do[5]
    Arrange-Act-Assert
    One assertion per test
    WebSearch framework patterns
    Mock external deps only
    Property-based for algorithms
  dont[4]
    Test implementation details
    Share state between tests
    Arbitrary sleeps
    Rely on test order

VALIDATE[4]{check,status}
  WebSearched current syntax,[ ]
  Pyramid distribution 60/30/10,[ ]
  No flaky tests,[ ]
  Coverage meets threshold,[ ]

FOOTER
  protocol: SP/1.0
  research_gate: enforced
```
