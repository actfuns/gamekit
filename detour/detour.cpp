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

// =============================================================================
// Navigation Mesh Generation Implementation (from OBJ files)
// =============================================================================

#include "Recast.h"
#include "RecastAlloc.h"
#include "RecastAssert.h"
#include "ChunkyTriMesh.h"

// Simple OBJ mesh loader (adapted from RecastDemo)
struct ObjMesh {
    float* verts;
    int nverts;
    int* tris;
    int ntris;
    
    ObjMesh() : verts(nullptr), nverts(0), tris(nullptr), ntris(0) {}
    ~ObjMesh() {
        if (verts) dtFree(verts);
        if (tris) dtFree(tris);
    }
};

static bool loadObj(const char* filename, ObjMesh& mesh)
{
    FILE* fp = fopen(filename, "r");
    if (!fp) return false;
    
    // Count vertices and faces
    int vertCount = 0, faceCount = 0;
    char line[256];
    while (fgets(line, sizeof(line), fp)) {
        if (line[0] == 'v' && line[1] == ' ') vertCount++;
        else if (line[0] == 'f') faceCount++;
    }
    rewind(fp);
    
    if (vertCount == 0 || faceCount == 0) {
        fclose(fp);
        return false;
    }
    
    mesh.verts = (float*)dtAlloc(sizeof(float) * vertCount * 3, DT_ALLOC_PERM);
    mesh.nverts = vertCount;
    mesh.tris = (int*)dtAlloc(sizeof(int) * faceCount * 3, DT_ALLOC_PERM);
    mesh.ntris = faceCount;
    
    int vidx = 0, tidx = 0;
    while (fgets(line, sizeof(line), fp)) {
        if (line[0] == 'v' && line[1] == ' ') {
            float x, y, z;
            if (sscanf(line, "v %f %f %f", &x, &y, &z) == 3) {
                mesh.verts[vidx++] = x;
                mesh.verts[vidx++] = y;
                mesh.verts[vidx++] = z;
            }
        }
        else if (line[0] == 'f') {
            int a, b, c;
            // Handle different face formats (f 1 2 3 or f 1/1/1 2/2/2 3/3/3)
            if (sscanf(line, "f %d %d %d", &a, &b, &c) == 3 ||
                sscanf(line, "f %d/%*d/%*d %d/%*d/%*d %d/%*d/%*d", &a, &b, &c) == 3) {
                mesh.tris[tidx++] = a - 1; // OBJ is 1-indexed
                mesh.tris[tidx++] = b - 1;
                mesh.tris[tidx++] = c - 1;
            }
        }
    }
    
    fclose(fp);
    return true;
}

static unsigned char* buildTileMesh(const float* verts, const int nverts,
                                   const int* tris, const int ntris,
                                   const float bmin[3], const float bmax[3],
                                   const float cellSize, const float cellHeight,
                                   const float agentHeight, const float agentRadius,
                                   const float agentMaxClimb, const float agentMaxSlope,
                                   const int regionMinSize, const int regionMergeSize,
                                   const float edgeMaxLen, const float edgeMaxError,
                                   const int vertsPerPoly, const float detailSampleDist,
                                   const float detailSampleMaxError,
                                   int* outDataSize)
{
    rcContext ctx;
    
    // Calculate bounds
    rcConfig cfg;
    memset(&cfg, 0, sizeof(cfg));
    cfg.cs = cellSize;
    cfg.ch = cellHeight;
    cfg.walkableHeight = (int)ceilf(agentHeight / cfg.ch);
    cfg.walkableClimb = (int)floorf(agentMaxClimb / cfg.ch);
    cfg.walkableRadius = (int)ceilf(agentRadius / cfg.cs);
    cfg.walkableSlopeAngle = agentMaxSlope;
    cfg.tileSize = 48;
    cfg.borderSize = cfg.walkableRadius + 3;
    cfg.width = cfg.tileSize + cfg.borderSize * 2;
    cfg.height = cfg.tileSize + cfg.borderSize * 2;
    cfg.maxEdgeLen = (int)(edgeMaxLen / cfg.cs);
    cfg.maxSimplificationError = edgeMaxError;
    cfg.minRegionArea = (int)rcSqr(regionMinSize);
    cfg.mergeRegionArea = (int)rcSqr(regionMergeSize);
    cfg.maxVertsPerPoly = vertsPerPoly;
    cfg.detailSampleDist = detailSampleDist < 0.9f ? 0 : cfg.cs * detailSampleDist;
    cfg.detailSampleMaxError = cfg.ch * detailSampleMaxError;
    
    rcVcopy(cfg.bmin, bmin);
    rcVcopy(cfg.bmax, bmax);
    
    // Allocate voxel heightfield
    rcHeightfield* hf = rcAllocHeightfield();
    if (!hf) return nullptr;
    
    if (!rcCreateHeightfield(&ctx, *hf, cfg.width, cfg.height, cfg.bmin, cfg.bmax, cfg.cs, cfg.ch)) {
        rcFreeHeightField(hf);
        return nullptr;
    }
    
    // Mark walkable triangles
    unsigned char* triAreas = new unsigned char[ntris];
    memset(triAreas, 0, ntris * sizeof(unsigned char));
    rcMarkWalkableTriangles(&ctx, cfg.walkableSlopeAngle, verts, nverts, tris, ntris, triAreas);
    
    // Rasterize triangles
    if (!rcRasterizeTriangles(&ctx, verts, nverts, tris, triAreas, ntris, *hf, cfg.walkableClimb)) {
        delete[] triAreas;
        rcFreeHeightField(hf);
        return nullptr;
    }
    
    delete[] triAreas;
    
    // Filter walkable areas
    rcFilterLowHangingWalkableObstacles(&ctx, cfg.walkableClimb, *hf);
    rcFilterLedgeSpans(&ctx, cfg.walkableHeight, cfg.walkableClimb, *hf);
    rcFilterWalkableLowHeightSpans(&ctx, cfg.walkableHeight, *hf);
    
    // Partition
    rcCompactHeightfield* chf = rcAllocCompactHeightfield();
    if (!chf) {
        rcFreeHeightField(hf);
        return nullptr;
    }
    
    if (!rcBuildCompactHeightfield(&ctx, cfg.walkableHeight, cfg.walkableClimb, *hf, *chf)) {
        rcFreeCompactHeightfield(chf);
        rcFreeHeightField(hf);
        return nullptr;
    }
    
    rcFreeHeightField(hf);
    
    if (!rcErodeWalkableArea(&ctx, cfg.walkableRadius, *chf)) {
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    if (!rcBuildDistanceField(&ctx, *chf)) {
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    if (!rcBuildRegions(&ctx, *chf, cfg.borderSize, cfg.minRegionArea, cfg.mergeRegionArea)) {
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    // Build contours
    rcContourSet* cset = rcAllocContourSet();
    if (!cset) {
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    if (!rcBuildContours(&ctx, *chf, cfg.maxSimplificationError, cfg.maxEdgeLen, *cset)) {
        rcFreeContourSet(cset);
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    if (cset->nconts == 0) {
        rcFreeContourSet(cset);
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    // Build polymesh
    rcPolyMesh* pmesh = rcAllocPolyMesh();
    if (!pmesh) {
        rcFreePolyMesh(pmesh);
        rcFreeContourSet(cset);
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    if (!rcBuildPolyMesh(&ctx, *cset, cfg.maxVertsPerPoly, *pmesh)) {
        rcFreePolyMesh(pmesh);
        rcFreeContourSet(cset);
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    // Build detail mesh
    rcPolyMeshDetail* dmesh = rcAllocPolyMeshDetail();
    if (!dmesh) {
        rcFreePolyMeshDetail(dmesh);
        rcFreePolyMesh(pmesh);
        rcFreeContourSet(cset);
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    if (!rcBuildPolyMeshDetail(&ctx, *pmesh, *chf, cfg.detailSampleDist, cfg.detailSampleMaxError, *dmesh)) {
        rcFreePolyMeshDetail(dmesh);
        rcFreePolyMesh(pmesh);
        rcFreeContourSet(cset);
        rcFreeCompactHeightfield(chf);
        return nullptr;
    }
    
    rcFreeCompactHeightfield(chf);
    rcFreeContourSet(cset);
    
    // Set flags
    for (int i = 0; i < pmesh->npolys; ++i) {
        pmesh->flags[i] = 1; // All polygons are walkable
    }
    
    // Create navmesh data
    dtNavMeshCreateParams params;
    memset(&params, 0, sizeof(params));
    params.verts = pmesh->verts;
    params.vertCount = pmesh->nverts;
    params.polys = pmesh->polys;
    params.polyAreas = pmesh->areas;
    params.polyFlags = pmesh->flags;
    params.polyCount = pmesh->npolys;
    params.nvp = pmesh->nvp;
    params.detailMeshes = dmesh->meshes;
    params.detailVerts = dmesh->verts;
    params.detailVertsCount = dmesh->nverts;
    params.detailTris = dmesh->tris;
    params.detailTriCount = dmesh->ntris;
    rcVcopy(params.bmin, pmesh->bmin);
    rcVcopy(params.bmax, pmesh->bmax);
    params.walkableHeight = agentHeight;
    params.walkableRadius = agentRadius;
    params.walkableClimb = agentMaxClimb;
    params.tileX = 0;
    params.tileY = 0;
    params.tileLayer = 0;
    params.cs = cfg.cs;
    params.ch = cfg.ch;
    params.buildBvTree = true;
    
    unsigned char* navData = 0;
    int navDataSize = 0;
    if (!dtCreateNavMeshData(&params, &navData, &navDataSize)) {
        rcFreePolyMeshDetail(dmesh);
        rcFreePolyMesh(pmesh);
        return nullptr;
    }
    
    rcFreePolyMeshDetail(dmesh);
    rcFreePolyMesh(pmesh);
    
    *outDataSize = navDataSize;
    return navData;
}

bool build_navmesh_from_obj(const char* obj_filename, const char* output_filename,
                           float cellSize, float cellHeight,
                           float agentHeight, float agentRadius, float agentMaxClimb, float agentMaxSlope,
                           int regionMinSize, int regionMergeSize,
                           float edgeMaxLen, float edgeMaxError,
                           int vertsPerPoly, float detailSampleDist, float detailSampleMaxError)
{
    // Load OBJ mesh
    ObjMesh mesh;
    if (!loadObj(obj_filename, mesh)) {
        return false;
    }
    
    if (mesh.nverts == 0 || mesh.ntris == 0) {
        return false;
    }
    
    // Calculate bounding box
    float bmin[3] = {FLT_MAX, FLT_MAX, FLT_MAX};
    float bmax[3] = {-FLT_MAX, -FLT_MAX, -FLT_MAX};
    for (int i = 0; i < mesh.nverts; ++i) {
        const float* v = &mesh.verts[i * 3];
        bmin[0] = rcMin(bmin[0], v[0]);
        bmin[1] = rcMin(bmin[1], v[1]);
        bmin[2] = rcMin(bmin[2], v[2]);
        bmax[0] = rcMax(bmax[0], v[0]);
        bmax[1] = rcMax(bmax[1], v[1]);
        bmax[2] = rcMax(bmax[2], v[2]);
    }
    
    // Expand bounds by 1 unit on each side
    bmin[0] -= 1.0f; bmin[1] -= 1.0f; bmin[2] -= 1.0f;
    bmax[0] += 1.0f; bmax[1] += 1.0f; bmax[2] += 1.0f;
    
    // Build navigation mesh tile
    int dataSize = 0;
    unsigned char* navData = buildTileMesh(mesh.verts, mesh.nverts, mesh.tris, mesh.ntris,
                                          bmin, bmax,
                                          cellSize, cellHeight,
                                          agentHeight, agentRadius, agentMaxClimb, agentMaxSlope,
                                          regionMinSize, regionMergeSize,
                                          edgeMaxLen, edgeMaxError,
                                          vertsPerPoly, detailSampleDist, detailSampleMaxError,
                                          &dataSize);
    
    if (!navData || dataSize == 0) {
        return false;
    }
    
    // Save to file
    FILE* fp = fopen(output_filename, "wb");
    if (!fp) {
        dtFree(navData);
        return false;
    }
    
    fwrite(navData, 1, dataSize, fp);
    fclose(fp);
    dtFree(navData);
    
    return true;
}

#ifdef __cplusplus
}
#endif