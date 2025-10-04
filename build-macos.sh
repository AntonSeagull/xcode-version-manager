#!/usr/bin/env bash
set -e

APP_NAME="xvm"
BUILD_DIR="build"

echo "🚀 Building ${APP_NAME} for macOS (arm64 + amd64)..."
mkdir -p "${BUILD_DIR}"

# macOS Apple Silicon (arm64)
echo "🔧 Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "${BUILD_DIR}/${APP_NAME}-darwin-arm64"

# macOS Intel (amd64)
echo "🔧 Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "${BUILD_DIR}/${APP_NAME}-darwin-amd64"

echo "✅ Build completed successfully!"
echo "📦 Binaries are in the '${BUILD_DIR}' directory:"
ls -lh "${BUILD_DIR}"