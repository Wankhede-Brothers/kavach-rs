---
name: sql
description: PostgreSQL & SQLx - Index Scan optimization
license: MIT
compatibility: claude-code
metadata:
  category: database
  triggers: [sql, postgresql, query, index, migration, sqlx]
  protocol: SP/3.0
---

```toon
# SQL Skill - SP/3.0 + DACE

SKILL:sql
  protocol: SP/3.0
  dace: lazy_load,skill_first,on_demand
  triggers[6]: sql,postgresql,query,index,migration,sqlx
  goal: Index Scan on all queries
  success: EXPLAIN shows Index Scan
  fail: Seq Scan on large tables, OFFSET, SELECT *

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get sql --inject
  references: ~/.claude/skills/sql/references.toon
  research: kavach memory bank | grep -i sql
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded SQL - WebSearch for current PostgreSQL patterns
  topics[6]
    POSTGRESQL: WebSearch "postgresql {YEAR}"
    INDEXING: WebSearch "postgresql index {YEAR}"
    SQLX: WebSearch "sqlx rust {YEAR}"
    TRANSACTIONS: WebSearch "postgresql transactions"
    CONNECTION: WebSearch "connection pooling {YEAR}"
    OPTIMIZATION: WebSearch "query optimization {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[5]
    TODAY=$(date +%Y-%m-%d)
    WebSearch "postgresql [pattern] $YEAR"
    WebSearch "sqlx rust $YEAR"
    WebFetch postgresql.org/docs/current
    Verify PG features
  forbidden[2]
    ✗ NEVER assume PG features
    ✓ ALWAYS EXPLAIN ANALYZE

GOLDEN_RULES[5]
  1. WHERE/JOIN/ORDER BY → Index
  2. NEVER OFFSET → Use keyset
  3. NEVER SELECT * → Specify columns
  4. Multi-statement → Transaction
  5. ALWAYS EXPLAIN ANALYZE

INDEXES[4]{type,use}
  B-Tree,Equality + range + ORDER BY
  Composite,Multi-column (order matters!)
  Partial,Common WHERE (smaller)
  GIN,JSONB + arrays + full-text

PAGINATION
  wrong: OFFSET 10000 (scans all)
  right: WHERE id > $cursor ORDER BY id LIMIT 20

SQLX[3]{feature,syntax}
  compile_check,sqlx::query_as!(Type "...")
  transaction,pool.begin() → tx.commit()
  offline,cargo sqlx prepare

RULES
  do[5]
    Keyset pagination
    Specify columns
    Index for WHERE/ORDER
    EXPLAIN ANALYZE
    Transactions for multi-op
  dont[4]
    OFFSET at scale
    SELECT *
    Unindexed large tables
    Trust without EXPLAIN

VALIDATE
  explain: EXPLAIN ANALYZE {query}
  success: Index Scan
  failure: Seq Scan → Add index

FOOTER
  protocol: SP/3.0
  research_gate: enforced
```
