#!/bin/bash

# Build RecastNavigation libraries and install to project directories
set -e

echo "Building RecastNavigation Library"
echo "=================================="

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture to common names
case $ARCH in
    x86_64|amd64)
        ARCH_NAME="x64"
        ;;
    i386|i686)
        ARCH_NAME="x86"
        ;;
    aarch64|arm64)
        ARCH_NAME="arm64"
        ;;
    armv7l|armv7)
        ARCH_NAME="arm"
        ;;
    *)
        ARCH_NAME="$ARCH"
        ;;
esac

PLATFORM_DIR="${OS}_${ARCH_NAME}"

echo "Detected platform: $OS ($ARCH_NAME)"
echo "Target directory: lib/$PLATFORM_DIR/"

# Create build directory
mkdir -p third_party/recastnavigation/build
cd third_party/recastnavigation/build

# Configure with CMake - only build core libraries, disable everything else
echo "Configuring RecastNavigation build..."
cmake .. \
    -DDISABLE_ASSERTS=ON \
    -DCMAKE_BUILD_TYPE=Release \
    -DRECASTNAVIGATION_DEMO=OFF \
    -DRECASTNAVIGATION_TESTS=OFF \
    -DRECASTNAVIGATION_EXAMPLES=OFF \
    -DRECASTNAVIGATION_STATIC=ON

# Build only the specific targets we need
echo "Building RecastNavigation libraries..."
make -j$(nproc) Detour DetourCrowd DetourTileCache

# Create target directories
mkdir -p ../../../lib/$PLATFORM_DIR
mkdir -p ../../../include/recastnavigation
mkdir -p ../../../include/recastnavigation/detour
mkdir -p ../../../include/recastnavigation/detourcrowd
mkdir -p ../../../include/recastnavigation/detourtilecache

# Copy static libraries (without platform suffix in filename)
echo "Installing libraries to project lib/$PLATFORM_DIR/ directory..."
cp Detour/libDetour.a ../../../lib/$PLATFORM_DIR/
cp DetourCrowd/libDetourCrowd.a ../../../lib/$PLATFORM_DIR/
cp DetourTileCache/libDetourTileCache.a ../../../lib/$PLATFORM_DIR/

# Copy headers (headers are platform-independent)
echo "Installing headers to project include/recastnavigation/ directory..."
cp ../Detour/Include/*.h ../../../include/recastnavigation/detour
cp ../DetourCrowd/Include/*.h ../../../include/recastnavigation/detourcrowd
cp ../DetourTileCache/Include/*.h ../../../include/recastnavigation/detourtilecache

echo "RecastNavigation libraries built and installed successfully!"
echo ""
echo "Platform: $OS ($ARCH_NAME)"
echo "Libraries installed to: ./lib/$PLATFORM_DIR/"
echo "Headers installed to: ./include/recastnavigation/"
echo ""
echo "Library files:"
ls -la ../../../lib/$PLATFORM_DIR/
echo ""
echo "To use in Go code, you may need to adjust CGO LDFLAGS based on platform."