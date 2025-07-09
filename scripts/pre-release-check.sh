#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Pre-Release Validation ===${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Go is installed: $(go version)${NC}"

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo -e "${RED}❌ Git is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Git is installed${NC}"

# Check if we're in a git repository
if ! git rev-parse --is-inside-work-tree &> /dev/null; then
    echo -e "${RED}❌ Not in a git repository${NC}"
    exit 1
fi
echo -e "${GREEN}✅ In git repository${NC}"

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${YELLOW}⚠️  Working directory has uncommitted changes:${NC}"
    git status --porcelain
    echo ""
else
    echo -e "${GREEN}✅ Working directory is clean${NC}"
fi

# Check current branch
current_branch=$(git branch --show-current)
if [ "$current_branch" != "main" ]; then
    echo -e "${YELLOW}⚠️  Current branch is not main: $current_branch${NC}"
else
    echo -e "${GREEN}✅ On main branch${NC}"
fi

# Check if there are any recent commits
commit_count=$(git rev-list --count HEAD ^origin/main 2>/dev/null || echo "0")
if [ "$commit_count" -gt 0 ]; then
    echo -e "${YELLOW}⚠️  There are $commit_count unpushed commits${NC}"
else
    echo -e "${GREEN}✅ No unpushed commits${NC}"
fi

# Run tests
echo -e "${BLUE}Running tests...${NC}"
if go test -v ./...; then
    echo -e "${GREEN}✅ All tests passed${NC}"
else
    echo -e "${RED}❌ Tests failed${NC}"
    exit 1
fi

# Check if code builds
echo -e "${BLUE}Testing build...${NC}"
if go build -o /tmp/cngt-cli-test ./cmd/main.go; then
    echo -e "${GREEN}✅ Build successful${NC}"
    rm -f /tmp/cngt-cli-test
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Check current version
current_version=$(grep 'Version.*=' internal/version/version.go | sed 's/.*"\(.*\)".*/\1/')
echo -e "${BLUE}Current version: $current_version${NC}"

# Check if there are any existing tags
tags=$(git tag -l | wc -l)
if [ "$tags" -gt 0 ]; then
    echo -e "${BLUE}Recent tags:${NC}"
    git tag -l --sort=-version:refname | head -5
fi

echo ""
echo -e "${GREEN}=== Pre-Release Check Complete ===${NC}"
echo ""
echo -e "${YELLOW}Ready to release? Use:${NC}"
echo -e "${BLUE}  make release VERSION=X.Y.Z${NC}"
echo -e "${BLUE}  ./scripts/tag-release.sh X.Y.Z${NC}"
echo ""