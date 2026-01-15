#!/bin/bash
# Platform and CLI detection - Kavach/Brahmastra Stack
# DACE: 45 lines

detect_platform() {
    case "$(uname -s)" in
        Linux*)     PLATFORM="linux";;
        Darwin*)    PLATFORM="darwin";;
        MINGW*|CYGWIN*|MSYS*) PLATFORM="windows";;
        *)          PLATFORM="unknown";;
    esac

    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)       ARCH="amd64";;
        aarch64|arm64) ARCH="arm64";;
        *)            ARCH="unknown";;
    esac

    export PLATFORM ARCH
}

detect_cli() {
    # CLI detection order: argument > environment > directory
    if [ -n "$1" ]; then
        CLI="$1"
    elif [ -n "$KAVACH_CLI" ]; then
        CLI="$KAVACH_CLI"
    elif [ -d "$HOME/.claude" ]; then
        CLI="claude-code"
    elif [ -n "$OPENCODE_HOME" ]; then
        CLI="opencode"
    else
        CLI="claude-code"
    fi
    export CLI
}

print_detect() {
    echo "[DETECT]"
    echo "  platform: $PLATFORM"
    echo "  arch: $ARCH"
    echo "  cli: $CLI"
}
