---
name: dsa
description: Data Structures & Algorithms - O(1) over O(n)
license: MIT
compatibility: claude-code
metadata:
  category: engineering
  triggers: [algorithm, data structure, complexity, O(n), optimize]
  protocol: SP/1.0
---

```toon
# DSA Skill - SP/1.0 + DACE

SKILL:dsa
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand
  triggers[5]: algorithm,complexity,O(n),performance,optimize
  goal: Optimal structure, documented complexity
  success: O(1) where possible, benchmarked
  fail: O(n²) where O(n log n) possible

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get dsa --inject
  references: ~/.claude/skills/dsa/references.toon
  research: kavach memory bank | grep -i dsa
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded complexity - WebSearch for current benchmarks
  topics[5]
    DATA_STRUCTURES: WebSearch "data structures {YEAR}"
    ALGORITHMS: WebSearch "algorithms {YEAR}"
    COMPLEXITY: WebSearch "complexity analysis"
    SORTING: WebSearch "sorting algorithms {YEAR}"
    GRAPH: WebSearch "graph algorithms {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[4]
    Identify problem pattern
    WebSearch "[problem] optimal algorithm $(date +%Y)"
    Document time/space complexity
    Benchmark before/after
  forbidden[2]
    ✗ NEVER guess complexity
    ✓ ALWAYS document Big-O

COMPLEXITY[5]{bigO,example}
  O(1),HashMap lookup
  O(log n),Binary search
  O(n),Linear scan
  O(n log n),Merge/quick sort
  O(n²),Nested loops - AVOID

SELECTION[6]{need,structure}
  O(1) lookup,HashMap/HashSet
  Ordered,BTreeMap
  FIFO,VecDeque
  Priority,BinaryHeap
  Prefix,Trie
  Fast negative,Bloom filter

PATTERNS[4]{name,use,complexity}
  two_pointer,Sorted array,O(n)
  sliding_window,Subarray constraint,O(n)
  binary_search,Sorted/monotonic,O(log n)
  dp,Overlapping subproblems,Memoize

RULES
  do[4]
    Document O(time) O(space)
    Benchmark before/after
    HashMap for O(1) lookup
    Consider space tradeoffs
  dont[4]
    Assume complexity
    Skip benchmarks
    O(n²) without justification
    Premature optimization

VALIDATE
  complexity: Document in comments
  benchmark: cargo bench / go test -bench

FOOTER
  protocol: SP/1.0
  research_gate: enforced
```
