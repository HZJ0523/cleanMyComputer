#!/bin/bash
set -e

echo "Packaging CleanMyComputer..."

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build"
DIST_DIR="dist"
APP_NAME="cleanMyComputer"

mkdir -p $DIST_DIR

./scripts/build.sh

cp $BUILD_DIR/${APP_NAME}.exe $DIST_DIR/
cp -r configs $DIST_DIR/
cp README.md $DIST_DIR/

echo "Package complete: $DIST_DIR/${APP_NAME}-${VERSION}-windows-amd64"

ls -la $DIST_DIR/
