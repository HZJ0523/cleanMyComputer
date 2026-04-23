#!/bin/bash
set -e

echo "Building CleanMyComputer..."

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build"
OUTPUT_NAME="cleanMyComputer"

mkdir -p $BUILD_DIR

CGO_ENABLED=1 go build -ldflags "-s -w -X main.version=$VERSION" \
  -o $BUILD_DIR/${OUTPUT_NAME}.exe \
  ./cmd/cleaner

echo "Build complete: $BUILD_DIR/${OUTPUT_NAME}.exe"
