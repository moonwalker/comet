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

# Construct download URL
# Adjust the filename pattern to match your release assets
# Example: comet_v1.0.0_darwin_amd64.tar.gz or comet_darwin_amd64.tar.gz
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"

echo "Downloading from: $DOWNLOAD_URL"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download and extract
curl -fsSL "$DOWNLOAD_URL" | tar -xz -C "$TMP_DIR"

# Determine installation directory
if [ -w "/usr/local/bin" ]; then
  INSTALL_DIR="/usr/local/bin"
elif [ -d "$HOME/.local/bin" ]; then
  INSTALL_DIR="$HOME/.local/bin"
else
  # Create ~/.local/bin if it doesn't exist
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
fi

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
