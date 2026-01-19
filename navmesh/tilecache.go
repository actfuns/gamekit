package navmesh

/*
#cgo CXXFLAGS: -std=c++11
#cgo linux CXXFLAGS: -I./third_party/recastnavigation/Detour/Include -I./third_party/recastnavigation/Recast/Include
#cgo linux LDFLAGS: -L./lib/linux_x64 -lRecast -lDetour -lDetourTileCache -lDetourCrowd
#include "navmesh.h"
*/
import "C"

// CreateTileCache creates a new tile cache for dynamic obstacles
func CreateTileCache(params TileCacheParams) (*TileCache, error) {
	cOrig := (*C.float)(&params.Orig[0])

	handle := C.create_tile_cache(
		cOrig,
		C.float(params.Cs), C.float(params.Ch),
		C.int(params.Width), C.int(params.Height),
		C.float(params.WalkableHeight), C.float(params.WalkableRadius), C.float(params.WalkableClimb),
		C.float(params.MaxSimplificationError), C.int(params.MaxTiles), C.int(params.MaxObstacles),
	)

	if handle == nil {
		return nil, ErrFailedToCreateTileCache
	}

	return &TileCache{handle: handle}, nil
}

// AddCylinderObstacle adds a cylindrical obstacle to the tile cache
func (tc *TileCache) AddCylinderObstacle(posX, posY, posZ, radius, height float32) (ObstacleRef, error) {
	var cRef C.uint

	result := C.add_cylinder_obstacle(
		tc.handle,
		C.float(posX), C.float(posY), C.float(posZ),
		C.float(radius), C.float(height),
		&cRef,
	)

	if result != 0 {
		return 0, ErrFailedToAddObstacle
	}

	return ObstacleRef(cRef), nil
}

// RemoveObstacle removes an obstacle from the tile cache
func (tc *TileCache) RemoveObstacle(ref ObstacleRef) error {
	result := C.remove_obstacle(tc.handle, C.uint(ref))
	if result != 0 {
		return ErrFailedToRemoveObstacle
	}
	return nil
}

// Update updates the tile cache and rebuilds affected tiles
func (tc *TileCache) Update(detour *NavMesh) error {
	result := C.update_tile_cache_with_navmesh(tc.handle, detour.handle)
	if result != 0 {
		return ErrFailedToUpdateTileCache
	}
	return nil
}

// Close cleans up the tile cache
func (tc *TileCache) Close() {
	if tc.handle != nil {
		C.destroy_tile_cache(tc.handle)
		tc.handle = nil
	}
}
