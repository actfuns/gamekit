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
│   └── recastnavigation/         # Header files (platform-independent)
│       ├── detour/               # Detour headers
│       │   ├── DetourNavMesh.h
│       │   ├── DetourNavMeshQuery.h
│       │   └── ...
│       ├── detourcrowd/          # DetourCrowd headers
│       │   └── DetourCrowd.h
│       └── detourtilecache/      # DetourTileCache headers
│           └── DetourTileCache.h
├── detour/                       # Go package with CGO bindings
│   ├── navmesh.go               # Core navigation mesh functions
│   ├── tilecache.go             # Dynamic obstacle management
│   ├── crowd.go                 # Multi-agent crowd simulation
│   ├── types.go                 # Common types and constants
│   ├── detour.h                 # C API bridge header
│   └── detour.cpp               # C++ implementation bridge
├── examples/                     # Usage examples
│   ├── findpath/                # Basic pathfinding example
│   ├── tilecache/               # Dynamic obstacle example
│   ├── crowd/                   # Crowd simulation example
│   └── dynamic_obstacles/       # Advanced dynamic obstacles
├── build_detour.sh              # Build RecastNavigation libraries
└── README.md                    # This documentation
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

The project uses pre-compiled static libraries stored in platform-specific directories under `lib/`:

```bash
# Build platform-specific libraries
./build_detour.sh
```

This script:
- Detects your platform (Linux/macOS/Windows + x86/x64/ARM)
- Compiles RecastNavigation core libraries only (no demos or tests)
- Installs libraries with standard names like `libDetour.a` in platform-specific directories
- Copies headers to `include/recastnavigation/` subdirectories

### 2. Build the Project

```bash
# Build all examples
go build ./...

# Or build specific components
go build ./detour                    # Build the Go package
cd examples/findpath && go build     # Build specific example
```

## Platform Support

The build system automatically detects and supports multiple platforms:

| Platform | Library Directory | Go Build Tags |
|----------|-------------------|---------------|
| Linux x64 | `linux_x64/` | `linux,amd64` |
| Linux x86 | `linux_x86/` | `linux,386` |
| macOS x64 | `darwin_x64/` | `darwin,amd64` |
| macOS ARM64 | `darwin_arm64/` | `darwin,arm64` |
| Windows x64 | `windows_x64/` | `windows,amd64` |

To add support for new platforms, extend the platform detection logic in `build_detour.sh` and add corresponding CGO LDFLAGS in the Go files.

## Usage Examples

### Basic Pathfinding
```go
package main

import (
    "fmt"
    "log"
    "os"
    "dr/detour"
)

func main() {
    // Load navigation mesh from file
    navMesh, err := detour.LoadNavMeshFromFile("detour.bin")
    if err != nil {
        log.Fatal(err)
    }
    defer navMesh.Close()
    
    // Create navigation mesh query
    query, err := navMesh.CreateQuery()
    if err != nil {
        log.Fatal(err)
    }
    defer query.Close()
    
    // Find path between two points
    path, err := query.FindStraightPath(0, 0, 0, 10, 0, 10)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found path with %d points\n", len(path))
    for i, point := range path {
        fmt.Printf("Point %d: (%.2f, %.2f, %.2f)\n", i, point[0], point[1], point[2])
    }
}
```

### Crowd Simulation
```go
// Create crowd manager
crowd, err := detour.CreateCrowd(navMesh, maxAgents, maxAgentRadius)
if err != nil {
    log.Fatal(err)
}
defer crowd.Close()

// Add agents and update simulation
agentID, err := crowd.AddAgent(positionX, positionY, positionZ, params)
if err != nil {
    log.Fatal(err)
}

crowd.Update(deltaTime)
```

## Directory Organization Benefits

- **Namespace Isolation**: `recastnavigation/` subdirectory prevents header conflicts
- **Platform Separation**: Each platform has its own library directory
- **Clean Integration**: Go code references libraries through standardized paths
- **Easy Maintenance**: Adding new platforms only requires updating build scripts

## Troubleshooting

### Common Issues

1. **"Cannot find library" errors**
   - Ensure you ran `./build_detour.sh` first
   - Verify library files exist in the correct platform directory under `lib/`
   - Check that your platform is supported in the Go CGO LDFLAGS

2. **Header file not found**
   - Confirm headers are in `include/recastnavigation/` subdirectories
   - Verify CGO CXXFLAGS include the correct include paths

3. **Linking errors**
   - Ensure all three libraries are present: Detour, DetourCrowd, DetourTileCache
   - Check that library linking order is correct in CGO LDFLAGS

### Adding New Platforms

To support additional platforms:

1. Update platform detection in `build_detour.sh`
2. Add corresponding CGO LDFLAGS in the Go files
3. Test compilation on the target platform

## License

- **This project**: MIT License
- **Recast Navigation**: zlib License

See individual license files for details.