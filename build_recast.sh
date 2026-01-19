#!/bin/bash

# Auto-fix CRLF line endings if present
if [ "$(head -n1 "$0" | cat -A | tail -c2)" = "^M$" ]; then
    # File has CRLF line endings, fix it
    sed -i 's/\r$//' "$0"
    echo "Fixed CRLF line endings in script, re-executing..."
    exec "$0" "$@"
fi

# Build RecastNavigation libraries and install to project directories
set -e

echo "Building RecastNavigation Library"
echo "=================================="

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

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

# Handle MSYS2 platform naming
if [[ "$OS" == mingw64* ]]; then
    PLATFORM_DIR="windows_${ARCH_NAME}"
else
    PLATFORM_DIR="${OS}_${ARCH_NAME}"
fi

echo "Detected platform: $OS ($ARCH_NAME)"
echo "Target directory: lib/$PLATFORM_DIR/"

# Function to get number of processors (fallback if nproc not available)
get_nproc() {
    if command -v nproc >/dev/null 2>&1; then
        nproc
    else
        # Fallback to 4 cores if nproc not available
        echo 4
    fi
}

# Clean up previous build
echo "Cleaning up previous build..."
rm -rf third_party/recastnavigation/build
rm -rf include/recastnavigation
rm -rf lib/$PLATFORM_DIR/detour*

# Create build directory
mkdir -p third_party/recastnavigation/build
cd third_party/recastnavigation/build

# Configure with CMake - only build core libraries, disable everything else
echo "Configuring RecastNavigation build..."
# Don't force generator on MSYS2, let CMake choose automatically
cmake .. \
    -DDISABLE_ASSERTS=ON \
    -DCMAKE_BUILD_TYPE=Release \
    -DRECASTNAVIGATION_DEMO=OFF \
    -DRECASTNAVIGATION_TESTS=OFF \
    -DRECASTNAVIGATION_EXAMPLES=OFF \
    -DRECASTNAVIGATION_STATIC=ON

# Detect which build system CMake generated
BUILD_CMD=""
if [ -f "build.ninja" ]; then
    # Ninja build system detected
    if command -v ninja >/dev/null 2>&1; then
        BUILD_CMD="ninja"
    elif command -v ninja-build >/dev/null 2>&1; then
        BUILD_CMD="ninja-build"
    else
        echo "Error: Ninja build system detected but ninja command not found!"
        echo "Please install ninja: pacman -S mingw-w64-x86_64-ninja"
        exit 1
    fi
elif [ -f "Makefile" ]; then
    # Make build system detected
    if command -v make >/dev/null 2>&1; then
        BUILD_CMD="make"
    elif command -v mingw32-make >/dev/null 2>&1; then
        BUILD_CMD="mingw32-make"
    else
        echo "Error: Make build system detected but make command not found!"
        echo "Please install make: pacman -S mingw-w64-x86_64-make"
        exit 1
    fi
else
    echo "Error: No recognized build system files found (neither build.ninja nor Makefile)!"
    echo "CMake configuration may have failed."
    exit 1
fi

# Build only the specific targets we need
echo "Building RecastNavigation libraries with $BUILD_CMD..."
if [ "$BUILD_CMD" = "ninja" ] || [ "$BUILD_CMD" = "ninja-build" ]; then
    # Ninja doesn't support -j option the same way, but uses parallel by default
    # We can use -j to specify explicit parallelism if needed
    $BUILD_CMD -j$(get_nproc) Detour DetourCrowd DetourTileCache
else
    # Make supports -j for parallel builds
    $BUILD_CMD -j$(get_nproc) Detour DetourCrowd DetourTileCache
fi

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