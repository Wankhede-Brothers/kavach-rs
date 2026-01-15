# Reusability Patterns - SP/3.0
# Lazy-loaded for reuse-first principles

REUSE:PRINCIPLE
  RULE: NEVER recreate, ALWAYS reuse
  priority: Search → Check → Import → Create

REUSE:BEFORE_CREATE
  1. Grep for existing function/type
  2. Check shared/, common/, utils/ dirs
  3. Check if similar implementation exists
  4. ONLY create if truly unique

REUSE:SEARCH_PATTERNS
  go:
    func: Grep "func.*{name}"
    type: Grep "type.*{name}"
  rust:
    fn: Grep "fn.*{name}"
    struct: Grep "struct.*{name}"
    trait: Grep "trait.*{name}"
  typescript:
    function: Grep "function.*{name}|const.*{name}"
    type: Grep "type.*{name}|interface.*{name}"
    component: Grep "export.*{name}"

REUSE:SHARED_LOCATIONS
  go: pkg/shared/, internal/common/
  rust: src/shared/, src/common/
  typescript: src/lib/, src/utils/, src/shared/
  python: src/common/, src/utils/

REUSE:FORBIDDEN
  - Duplicate functions (different names, same logic)
  - Copy-paste code blocks (>5 lines)
  - Recreating existing utilities
  - Inline implementations of shared logic

REUSE:REQUIRED
  - Import from shared modules
  - Extend existing types (composition)
  - Use type aliases for clarity
  - Extract repeated patterns to utils
