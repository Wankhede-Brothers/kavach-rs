#!/bin/bash
# Kavach Installer - Brahmastra Stack
# Usage: curl -fsSL https://raw.githubusercontent.com/Wankhede-Brothers/kavach-rs/main/install/install.sh | bash
# Or: curl ... | bash -s -- --cli opencode --rust-cli
# DACE: Single responsibility - kavach installation

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REMOTE_BASE="https://raw.githubusercontent.com/Wankhede-Brothers/kavach-rs/main/install"
INSTALL_RUST_CLI=false

# Source libraries (local or remote)
source_lib() {
    if [ -f "$SCRIPT_DIR/lib/$1" ]; then
        source "$SCRIPT_DIR/lib/$1"
    else
        source <(curl -fsSL "$REMOTE_BASE/lib/$1")
    fi
}

print_header() {
    echo ""
    echo "╔════════════════════════════════════════════╗"
    echo "║  KAVACH - Brahmastra Stack Installer       ║"
    echo "║  Protocol: SP/1.0 (Sutra Protocol)        ║"
    echo "║  DACE: Rust CLI + Go Binary               ║"
    echo "╚════════════════════════════════════════════╝"
    echo ""
}

print_success() {
    echo ""
    echo "[COMPLETE]"
    echo "  binary: $BIN_DIR/kavach"
    echo "  memory: $MEMORY_DIR"
    echo "  prompt: $SYSTEM_PROMPT"
    echo ""
    if [ "$INSTALL_RUST_CLI" = true ]; then
        echo "[RUST_CLI]"
        echo "  bat, eza, fd, rg, sd, procs, dust, btm, delta"
        echo ""
    fi
    echo "[NEXT]"
    echo "  1. Restart terminal or: source ~/.bashrc"
    echo "  2. Run: kavach status"
    echo "  3. Run: kavach memory bank"
    echo "  4. Run: kavach memory view --tree"
    echo ""
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --cli) CLI_ARG="$2"; shift 2;;
            --rust-cli) INSTALL_RUST_CLI=true; shift;;
            --full) INSTALL_RUST_CLI=true; shift;;
            --help)
                echo "Usage: install.sh [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --cli <name>    Target CLI (claude-code|opencode)"
                echo "  --rust-cli      Install Rust CLI tools (bat, eza, fd, rg, etc.)"
                echo "  --full          Full install (kavach + Rust CLI tools)"
                echo "  --help          Show this help"
                exit 0;;
            *) shift;;
        esac
    done
}

main() {
    parse_args "$@"
    print_header

    source_lib "detect.sh"
    source_lib "paths.sh"
    source_lib "download.sh"
    source_lib "config.sh"

    detect_platform
    detect_cli "$CLI_ARG"
    print_detect

    if [ "$PLATFORM" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        echo "[ERROR] Unsupported platform: $PLATFORM/$ARCH"
        exit 1
    fi

    set_paths
    print_paths
    create_directories

    install_binary
    install_system_prompt
    install_memory_templates
    create_symlinks

    # Optional: Install Rust CLI tools
    if [ "$INSTALL_RUST_CLI" = true ]; then
        source_lib "rust-cli.sh"
        install_rust_tools
        create_rust_aliases
    fi

    print_success
}

main "$@"
