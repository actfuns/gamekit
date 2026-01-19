package main

import (
	"fmt"
	"log"
	"os"

	"dr/detour"
)

func main() {
	// 使用当前工作目录下的detour.bin
	navMeshFile := "../assets/navmesh.bin"

	// 检查文件是否存在
	if _, err := os.Stat(navMeshFile); os.IsNotExist(err) {
		log.Fatalf("Navigation mesh file %s does not exist in current directory", navMeshFile)
	}

	// Load the navigation mesh
	detour, err := detour.LoadNavMeshFromFile(navMeshFile)
	if err != nil {
		log.Fatal("Failed to load detour:", err)
	}
	defer detour.Close()

	// Create a query object for pathfinding
	query, err := detour.CreateQuery()
	if err != nil {
		log.Fatal("Failed to create query:", err)
	}
	defer query.Close()

	// 定义多组寻路测试点
	testCases := []struct {
		name  string
		start [3]float32
		end   [3]float32
	}{
		{"Test 1", [3]float32{0.0, 0.0, 0.0}, [3]float32{10.0, 0.0, 10.0}},
		{"Test 2", [3]float32{5.0, 0.0, 5.0}, [3]float32{15.0, 0.0, 15.0}},
		{"Test 3", [3]float32{0.0, 0.0, 10.0}, [3]float32{10.0, 0.0, 0.0}},
		{"Test 4", [3]float32{2.0, 0.0, 2.0}, [3]float32{8.0, 0.0, 8.0}},
		{"Test 5", [3]float32{1.0, 0.0, 1.0}, [3]float32{9.0, 0.0, 9.0}},
	}

	// 执行多次寻路
	for _, tc := range testCases {
		fmt.Printf("\n=== %s ===\n", tc.name)
		fmt.Printf("Finding path from (%.2f, %.2f, %.2f) to (%.2f, %.2f, %.2f)\n",
			tc.start[0], tc.start[1], tc.start[2],
			tc.end[0], tc.end[1], tc.end[2])

		path, err := query.FindStraightPath(tc.start[0], tc.start[1], tc.start[2],
			tc.end[0], tc.end[1], tc.end[2])
		if err != nil {
			fmt.Printf("Failed to find path: %v\n", err)
			continue
		}

		fmt.Printf("Found path with %d points:\n", len(path))
		for i, point := range path {
			fmt.Printf("  Point %d: (%.2f, %.2f, %.2f)\n", i, point[0], point[1], point[2])
		}
	}

	fmt.Println("\nMultiple pathfinding tests completed!")
}
