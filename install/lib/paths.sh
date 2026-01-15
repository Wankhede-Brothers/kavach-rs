#!/bin/bash
# Path configuration - Kavach/Brahmastra Stack
# DACE: 55 lines

set_paths() {
    case "$PLATFORM" in
        linux)
            DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}"
            CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"
            MEMORY_DIR="$DATA_HOME/shared-ai/memory"
            ;;
        darwin)
            DATA_HOME="$HOME/Library/Application Support"
            CONFIG_HOME="$DATA_HOME"
            MEMORY_DIR="$DATA_HOME/shared-ai/memory"
            ;;
    esac

    case "$CLI" in
        claude-code)
            BIN_DIR="$HOME/.claude/bin"
            SETTINGS_PATH="$HOME/.claude/settings.json"
            SYSTEM_PROMPT="$HOME/.claude/CLAUDE.md"
            ;;
        opencode)
            BIN_DIR="$HOME/.local/bin"
            SETTINGS_PATH="$CONFIG_HOME/opencode/settings.json"
            SYSTEM_PROMPT="$CONFIG_HOME/opencode/AGENTS.md"
            ;;
        *)
            BIN_DIR="$HOME/.local/bin"
            SETTINGS_PATH="$CONFIG_HOME/$CLI/settings.json"
            SYSTEM_PROMPT="$CONFIG_HOME/$CLI/AGENTS.md"
            ;;
    esac

    export DATA_HOME CONFIG_HOME MEMORY_DIR BIN_DIR SETTINGS_PATH SYSTEM_PROMPT
}

create_directories() {
    mkdir -p "$BIN_DIR"
    mkdir -p "$MEMORY_DIR"/{decisions,graph,kanban,patterns,proposals,research,roadmaps,STM}
    mkdir -p "$(dirname "$SETTINGS_PATH")"
}

print_paths() {
    echo "[PATHS]"
    echo "  bin: $BIN_DIR/kavach"
    echo "  memory: $MEMORY_DIR"
    echo "  settings: $SETTINGS_PATH"
    echo "  prompt: $SYSTEM_PROMPT"
}
