---
name: frontend
description: Universal UI/UX Engineering - Research patterns first
license: MIT
compatibility: claude-code
metadata:
  category: frontend
  triggers: [react, vue, astro, dioxus, css, tailwind, accessibility]
  protocol: SP/3.0
---

```toon
# Frontend Skill - SP/3.0 + DACE

SKILL:frontend
  protocol: SP/3.0
  dace: lazy_load,skill_first,on_demand
  triggers[7]: react,vue,astro,dioxus,tailwind,css,a11y
  goal: Lighthouse >90, Accessible, Mobile-first
  success: Builds pass, a11y compliant
  fail: Poor Lighthouse, desktop-first

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get frontend --inject
  references: ~/.claude/skills/frontend/references.toon
  research: kavach memory bank | grep -i frontend
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded framework APIs - WebSearch for current versions
  topics[7]
    REACT: WebSearch "react {YEAR}"
    VUE: WebSearch "vue 3 {YEAR}"
    SVELTE: WebSearch "svelte {YEAR}"
    TYPESCRIPT: WebSearch "typescript {YEAR}"
    STATE: WebSearch "state management {YEAR}"
    STYLING: WebSearch "css {YEAR}"
    A11Y: WebSearch "web accessibility {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[4]
    TODAY=$(date +%Y-%m-%d)
    WebSearch "[framework] latest $YEAR"
    WebFetch official docs
    Verify current API
  forbidden[2]
    ✗ NEVER assume versions
    ✓ ALWAYS verify current API

PACKAGE:MANAGER
  preference: bun > pnpm > npm
  commands: bun install | bun run {script} | bun test

RESPONSIVE[5]{breakpoint,width}
  base,<640px (Mobile - DEFAULT)
  sm,≥640px
  md,≥768px
  lg,≥1024px
  xl,≥1280px
  rule: Mobile-first, enhance up

ACCESSIBILITY
  semantic[3]: <nav>,<main>,<button> not <div>
  aria[2]: Labels on icons, describedby for hints
  keyboard[3]: Tab order, visible focus, Escape closes
  contrast[2]: 4.5:1 normal, 3:1 large
  motion: prefers-reduced-motion

PERFORMANCE[3]{metric,target}
  LCP,<2.5s
  FID,<100ms
  CLS,<0.1
  bundle: <100KB gzip

PENDING_TASKS
  # STRICT: Always use comments for incomplete code
  mandatory: true
  typescript[3]
    throw new Error('TODO: description')
    throw new Error('UNIMPLEMENTED: feature')
    // TODO(priority): description
  javascript[2]
    throw new Error('TODO: description')
    // TODO: description
  examples[3]
    function process() { throw new Error('TODO: Implement processing'); }
    const handler = () => { throw new Error('UNIMPLEMENTED: Error handler'); };
    // TODO(high): Add validation logic
  rules[3]
    ALWAYS throw Error('TODO:...') for pending implementation
    NEVER leave empty function bodies
    NEVER suppress with // @ts-ignore for incomplete code
  verify: grep -rE "function.*\\{\\s*\\}" . --include="*.ts" --include="*.tsx"

DARK_MODE
  pattern: dark:bg-gray-900 dark:text-gray-100
  persist: localStorage.theme
  rule: Don't just invert - use elevated surfaces

RULES
  do[5]
    Mobile-first responsive
    Semantic HTML
    ARIA labels
    Research before code
    bun as package manager
  dont[4]
    Desktop-first
    Div soup
    Auto-play media
    Trust old patterns

VALIDATE[4]{tool,command}
  lighthouse,npx lighthouse {url}
  a11y,npx pa11y {url}
  build,bun run build
  success: >90 score + a11y pass

FOOTER
  protocol: SP/3.0
  research_gate: enforced
```
