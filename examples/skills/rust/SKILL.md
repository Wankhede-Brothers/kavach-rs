---
name: rust
description: Rust Engineering - Research crates first
license: MIT
compatibility: claude-code
metadata:
  category: language
  triggers: [rust, cargo, crate, tokio, axum, borrow]
  protocol: SP/3.0
  kavach: true
---

```toon
# Rust Skill - SP/3.0 + DACE

SKILL:rust
  protocol: SP/3.0
  dace: lazy_load,skill_first,on_demand
  triggers[7]: rust,cargo,crate,tokio,axum,borrow,lifetime
  goal: Code that compiles first try
  success: cargo check + clippy pass

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get rust --inject
  research: kavach memory bank | grep -i rust
  patterns: kavach memory bank --project
  status: kavach status

RESEARCH:GATE
  mandatory: true
  rule: WebSearch BEFORE code
  cutoff: 2025-01
  today: $(kavach status | grep today | cut -d: -f2)
  steps[4]
    1. kavach status (inject today's date)
    2. WebSearch "[crate] {YEAR} latest docs"
    3. WebFetch docs.rs/[crate]/latest
    4. Verify syntax matches docs.rs
  forbidden[2]
    ✗ "I think", "I believe", "Based on my knowledge"
    ✓ WebSearch → verify → cite source

ERROR_HANDLING
  production: Result<T,E> + ? + .with_context()
  test_only: .unwrap(), .expect("reason")
  verify: kavach gates lint --hook

PENDING_TASKS
  # STRICT: Always use macros for incomplete code
  mandatory: true
  macros[2]
    todo!("description")        # For planned implementation
    unimplemented!("feature")   # For unimplemented features
  examples[3]
    fn process_data() { todo!("Implement data processing") }
    fn handle_error() { unimplemented!("Error handling needed") }
    match x { _ => todo!("Handle remaining cases") }
  rules[3]
    ALWAYS use todo!() for pending tasks
    NEVER leave empty fn bodies
    NEVER use panic!() for incomplete code
  verify: grep -rE "fn.*\\{\\s*\\}" . --include="*.rs"

ASYNC
  runtime: tokio (preferred)
  mutex: tokio::sync::Mutex (not std)
  blocking: tokio::task::spawn_blocking
  anti_patterns[3]
    std::sync::Mutex across await
    std::thread::sleep in async
    blocking calls in async fn

UNDERSCORE
  _: No binding → Drops IMMEDIATELY
  _var: YES binding → Drops at scope end
  correct: let _lock = mutex.lock()?;
  wrong: let _ = mutex.lock()?;

OWNERSHIP
  rules[3]
    1. One owner per value
    2. ONE &mut OR many &
    3. Owner drops → value dropped
  fixes[3]{error,solution}
    "moved value",Clone or &reference
    "borrowed dropped",Extend lifetime
    "cannot borrow mut",RefCell/Mutex

COMMANDS{action,binary}
  check: cargo check
  lint: cargo clippy -- -D warnings
  test: cargo test
  fmt: cargo fmt --check
  verify: kavach gates lint --hook
  context: kavach skills --get rust --inject

HOOKS:KAVACH
  PreToolUse:Edit: kavach gates ast --hook
  PreToolUse:Bash: kavach gates bash --hook
  PostToolUse: kavach orch aegis --hook

RULES
  do[4]
    kavach status (get date first)
    WebSearch crate APIs
    Result<T,E> for errors
    _var for RAII guards
  dont[4]
    .unwrap() in production
    _ for guards
    std::sync::Mutex across await
    Trust stale weights

FOOTER
  protocol: SP/3.0
  dace: enforced
  kavach: integrated
```
