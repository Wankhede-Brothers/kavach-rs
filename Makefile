# Kavach - Brahmastra Stack
# ==========================

# Go binary to use
GO ?= go

# Installation directory
PREFIX ?= $(HOME)/.local
BIN_DIR = $(PREFIX)/bin

# Binary name
BINARY = kavach

# Version (from git tag or default)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build flags
LDFLAGS = -s -w -X main.version=$(VERSION)

# Symlinks for backward compatibility
SYMLINKS = ceo-gate bash-sanitizer read-blocker session-init memory-bank enforcer

.PHONY: all build install uninstall clean test lint fmt release help workspace

# Default target
all: build

## build: Build the kavach binary
build: workspace
	@echo "[BUILD]"
	@echo "  binary: $(BINARY)"
	@echo "  version: $(VERSION)"
	@mkdir -p $(BIN_DIR)
	@cd cmd/kavach && $(GO) build -ldflags "$(LDFLAGS)" -o "$(BIN_DIR)/$(BINARY)" .
	@echo "  size: $$(du -h $(BIN_DIR)/$(BINARY) | cut -f1)"
	@echo "  path: $(BIN_DIR)/$(BINARY)"
	@echo "[COMPLETE]"

## workspace: Initialize Go workspace
workspace:
	@$(GO) work init ./cmd/kavach ./shared 2>/dev/null || true
	@$(GO) work sync 2>/dev/null || true

## install: Build and install with symlinks
install: build symlinks
	@echo "[INSTALL]"
	@echo "  binary: $(BIN_DIR)/$(BINARY)"
	@echo "  symlinks: $(words $(SYMLINKS))"
	@echo "[COMPLETE]"

## symlinks: Create backward-compatible symlinks
symlinks:
	@echo "[SYMLINKS]"
	@for name in $(SYMLINKS); do \
		ln -sf $(BINARY) "$(BIN_DIR)/$$name" 2>/dev/null || true; \
	done
	@echo "  created: $(SYMLINKS)"

## uninstall: Remove kavach binary and symlinks
uninstall:
	@echo "[UNINSTALL]"
	@rm -f "$(BIN_DIR)/$(BINARY)"
	@for name in $(SYMLINKS); do \
		rm -f "$(BIN_DIR)/$$name"; \
	done
	@echo "[COMPLETE]"

## clean: Remove build artifacts
clean:
	@echo "[CLEAN]"
	@rm -rf dist/
	@cd cmd/kavach && $(GO) clean 2>/dev/null || true
	@cd shared && $(GO) clean 2>/dev/null || true
	@echo "[COMPLETE]"

## test: Run all tests
test: workspace
	@echo "[TEST]"
	@echo "  cmd/kavach..."
	@cd cmd/kavach && $(GO) test ./...
	@echo "  shared..."
	@cd shared && $(GO) test ./...
	@echo "[COMPLETE]"

## lint: Run linter
lint: workspace
	@echo "[LINT]"
	@cd cmd/kavach && $(GO) vet ./...
	@cd shared && $(GO) vet ./...
	@echo "[COMPLETE]"

## fmt: Format all Go code
fmt:
	@echo "[FORMAT]"
	@cd cmd/kavach && $(GO) fmt ./...
	@cd shared && $(GO) fmt ./...
	@echo "[COMPLETE]"

## release: Build release binaries for all platforms
release: workspace
	@echo "[RELEASE]"
	@echo "  version: $(VERSION)"
	@mkdir -p dist
	@echo ""
	@echo "  Building Linux amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o dist/kavach-linux-amd64 ./cmd/kavach
	@echo "  Building Linux arm64..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o dist/kavach-linux-arm64 ./cmd/kavach
	@echo "  Building macOS amd64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o dist/kavach-darwin-amd64 ./cmd/kavach
	@echo "  Building macOS arm64..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o dist/kavach-darwin-arm64 ./cmd/kavach
	@echo "  Building Windows amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o dist/kavach-windows-amd64.exe ./cmd/kavach
	@echo ""
	@echo "[RELEASE_COMPLETE]"
	@ls -lh dist/
	@echo ""

## checksums: Generate SHA256 checksums for release
checksums:
	@echo "[CHECKSUMS]"
	@cd dist && sha256sum kavach-* > SHA256SUMS.txt
	@cat dist/SHA256SUMS.txt

## deps: Download dependencies
deps: workspace
	@echo "[DEPS]"
	@cd cmd/kavach && $(GO) mod download
	@cd shared && $(GO) mod download
	@echo "[COMPLETE]"

## tidy: Tidy go.mod files
tidy:
	@echo "[TIDY]"
	@cd cmd/kavach && $(GO) mod tidy
	@cd shared && $(GO) mod tidy
	@echo "[COMPLETE]"

## status: Show build status
status:
	@echo "[STATUS]"
	@echo "  binary: $(BINARY)"
	@echo "  version: $(VERSION)"
	@echo "  go: $$($(GO) version | cut -d' ' -f3)"
	@echo "  installed: $$(which $(BINARY) 2>/dev/null || echo 'not found')"
	@if [ -f "$(BIN_DIR)/$(BINARY)" ]; then \
		echo "  size: $$(du -h $(BIN_DIR)/$(BINARY) | cut -f1)"; \
	fi

## help: Show this help message
help:
	@echo "Kavach - Brahmastra Stack"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "Examples:"
	@echo "  make build      # Build kavach binary"
	@echo "  make install    # Build + install + symlinks"
	@echo "  make release    # Cross-platform release builds"
	@echo "  make test       # Run all tests"
	@echo ""
