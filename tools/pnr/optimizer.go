package main

import "slices"

type LayerIndex int

type NodeIndex int

type Node struct {
	Id                  NodeIndex
	Label               string
	ForwardConnections  []NodeIndex
	BackwardConnections []NodeIndex
}

type NodeGraph struct {
	// In these arrays, nodes are stored from left to right
	InputNodes  []NodeIndex
	OutputNodes []NodeIndex

	// Layers are stored from input to output, nodes within layers are stored from left to right
	Layers [][]NodeIndex

	Nodes []Node
}

func NewNodeGraph(layers int) *NodeGraph {
	return &NodeGraph{
		InputNodes:  []NodeIndex{},
		OutputNodes: []NodeIndex{},
		Layers:      make([][]NodeIndex, layers),
		Nodes:       []Node{},
	}
}

func (ng *NodeGraph) AddInputNode() NodeIndex {
	nodeId := ng.addNode("", []NodeIndex{})
	ng.InputNodes = append(ng.InputNodes, nodeId)
	return nodeId
}

// AddLayerNode puts a new node in a layer
func (ng *NodeGraph) AddLayerNode(label string, layer LayerIndex, backwardConnections []NodeIndex) NodeIndex {
	nodeId := ng.addNode(label, backwardConnections)
	ng.Layers[layer] = append(ng.Layers[layer], nodeId)
	return nodeId
}

func (ng *NodeGraph) AddOutputNode(backwardConnections []NodeIndex) NodeIndex {
	nodeId := ng.addNode("", backwardConnections)
	ng.OutputNodes = append(ng.OutputNodes, nodeId)
	return nodeId
}

func (ng *NodeGraph) addNode(label string, backwardConnections []NodeIndex) NodeIndex {
	nodeId := NodeIndex(len(ng.Nodes))
	node := Node{Id: nodeId, Label: label, ForwardConnections: []NodeIndex{}, BackwardConnections: backwardConnections}
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

type Optimizer interface {
	// Optimize inner layers for shortest connection distance
	Optimize(graph *NodeGraph) [][]NodeIndex
}
