#!/usr/bin/env bash
set -e

REPO="f24aalam/dbmcp"
VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep -o '"tag_name": *"[^"]*"' | cut -d'"' -f4)
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

BINARY="dbmcp_${OS}_${ARCH}"
URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY"

echo "📦 Installing dbmcp ($OS/$ARCH)"
echo "⬇️  $URL"

mkdir -p "$INSTALL_DIR"
echo "Downloading..."
curl -#L "$URL" -o "${INSTALL_DIR}/dbmcp"
chmod +x "${INSTALL_DIR}/dbmcp"

if [ -w "$FALLBACK_DIR" ] 2>/dev/null; then
    mv "${INSTALL_DIR}/dbmcp" "${FALLBACK_DIR}/dbmcp"
    echo "✅ Installed to ${FALLBACK_DIR}/dbmcp"
else
    echo "✅ Installed to ${INSTALL_DIR}/dbmcp"
    echo "Add to PATH: export PATH=\"\$PATH:${INSTALL_DIR}\""
fi

echo "Run: dbmcp --help"
