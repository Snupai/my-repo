#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if version is provided
if [ -z "$1" ]; then
    echo -e "${RED}Usage: $0 <version>${NC}"
    echo -e "${YELLOW}Example: $0 1.0.1${NC}"
    echo ""
    echo "This script will:"
    echo "1. Validate the version format"
    echo "2. Update version in source code"
    echo "3. Create a git commit and tag"
    echo "4. Push to GitHub (triggers automated release)"
    exit 1
fi

VERSION=$1

# Validate version format (semantic versioning)
if [[ ! $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}Error: Version must be in format X.Y.Z (e.g., 1.0.1)${NC}"
    exit 1
fi

# Check if tag already exists
if git tag -l | grep -q "^v$VERSION$"; then
    echo -e "${RED}Error: Tag v$VERSION already exists${NC}"
    exit 1
fi

# Check if we're on main branch
current_branch=$(git branch --show-current)
if [ "$current_branch" != "main" ]; then
    echo -e "${YELLOW}Warning: You're not on the main branch. Current branch: $current_branch${NC}"
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}Error: Working directory is not clean. Please commit or stash changes.${NC}"
    exit 1
fi

# Check if tests pass
echo -e "${BLUE}Running tests...${NC}"
if ! go test -v ./...; then
    echo -e "${RED}Tests failed. Please fix them before releasing.${NC}"
    exit 1
fi

# Show current version
current_version=$(grep 'Version.*=' internal/version/version.go | sed 's/.*"\(.*\)".*/\1/')
echo -e "${BLUE}Current version: $current_version${NC}"
echo -e "${BLUE}New version: $VERSION${NC}"

# Confirm release
read -p "Are you sure you want to release version $VERSION? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Release cancelled.${NC}"
    exit 1
fi

# Update version in version.go
echo -e "${BLUE}Updating version in source code...${NC}"
sed -i.bak "s/Version   = \".*\"/Version   = \"$VERSION\"/" internal/version/version.go
rm -f internal/version/version.go.bak

# Commit version update
echo -e "${BLUE}Creating commit and tag...${NC}"
git add internal/version/version.go
git commit -m "chore: bump version to $VERSION"

# Create and push tag
git tag -a "v$VERSION" -m "Release version $VERSION"

# Push to GitHub
echo -e "${BLUE}Pushing to GitHub...${NC}"
git push origin main
git push origin "v$VERSION"

echo ""
echo -e "${GREEN}✅ Successfully created and pushed tag v$VERSION${NC}"
echo -e "${GREEN}🚀 GitHub Actions will now build and create the release automatically${NC}"
echo -e "${GREEN}📦 Release will be available at: https://github.com/snupai/cngt-cli/releases/tag/v$VERSION${NC}"
echo ""
echo -e "${YELLOW}You can monitor the release process at:${NC}"
echo -e "${BLUE}https://github.com/snupai/cngt-cli/actions${NC}"