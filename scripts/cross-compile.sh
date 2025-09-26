#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

echo "🚀 Starting K2Ray cross-compilation process..."

# Get the root directory of the project
ROOT_DIR=$(git rev-parse --show-toplevel)
RELEASE_DIR="$ROOT_DIR/dist/release"

# Target platforms (OS/Architecture)
PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64"

# Clean up previous release builds
echo "🧹 Cleaning up old release artifacts..."
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

# Get project version from git tag or use a default
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0-dev")
echo "📦 Building version: $VERSION"

# Loop through platforms and build
for PLATFORM in $PLATFORMS; do
    # Split the platform string into OS and ARCH
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}

    OUTPUT_NAME="k2ray-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi

    echo "🛠️ Building for $GOOS/$GOARCH..."

    # Set environment variables and build
    env GOOS=$GOOS GOARCH=$GOARCH go build -v -o "$RELEASE_DIR/$OUTPUT_NAME" -ldflags="-s -w -X main.Version=$VERSION" "$ROOT_DIR/cmd/k2ray"

    # Create an archive for the binary
    echo "🗜️ Compressing $OUTPUT_NAME..."
    if [ "$GOOS" = "windows" ]; then
        zip -j "$RELEASE_DIR/k2ray-${VERSION}-${GOOS}-${GOARCH}.zip" "$RELEASE_DIR/$OUTPUT_NAME"
    else
        tar -czf "$RELEASE_DIR/k2ray-${VERSION}-${GOOS}-${GOARCH}.tar.gz" -C "$RELEASE_DIR" "$OUTPUT_NAME"
    fi

    # Clean up the raw binary
    rm "$RELEASE_DIR/$OUTPUT_NAME"
done

echo "✅ Cross-compilation successful! Artifacts are in the $RELEASE_DIR directory."