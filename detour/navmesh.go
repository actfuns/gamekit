package detour

import (
	"fmt"
	"unsafe"
)

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

	if result <= 0 {
		return nil, ErrFailedToFindPath
	}

	// Convert C array to Go slice
	path := make([][3]float32, result)
	ptr := uintptr(unsafe.Pointer(cPath))
	for i := 0; i < int(result); i++ {
		path[i][0] = *(*float32)(unsafe.Pointer(ptr))
		path[i][1] = *(*float32)(unsafe.Pointer(ptr + 4))
		path[i][2] = *(*float32)(unsafe.Pointer(ptr + 8))
		ptr += 12 // 3 * 4 bytes
	}

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

// BuildNavMeshFromObj builds a navigation mesh from an OBJ file and saves it to a binary file
func BuildNavMeshFromObj(objFilename, outputFilename string,
	cellSize, cellHeight float32,
	agentHeight, agentRadius, agentMaxClimb, agentMaxSlope float32,
	regionMinSize, regionMergeSize int,
	edgeMaxLen, edgeMaxError float32,
	vertsPerPoly int, detailSampleDist, detailSampleMaxError float32) error {
	
	cObjFilename := C.CString(objFilename)
	defer C.free(unsafe.Pointer(cObjFilename))
	
	cOutputFilename := C.CString(outputFilename)
	defer C.free(unsafe.Pointer(cOutputFilename))
	
	success := C.build_navmesh_from_obj(
		cObjFilename, cOutputFilename,
		C.float(cellSize), C.float(cellHeight),
		C.float(agentHeight), C.float(agentRadius), C.float(agentMaxClimb), C.float(agentMaxSlope),
		C.int(regionMinSize), C.int(regionMergeSize),
		C.float(edgeMaxLen), C.float(edgeMaxError),
		C.int(vertsPerPoly), C.float(detailSampleDist), C.float(detailSampleMaxError))
	
	if !success {
		return fmt.Errorf("failed to build navigation mesh from OBJ file")
	}
	
	return nil
}
