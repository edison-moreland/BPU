package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"slices"
)

const nodeWidth = 10
const nodeDistance = nodeWidth
const aspectRatio = 0.5 // layerDistance/nodeDistance

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
}

func NewLayout(layers int) *Layout {
	return &Layout{
		inputNodes:  []NodeId{},
		outputNodes: []NodeId{},
		middleNodes: make([][]NodeId, layers-1),
		layerCount:  layers,
		//	nodePositions intentionally left nil, it will be calculated lazily
	}
}

func (l *Layout) calculateNodePositions() {
	nodePositions := make([][][]rl.Vector3, l.layerCount)
	layerDistance := float32(nodeDistance) * aspectRatio

	for layerIdx := 0; layerIdx < l.layerCount; layerIdx++ {
		layer := l.getLayer(layerIdx)
		layerNodeCount := len(layer)
		layerWidth := min(l.circuitWidth, layerNodeCount)

		layerPositions := [][]rl.Vector3{}
		for i := 0; i < layerNodeCount; i++ {
			x := i % layerWidth
			y := i / layerWidth

			//LayoutIndex{layerIdx, x, y}

			position := rl.NewVector3(
				float32(layerIdx)*layerDistance,  // front to back
				float32(y)*float32(nodeDistance), // up down
				float32(x)*float32(nodeDistance), // side to side
			)

			if len(layerPositions)-1 < 0 {
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

			f(LayoutIndex{layerIdx, x, y}, layer[i])
		}

	}
}

func (l *Layout) AddInputNode(ng *NodeGraph) NodeId {
	nodeId := ng.addNode("", []NodeId{})
	l.inputNodes = append(l.inputNodes, nodeId)
	l.circuitWidth = max(len(l.inputNodes), len(l.outputNodes))
	return nodeId
}

// AddMiddleNode puts a new node in a layer
func (l *Layout) AddMiddleNode(ng *NodeGraph, label string, layer int, backwardConnections []NodeId) NodeId {
	nodeId := ng.addNode(label, backwardConnections)
	l.middleNodes[layer] = append(l.middleNodes[layer-1], nodeId)
	return nodeId
}

func (l *Layout) AddOutputNode(ng *NodeGraph, backwardConnections []NodeId) NodeId {
	nodeId := ng.addNode("", backwardConnections)
	l.outputNodes = append(l.outputNodes, nodeId)
	l.circuitWidth = max(len(l.inputNodes), len(l.outputNodes))
	return nodeId
}

type NodeId int

type Node struct {
	Id                  NodeId
	Label               string
	ForwardConnections  []NodeId
	BackwardConnections []NodeId
}

type NodeGraph struct {
	Nodes []Node
}

func NewNodeGraph() *NodeGraph {
	return &NodeGraph{
		Nodes: []Node{},
	}
}

func (ng *NodeGraph) addNode(label string, backwardConnections []NodeId) NodeId {
	nodeId := NodeId(len(ng.Nodes))
	node := Node{Id: nodeId, Label: label, ForwardConnections: []NodeId{}, BackwardConnections: backwardConnections}
	ng.Nodes = append(ng.Nodes, node)

	if len(backwardConnections) != 0 {
		for _, previousNodeId := range backwardConnections {
			previousNodeForwardConnections := ng.Nodes[previousNodeId].ForwardConnections
			if !slices.Contains(previousNodeForwardConnections, nodeId) {
				ng.Nodes[previousNodeId].ForwardConnections = append(previousNodeForwardConnections, nodeId)
			}
		}
	}

	return nodeId
}
