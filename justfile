# Kavach - Brahmastra Stack
# ==========================

# Variables
go := env("GO", "go")
prefix := env("PREFIX", env("HOME") / ".local")
bin_dir := prefix / "bin"
binary := "kavach"
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
ldflags := "-s -w -X main.version=" + version
symlinks := "ceo-gate bash-sanitizer read-blocker session-init memory-bank enforcer"

# Default recipe
default: build

# Build the Go kavach binary
build: workspace
    @echo "[BUILD]"
    @echo "  binary: {{binary}}"
    @echo "  version: {{version}}"
    @mkdir -p {{bin_dir}}
    cd cmd/kavach && {{go}} build -ldflags "{{ldflags}}" -o "{{bin_dir}}/{{binary}}" .
    @echo "  size: $(du -h {{bin_dir}}/{{binary}} | cut -f1)"
    @echo "  path: {{bin_dir}}/{{binary}}"
    @echo "[COMPLETE]"

# Initialize Go workspace
workspace:
    -{{go}} work init ./cmd/kavach ./shared 2>/dev/null
    -{{go}} work sync 2>/dev/null

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

# Build the Rust kavach binary and install to ~/.local/bin/kavach
build-rust:
    @echo "[BUILD:RUST]"
    @mkdir -p {{bin_dir}}
    cd crates/kavach-cli && cargo build --release
    cp crates/kavach-cli/target/release/kavach "{{bin_dir}}/{{binary}}"
    @echo "  size: $(du -h {{bin_dir}}/{{binary}} | cut -f1)"
    @echo "  path: {{bin_dir}}/{{binary}}"
    @echo "[COMPLETE]"

# Remove build artifacts
clean:
    @echo "[CLEAN]"
    rm -rf dist/
    -cd cmd/kavach && {{go}} clean 2>/dev/null
    -cd shared && {{go}} clean 2>/dev/null
    @echo "[COMPLETE]"

# Run all tests
test: workspace
    @echo "[TEST]"
    @echo "  cmd/kavach..."
    cd cmd/kavach && {{go}} test ./...
    @echo "  shared..."
    cd shared && {{go}} test ./...
    @echo "[COMPLETE]"

# Run linter
lint: workspace
    @echo "[LINT]"
    cd cmd/kavach && {{go}} vet ./...
    cd shared && {{go}} vet ./...
    @echo "[COMPLETE]"

# Format all Go code
fmt:
    @echo "[FORMAT]"
    cd cmd/kavach && {{go}} fmt ./...
    cd shared && {{go}} fmt ./...
    @echo "[COMPLETE]"

# Build release binaries for all platforms
release: workspace
    @echo "[RELEASE]"
    @echo "  version: {{version}}"
    mkdir -p dist
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 {{go}} build -ldflags "{{ldflags}}" -o dist/kavach-linux-amd64 ./cmd/kavach
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 {{go}} build -ldflags "{{ldflags}}" -o dist/kavach-linux-arm64 ./cmd/kavach
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 {{go}} build -ldflags "{{ldflags}}" -o dist/kavach-darwin-amd64 ./cmd/kavach
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 {{go}} build -ldflags "{{ldflags}}" -o dist/kavach-darwin-arm64 ./cmd/kavach
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 {{go}} build -ldflags "{{ldflags}}" -o dist/kavach-windows-amd64.exe ./cmd/kavach
    @echo "[RELEASE_COMPLETE]"

# Generate SHA256 checksums for release
checksums:
    @echo "[CHECKSUMS]"
    cd dist && sha256sum kavach-* > SHA256SUMS.txt

# Download dependencies
deps: workspace
    @echo "[DEPS]"
    cd cmd/kavach && {{go}} mod download
    cd shared && {{go}} mod download
    @echo "[COMPLETE]"

# Tidy go.mod files
tidy:
    @echo "[TIDY]"
    cd cmd/kavach && {{go}} mod tidy
    cd shared && {{go}} mod tidy
    @echo "[COMPLETE]"

# Show build status
status:
    @echo "[STATUS]"
    @echo "  binary: {{binary}}"
    @echo "  version: {{version}}"
    @echo "  go: $({{go}} version | cut -d' ' -f3)"
    @echo "  installed: $(which {{binary}} 2>/dev/null || echo 'not found')"
    @if [ -f "{{bin_dir}}/{{binary}}" ]; then echo "  size: $(du -h {{bin_dir}}/{{binary}} | cut -f1)"; fi
