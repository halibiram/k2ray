#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

echo "🚀 Starting K2Ray build process..."

# Get the root directory of the project
ROOT_DIR=$(git rev-parse --show-toplevel)
DIST_DIR="$ROOT_DIR/dist"

# Clean up previous builds
echo "🧹 Cleaning up old build artifacts..."
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

# Build Go backend
echo "🛠️ Building Go backend..."
go build -v -o "$DIST_DIR/k2ray" "$ROOT_DIR/cmd/k2ray"

# Build frontend
echo "🎨 Building frontend..."
cd "$ROOT_DIR/web"
npm install
npm run build
cd "$ROOT_DIR"

# Move frontend build to dist
echo "🚚 Moving frontend assets to dist..."
mv "$ROOT_DIR/web/dist" "$DIST_DIR/web"

# Copy configuration files
echo "⚙️ Copying configuration files..."
cp -r "$ROOT_DIR/configs" "$DIST_DIR/"

echo "✅ Build successful! Output is in the $DIST_DIR directory."