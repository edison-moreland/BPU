package main

import (
	"context"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Visualize a layout being optimized. Blocks until window is closed, or context is ended
func Visualize(ctx context.Context, nodes *NodeGraph, layout *Layout) {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(700, 500, "PnR Visualization")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(10.0, 10.0, 10.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	done := false
	for !rl.WindowShouldClose() && !done {
		rl.UpdateCamera(&camera, rl.CameraFree)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawGrid(100, 2.0)

		drawNodeLayout(nodes, layout)

		rl.EndMode3D()

		rl.EndDrawing()

		select {
		case <-ctx.Done():
			done = true
		default:
			continue
		}
	}
}

func drawNodeLayout(nodes *NodeGraph, layout *Layout) {
	nWidth := layout.GetNodeDimensions()
	nodeSize := rl.NewVector3(nWidth, nWidth, nWidth)
	wiresSize := rl.Vector3AddValue(nodeSize, 0.01)
	connRadius := nodeSize.X / 16

	// Position circuit in center of world
	width, depth := layout.GetCircuitDimensions()
	offset := rl.NewVector3(-(depth / 2), 0, -(width / 2))

	layout.ForEachNode(func(idx LayoutIndex, node NodeId) {
		position := rl.Vector3Add(offset, layout.GetNodePosition(idx))

		rl.DrawCubeV(position, nodeSize, rl.Blue)
		rl.DrawCubeWiresV(position, wiresSize, rl.Black)

		for _, forwardConn := range nodes.ForwardConnections(node) {
			forwardPosition := layout.GetNodePosition(layout.GetNodeIndex(forwardConn))

			//	Offset connection slightly
			forwardConnPosition := rl.Vector3Add(offset, forwardPosition)
			forwardConnPosition.X -= nodeSize.X / 2

			connPosition := position
			connPosition.X += nodeSize.X / 2

			rl.DrawCylinderEx(forwardConnPosition, connPosition, connRadius, connRadius, 4, rl.Red)
			rl.DrawCylinderWiresEx(forwardConnPosition, connPosition, connRadius+0.02, connRadius, 4, rl.Black)
		}
	})
}
