# Output Style - SP/3.0
# Visual dashboard templates

OUTPUT:RULES
  ALWAYS: ASCII boxes for headers
  ALWAYS: Tables for structured data
  ALWAYS: Progress bars for completion
  ALWAYS: Tree diagrams for hierarchies
  NEVER: Paragraphs or text walls
  NEVER: Raw JSON without formatting

OUTPUT:HEADER
  ╔═══════════════════════════════════════╗
  ║           HEADER TITLE                ║
  ╠═══════════════════════════════════════╣
  ║  Content with consistent formatting   ║
  ╚═══════════════════════════════════════╝

OUTPUT:TABLE
  ┌─────────────┬─────────────┬───────────┐
  │  Column 1   │  Column 2   │  Column 3 │
  ├─────────────┼─────────────┼───────────┤
  │  Value      │  Value      │  Value    │
  └─────────────┴─────────────┴───────────┘

OUTPUT:PROGRESS
  [PROGRESS]
  task: Task Name
  status: ████████████░░░░░░░░ 60%

OUTPUT:TREE
  root/
  ├── branch1/
  │   ├── leaf1
  │   └── leaf2
  ├── branch2/
  │   └── leaf3
  └── branch3/

OUTPUT:STATUS
  ✓ PASS    ✗ FAIL    ○ PENDING
  ◐ RUNNING ⊘ BLOCKED ⚠ WARNING

OUTPUT:TOON_FORMAT
  SECTION:NAME
    key: value
    list[]: item1, item2
    nested:
      sub_key: sub_value
