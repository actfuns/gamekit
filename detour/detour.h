#ifndef detour_H
#define detour_H

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// =============================================================================
// Navigation Mesh System C API Bridge
// 
// This header provides C-compatible interfaces for the Recast Navigation library,
// enabling Go code to interact with the following components:
// - Detour: Navigation mesh loading and pathfinding
// - DetourTileCache: Dynamic obstacle management  
// - DetourCrowd: Multi-agent crowd simulation
// =============================================================================

// Handle types for opaque pointers
typedef void* NavMeshHandle;
typedef void* NavMeshQueryHandle;
typedef void* TileCacheHandle;
typedef void* CrowdHandle;

// =============================================================================
// Navigation Mesh Functions
// =============================================================================

// Load navigation mesh from file
NavMeshHandle load_navmesh_from_file(const char* filename);

// Create navigation mesh query instance
NavMeshQueryHandle create_navmesh_query(NavMeshHandle navMesh);

// Find straight path between two points
int find_straight_path(NavMeshQueryHandle query, 
                       float start_x, float start_y, float start_z,
                       float end_x, float end_y, float end_z,
                       float* path, int max_points);

// Cleanup functions
void destroy_navmesh_query(NavMeshQueryHandle query);
void destroy_navmesh(NavMeshHandle navMesh);

// =============================================================================
// Tile Cache Functions (Dynamic Obstacles)
// =============================================================================

// Create tile cache for dynamic obstacle management
TileCacheHandle create_tile_cache(const float* orig, float cs, float ch, 
                                  int width, int height,
                                  float walkableHeight, float walkableRadius, float walkableClimb,
                                  float maxSimplificationError, int maxTiles, int maxObstacles);

// Add cylindrical obstacle to tile cache
int add_cylinder_obstacle(TileCacheHandle tileCache, 
                          float pos_x, float pos_y, float pos_z,
                          float radius, float height,
                          unsigned int* obstacleRef);

// Remove obstacle from tile cache
int remove_obstacle(TileCacheHandle tileCache, unsigned int obstacleRef);

// Update tile cache and rebuild affected navigation mesh tiles
int update_tile_cache_with_navmesh(TileCacheHandle tileCache, NavMeshHandle navMesh);

// Cleanup tile cache
void destroy_tile_cache(TileCacheHandle tileCache);

// =============================================================================
// Crowd Simulation Functions
// =============================================================================

// Create crowd simulation instance
CrowdHandle create_crowd(int maxAgents, float maxAgentRadius, NavMeshHandle navMesh);

// Cleanup crowd simulation
void destroy_crowd(CrowdHandle handle);

// Add agent to crowd simulation
int add_crowd_agent(CrowdHandle handle, float posX, float posY, float posZ, 
                    float radius, float height, float maxAcceleration, float maxSpeed,
                    float collisionQueryRange, float pathOptimizationRange, float separationWeight,
                    unsigned char updateFlags, unsigned char obstacleAvoidanceType, 
                    unsigned char queryFilterType);

// Remove agent from crowd simulation
void remove_crowd_agent(CrowdHandle handle, int agentIdx);

// Request agent to move to target position
bool request_crowd_agent_move_target(CrowdHandle handle, int agentIdx, 
                                    uint64_t polyRef, float targetX, float targetY, float targetZ);

// Update crowd simulation state
void update_crowd(CrowdHandle handle, float dt);

// Get active agents' positions and velocities
int get_crowd_active_agents(CrowdHandle handle, float* positions, float* velocities, int maxAgents);

#ifdef __cplusplus
}
#endif

#endif // detour_H