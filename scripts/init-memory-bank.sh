#!/bin/bash
# Memory Bank Initialization Script
# Part of Kavach - Brahmastra Stack
# https://github.com/Wankhede-Brothers/kavach-rs

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
MEMORY_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/shared-ai/memory"
SCHEMA_VERSION="3.0"

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  Memory Bank Initialization v${SCHEMA_VERSION}${NC}"
echo -e "${BLUE}  Kavach - Brahmastra Stack${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""

# Check if kavach binary exists
if ! command -v kavach &> /dev/null; then
    echo -e "${RED}Error: kavach binary not found${NC}"
    echo -e "${YELLOW}Install: curl -fsSL https://raw.githubusercontent.com/Wankhede-Brothers/kavach-rs/main/install/install.sh | bash${NC}"
    exit 1
fi

# Check if already initialized
if [ -f "${MEMORY_DIR}/index.toon" ]; then
    echo -e "${YELLOW}Memory Bank already exists at ${MEMORY_DIR}${NC}"
    echo ""
    read -p "Do you want to reinitialize? This will wipe existing data. [y/N] " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 0
    fi
    echo -e "${YELLOW}Backing up existing Memory Bank...${NC}"
    mv "${MEMORY_DIR}" "${MEMORY_DIR}.backup.$(date +%s)"
fi

# Create directory structure
echo -e "${BLUE}Initializing Memory Bank...${NC}"
mkdir -p "${MEMORY_DIR}"/{decisions,patterns,research,kanban,proposals,roadmaps,graph,STM}

# Create index.toon
DATE=$(date +%Y-%m-%d)
cat > "${MEMORY_DIR}/index.toon" << EOF
# Memory Bank Index - SP/1.0
INDEX:memory-bank
  version: 1.0
  created: ${DATE}
  protocol: SP/1.0

STRUCTURE[8]{dir,purpose}
  decisions/,Architecture decisions
  research/,Research findings
  roadmaps/,Project roadmaps
  patterns/,Code patterns
  proposals/,Feature proposals
  kanban/,Task boards
  graph/,Knowledge graphs
  STM/,Session context
EOF

# Create volatile.toon
cat > "${MEMORY_DIR}/volatile.toon" << EOF
# Volatile Session State - SP/1.0
VOLATILE:session
  created: ${DATE}
  ttl: session
  persist: false
EOF

# Create GOVERNANCE.toon
cat > "${MEMORY_DIR}/GOVERNANCE.toon" << EOF
# Memory Bank Governance - SP/1.0
GOVERNANCE:memory-bank
  version: 1.0
  created: ${DATE}

RULES[4]{name,description}
  project_isolation,Scope queries to active project
  no_hardcode,WebSearch for current patterns
  toon_format,Use TOON for all memory files
  file_locking,Prevent concurrent write conflicts
EOF

echo -e "${GREEN}Memory Bank initialized successfully${NC}"
echo ""
echo -e "${GREEN}Memory Bank structure created:${NC}"
echo "  ${MEMORY_DIR}/"
echo "     ├── decisions/    (Architecture decisions)"
echo "     ├── patterns/     (Code patterns)"
echo "     ├── research/     (Research findings)"
echo "     ├── kanban/       (Task boards)"
echo "     ├── proposals/    (Feature proposals)"
echo "     ├── roadmaps/     (Project roadmaps)"
echo "     ├── graph/        (Knowledge graphs)"
echo "     ├── STM/          (Session context)"
echo "     ├── GOVERNANCE.toon"
echo "     ├── index.toon"
echo "     └── volatile.toon"
echo ""
echo -e "${GREEN}============================================${NC}"
echo -e "${GREEN}Memory Bank initialization complete!${NC}"
echo -e "${GREEN}============================================${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "  1. Verify:  kavach status"
echo "  2. Check:   kavach memory bank"
echo "  3. Kanban:  kavach memory kanban"
echo ""
