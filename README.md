# Recast Navigation Go Binding

This project provides Go bindings for the [Recast Navigation](https://github.com/recastnavigation/recastnavigation) library using CGO. It enables pathfinding, dynamic obstacle management, and crowd simulation capabilities in Go applications.

## Project Structure

```
dr/
├── lib/
│   ├── linux_x64/                # Linux x64 static libraries
│   ├── darwin_x64/               # macOS x64 static libraries  
│   ├── darwin_arm64/             # macOS ARM64 static libraries
│   └── windows_x64/              # Windows x64 static libraries
│       ├── libDetour.a
│       ├── libDetourCrowd.a  
│       └── libDetourTileCache.a
├── include/
│   └── recastnavigation/           # Header files (platform-independent)
│       ├── Detourdetour.h
│       ├── DetourCrowd.h
│       └── ...
├── detour/                       # Go package with CGO bindings
│   ├── detour.go                # Core navigation mesh functions
│   ├── tilecache.go              # Dynamic obstacle management
│   ├── crowd.go                  # Multi-agent crowd simulation
│   ├── types.go                  # Common types and constants
│   ├── detour.h          # C API bridge header
│   └── detour.cpp        # C++ implementation bridge
├── examples/                      # Usage examples
│   ├── findpath/                 # Basic pathfinding example
│   ├── tilecache/                # Dynamic obstacle example
│   ├── crowd/                    # Crowd simulation example
│   └── dynamic_obstacles/        # Advanced dynamic obstacles
├── build_recast.sh               # Build RecastNavigation libraries
├── build.sh                      # Build entire project
└── README.md                     # This documentation
```

## Features

- **Pathfinding**: Load navigation meshes and find paths between points
- **Dynamic Obstacles**: Add/remove obstacles at runtime using TileCache
- **Crowd Simulation**: Manage multiple agents with automatic pathfinding and collision avoidance
- **Cross-platform**: Supports Linux, macOS, and Windows with platform-specific libraries

## Prerequisites

### Build Dependencies
- **Go** 1.24.1 or later
- **CMake** 3.10 or later
- **GCC/Clang** with C++11 support
- **Make** (for building RecastNavigation)

### Runtime Dependencies
- None (statically linked libraries)

## Building

### 1. Build RecastNavigation Libraries

The project uses pre-compiled static libraries stored in `lib/recastnavigation/` with platform-specific naming:

```bash
# Build platform-specific libraries
./build_recast.sh
```

This script:
- Detects your platform (Linux/macOS/Windows + x86/x64/ARM)
- Compiles RecastNavigation core libraries only (no demos or tests)
- Installs libraries with names like `libDetour_linux_x64.a`
- Copies headers to `include/recastnavigation/`

### 2. Build the Project

```bash
# Build all examples
./build.sh

# Or build specific components
go build ./detour                    # Build the Go package
go build -o examples/findpath/findpath_demo examples/findpath/main.go
```

## Platform Support

The build system automatically detects and supports multiple platforms:

| Platform | Library Suffix | Go Build Tags |
|----------|----------------|---------------|
| Linux x64 | `_linux_x64` | `linux,amd64` |
| Linux x86 | `_linux_x86` | `linux,386` |
| macOS x64 | `_darwin_x64` | `darwin,amd64` |
| macOS ARM64 | `_darwin_arm64` | `darwin,arm64` |
| Windows x64 | `_windows_x64` | `windows,amd64` |

To add support for new platforms, extend the platform detection logic in `build_recast.sh` and add corresponding CGO LDFLAGS in `detour/detour.go`.

## Usage Examples

### Basic Pathfinding
```go
package main

import (
    "fmt"
    "dr/detour"
)

func main() {
    // Load navigation mesh
    detourHandle := detour.LoaddetourFromFile("mesh.nav")
    defer detour.Destroydetour(detourHandle)
    
    // Find path between two points
    start := detour.Vector3{X: 0, Y: 0, Z: 0}
    end := detour.Vector3{X: 10, Y: 0, Z: 10}
    
    path, success := detour.FindStraightPath(detourHandle, start, end)
    if success {
        fmt.Printf("Found path with %d points\n", len(path))
    }
}
```

### Crowd Simulation
```go
// Create crowd manager
crowd := detour.CreateCrowd(detourHandle, maxAgents, maxAgentRadius)
defer detour.DestroyCrowd(crowd)

// Add agents and update simulation
agentID := detour.AddAgent(crowd, position, params)
detour.UpdateCrowd(crowd, deltaTime)
```

## Directory Organization Benefits

- **Namespace Isolation**: `recastnavigation/` subdirectory prevents header conflicts
- **Platform Separation**: Each platform has its own library files
- **Clean Integration**: Go code references libraries through standardized paths
- **Easy Maintenance**: Adding new platforms only requires updating build scripts

## Troubleshooting

### Common Issues

1. **"Cannot find library" errors**
   - Ensure you ran `./build_recast.sh` first
   - Verify library files exist in `lib/recastnavigation/`
   - Check that your platform is supported in `detour.go`

2. **Header file not found**
   - Confirm headers are in `include/recastnavigation/`
   - Verify `detour.cpp` includes use correct paths

3. **Linking errors**
   - Ensure all three libraries are present: Detour, DetourCrowd, DetourTileCache
   - Check that library names match your platform suffix

### Adding New Platforms

To support additional platforms:

1. Update platform detection in `build_recast.sh`
2. Add corresponding CGO LDFLAGS in `detour/detour.go`
3. Test compilation on the target platform

## License

- **This project**: MIT License
- **Recast Navigation**: zlib License

See individual license files for details.