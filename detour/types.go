// package detour provides Go bindings for Recast Navigation library,
// enabling pathfinding, dynamic obstacle handling, and crowd simulation
// in Go applications.
//
// The package wraps the following Recast Navigation components:
//   - Detour: Navigation mesh loading and pathfinding
//   - DetourTileCache: Dynamic obstacle support
//   - DetourCrowd: Multi-agent crowd simulation
package detour

/*
#cgo CXXFLAGS: -std=c++11 -I${SRCDIR}/../include/recastnavigation/detour -I${SRCDIR}/../include/recastnavigation/detourcrowd -I${SRCDIR}/../include/recastnavigation/detourtilecache
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../lib/linux_x64 -lDetourTileCache -lDetourCrowd -lDetour -lstdc++
#cgo linux,386 LDFLAGS: -L${SRCDIR}/../lib/linux_x86 -lDetourTileCache -lDetourCrowd -lDetour -lstdc++
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../lib/darwin_x64 -lDetourTileCache -lDetourCrowd -lDetour -lstdc++
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../lib/darwin_arm64 -lDetourTileCache -lDetourCrowd -lDetour -lstdc++
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../lib/windows_x64 -lDetourTileCache -lDetourCrowd -lDetour -lstdc++

#include "detour.h"
#include <stdlib.h>
*/
import "C"

// NavMesh represents a navigation mesh
type NavMesh struct {
	handle C.NavMeshHandle
}

// NavMeshQuery represents a navigation mesh query
type NavMeshQuery struct {
	handle C.NavMeshQueryHandle
}

// TileCache represents a dynamic navigation mesh cache that supports obstacles
type TileCache struct {
	handle C.TileCacheHandle
}

// Crowd represents a crowd simulation for multiple agents
type Crowd struct {
	handle C.CrowdHandle
}

// ObstacleRef represents a reference to a dynamic obstacle
type ObstacleRef uint32

// TileCacheParams contains parameters for creating a tile cache
type TileCacheParams struct {
	Orig                   [3]float32
	Cs, Ch                 float32
	Width, Height          int32
	WalkableHeight         float32
	WalkableRadius         float32
	WalkableClimb          float32
	MaxSimplificationError float32
	MaxTiles               int32
	MaxObstacles           int32
}

// CrowdAgentParams contains parameters for crowd agents
type CrowdAgentParams struct {
	Radius                float32
	Height                float32
	MaxAcceleration       float32
	MaxSpeed              float32
	CollisionQueryRange   float32
	PathOptimizationRange float32
	SeparationWeight      float32
	UpdateFlags           uint8
	ObstacleAvoidanceType uint8
	QueryFilterType       uint8
}

// Error represents a detour error
type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

// Errors for navigation mesh operations
var (
	ErrFailedToLoadNavMesh = &Error{"failed to load navigation mesh"}
	ErrFailedToCreateQuery = &Error{"failed to create navigation mesh query"}
	ErrNoPathFound         = &Error{"no path found"}
)

// Errors for tile cache operations
var (
	ErrFailedToCreateTileCache = &Error{"failed to create tile cache"}
	ErrFailedToAddObstacle     = &Error{"failed to add obstacle"}
	ErrFailedToRemoveObstacle  = &Error{"failed to remove obstacle"}
	ErrFailedToUpdateTileCache = &Error{"failed to update tile cache"}
)

// Errors for crowd operations
var (
	ErrFailedToCreateCrowd       = &Error{"failed to create crowd"}
	ErrFailedToAddAgent          = &Error{"failed to add agent to crowd"}
	ErrFailedToRequestMoveTarget = &Error{"failed to request move target for agent"}
	ErrFailedToGetActiveAgents   = &Error{"failed to get active agents from crowd"}
	ErrInvalidMaxAgents          = &Error{"invalid max agents value"}
)
