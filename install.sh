#!/usr/bin/env bash
set -e

REPO="f24aalam/dbmcp"
VERSION="v0.1.0"
INSTALL_DIR="/usr/local/bin"

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

curl -fsSL "$URL" -o /tmp/godbmcp
chmod +x /tmp/godbmcp
sudo mv /tmp/godbmcp "$INSTALL_DIR/godbmcp"

echo "✅ Installed successfully"
echo "Run: godbmcp --help"
