package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"slices"
)

const nodeWidth = 4
const nodeDistance = nodeWidth * 1.5
const aspectRatio = 2 // layerDistance/nodeDistance

type Optimizer interface {
	// Optimize layout for shortest connection distance
	Optimize(graph *NodeGraph, layout *Layout)
}

type LayoutIndex struct {
	layer int // front to back
	x     int // left to right
	y     int // bottom to top
}

type Layout struct {
	// In these arrays, nodes are stored from left to right
	inputNodes  []NodeId
	outputNodes []NodeId

	// middleNodes are stored from input to output
	middleNodes [][]NodeId

	layerCount int

	circuitWidth  int
	nodePositions [][][]rl.Vector3 // [layer][y][x]
	nodeIndices   []LayoutIndex
}

func NewLayout(layers int) *Layout {
	return &Layout{
		inputNodes:  []NodeId{},
		outputNodes: []NodeId{},
		middleNodes: make([][]NodeId, layers-1),
		layerCount:  layers,
		nodeIndices: []LayoutIndex{},
		//	nodePositions intentionally left nil, it will be calculated lazily
	}
}

func (l *Layout) calculateNodePositions() {
	fmt.Println(l.circuitWidth, l.layerCount)

	nodePositions := make([][][]rl.Vector3, l.layerCount)
	layerDistance := float32(nodeDistance) * aspectRatio

	for layerIdx := 0; layerIdx < l.layerCount; layerIdx++ {
		layer := l.getLayer(layerIdx)
		layerNodeCount := len(layer)
		layerWidth := min(l.circuitWidth, layerNodeCount)
		fmt.Println(layerIdx, layerNodeCount, l.circuitWidth)

		//if layerNodeCount > l.circuitWidth {
		//	fmt.Println(layerIdx, layerNodeCount, l.circuitWidth)
		//}

		layerPositions := [][]rl.Vector3{}
		for i := 0; i < layerNodeCount; i++ {
			x := i % layerWidth
			y := i / layerWidth

			l.nodeIndices[layer[i]] = LayoutIndex{layerIdx, x, y}

			position := rl.NewVector3(
				float32(layerIdx)*layerDistance,  // front to back
				float32(y)*float32(nodeDistance), // up down
				float32(x)*float32(nodeDistance), // side to side
			)

			if len(layerPositions)-1 < y {
				layerPositions = append(layerPositions, []rl.Vector3{})
			}
			layerPositions[y] = append(layerPositions[y], position)

			if layerPositions[y][x] != position {
				panic("Uh oh!!")
			}

		}
		nodePositions[layerIdx] = layerPositions
	}
	l.nodePositions = nodePositions
}

func (l *Layout) forEachLayer(f func([]NodeId)) {
	f(l.inputNodes)
	for _, layer := range l.middleNodes {
		f(layer)
	}
	f(l.outputNodes)
}

func (l *Layout) getLayer(layer int) []NodeId {
	switch layer {
	case 0:
		return l.inputNodes
	case l.layerCount - 1:
		return l.outputNodes
	default:
		return l.middleNodes[layer-1]
	}
}

func (l *Layout) GetNodePosition(idx LayoutIndex) rl.Vector3 {
	if l.nodePositions == nil {
		l.calculateNodePositions()
	}

	return l.nodePositions[idx.layer][idx.y][idx.x]
}

func (l *Layout) GetNodeIndex(node NodeId) LayoutIndex {
	return l.nodeIndices[node]
}

func (l *Layout) LayerWidth(layer int) int {
	return min(l.circuitWidth, len(l.getLayer(layer)))
}

func (l *Layout) ForEachNode(f func(LayoutIndex, NodeId)) {
	for layerIdx := 0; layerIdx < l.layerCount; layerIdx++ {
		layer := l.getLayer(layerIdx)
		layerNodeCount := len(layer)
		layerWidth := min(l.circuitWidth, layerNodeCount)

		for i := 0; i < layerNodeCount; i++ {
			x := i % layerWidth
			y := i / layerWidth

			//fmt.Printf("%d %d\n", x, y)

			f(LayoutIndex{layerIdx, x, y}, layer[i])
		}

	}
}

func (l *Layout) AddInputNode(ng *NodeGraph) NodeId {
	nodeId := ng.addNode("", []NodeId{})
	l.inputNodes = append(l.inputNodes, nodeId)
	l.circuitWidth = max(len(l.inputNodes), len(l.outputNodes))
	l.nodeIndices = append(l.nodeIndices, LayoutIndex{})
	return nodeId
}

// AddMiddleNode puts a new node in a layer
func (l *Layout) AddMiddleNode(ng *NodeGraph, label string, layer int, backwardConnections []NodeId) NodeId {
	nodeId := ng.addNode(label, backwardConnections)
	l.middleNodes[layer-1] = append(l.middleNodes[layer-1], nodeId)
	l.nodeIndices = append(l.nodeIndices, LayoutIndex{})
	return nodeId
}

func (l *Layout) AddOutputNode(ng *NodeGraph, backwardConnections []NodeId) NodeId {
	nodeId := ng.addNode("", backwardConnections)
	l.outputNodes = append(l.outputNodes, nodeId)
	l.circuitWidth = max(len(l.inputNodes), len(l.outputNodes))
	l.nodeIndices = append(l.nodeIndices, LayoutIndex{})
	return nodeId
}

func (l *Layout) GetCircuitDimensions() (float32, float32) {
	// width,depth

	width := (nodeDistance * float32(l.circuitWidth)) + nodeWidth
	depth := (nodeDistance * aspectRatio * float32(l.layerCount)) + nodeWidth

	return width, depth
}

func (l *Layout) GetNodeDimensions() float32 {
	return nodeWidth
}

type NodeId int

type Node struct {
	Id                  NodeId
	Label               string
	ForwardConnections  []NodeId
	BackwardConnections []NodeId
}

type NodeGraph struct {
	nodes []Node
}

func NewNodeGraph() *NodeGraph {
	return &NodeGraph{
		nodes: []Node{},
	}
}

func (ng *NodeGraph) addNode(label string, backwardConnections []NodeId) NodeId {
	nodeId := NodeId(len(ng.nodes))
	node := Node{Id: nodeId, Label: label, ForwardConnections: []NodeId{}, BackwardConnections: backwardConnections}
	ng.nodes = append(ng.nodes, node)

	if len(backwardConnections) != 0 {
		for _, previousNodeId := range backwardConnections {
			previousNodeForwardConnections := ng.nodes[previousNodeId].ForwardConnections
			if !slices.Contains(previousNodeForwardConnections, nodeId) {
				ng.nodes[previousNodeId].ForwardConnections = append(previousNodeForwardConnections, nodeId)
			}
		}
	}

	return nodeId
}

func (ng *NodeGraph) ForwardConnections(nodeId NodeId) []NodeId {
	return ng.nodes[nodeId].ForwardConnections
}

func (ng *NodeGraph) BackwardConnections(nodeId NodeId) []NodeId {
	return ng.nodes[nodeId].BackwardConnections
}

func (ng *NodeGraph) Label(nodeId NodeId) string {
	return ng.nodes[nodeId].Label
}
