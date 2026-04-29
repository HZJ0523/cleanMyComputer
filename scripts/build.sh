#!/bin/bash
set -e

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build"
DIST_DIR="dist/cleanMyComputer-v${VERSION}"

echo "Building CleanMyComputer v${VERSION}..."

# Build
CGO_ENABLED=1 go build -ldflags "-s -w" \
  -o ${BUILD_DIR}/cleanMyComputer.exe \
  ./cmd/cleaner

echo "Build complete."

# Package for distribution
echo "Packaging..."
rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}/configs/rules"

cp ${BUILD_DIR}/cleanMyComputer.exe "${DIST_DIR}/"
cp configs/rules/level1_safe.json "${DIST_DIR}/configs/rules/"
cp configs/rules/level2_deep.json "${DIST_DIR}/configs/rules/"
cp configs/rules/level3_advanced.json "${DIST_DIR}/configs/rules/"

cd dist
powershell -Command "Compress-Archive -Path 'cleanMyComputer-v${VERSION}/*' -DestinationPath 'cleanMyComputer-v${VERSION}.zip' -Force"
cd ..

echo "Distribution package: dist/cleanMyComputer-v${VERSION}.zip"
