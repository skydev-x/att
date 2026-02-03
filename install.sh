#!/bin/bash
# Installation script - builds from source
set -e

# Configuration - UPDATE THESE
APP_NAME="at"           # What users will type to run it
GITHUB_USER="skydev-x"   # Your GitHub username
GITHUB_REPO="at"       # Your repo name
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Installing $APP_NAME...${NC}\n"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓ Go found: $(go version)${NC}"

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: git is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Git found${NC}\n"

# Clone or download the repo
TMP_DIR=$(mktemp -d)
echo -e "${YELLOW}Downloading source code...${NC}"

git clone --depth 1 "https://github.com/$GITHUB_USER/$GITHUB_REPO.git" "$TMP_DIR" 2>/dev/null || {
    echo -e "${RED}Failed to clone repository${NC}"
    exit 1
}

echo -e "${GREEN}✓ Source downloaded${NC}\n"

# Build the binary
echo -e "${YELLOW}Building $APP_NAME...${NC}"
cd "$TMP_DIR"

go build -o "$APP_NAME" -ldflags="-s -w" . || {
    echo -e "${RED}Build failed${NC}"
    exit 1
}

echo -e "${GREEN}✓ Build successful${NC}\n"

# Install the binary
echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"

if [ -w "$INSTALL_DIR" ]; then
    mv "$APP_NAME" "$INSTALL_DIR/$APP_NAME"
else
    sudo mv "$APP_NAME" "$INSTALL_DIR/$APP_NAME"
fi

# Cleanup
cd - > /dev/null
rm -rf "$TMP_DIR"

# Verify installation
if command -v "$APP_NAME" &> /dev/null; then
    echo -e "${GREEN}✓ $APP_NAME installed successfully!${NC}\n"
    echo -e "Run ${GREEN}'$APP_NAME'${NC} to start\n"
else
    echo -e "${RED}Installation verification failed${NC}"
    echo "The binary was moved to $INSTALL_DIR but is not in PATH"
    exit 1
fi
