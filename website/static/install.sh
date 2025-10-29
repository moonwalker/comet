#!/bin/sh
set -e

# Comet installer script
# Usage: curl -fsSL https://moonwalker.github.io/comet/install.sh | sh

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Darwin)
    OS="darwin"
    ;;
  Linux)
    OS="linux"
    ;;
  *)
    echo "Unsupported operating system: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64|arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# GitHub repository
REPO="moonwalker/comet"
BINARY_NAME="comet"

# Get latest release tag
echo "Fetching latest release..."
LATEST_RELEASE=$(curl -sL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
  echo "Error: Could not fetch latest release"
  exit 1
fi

echo "Latest release: $LATEST_RELEASE"

# Extract version without 'v' prefix for asset filename
VERSION="${LATEST_RELEASE#v}"

# Construct download URL
# GoReleaser creates archives like: comet_0.6.7_darwin_arm64.tar.gz
# (version without 'v' prefix, even though tag is v0.6.7)
# Windows uses .tar.gz for GoReleaser v2
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"

echo "Downloading from: $DOWNLOAD_URL"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download and extract
curl -fsSL "$DOWNLOAD_URL" | tar -xz -C "$TMP_DIR"

# Determine installation directory
# Prefer ~/.local/bin (modern standard, no sudo needed)
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

echo "Installing to: $INSTALL_DIR"

# Install binary
mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo ""
echo "✓ Comet installed successfully!"
echo ""
echo "Installation location: $INSTALL_DIR/$BINARY_NAME"
echo ""

# Check if installation directory is in PATH
case ":$PATH:" in
  *":$INSTALL_DIR:"*)
    # Already in PATH
    echo "Run 'comet version' to verify the installation."
    ;;
  *)
    # Not in PATH
    echo "⚠️  Note: $INSTALL_DIR is not in your PATH"
    echo ""
    echo "To add it to your PATH, add this line to your shell profile:"
    echo ""
    if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
      echo "  # /usr/local/bin should already be in PATH"
    else
      echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
    echo ""
    echo "Then run: source ~/.bashrc (or ~/.zshrc, ~/.profile, etc.)"
    echo ""
    echo "Or run comet directly: $INSTALL_DIR/comet version"
    ;;
esac
