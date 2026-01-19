package navmesh

/*
#cgo CXXFLAGS: -std=c++11
#cgo linux CXXFLAGS: -I./third_party/recastnavigation/Detour/Include -I./third_party/recastnavigation/Recast/Include
#cgo linux LDFLAGS: -L./lib/linux_x64 -lRecast -lDetour -lDetourTileCache -lDetourCrowd
#include "navmesh.h"
*/
import "C"
import (
	"unsafe"
)

// CreateCrowd creates a new crowd simulation for multiple agents
func CreateCrowd(maxAgents int, maxAgentRadius float32, detour *NavMesh) (*Crowd, error) {
	if maxAgents <= 0 {
		return nil, ErrInvalidMaxAgents
	}

	handle := C.create_crowd(
		C.int(maxAgents),
		C.float(maxAgentRadius),
		detour.handle,
	)

	if handle == nil {
		return nil, ErrFailedToCreateCrowd
	}

	return &Crowd{handle: handle}, nil
}

// AddAgent adds an agent to the crowd with specified parameters and position
func (crowd *Crowd) AddAgent(params CrowdAgentParams, posX, posY, posZ float32) (int, error) {
	result := C.add_crowd_agent(
		crowd.handle,
		C.float(posX), C.float(posY), C.float(posZ),
		C.float(params.Radius), C.float(params.Height),
		C.float(params.MaxAcceleration), C.float(params.MaxSpeed),
		C.float(params.CollisionQueryRange), C.float(params.PathOptimizationRange),
		C.float(params.SeparationWeight),
		C.uint8_t(params.UpdateFlags), C.uint8_t(params.ObstacleAvoidanceType), C.uint8_t(params.QueryFilterType),
	)

	if result < 0 {
		return -1, ErrFailedToAddAgent
	}

	return int(result), nil
}

// RemoveAgent removes an agent from the crowd
func (crowd *Crowd) RemoveAgent(agentIdx int) {
	C.remove_crowd_agent(crowd.handle, C.int(agentIdx))
}

// RequestMoveTarget requests an agent to move to a target position on the navigation mesh
func (crowd *Crowd) RequestMoveTarget(agentIdx int, endX, endY, endZ float32) error {
	result := C.request_crowd_agent_move_target(crowd.handle, C.int(agentIdx), 0, C.float(endX), C.float(endY), C.float(endZ))
	if result == false {
		return ErrFailedToRequestMoveTarget
	}

	return nil
}

// Update updates the crowd simulation
// This function has been removed as it's no longer needed with the new get_active_agents implementation

// GetActiveAgents retrieves information about all active agents in the crowd
func (crowd *Crowd) GetActiveAgents() ([][3]float32, error) {
	// We need to allocate memory for positions
	// For simplicity, let's assume a maximum of 1000 agents
	maxAgents := 1000
	cPositions := (*C.float)(C.malloc(C.size_t(maxAgents * 3 * 4))) // 3 floats per agent, 4 bytes per float
	defer C.free(unsafe.Pointer(cPositions))

	cCount := C.get_crowd_active_agents(crowd.handle, cPositions, nil, C.int(maxAgents))
	if cCount < 0 {
		return nil, ErrFailedToGetActiveAgents
	}

	if cCount == 0 {
		return [][3]float32{}, nil
	}

	// Convert C array to Go slice
	positions := make([][3]float32, int(cCount))
	for i := 0; i < int(cCount); i++ {
		positions[i][0] = float32(*(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(cPositions)) + uintptr(i*3*4))))
		positions[i][1] = float32(*(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(cPositions)) + uintptr(i*3*4+4))))
		positions[i][2] = float32(*(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(cPositions)) + uintptr(i*3*4+8))))
	}

	return positions, nil
}

// Update updates the crowd simulation with the given time step
func (crowd *Crowd) Update(dt float32) {
	C.update_crowd(crowd.handle, C.float(dt))
}

// Close cleans up the crowd
func (crowd *Crowd) Close() {
	if crowd.handle != nil {
		C.destroy_crowd(crowd.handle)
		crowd.handle = nil
	}
}
