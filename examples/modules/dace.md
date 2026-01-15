# DACE - Dynamic Agentic Context Engineering
# Lazy-loaded for detailed DACE rules

DACE:MICRO_MODULAR
  principle: Smallest files, deepest structure
  applies_to: ALL projects (Go, Rust, TS, Python)
  max_lines: 100 (hard block by enforcer)
  warn_lines: 50 (suggest split)
  ideal_lines: 30-50
  depth: min=3, avg=5, max=7

DACE:BENEFITS
  - LLM reads only 50 lines, not 500
  - Auto-doc comments fit naturally
  - Single responsibility per file
  - Tree scan shows structure without content
  - Faster context loading
  - Better code navigation

DACE:FILE_NAMING
  types: types.{ext}, models.{ext}
  impl: impl.{ext}, service.{ext}
  traits: traits.{ext}, interfaces.{ext}
  handlers: handlers.{ext}, routes.{ext}
  utils: utils.{ext}, helpers.{ext}
  mod: mod.rs, index.ts (<10 lines)

DACE:AUTO_DOC
  line_1-5: Module/file purpose
  line_6-10: DACE compliance note
  line_11+: Imports then code

DACE:SCAN_FIRST
  BEFORE writing any file:
  1. kavach scan [path]
  2. Identify existing modules
  3. Check for similar files
  4. Decide: extend or create new
