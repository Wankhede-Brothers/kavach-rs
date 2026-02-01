# Kavach - Brahmastra Stack (Rust)
# =================================

# Variables
prefix := env("PREFIX", env("HOME") / ".local")
bin_dir := prefix / "bin"
binary := "kavach"
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
symlinks := "ceo-gate bash-sanitizer read-blocker session-init memory-bank enforcer"

# Default recipe
default: build

# Build the Rust kavach binary and install to ~/.local/bin/kavach
build:
    @echo "[BUILD]"
    @echo "  binary: {{binary}}"
    @echo "  version: {{version}}"
    @mkdir -p {{bin_dir}}
    cd crates/kavach-cli && cargo build --release
    cp crates/kavach-cli/target/release/kavach "{{bin_dir}}/{{binary}}"
    @echo "  size: $(du -h {{bin_dir}}/{{binary}} | cut -f1)"
    @echo "  path: {{bin_dir}}/{{binary}}"
    @echo "[COMPLETE]"

# Build debug binary (faster compile, no install)
build-debug:
    cd crates/kavach-cli && cargo build

# Build and install with symlinks
install: build symlinks
    @echo "[INSTALL]"
    @echo "  binary: {{bin_dir}}/{{binary}}"
    @echo "[COMPLETE]"

# Create backward-compatible symlinks
symlinks:
    @echo "[SYMLINKS]"
    for name in {{symlinks}}; do ln -sf {{binary}} "{{bin_dir}}/$name" 2>/dev/null || true; done
    @echo "  created: {{symlinks}}"

# Remove kavach binary and symlinks
uninstall:
    @echo "[UNINSTALL]"
    rm -f "{{bin_dir}}/{{binary}}"
    for name in {{symlinks}}; do rm -f "{{bin_dir}}/$name"; done
    @echo "[COMPLETE]"

# Remove build artifacts
clean:
    @echo "[CLEAN]"
    rm -rf dist/
    cd crates/kavach-cli && cargo clean
    @echo "[COMPLETE]"

# Run Rust tests
test:
    @echo "[TEST]"
    cd crates/kavach-cli && cargo test
    @echo "[COMPLETE]"

# Run clippy linter
lint:
    @echo "[LINT]"
    cd crates/kavach-cli && cargo clippy -- -D warnings
    @echo "[COMPLETE]"

# Format Rust code
fmt:
    @echo "[FORMAT]"
    cd crates/kavach-cli && cargo fmt
    @echo "[COMPLETE]"

# Check formatting without modifying
fmt-check:
    cd crates/kavach-cli && cargo fmt -- --check

# Build release binaries for all platforms
release:
    @echo "[RELEASE]"
    @echo "  version: {{version}}"
    mkdir -p dist
    cd crates/kavach-cli && cargo build --release --target aarch64-apple-darwin 2>/dev/null || true
    cd crates/kavach-cli && cargo build --release
    cp crates/kavach-cli/target/release/kavach dist/kavach-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)
    @echo "[RELEASE_COMPLETE]"

# Show build status
status:
    @echo "[STATUS]"
    @echo "  binary: {{binary}}"
    @echo "  version: {{version}}"
    @echo "  rustc: $(rustc --version)"
    @echo "  installed: $(which {{binary}} 2>/dev/null || echo 'not found')"
    @if [ -f "{{bin_dir}}/{{binary}}" ]; then echo "  size: $(du -h {{bin_dir}}/{{binary}} | cut -f1)"; fi
