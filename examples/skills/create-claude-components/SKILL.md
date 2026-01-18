---
name: create-claude-components
description: Claude Code Extension Builder - Skills, Hooks, Agents, Commands
license: MIT
compatibility: claude-code
metadata:
  category: meta
  triggers: [create skill, create hook, create agent, create command, SKILL.md]
  protocol: SP/1.0
---

```toon
# Create Claude Components Skill - SP/1.0 + DACE

SKILL:create-claude-components
  protocol: SP/1.0
  dace: lazy_load,skill_first,on_demand
  triggers[5]: create skill,create hook,create agent,create command,SKILL.md
  goal: Generate valid Claude Code extension components
  success: Component created + validated + loads correctly
  fail: Invalid structure, missing fields, won't load

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get create-claude-components --inject
  references: ~/.claude/skills/create-claude-components/references.toon
  agents: kavach agents
  skills: kavach skills
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded formats - WebSearch for current Claude Code docs
  topics[5]
    HOOKS: WebSearch "claude code hooks {YEAR}"
    SKILLS: WebSearch "claude code skills {YEAR}"
    AGENTS: WebSearch "claude code agents {YEAR}"
    MCP: WebSearch "model context protocol {YEAR}"
    SETTINGS: WebSearch "claude code settings {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[4]
    TODAY=$(date +%Y-%m-%d)
    WebSearch "Claude Code skills hooks $YEAR"
    WebFetch official Claude Code docs
    Verify current component structure
  forbidden[2]
    ✗ NEVER hardcode versions
    ✓ ALWAYS validate component loads

COMPONENTS:PATHS[4]{type,location,format}
  Skill,~/.claude/skills/[name]/skill.md,YAML frontmatter + SP/1.0 TOON
  Hook,~/.claude/hooks.json,JSON matcher + command
  Agent,~/.claude/agents/[name].md,YAML frontmatter + SP/1.0 TOON
  Command,~/.claude/commands/[name].md,YAML frontmatter + prompt

TEMPLATE:SKILL
  structure: |
    ---
    name: {skill-name}
    description: {one-line}
    license: MIT
    compatibility: claude-code
    metadata:
      category: {category}
      triggers: [{keywords}]
      protocol: SP/1.0
    ---
    ```toon
    SKILL:{name}
      protocol: SP/1.0
      triggers[N]: {comma-separated}
      goal: {what it achieves}
      success: {acceptance criteria}
      fail: {anti-patterns}
    RESEARCH:GATE
      mandatory: true
    RULES
      do[N]: {best practices}
      dont[N]: {anti-patterns}
    ```
  max_lines: 150
  max_chars: 2000

TEMPLATE:HOOK
  format: |
    {
      "hooks": {
        "{EventType}": [{
          "matcher": "{tool pattern}",
          "command": "{executable path}"
        }]
      }
    }
  events[5]: PreToolUse,PostToolUse,SessionStart,UserPromptSubmit,Stop
  rule: command must be executable (chmod +x)

RULES
  do[5]
    Use YAML frontmatter
    Keep under 150 lines
    Include RESEARCH:GATE
    Use $(date +%Y) not hardcoded years
    Validate component loads
  dont[4]
    Hardcode versions/syntax
    Exceed 2000 characters
    Skip research gate
    Create without testing

VALIDATE:COMPONENT[5]{check,status}
  YAML frontmatter valid,[ ]
  Under 150 lines,[ ]
  Has RESEARCH:GATE,[ ]
  No hardcoded dates,[ ]
  Component loads in Claude Code,[ ]

FOOTER
  protocol: SP/1.0
  max_lines: 150
```
