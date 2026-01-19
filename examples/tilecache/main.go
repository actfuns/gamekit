package main

import (
	"fmt"
	"log"
	"os"

	"dr/navmesh"
)

func main() {
	fmt.Println("Recast Navigation TileCache Demo with Go Interface")

	navMeshFile := "../assets/navmesh.bin"

	// 检查文件是否存在
	if _, err := os.Stat(navMeshFile); os.IsNotExist(err) {
		log.Fatalf("Navigation mesh file %s does not exist in current directory", navMeshFile)
	}

	// 加载导航网格
	nm, err := navmesh.LoadNavMeshFromFile(navMeshFile)
	if err != nil {
		log.Fatal("Failed to load navigation mesh:", err)
	}
	defer nm.Close()

	// 创建TileCache参数
	params := navmesh.TileCacheParams{
		Orig:                   [3]float32{0, 0, 0},
		Cs:                     0.3,
		Ch:                     0.2,
		Width:                  256,
		Height:                 256,
		WalkableHeight:         2.0,
		WalkableRadius:         0.6,
		WalkableClimb:          0.9,
		MaxSimplificationError: 1.3,
		MaxTiles:               128,
		MaxObstacles:           128,
	}

	// 创建TileCache
	tileCache, err := navmesh.CreateTileCache(params)
	if err != nil {
		log.Fatal("Failed to create tile cache:", err)
	}
	defer tileCache.Close()

	// 添加动态障碍物
	obstacleRef, err := tileCache.AddCylinderObstacle(5.0, 0.0, 5.0, 1.0, 2.0)
	if err != nil {
		log.Fatal("Failed to add obstacle:", err)
	}

	fmt.Printf("Added obstacle with ref: %d\n", obstacleRef)

	// 更新tile cache
	err = tileCache.Update(nm)
	if err != nil {
		log.Printf("Warning: Failed to update tile cache: %v", err)
	} else {
		fmt.Println("Tile cache updated successfully")
	}

	// 移除障碍物
	err = tileCache.RemoveObstacle(obstacleRef)
	if err != nil {
		log.Printf("Warning: Failed to remove obstacle: %v", err)
	} else {
		fmt.Println("Obstacle removed successfully")
	}

	// 再次更新tile cache
	err = tileCache.Update(nm)
	if err != nil {
		log.Printf("Warning: Failed to update tile cache after removal: %v", err)
	} else {
		fmt.Println("Tile cache updated after obstacle removal")
	}

	fmt.Println("TileCache demo completed successfully!")
}
