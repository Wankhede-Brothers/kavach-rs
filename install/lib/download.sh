#!/bin/bash
# Binary download - Kavach/Brahmastra Stack
# DACE: 60 lines

REPO="Wankhede-Brothers/kavach-rs"
BINARY_NAME="kavach"

get_latest_version() {
    curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null | \
        grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/' || echo "v0.1.0"
}

download_binary() {
    VERSION=$(get_latest_version)

    # Map platform names to release artifact names
    case "$PLATFORM" in
        linux)  RELEASE_PLATFORM="linux" ;;
        darwin) RELEASE_PLATFORM="darwin" ;;
        *)      RELEASE_PLATFORM="$PLATFORM" ;;
    esac

    # Map architecture names
    case "$ARCH" in
        x86_64|amd64) RELEASE_ARCH="amd64" ;;
        aarch64|arm64) RELEASE_ARCH="arm64" ;;
        *)            RELEASE_ARCH="$ARCH" ;;
    esac

    BINARY_URL="https://github.com/$REPO/releases/download/$VERSION/${BINARY_NAME}-${RELEASE_PLATFORM}-${RELEASE_ARCH}"

    echo "[DOWNLOAD]"
    echo "  version: $VERSION"
    echo "  platform: $RELEASE_PLATFORM/$RELEASE_ARCH"
    echo "  url: $BINARY_URL"

    if command -v curl &>/dev/null; then
        curl -fsSL "$BINARY_URL" -o "$BIN_DIR/$BINARY_NAME" 2>/dev/null && {
            chmod +x "$BIN_DIR/$BINARY_NAME"
            echo "  status: ok"
            return 0
        }
    elif command -v wget &>/dev/null; then
        wget -q "$BINARY_URL" -O "$BIN_DIR/$BINARY_NAME" 2>/dev/null && {
            chmod +x "$BIN_DIR/$BINARY_NAME"
            echo "  status: ok"
            return 0
        }
    fi

    echo "  status: failed (will build from source)"
    return 1
}

build_from_source() {
    echo "[BUILD]"

    if ! command -v cargo &>/dev/null; then
        echo "  error: Rust not installed (install via https://rustup.rs)"
        return 1
    fi

    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    git clone --depth 1 "https://github.com/$REPO.git" .
    cd crates/kavach-cli
    cargo build --release
    cp target/release/"$BINARY_NAME" "$BIN_DIR/$BINARY_NAME"

    cd /
    rm -rf "$TEMP_DIR"

    echo "  status: ok (built from source)"
}

install_binary() {
    download_binary || build_from_source
}
