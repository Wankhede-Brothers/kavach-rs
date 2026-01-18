#!/bin/bash
# rust-cli.sh - Rust CLI tools installation for DACE
# DACE: Single responsibility - Rust tool installation only
# Protocol: SP/1.0

# Rust CLI tool specifications
# Format: tool_name|cargo_name|description
RUST_TOOLS=(
    "bat|bat|cat replacement with syntax highlighting"
    "eza|eza|ls replacement with icons and git"
    "fd|fd-find|find replacement (faster)"
    "rg|ripgrep|grep replacement (ripgrep)"
    "sd|sd|sed replacement (simpler)"
    "procs|procs|ps replacement"
    "dust|du-dust|du replacement"
    "btm|bottom|top replacement"
    "delta|git-delta|diff replacement"
)

# Check if cargo is installed
check_cargo() {
    if command -v cargo &>/dev/null; then
        echo "[CARGO] $(cargo --version)"
        return 0
    fi
    return 1
}

# Install cargo via rustup
install_cargo() {
    echo "[INSTALL] Installing Rust toolchain..."
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    source "$HOME/.cargo/env"
}

# Check if a tool is installed
tool_installed() {
    command -v "$1" &>/dev/null
}

# Install a single Rust tool
install_rust_tool() {
    local tool="$1"
    local cargo_name="$2"
    local desc="$3"

    if tool_installed "$tool"; then
        echo "[OK] $tool already installed"
        return 0
    fi

    echo "[INSTALL] $tool ($desc)..."
    cargo install "$cargo_name" --locked 2>/dev/null || cargo install "$cargo_name"
}

# Install all Rust CLI tools
install_rust_tools() {
    echo ""
    echo "[RUST_CLI]"
    echo "desc: Installing Rust-based CLI tools for DACE"
    echo ""

    # Ensure cargo is available
    if ! check_cargo; then
        install_cargo
    fi

    # Install each tool
    for spec in "${RUST_TOOLS[@]}"; do
        IFS='|' read -r tool cargo_name desc <<< "$spec"
        install_rust_tool "$tool" "$cargo_name" "$desc"
    done

    echo ""
    echo "[COMPLETE] Rust CLI tools installed"
}

# Print Rust CLI status
print_rust_status() {
    echo ""
    echo "[RUST_CLI_STATUS]"
    printf "%-8s %-12s %s\n" "Legacy" "Rust" "Status"
    printf "%-8s %-12s %s\n" "------" "----" "------"

    local legacy_tools=("cat" "ls" "find" "grep" "sed" "ps" "du" "top" "diff")
    local rust_tools=("bat" "eza" "fd" "rg" "sd" "procs" "dust" "btm" "delta")

    for i in "${!legacy_tools[@]}"; do
        local legacy="${legacy_tools[$i]}"
        local rust="${rust_tools[$i]}"
        local status="not installed"

        if tool_installed "$rust"; then
            status="available"
        fi

        printf "%-8s %-12s %s\n" "$legacy" "$rust" "$status"
    done
    echo ""
}

# Create shell aliases for Rust tools
create_rust_aliases() {
    local shell_rc="$HOME/.bashrc"
    [ -f "$HOME/.zshrc" ] && shell_rc="$HOME/.zshrc"

    local alias_block='
# DACE Rust CLI Aliases (added by kavach)
if command -v bat &>/dev/null; then alias cat="bat --plain"; fi
if command -v eza &>/dev/null; then
    alias ls="eza --icons"
    alias ll="eza -la --icons --git"
    alias lt="eza --tree --icons --level=2"
fi
if command -v fd &>/dev/null; then alias find="fd"; fi
if command -v rg &>/dev/null; then alias grep="rg"; fi
if command -v sd &>/dev/null; then alias sed="sd"; fi
if command -v procs &>/dev/null; then alias ps="procs"; alias psa="procs -a"; fi
if command -v dust &>/dev/null; then alias du="dust"; fi
if command -v btm &>/dev/null; then alias top="btm"; fi
if command -v delta &>/dev/null; then alias diff="delta"; fi
# End DACE Rust CLI Aliases
'

    # Check if already added
    if grep -q "DACE Rust CLI Aliases" "$shell_rc" 2>/dev/null; then
        echo "[ALIASES] Already configured in $shell_rc"
        return 0
    fi

    echo "$alias_block" >> "$shell_rc"
    echo "[ALIASES] Added to $shell_rc"
    echo "[NOTE] Run: source $shell_rc"
}

# Main function for standalone execution
rust_cli_main() {
    case "${1:-install}" in
        install)
            install_rust_tools
            create_rust_aliases
            print_rust_status
            ;;
        status)
            print_rust_status
            ;;
        aliases)
            create_rust_aliases
            ;;
        *)
            echo "Usage: rust-cli.sh [install|status|aliases]"
            ;;
    esac
}

# Run if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    rust_cli_main "$@"
fi
