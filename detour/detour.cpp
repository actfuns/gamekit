/**
 * @file detour.cpp
 * @brief Navigation Mesh C++ Bridge Implementation
 * 
 * This file provides the C++ implementation for the Go CGO bridge interface
 * that wraps the Recast Navigation library components:
 * - Detour (Navigation Mesh)
 * - TileCache (Dynamic Obstacle Management)  
 * - Crowd (Agent Simulation)
 * 
 * The bridge exposes C-compatible functions that can be called from Go code
 * through CGO, providing a clean interface to the underlying C++ functionality.
 */

#include "detour.h"
#include <stdio.h>
#include <stdlib.h>
#include <cstring>

#include "DetourNavMesh.h"
#include "DetourNavMeshQuery.h"
#include "DetourAlloc.h"
#include "DetourCommon.h"
#include "DetourTileCache.h"
#include "DetourCrowd.h"

struct NavMeshWrapper {
    dtNavMesh* navMesh;
};

struct NavMeshQueryWrapper {
    dtNavMeshQuery* query;
};

struct TileCacheWrapper {
    dtTileCache* tileCache;
};

struct CrowdWrapper {
    dtCrowd* crowd;
    
    CrowdWrapper() : crowd(nullptr) {}
    ~CrowdWrapper() {
        if (crowd) {
            dtFreeCrowd(crowd);
        }
    }
};

#ifdef __cplusplus
extern "C" {
#endif

NavMeshHandle load_navmesh_from_file(const char* filename) {
    if (!filename) return nullptr;
    
    // Open file
    FILE* file = fopen(filename, "rb");
    if (!file) {
        printf("Error: Cannot open file %s\n", filename);
        return nullptr;
    }
    
    // Get file size
    fseek(file, 0, SEEK_END);
    long fileSize = ftell(file);
    fseek(file, 0, SEEK_SET);
    
    printf("Loading navmesh from %s, file size: %ld bytes\n", filename, fileSize);
    
    if (fileSize <= 0) {
        fclose(file);
        printf("Error: File size is %ld, expected positive size\n", fileSize);
        return nullptr;
    }
    
    // Read file content
    unsigned char* data = new unsigned char[fileSize];
    size_t bytesRead = fread(data, 1, fileSize, file);
    fclose(file);
    
    if (bytesRead != (size_t)fileSize) {
        delete[] data;
        printf("Error: Read %zu bytes, expected %ld bytes\n", bytesRead, fileSize);
        return nullptr;
    }
    
    NavMeshWrapper* wrapper = new NavMeshWrapper();
    wrapper->navMesh = dtAllocNavMesh();
    if (!wrapper->navMesh) {
        delete[] data;
        delete wrapper;
        printf("Error: Failed to allocate dtNavMesh\n");
        return nullptr;
    }
    
    dtStatus status = wrapper->navMesh->init(data, (int)fileSize, DT_TILE_FREE_DATA);
    if (dtStatusFailed(status)) {
        dtFreeNavMesh(wrapper->navMesh);
        delete[] data;
        delete wrapper;
        printf("Error: dtNavMesh::init failed with status: %d\n", status);
        return nullptr;
    }
    
    delete[] data;
    
    printf("Successfully loaded navmesh from %s\n", filename);
    return static_cast<NavMeshHandle>(wrapper);
}

TileCacheHandle create_tile_cache(const float* orig, float cs, float ch, 
                                  int width, int height,
                                  float walkableHeight, float walkableRadius, float walkableClimb,
                                  float maxSimplificationError, int maxTiles, int maxObstacles) {
    if (!orig) return nullptr;
    
    dtTileCacheParams params;
    memcpy(params.orig, orig, sizeof(float)*3);
    params.cs = cs;
    params.ch = ch;
    params.width = width;
    params.height = height;
    params.walkableHeight = walkableHeight;
    params.walkableRadius = walkableRadius;
    params.walkableClimb = walkableClimb;
    params.maxSimplificationError = maxSimplificationError;
    params.maxTiles = maxTiles;
    params.maxObstacles = maxObstacles;
    
    TileCacheWrapper* wrapper = new TileCacheWrapper();
    wrapper->tileCache = dtAllocTileCache();
    if (!wrapper->tileCache) {
        delete wrapper;
        return nullptr;
    }
    
    if (dtStatusFailed(wrapper->tileCache->init(&params, nullptr, nullptr, nullptr))) {
        dtFreeTileCache(wrapper->tileCache);
        delete wrapper;
        return nullptr;
    }
    
    return static_cast<TileCacheHandle>(wrapper);
}

int add_cylinder_obstacle(TileCacheHandle tileCache, 
                          float pos_x, float pos_y, float pos_z,
                          float radius, float height,
                          unsigned int* obstacleRef) {
    if (!tileCache || !obstacleRef) return -1;
    
    TileCacheWrapper* wrapper = static_cast<TileCacheWrapper*>(tileCache);
    float pos[3] = {pos_x, pos_y, pos_z};
    dtObstacleRef ref;
    
    dtStatus status = wrapper->tileCache->addObstacle(pos, radius, height, &ref);
    if (dtStatusFailed(status)) {
        return -1;
    }
    
    *obstacleRef = static_cast<unsigned int>(ref);
    return 0;
}

int remove_obstacle(TileCacheHandle tileCache, unsigned int obstacleRef) {
    if (!tileCache) return -1;
    
    TileCacheWrapper* wrapper = static_cast<TileCacheWrapper*>(tileCache);
    dtStatus status = wrapper->tileCache->removeObstacle(static_cast<dtObstacleRef>(obstacleRef));
    if (dtStatusFailed(status)) {
        return -1;
    }
    
    return 0;
}

int update_tile_cache_with_navmesh(TileCacheHandle tileCache, NavMeshHandle navMesh) {
    if (!tileCache || !navMesh) return -1;
    
    TileCacheWrapper* tcWrapper = static_cast<TileCacheWrapper*>(tileCache);
    NavMeshWrapper* nmWrapper = static_cast<NavMeshWrapper*>(navMesh);
    bool upToDate = false;
    
    dtStatus status = tcWrapper->tileCache->update(0.0f, nmWrapper->navMesh, &upToDate);
    if (dtStatusFailed(status)) {
        return -1;
    }
    
    return upToDate ? 1 : 0; // 1 indicates update completed, 0 indicates still processing
}

void destroy_tile_cache(TileCacheHandle tileCache) {
    if (!tileCache) return;
    
    TileCacheWrapper* wrapper = static_cast<TileCacheWrapper*>(tileCache);
    if (wrapper->tileCache) {
        dtFreeTileCache(wrapper->tileCache);
    }
    delete wrapper;
}

NavMeshQueryHandle create_navmesh_query(NavMeshHandle navMesh) {
    if (!navMesh) return nullptr;
    
    NavMeshWrapper* wrapper = static_cast<NavMeshWrapper*>(navMesh);
    NavMeshQueryWrapper* queryWrapper = new NavMeshQueryWrapper();
    
    queryWrapper->query = dtAllocNavMeshQuery();
    if (!queryWrapper->query) {
        delete queryWrapper;
        return nullptr;
    }
    
    if (dtStatusFailed(queryWrapper->query->init(wrapper->navMesh, 2048))) {
        dtFreeNavMeshQuery(queryWrapper->query);
        delete queryWrapper;
        return nullptr;
    }
    
    return static_cast<NavMeshQueryHandle>(queryWrapper);
}

int find_straight_path(NavMeshQueryHandle query, 
                       float start_x, float start_y, float start_z,
                       float end_x, float end_y, float end_z,
                       float* path, int max_points) {
    if (!query || !path) return 0;
    
    NavMeshQueryWrapper* queryWrapper = static_cast<NavMeshQueryWrapper*>(query);
    
    float startPos[3] = {start_x, start_y, start_z};
    float endPos[3] = {end_x, end_y, end_z};
    
    dtPolyRef startRef, endRef;
    float nearestStart[3], nearestEnd[3];
    
    static const dtQueryFilter filter;
    
    queryWrapper->query->findNearestPoly(startPos, nullptr, &filter, &startRef, nearestStart);
    queryWrapper->query->findNearestPoly(endPos, nullptr, &filter, &endRef, nearestEnd);
    
    if (!startRef || !endRef) {
        return 0;
    }
    
    dtPolyRef polys[256];
    int npolys = 0;
    queryWrapper->query->findPath(startRef, endRef, nearestStart, nearestEnd, &filter, polys, &npolys, 256);
    
    if (npolys == 0) {
        return 0;
    }
    
    int straightPathCount = 0;
    float straightPath[256*3];
    unsigned char straightPathFlags[256];
    dtPolyRef straightPathPolys[256];
    
    queryWrapper->query->findStraightPath(nearestStart, nearestEnd, polys, npolys,
                                         straightPath, straightPathFlags, straightPathPolys,
                                         &straightPathCount, 256);
    
    int points_to_copy = (straightPathCount < max_points) ? straightPathCount : max_points;
    for (int i = 0; i < points_to_copy; ++i) {
        path[i*3] = straightPath[i*3];
        path[i*3+1] = straightPath[i*3+1];
        path[i*3+2] = straightPath[i*3+2];
    }
    
    return points_to_copy;
}

void destroy_navmesh_query(NavMeshQueryHandle query) {
    if (query) {
        NavMeshQueryWrapper* wrapper = static_cast<NavMeshQueryWrapper*>(query);
        if (wrapper->query) {
            dtFreeNavMeshQuery(wrapper->query);
        }
        delete wrapper;
    }
}

void destroy_navmesh(NavMeshHandle navMesh) {
    if (navMesh) {
        NavMeshWrapper* wrapper = static_cast<NavMeshWrapper*>(navMesh);
        if (wrapper->navMesh) {
            dtFreeNavMesh(wrapper->navMesh);
        }
        delete wrapper;
    }
}

CrowdHandle create_crowd(int maxAgents, float maxAgentRadius, NavMeshHandle navMesh) {
    if (!navMesh) return nullptr;
    
    CrowdWrapper* wrapper = new CrowdWrapper();
    wrapper->crowd = dtAllocCrowd();
    if (!wrapper->crowd) {
        delete wrapper;
        return nullptr;
    }
    
    NavMeshWrapper* navMeshWrapper = static_cast<NavMeshWrapper*>(navMesh);
    if (!wrapper->crowd->init(maxAgents, maxAgentRadius, navMeshWrapper->navMesh)) {
        dtFreeCrowd(wrapper->crowd);
        delete wrapper;
        return nullptr;
    }
    
    return static_cast<CrowdHandle>(wrapper);
}

void destroy_crowd(CrowdHandle handle) {
    if (handle) {
        CrowdWrapper* wrapper = static_cast<CrowdWrapper*>(handle);
        delete wrapper;
    }
}

int add_crowd_agent(CrowdHandle handle, float posX, float posY, float posZ, 
                    float radius, float height, float maxAcceleration, float maxSpeed,
                    float collisionQueryRange, float pathOptimizationRange, float separationWeight,
                    unsigned char updateFlags, unsigned char obstacleAvoidanceType, 
                    unsigned char queryFilterType) {
    if (!handle) return -1;
    
    CrowdWrapper* wrapper = static_cast<CrowdWrapper*>(handle);
    dtCrowdAgentParams params;
    params.radius = radius;
    params.height = height;
    params.maxAcceleration = maxAcceleration;
    params.maxSpeed = maxSpeed;
    params.collisionQueryRange = collisionQueryRange;
    params.pathOptimizationRange = pathOptimizationRange;
    params.separationWeight = separationWeight;
    params.updateFlags = updateFlags;
    params.obstacleAvoidanceType = obstacleAvoidanceType;
    params.queryFilterType = queryFilterType;
    params.userData = nullptr;
    
    float pos[3] = {posX, posY, posZ};
    return wrapper->crowd->addAgent(pos, &params);
}

void remove_crowd_agent(CrowdHandle handle, int agentIdx) {
    if (!handle) return;
    CrowdWrapper* wrapper = static_cast<CrowdWrapper*>(handle);
    wrapper->crowd->removeAgent(agentIdx);
}

bool request_crowd_agent_move_target(CrowdHandle handle, int agentIdx, 
                                    uint64_t polyRef, float targetX, float targetY, float targetZ) {
    if (!handle) return false;
    CrowdWrapper* wrapper = static_cast<CrowdWrapper*>(handle);
    float targetPos[3] = {targetX, targetY, targetZ};
    return wrapper->crowd->requestMoveTarget(agentIdx, static_cast<dtPolyRef>(polyRef), targetPos);
}

void update_crowd(CrowdHandle handle, float dt) {
    if (!handle) return;
    CrowdWrapper* wrapper = static_cast<CrowdWrapper*>(handle);
    wrapper->crowd->update(dt, nullptr);
}

int get_crowd_active_agents(CrowdHandle handle, float* positions, float* velocities, int maxAgents) {
    if (!handle || !positions || !velocities) return 0;
    CrowdWrapper* wrapper = static_cast<CrowdWrapper*>(handle);
    
    dtCrowdAgent** agents = new dtCrowdAgent*[maxAgents];
    int numAgents = wrapper->crowd->getActiveAgents(agents, maxAgents);
    
    for (int i = 0; i < numAgents; i++) {
        positions[i * 3] = agents[i]->npos[0];
        positions[i * 3 + 1] = agents[i]->npos[1];
        positions[i * 3 + 2] = agents[i]->npos[2];
        
        velocities[i * 3] = agents[i]->vel[0];
        velocities[i * 3 + 1] = agents[i]->vel[1];
        velocities[i * 3 + 2] = agents[i]->vel[2];
    }
    
    delete[] agents;
    return numAgents;
}

#ifdef __cplusplus
}
#endif