package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"dr/navmesh"
)

func main() {
	// Command line flags
	objFile := flag.String("obj", "", "Input OBJ file path")
	outputFile := flag.String("output", "navmesh.bin", "Output navigation mesh binary file path")
	
	// Recast parameters with reasonable defaults
	cellSize := flag.Float64("cellSize", 0.3, "Rasterization cell size")
	cellHeight := flag.Float64("cellHeight", 0.2, "Rasterization cell height")
	agentHeight := flag.Float64("agentHeight", 2.0, "Agent height")
	agentRadius := flag.Float64("agentRadius", 0.6, "Agent radius")
	agentMaxClimb := flag.Float64("agentMaxClimb", 0.9, "Agent max climb")
	agentMaxSlope := flag.Float64("agentMaxSlope", 45.0, "Agent max slope angle in degrees")
	regionMinSize := flag.Int("regionMinSize", 8, "Minimum region size")
	regionMergeSize := flag.Int("regionMergeSize", 20, "Region merge size")
	edgeMaxLen := flag.Float64("edgeMaxLen", 12.0, "Maximum edge length")
	edgeMaxError := flag.Float64("edgeMaxError", 1.3, "Maximum edge simplification error")
	vertsPerPoly := flag.Int("vertsPerPoly", 6, "Maximum vertices per polygon")
	detailSampleDist := flag.Float64("detailSampleDist", 6.0, "Detail mesh sample distance")
	detailSampleMaxError := flag.Float64("detailSampleMaxError", 1.0, "Detail mesh maximum sample error")
	
	flag.Parse()
	
	// Validate input
	if *objFile == "" {
		fmt.Println("Usage: obj2navmesh -obj <input.obj> [-output <output.bin>] [parameters...]")
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	// Check if OBJ file exists
	if _, err := os.Stat(*objFile); os.IsNotExist(err) {
		log.Fatalf("OBJ file does not exist: %s", *objFile)
	}
	
	fmt.Printf("Building navigation mesh from: %s\n", *objFile)
	fmt.Printf("Output file: %s\n", *outputFile)
	
	// Build navigation mesh
	err := navmesh.BuildNavMeshFromObj(
		*objFile, *outputFile,
		float32(*cellSize), float32(*cellHeight),
		float32(*agentHeight), float32(*agentRadius), float32(*agentMaxClimb), float32(*agentMaxSlope),
		*regionMinSize, *regionMergeSize,
		float32(*edgeMaxLen), float32(*edgeMaxError),
		*vertsPerPoly, float32(*detailSampleDist), float32(*detailSampleMaxError))
	
	if err != nil {
		log.Fatalf("Failed to build navigation mesh: %v", err)
	}
	
	fmt.Printf("Navigation mesh successfully built and saved to: %s\n", *outputFile)
}