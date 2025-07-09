#!/bin/bash

set -e

VERSION=${1:-"1.0.0"}
OUTPUT_DIR="dist"

echo "Building cngt-cli version $VERSION"

# Clean previous builds
rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR

# Build for different platforms
PLATFORMS=(
    "windows/amd64"
    "linux/amd64"
    "darwin/amd64"
    "linux/arm64"
    "darwin/arm64"
)

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    echo "Building for $GOOS/$GOARCH..."
    
    output_name="cngt-cli-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-X github.com/snupai/cngt-cli/internal/version.Version=$VERSION -X github.com/snupai/cngt-cli/internal/version.GitCommit=$(git rev-parse HEAD) -X github.com/snupai/cngt-cli/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -s -w" \
        -o "$OUTPUT_DIR/$output_name" \
        ./cmd/main.go
    
    echo "Built: $OUTPUT_DIR/$output_name"
done

echo "Build completed successfully!"
echo "Binaries available in $OUTPUT_DIR/"