package main

import (
	"context"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Visualize a layout being optimized. Blocks until window is closed, or context is ended
func Visualize(ctx context.Context, nodes *NodeGraph, layout *Layout) {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(700, 500, "annealing")
	defer rl.CloseWindow()

	rl.SetTargetFPS(10)
	done := false
	for !rl.WindowShouldClose() && !done {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		//drawNodes(graph, visualizedLayers)

		rl.EndDrawing()

		select {
		case <-ctx.Done():
			done = true
		default:
			continue
		}
	}
}
