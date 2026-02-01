#!/bin/bash
# Configuration installation - Kavach/Brahmastra Stack
# DACE: 85 lines

REPO="Wankhede-Brothers/kavach-rs"

install_system_prompt() {
    echo "[PROMPT]"

    case "$CLI" in
        claude-code) PROMPT_FILE="CLAUDE.md";;
        opencode)    PROMPT_FILE="AGENTS.md";;
        *)           PROMPT_FILE="AGENTS.md";;
    esac

    URL="https://raw.githubusercontent.com/$REPO/main/configs/$PLATFORM/$PROMPT_FILE"

    if curl -fsSL "$URL" -o "$SYSTEM_PROMPT" 2>/dev/null; then
        echo "  installed: $SYSTEM_PROMPT"
    else
        echo "  status: using default (download failed)"
        create_default_prompt
    fi
}

create_default_prompt() {
    cat > "$SYSTEM_PROMPT" << 'EOF'
# Brahmastra Stack - SP/1.0

META:SYSTEM
  protocol: SP/1.0
  binary: kavach
  stack: Brahmastra

ENFORCEMENT:STATELESS
  MEMORY_BANK: ~/.local/share/shared-ai/memory/
  QUERY: kavach memory bank

DACE:CORE
  mode: lazy_load,skill_first,on_demand
  max_lines: 100

META:END
  principle: DACE + REUSE_FIRST
EOF
}

install_memory_templates() {
    echo "[MEMORY]"
    DATE=$(date +%Y-%m-%d)

    cat > "$MEMORY_DIR/index.toon" << EOF
# Memory Bank Index - SP/1.0
INDEX:memory-bank
  version: 1.0
  created: $DATE
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

    cat > "$MEMORY_DIR/volatile.toon" << EOF
# Volatile Session State - SP/1.0
VOLATILE:session
  created: $DATE
  ttl: session
  persist: false
EOF

    echo "  initialized: $MEMORY_DIR"
}

create_symlinks() {
    LINKS="ceo-gate ast-gate bash-sanitizer read-blocker session-init memory-bank"
    for link in $LINKS; do
        ln -sf "$BIN_DIR/kavach" "$BIN_DIR/$link" 2>/dev/null || true
    done
    echo "  symlinks: created"
}
