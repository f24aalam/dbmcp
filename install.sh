#!/usr/bin/env bash
set -e

REPO="f24aalam/dbmcp"
VERSION="v0.1.0"
INSTALL_DIR="${HOME}/.local/bin"
FALLBACK_DIR="/usr/local/bin"

OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux|darwin) ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

BINARY="godbmcp_${OS}_${ARCH}"
URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY"

echo "📦 Installing godbmcp ($OS/$ARCH)"
echo "⬇️  $URL"

mkdir -p "$INSTALL_DIR"
curl -fsSL "$URL" -o "${INSTALL_DIR}/godbmcp"
chmod +x "${INSTALL_DIR}/godbmcp"

if [ -w "$FALLBACK_DIR" ] 2>/dev/null; then
    mv "${INSTALL_DIR}/godbmcp" "${FALLBACK_DIR}/godbmcp"
    echo "✅ Installed to ${FALLBACK_DIR}/godbmcp"
else
    echo "✅ Installed to ${INSTALL_DIR}/godbmcp"
    echo "Add to PATH: export PATH=\"\$PATH:${INSTALL_DIR}\""
fi

echo "Run: godbmcp --help"
