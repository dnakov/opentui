#!/bin/bash
set -e

# Download header and libraries from GitHub releases
echo "Downloading OpenTUI assets..."

REPO="sst/opentui"
RELEASE_URL="https://github.com/$REPO/releases/latest/download"

# Create directories
mkdir -p lib/aarch64-linux lib/aarch64-macos lib/aarch64-windows lib/x86_64-linux lib/x86_64-macos lib/x86_64-windows

# Download header
curl -L -o opentui.h "$RELEASE_URL/opentui.h"

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="x86_64" ;;
    arm64|aarch64) ARCH="aarch64" ;;
esac

case "$OS" in
    darwin) OS="macos" ;;
    linux) OS="linux" ;;
    mingw*|cygwin*|msys*) OS="windows" ;;
esac

PLATFORM="${ARCH}-${OS}"

# Download appropriate library
case "$OS" in
    macos)
        curl -L -o "lib/$PLATFORM/libopentui.dylib" "$RELEASE_URL/libopentui-$PLATFORM.dylib"
        ;;
    linux)
        curl -L -o "lib/$PLATFORM/libopentui.so" "$RELEASE_URL/libopentui-$PLATFORM.so"
        ;;
    windows)
        curl -L -o "lib/$PLATFORM/opentui.dll" "$RELEASE_URL/opentui-$PLATFORM.dll"
        ;;
esac

echo "OpenTUI assets downloaded successfully for $PLATFORM"