package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const nodeSize = 10
const nodeDistance = 60
const aspectRatio = 0.5
const maxSteps = 99_9999
const progressUpdates = 100

type annealingOptimizer struct {
}

func NewAnnealingOptimizer() Optimizer {
	return &annealingOptimizer{}
}

func (a annealingOptimizer) Optimize(graph *NodeGraph) [][]NodeIndex {
	layers := graph.Layers

	visualizedLayers := graph.Layers

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(700, 500, "annealing")
	defer rl.CloseWindow()

	rl.SetTargetFPS(10)

	done := false
	go func() {
		defer func() {
			done = true
		}()

		layers = anneal[[][]NodeIndex](
			layers,
			energy(graph),
			neighbour(graph),
			temperature,
			maxSteps,
			func(state [][]NodeIndex) {
				visualizedLayers = state
			},
		)
		//
		//time.Sleep(time.Second * 10)
	}()

	for !rl.WindowShouldClose() && !done {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		visualizeNodeGraph(graph, visualizedLayers)

		rl.EndDrawing()
	}

	//for !done {
	//	//wait
	//	time.Sleep(time.Millisecond)
	//}

	return layers
}

func visualizeNodeGraph(graph *NodeGraph, middleLayers [][]NodeIndex) {

	nodePositions := make([]rl.Vector2, len(graph.Nodes))
	placeNodes(graph, middleLayers, nodeDistance, aspectRatio, func(nodeId NodeIndex, position rl.Vector2) {

		nodePosition := rl.Vector2Add(position, rl.NewVector2(nodeDistance/2, nodeDistance/2))

		rl.DrawCircleV(nodePosition, nodeSize, rl.Black)
		nodePositions[nodeId] = nodePosition

		for _, previousNodeId := range graph.Nodes[nodeId].BackwardConnections {
			previousNodePosition := nodePositions[previousNodeId]
			rl.DrawLineV(nodePosition, previousNodePosition, rl.Black)
		}
	})

	for nodeId, position := range nodePositions {

		drawCenteredText(
			strconv.Itoa(int(nodeId)),
			nodeSize/1.5,
			rl.Vector2Add(position, rl.NewVector2(0, -(nodeSize/2.5))),
			rl.White,
		)

		drawCenteredText(
			graph.Nodes[nodeId].Label,
			nodeSize/1.5,
			rl.Vector2Add(position, rl.NewVector2(0, nodeSize/2.5)),
			rl.White,
		)
	}
}

func drawCenteredText(text string, textSize float32, position rl.Vector2, color rl.Color) {
	font := rl.GetFontDefault()
	const spacing = 1

	stringSize := rl.MeasureTextEx(font, text, textSize, spacing)

	position = rl.Vector2Subtract(position, rl.Vector2Scale(stringSize, 0.5))

	rl.DrawTextEx(font, text, position, textSize, spacing, color)

}

func anneal[T any](state T, energy func(T) float32, neighbour func(T, float32) T, temperature func(t float32) float32, maxSteps int, visualizeStep func(T)) T {
	p := func(e, ep, temp float32) float32 {
		if ep <= e {
			return 1
		}

		return float32(math.Exp(float64(-(ep - e) / temp)))
	}

	eLow := float32(math.Inf(1))
	eHigh := float32(math.Inf(-1))

	startTime := time.Now()

	s := state
	for i := 0; i < maxSteps; i++ {
		temp := temperature(1.0 - float32(i+1)/float32(maxSteps))
		newS := neighbour(s, temp)

		sEnergy := energy(s)
		newSEnergy := energy(newS)

		eLow = min(eLow, sEnergy)
		eHigh = max(eHigh, sEnergy)

		prob := p(sEnergy, newSEnergy, temp)

		if i%(maxSteps/progressUpdates) == 0 {
			elapsed := time.Since(startTime)
			fmt.Printf("%s: T = %f, E(s) = %f, E(s') = %f, P(E(s), E(s')) = %f, eLow = %f, eHigh = %f\n", elapsed.String(), temp, sEnergy, newSEnergy, prob, eLow, eHigh)
		}

		if prob >= rand.Float32() {
			s = newS
			visualizeStep(s)
			// fmt.Println("chose new state")
		}
	}

	return s
}

func energy(graph *NodeGraph) func([][]NodeIndex) float32 {
	nodePositions := make([]rl.Vector2, len(graph.Nodes))
	return func(layers [][]NodeIndex) float32 {
		placeNodes(graph, layers, nodeDistance, aspectRatio, func(nodeId NodeIndex, position rl.Vector2) {
			nodePositions[nodeId] = position
		})

		foreachForwardConnection := func(layerA []NodeIndex, f func(nodeA, nodeB NodeIndex)) {
			for _, nodeA := range layerA {
				for _, nodeB := range graph.Nodes[nodeA].ForwardConnections {
					f(nodeA, nodeB)
				}
			}
		}

		totalLength := float32(0.0)
		highestLength := float32(0.0)
		connections := 0

		measureConnections := func(layerA []NodeIndex) {
			foreachForwardConnection(layerA, func(nodeA, nodeB NodeIndex) {
				length := rl.Vector2Distance(nodePositions[nodeA], nodePositions[nodeB])

				totalLength += length
				highestLength = max(highestLength, length)
				connections += 1
			})
		}

		measureConnections(graph.InputNodes)
		for _, layer := range layers {
			measureConnections(layer)
		}
		averageLength := totalLength / float32(connections)

		//return (averageLength) / float32(len(graph.Nodes))

		//Count intersections
		countIntersections := func(layerA []NodeIndex) int {
			intersections := 0

			foreachForwardConnection(layerA, func(nodeAA, nodeAB NodeIndex) {
				foreachForwardConnection(layerA, func(nodeBA, nodeBB NodeIndex) {
					if nodeAA == nodeBA && nodeAB == nodeBB {
						return
					}

					if linesIntersect(
						nodePositions[nodeAA], nodePositions[nodeAB],
						nodePositions[nodeBA], nodePositions[nodeBB],
					) {
						intersections += 1
					}
				})
			})

			return intersections
		}

		intersections := 0
		intersections += countIntersections(graph.InputNodes)
		for _, layer := range layers {
			intersections += countIntersections(layer)
		}

		//
		//return (float32(intersections) / float32(len(graph.Nodes)))
		return (float32(intersections) / float32(len(graph.Nodes))) * (averageLength / float32(len(graph.Nodes)))
	}
}

func neighbour(graph *NodeGraph) func(layers [][]NodeIndex, temp float32) [][]NodeIndex {
	return func(layers [][]NodeIndex, temp float32) [][]NodeIndex {
		newLayers := make([][]NodeIndex, len(layers))
		for i := range newLayers {
			newLayers[i] = make([]NodeIndex, len(layers[i]))

			for j := range newLayers[i] {
				newLayers[i][j] = layers[i][j]
			}

			//	Shuffle every layer
			//rand.Shuffle(len(newLayers[i]), func(k, l int) {
			//	newLayers[i][k], newLayers[i][l] = newLayers[i][l], newLayers[i][k]
			//})

			// Shuffle a single node
			//if temp >= rand.Float32() {
			//	firstNode := rand.IntN(len(newLayers[i]))
			//	secondNode := firstNode
			//	for firstNode == secondNode {
			//		secondNode = rand.IntN(len(newLayers[i]))
			//	}
			//
			//	newLayers[i][firstNode], newLayers[i][secondNode] = newLayers[i][secondNode], newLayers[i][firstNode]
			//}

			//swaps := int(float32(len(newLayers[i])) * temp)

			//for j := 0; j < swaps; j++ {
			//	firstNode := rand.IntN(len(newLayers[i]))
			//	secondNode := firstNode
			//	for firstNode == secondNode {
			//		secondNode = rand.IntN(len(newLayers[i]))
			//	}
			//
			//	newLayers[i][firstNode], newLayers[i][secondNode] = newLayers[i][secondNode], newLayers[i][firstNode]
			//
			//}
		}

		// Pick a random layer then swap two nodes
		//swaps := max(1.0, min(int(2*temp), len(graph.Nodes)))
		swaps := 1

		for i := 0; i < swaps; i++ {
			layer := newLayers[rand.IntN(len(newLayers))]

			// shuffle layer
			//rand.Shuffle(len(newLayers[layer]), func(k, l int) {
			//	newLayers[layer][k], newLayers[layer][l] = newLayers[layer][l], newLayers[layer][k]
			//})

			// swap two nodes
			firstNode, secondNode := pickNodePair(layer)
			swapNodes(layer, firstNode, secondNode)
		}

		return newLayers
	}
}

func pickNodePair(layer []NodeIndex) (int, int) {
	a := rand.IntN(len(layer))
	b := a
	for a == b {
		b = rand.IntN(len(layer))
	}
	return a, b
}

func swapNodes(layer []NodeIndex, a, b int) {
	// Swap a and b and all nodes between
	if a < b {
		a, b = b, a
	}

	if rand.IntN(1) == 1 {
		for a > b {
			layer[a], layer[b] = layer[b], layer[a]
			a -= 1
			b += 1
		}
	} else {
		layer[a], layer[b] = layer[b], layer[a]
	}
}

func temperature(temp float32) float32 {
	// TODO: This could be more advanced
	return easeInSine(temp)
}

func easeInSine(x float32) float32 {
	return float32(1.0 - math.Cos((float64(x)*math.Pi)/2))
}

func placeNodes(graph *NodeGraph, layers [][]NodeIndex, nodeDistance, aspectRatio float32, place func(node NodeIndex, position rl.Vector2)) {
	widestLayerLen := len(graph.InputNodes)
	widestLayerLen = max(widestLayerLen, len(graph.OutputNodes))
	for _, layer := range layers {
		widestLayerLen = max(widestLayerLen, len(layer))
	}

	currentY := float32(0.0)

	placeLayer := func(y float32, nodes []NodeIndex) {
		currentX := float32(0.0)
		if len(nodes) < widestLayerLen {
			offset := (widestLayerLen - len(nodes)) / 2

			currentX += (nodeDistance * aspectRatio) * float32(offset)
		}

		for _, nodeId := range nodes {
			place(nodeId, rl.NewVector2(currentX, y))
			currentX += nodeDistance * aspectRatio
		}
	}

	placeLayer(currentY, graph.InputNodes)
	currentY += nodeDistance

	for _, layer := range layers {
		placeLayer(currentY, layer)

		currentY += nodeDistance
	}

	placeLayer(currentY, graph.OutputNodes)
}

func linesIntersect(p1, q1, p2, q2 rl.Vector2) bool {
	// Helper function to calculate orientation
	orientation := func(p, q, r rl.Vector2) int {
		val := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
		if val == 0 {
			return 0 // Collinear
		} else if val > 0 {
			return 1 // Clockwise
		} else {
			return 2 // Counterclockwise
		}
	}

	// Calculate orientations
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	// General case
	if o1 != o2 && o3 != o4 {
		return true
	}

	return false
}
