#!/bin/bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building DR Navigation Mesh Project${NC}"
echo "========================================"

# 检查是否在正确的目录
if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: This script must be run from the project root directory${NC}"
    exit 1
fi

# 构建主detour包
echo -e "${YELLOW}Building detour package...${NC}"
go build -v ./detour
echo -e "${GREEN}✓ detour package built successfully${NC}"

echo ""

# 构建所有示例
EXAMPLES=("findpath" "tilecache" "dynamic_obstacles" "crowd")

for example in "${EXAMPLES[@]}"; do
    if [ -d "examples/$example" ]; then
        echo -e "${YELLOW}Building $example example...${NC}"
        go build -o "examples/$example/${example}_demo" "examples/$example/main.go"
        echo -e "${GREEN}✓ $example example built successfully${NC}"
        echo ""
    else
        echo -e "${RED}Warning: examples/$example directory not found${NC}"
    fi
done

echo -e "${GREEN}All builds completed successfully!${NC}"
echo ""
echo "To run examples:"
echo "  ./examples/findpath/findpath_demo"
echo "  ./examples/tilecache/tilecache_demo"
echo "  ./examples/dynamic_obstacles/dynamic_obstacles_demo"
echo "  ./examples/crowd/crowd_demo"