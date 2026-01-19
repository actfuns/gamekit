package main

import (
	"fmt"
	"log"
	"os"

	"dr/detour"
)

func main() {
	// 使用当前工作目录下的detour.bin
	detourFile := "detour.bin"

	// 检查文件是否存在
	if _, err := os.Stat(detourFile); os.IsNotExist(err) {
		log.Fatalf("Navigation mesh file %s does not exist in current directory", detourFile)
	}

	// Load navigation mesh
	nm, err := detour.LoadNavMeshFromFile(detourFile)
	if err != nil {
		log.Fatalf("Failed to load detour: %v", err)
	}
	defer nm.Close()

	// Create tile cache with parameters
	params := detour.TileCacheParams{
		Orig:                   [3]float32{0, 0, 0},
		Cs:                     0.3,
		Ch:                     0.2,
		Width:                  64,
		Height:                 64,
		WalkableHeight:         2.0,
		WalkableRadius:         0.6,
		WalkableClimb:          0.9,
		MaxSimplificationError: 1.3,
		MaxTiles:               128,
		MaxObstacles:           128,
	}

	tileCache, err := detour.CreateTileCache(params)
	if err != nil {
		log.Fatalf("Failed to create tile cache: %v", err)
	}
	defer tileCache.Close()

	// Add a cylindrical obstacle at position (0, 0, 0) with radius 1.0 and height 2.0
	obstacleRef, err := tileCache.AddCylinderObstacle(0, 0, 0, 1.0, 2.0)
	if err != nil {
		log.Fatalf("Failed to add obstacle: %v", err)
	}

	fmt.Printf("Added obstacle with ref: %d\n", obstacleRef)

	// Update the tile cache to rebuild affected navigation mesh tiles
	err = tileCache.Update(nm)
	if err != nil {
		log.Printf("Warning: Failed to update tile cache: %v", err)
	} else {
		fmt.Println("Tile cache update completed")
	}

	// Remove the obstacle
	err = tileCache.RemoveObstacle(obstacleRef)
	if err != nil {
		log.Printf("Warning: Failed to remove obstacle: %v", err)
	} else {
		fmt.Println("Removed obstacle")
	}

	// Update again to rebuild without the obstacle
	err = tileCache.Update(nm)
	if err != nil {
		log.Printf("Warning: Failed to update tile cache after removal: %v", err)
	} else {
		fmt.Println("Tile cache update completed after removal")
	}

	fmt.Println("Dynamic obstacles example completed successfully!")
}
