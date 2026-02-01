#!/usr/bin/env bash
# P0 Audit: find missing parts and bugs in kavach binary
set -euo pipefail

KAVACH="$HOME/.local/bin/kavach"
PASS=0
FAIL=0
BUGS=""

check() {
    local name="$1" expect="$2" output="$3"
    if echo "$output" | grep -q "$expect"; then
        PASS=$((PASS + 1))
    else
        FAIL=$((FAIL + 1))
        BUGS="$BUGS\n  BUG: $name (expected: $expect)"
        echo "  BUG: $name"
        echo "       expected: $expect"
        echo "       got: ${output:0:120}"
    fi
}

# Helper: pipe JSON via temp file to avoid bash escaping issues
run_hook() {
    local subcmd="$1" json="$2"
    local tmp
    tmp=$(mktemp)
    echo "$json" > "$tmp"
    $KAVACH $subcmd --hook < "$tmp" 2>&1
    rm -f "$tmp"
}

echo "=== STUBS vs REAL ==="
echo "--- Gates ---"
for g in pre-write post-write pre-tool post-tool intent ceo ast bash read skill lint research content quality enforcer context dag task code-guard chain subagent failure mockdata; do
    tmp=$(mktemp)
    echo '{"toolName":"Read","toolInput":{"file_path":"/tmp/x"}}' > "$tmp"
    out=$($KAVACH gates "$g" --hook < "$tmp" 2>&1 || true)
    rm -f "$tmp"
    if echo "$out" | grep -q STUB; then
        echo "  STUB  gates $g"
    else
        echo "  REAL  gates $g"
    fi
done

echo "--- Session ---"
for s in init validate end compact resume land end-hook; do
    out=$($KAVACH session "$s" 2>&1 || true)
    if echo "$out" | grep -q STUB; then
        echo "  STUB  session $s"
    else
        echo "  REAL  session $s"
    fi
done

echo "--- Memory ---"
for m in bank write stm kanban view sync; do
    out=$($KAVACH memory "$m" 2>&1 || true)
    if echo "$out" | grep -q STUB; then
        echo "  STUB  memory $m"
    else
        echo "  REAL  memory $m"
    fi
done

echo "--- Orch ---"
for o in aegis verify task-health dag; do
    out=$($KAVACH orch "$o" 2>&1 || true)
    if echo "$out" | grep -q STUB; then
        echo "  STUB  orch $o"
    else
        echo "  REAL  orch $o"
    fi
done

echo "--- Top-Level ---"
for t in status agents skills lint quality telemetry; do
    out=$($KAVACH "$t" 2>&1 || true)
    if echo "$out" | grep -q STUB; then
        echo "  STUB  $t"
    else
        echo "  REAL  $t"
    fi
done

echo ""
echo "=== BUG HUNT: PRE-TOOL ==="

out=$(run_hook "gates pre-tool" '{"toolName":"Bash","toolInput":{"command":""}}')
check "bash: block empty command" "deny" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"Bash","toolInput":{"command":"cargo build"}}')
check "bash: approve cargo build" "approve" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"Bash","toolInput":{"command":"git status"}}')
check "bash: approve git status" "approve" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"Read","toolInput":{"file_path":"/etc/shadow"}}')
check "read: block /etc/shadow" "deny" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"Read","toolInput":{"file_path":"/tmp/key.pem"}}')
check "read: block .pem extension" "deny" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"Read","toolInput":{"file_path":"/tmp/foo.txt"}}')
check "read: approve normal file" "approve" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"Task","toolInput":{"prompt":"do something"}}')
check "task: block missing subagent_type" "deny" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"Task","toolInput":{"subagent_type":"fake-agent"}}')
check "task: block unknown agent" "deny" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"Task","toolInput":{"subagent_type":"Explore","prompt":"find"}}')
check "task: approve Explore" "approve" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"TaskCreate","toolInput":{"description":"foo"}}')
check "taskcreate: block missing subject" "deny" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"Skill","toolInput":{}}')
check "skill: block empty skill" "deny" "$out"

out=$(run_hook "gates pre-tool" '{"toolName":"AskUserQuestion","toolInput":{}}')
check "unknown: approve AskUserQuestion" "approve" "$out"

echo ""
echo "=== BUG HUNT: PRE-WRITE ==="

out=$(run_hook "gates pre-write" '{"toolName":"Write","toolInput":{"file_path":"/etc/shadow","content":"x"}}')
check "prewrite: block /etc/ write" "deny" "$out"

out=$(run_hook "gates pre-write" '{"toolName":"Write","toolInput":{"file_path":"/bin/evil","content":"x"}}')
check "prewrite: block /bin/ write" "deny" "$out"

out=$(run_hook "gates pre-write" '{"toolName":"Write","toolInput":{"file_path":"/tmp/readme.md","content":"hello"}}')
check "prewrite: approve non-code .md" "approve" "$out"

out=$(run_hook "gates pre-write" '{"toolName":"Edit","toolInput":{"file_path":"/tmp/a.rs","old_string":"impl Foo { fn bar() {} }","new_string":"struct Foo;"}}')
check "prewrite: block impl removal" "deny" "$out"

out=$(run_hook "gates pre-write" '{"toolName":"Edit","toolInput":{"file_path":"/tmp/a.rs","old_string":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","new_string":""}}')
check "prewrite: block complete deletion" "deny" "$out"

echo ""
echo "=== BUG HUNT: POST-WRITE ==="

out=$(run_hook "gates post-write" '{"toolName":"Write","toolInput":{"file_path":"/tmp/app.ts","content":"console.log(\"hi\")"}}')
check "postwrite: block console.log in .ts" "deny" "$out"

out=$(run_hook "gates post-write" '{"toolName":"Write","toolInput":{"file_path":"/tmp/app.tsx","content":"const x = foo as any"}}')
check "postwrite: block as any in .tsx" "deny" "$out"

out=$(run_hook "gates post-write" '{"toolName":"Write","toolInput":{"file_path":"/tmp/app.test.ts","content":"console.log(\"hi\")"}}')
check "postwrite: allow test file" "approve" "$out"

out=$(run_hook "gates post-write" '{"toolName":"Write","toolInput":{"file_path":"/tmp/Dockerfile","content":"FROM node:latest\nRUN echo hi"}}')
check "postwrite: block FROM :latest" "deny" "$out"

out=$(run_hook "gates post-write" '{"toolName":"Write","toolInput":{"file_path":"/tmp/app.py","content":"import logging\nlogging.info(\"ok\")"}}')
check "postwrite: approve clean .py" "approve" "$out"

out=$(run_hook "gates post-write" '{"toolName":"Write","toolInput":{"file_path":"/tmp/app.rs","content":"fn main() -> Result<()> { Ok(()) }"}}')
check "postwrite: approve clean .rs" "approve" "$out"

out=$(run_hook "gates post-write" '{"toolName":"Write","toolInput":{"file_path":"/tmp/app.rs","content":"unsafe { core::ptr::null() }"}}')
check "postwrite: block unsafe without SAFETY" "deny" "$out"

echo ""
echo "=== BUG HUNT: POST-TOOL ==="

out=$(run_hook "gates post-tool" '{"toolName":"WebSearch","toolInput":{"query":"rust patterns"}}')
check "posttool: approve WebSearch" "approve" "$out"

out=$(run_hook "gates post-tool" '{"toolName":"Task","toolInput":{"subagent_type":"Explore"}}')
check "posttool: Task agent complete" "AGENT_COMPLETE" "$out"

out=$(run_hook "gates post-tool" '{"toolName":"EnterPlanMode","toolInput":{}}')
check "posttool: approve unknown tool" "approve" "$out"

echo ""
echo "=== BUG HUNT: INTENT ==="

out=$(run_hook "gates intent" '{"prompt":"hello"}')
check "intent: trivial hello" "UserPromptSubmit" "$out"

out=$(run_hook "gates intent" '{"prompt":"status"}')
check "intent: status query" "BINARY_FIRST" "$out"

out=$(run_hook "gates intent" '{"prompt":"implement user auth"}')
check "intent: implement task" "UserPromptSubmit" "$out"

echo ""
echo "=== BUG HUNT: SESSION ==="

out=$($KAVACH session init 2>&1)
check "session init: has META" "META" "$out"
check "session init: has SESSION" "SESSION" "$out"

out=$($KAVACH session end 2>&1)
check "session end: has END" "END" "$out"
check "session end: has STATE" "STATE" "$out"

out=$($KAVACH session compact 2>&1)
check "session compact: has COMPACT" "COMPACT" "$out"

echo ""
echo "================================================="
echo "PASS: $PASS  FAIL: $FAIL  TOTAL: $((PASS + FAIL))"
if [ "$FAIL" -gt 0 ]; then
    echo -e "\nBUGS:$BUGS"
fi
