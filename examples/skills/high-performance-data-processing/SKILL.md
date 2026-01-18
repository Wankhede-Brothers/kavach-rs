---
name: high-performance-data-processing
description: Data Pipeline Optimization - Parquet, Arrow, Polars, Rayon
license: MIT
compatibility: claude-code
metadata:
  category: data
  triggers: [parquet, arrow, polars, data pipeline, 100MB+, batch processing]
  protocol: SP/1.0
---

```toon
# High Performance Data Skill - SP/1.0 + DACE

SKILL:data-processing
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand
  triggers[6]: parquet,arrow,polars,large dataset,pipeline,batch
  goal: 10-100x faster data pipelines
  success: 500GB in <5min, predicate pushdown, parallel
  fail: CSV at scale, row-oriented reads, single-threaded

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get high-performance-data-processing --inject
  references: ~/.claude/skills/high-performance-data-processing/references.toon
  research: kavach memory bank | grep -i data
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded crate versions - WebSearch for current APIs
  topics[6]
    ARROW: WebSearch "apache arrow {YEAR}"
    POLARS: WebSearch "polars {YEAR}"
    PARQUET: WebSearch "parquet {YEAR}"
    RAYON: WebSearch "rayon rust {YEAR}"
    SIMD: WebSearch "SIMD {YEAR}"
    STREAMING: WebSearch "stream processing {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[5]
    TODAY=$(date +%Y-%m-%d)
    WebSearch "polars parquet $YEAR best practices"
    WebSearch "arrow rust $YEAR latest API"
    WebFetch docs.rs/polars, arrow.apache.org
    Verify latest crate versions
  forbidden[2]
    ✗ NEVER hardcode crate versions
    ✓ ALWAYS check latest crate versions

PRINCIPLE:COLUMNAR
  rule: Data > 100MB + analytical → Parquet
  csv_problems[4]
    Reads ALL columns (wastes I/O)
    No predicate pushdown
    Poor compression
    Cannot parallelize
  parquet_wins[4]
    Column pruning (75% I/O savings)
    Predicate pushdown (skip 90%+ rows)
    14x compression (Zstd)
    Parallel row groups

STACK:SELECTION[4]{size,stack,research}
  <100MB,CSV/JSON,N/A
  <50GB,Polars + Rayon,WebSearch "polars $YEAR"
  50GB-1TB,K8s + Polars,WebSearch "polars distributed"
  >1TB,Ballista/DataFusion,WebSearch "ballista $YEAR"

PARALLELISM:LAYERS[4]{layer,speedup}
  1. Multi-file,Distribute across nodes (40x)
  2. Row-group,Rayon work-stealing (10x)
  3. Column,Decode concurrently (8x)
  4. SIMD,Vectorized ops (8x)

DECISION[3]{bound,runtime}
  CPU-bound,Rayon
  I/O-bound,Tokio
  Both,Rayon + Tokio

RULES
  do[6]
    WebSearch crate versions before Cargo.toml
    Parquet for >100MB analytical
    Predicate pushdown
    Rayon for CPU parallelism
    R2 for zero-egress storage
    Profile before optimizing
  dont[5]
    Hardcode crate versions
    CSV at scale
    Load all into RAM
    Single-threaded processing
    Assume improvements (benchmark)

VALIDATE:PIPELINE[5]{check,status}
  WebSearched latest crate versions,[ ]
  10-100x faster than naive,[ ]
  <8x memory vs Pandas,[ ]
  Predicate pushdown enabled,[ ]
  cargo check passes,[ ]

FOOTER
  protocol: SP/1.0
  research_gate: enforced
```
