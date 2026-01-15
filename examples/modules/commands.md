# Binary Commands - SP/3.0
# Reference for kavach CLI commands

COMMANDS:SESSION
  kavach session init      # Initialize session
  kavach session validate  # Validate state
  kavach session end       # Save end state
  kavach session compact   # Pre-compact save

COMMANDS:MEMORY
  kavach memory bank             # Project context
  kavach memory bank --all       # All projects
  kavach memory bank --status    # Memory health
  kavach memory kanban           # Visual kanban
  kavach memory kanban --status  # Quick status
  kavach memory stm              # Short-term memory
  kavach memory write            # Write to bank

COMMANDS:RESOURCES
  kavach agents              # List agents
  kavach agents --get <name> # Get agent
  kavach skills              # List skills
  kavach skills --get <name> # Get skill
  kavach status              # System status
  kavach tools               # Rust/Zig tools

COMMANDS:GATES
  kavach gates enforcer --hook  # Pipeline enforcement
  kavach gates bash --hook      # Bash sanitization
  kavach gates read --hook      # File access control
  kavach gates ceo --hook       # CEO validation
  kavach gates intent --hook    # Intent classification

COMMANDS:SCAN
  kavach scan [path]         # DACE tree scanner
  kavach scan --depth 3      # Limit depth
  kavach scan --ext go,rs    # Filter extensions
