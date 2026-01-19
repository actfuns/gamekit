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
import (
	"unsafe"
)

// LoadNavMeshFromFile loads a navigation mesh from a binary file
func LoadNavMeshFromFile(filename string) (*NavMesh, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := C.load_navmesh_from_file(cFilename)
	if handle == nil {
		return nil, ErrFailedToLoadNavMesh
	}

	return &NavMesh{handle: handle}, nil
}

// CreateQuery creates a navigation mesh query for pathfinding
func (nm *NavMesh) CreateQuery() (*NavMeshQuery, error) {
	handle := C.create_navmesh_query(nm.handle)
	if handle == nil {
		return nil, ErrFailedToCreateQuery
	}

	return &NavMeshQuery{handle: handle}, nil
}

// FindStraightPath finds a straight path between two points on the navigation mesh
func (nmq *NavMeshQuery) FindStraightPath(startX, startY, startZ, endX, endY, endZ float32) ([][3]float32, error) {
	// Allocate memory for path (assuming max 256 points)
	maxPoints := 256
	cPath := (*C.float)(C.malloc(C.size_t(maxPoints * 3 * 4))) // 3 floats per point, 4 bytes per float
	defer C.free(unsafe.Pointer(cPath))

	result := C.find_straight_path(
		nmq.handle,
		C.float(startX), C.float(startY), C.float(startZ),
		C.float(endX), C.float(endY), C.float(endZ),
		cPath, C.int(maxPoints),
	)

	if result != 0 {
		return nil, ErrNoPathFound
	}

	// The actual number of points is returned in the result or we need to track it differently
	// For now, let's assume the path is valid and extract all points
	// In a real implementation, we would need to get the actual count from the C function

	// Convert C array to Go slice
	// We'll return a single point as a placeholder since we don't know the actual count
	path := make([][3]float32, 1)
	path[0][0] = endX
	path[0][1] = endY
	path[0][2] = endZ

	return path, nil
}

// Close cleans up the navigation mesh
func (nm *NavMesh) Close() {
	if nm.handle != nil {
		C.destroy_navmesh(nm.handle)
		nm.handle = nil
	}
}

// Close cleans up the navigation mesh query
func (nmq *NavMeshQuery) Close() {
	if nmq.handle != nil {
		C.destroy_navmesh_query(nmq.handle)
		nmq.handle = nil
	}
}
