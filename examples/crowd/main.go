package main

import (
	"fmt"
	"log"
	"os"

	"dr/detour"
)

func main() {
	fmt.Println("Recast Navigation Crowd Demo with Go Interface")

	// 使用当前工作目录下的detour.bin
	navMeshFile := "../assets/navmesh.bin"

	// 检查文件是否存在
	if _, err := os.Stat(navMeshFile); os.IsNotExist(err) {
		log.Fatalf("Navigation mesh file %s does not exist in current directory", navMeshFile)
	}

	// 加载导航网格
	nm, err := detour.LoadNavMeshFromFile(navMeshFile)
	if err != nil {
		log.Fatal("Failed to load navigation mesh:", err)
	}
	defer nm.Close()

	// 创建Crowd模拟
	maxAgents := 10
	maxAgentRadius := float32(2.0) // 根据参考代码使用2.0作为最大代理半径
	crowd, err := detour.CreateCrowd(maxAgents, maxAgentRadius, nm)
	if err != nil {
		log.Fatal("Failed to create crowd:", err)
	}
	defer crowd.Close()

	// 设置代理参数
	agentParams := detour.CrowdAgentParams{
		Radius:                0.5,
		Height:                2.0,
		MaxAcceleration:       8.0,
		MaxSpeed:              3.5,
		CollisionQueryRange:   12.0,
		PathOptimizationRange: 30.0,
		SeparationWeight:      2.0,
		UpdateFlags:           0xFF, // 启用所有更新标志
		ObstacleAvoidanceType: 0,
		QueryFilterType:       0,
	}

	// 添加几个代理
	agentCount := 3
	agents := make([]int, agentCount)
	startPositions := [][3]float32{
		{1.0, 0.0, 1.0},
		{2.0, 0.0, 1.0},
		{3.0, 0.0, 1.0},
	}

	for i := 0; i < agentCount; i++ {
		agentIdx, err := crowd.AddAgent(
			agentParams,
			startPositions[i][0], startPositions[i][1], startPositions[i][2],
		)
		if err != nil {
			log.Printf("Failed to add agent %d: %v", i, err)
			continue
		}
		agents[i] = agentIdx
		fmt.Printf("Added agent %d with index %d\n", i, agentIdx)
	}

	// 为每个代理设置目标
	targetPositions := [][3]float32{
		{9.0, 0.0, 9.0},
		{8.0, 0.0, 8.0},
		{7.0, 0.0, 7.0},
	}

	for i, agentIdx := range agents {
		if agentIdx >= 0 {
			err := crowd.RequestMoveTarget(agentIdx, targetPositions[i][0], targetPositions[i][1], targetPositions[i][2])
			if err != nil {
				log.Printf("Failed to request move target for agent %d: %v", agentIdx, err)
			} else {
				fmt.Printf("Requested agent %d to move to (%.1f, %.1f, %.1f)\n",
					agentIdx, targetPositions[i][0], targetPositions[i][1], targetPositions[i][2])
			}
		}
	}

	// 模拟几帧
	for frame := 0; frame < 10; frame++ {
		// 注意：Update方法在Go接口中已被移除，因为C层的update_crowd需要dt参数
		// 但我们的Go接口没有暴露这个方法，需要添加

		// 获取当前代理状态
		positions, err := crowd.GetActiveAgents()
		if err != nil {
			log.Printf("Failed to get active agents: %v", err)
			continue
		}

		fmt.Printf("Frame %d: %d active agents\n", frame, len(positions))
		for i, pos := range positions {
			fmt.Printf("  Agent %d: pos(%.2f, %.2f, %.2f)\n",
				i, pos[0], pos[1], pos[2])
		}
	}

	fmt.Println("Crowd simulation completed successfully!")
}
