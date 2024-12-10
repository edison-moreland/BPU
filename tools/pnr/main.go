package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"os"
	"pnr/netlist"
	"strconv"
	"strings"
)

func main() {
	netlistRaw, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	nt, err := netlist.UnmarshalNetlist(netlistRaw)
	if err != nil {
		panic(err)
	}

	_, top, found := findTop(nt)
	if !found {
		panic("Top module not found")
	}

	ng, layout, _, cellMap := moduleToNodeGraph(top)

	NewAnnealingOptimizer().Optimize(ng, layout)

	// Update nt with attributes for positions
	maxLayerWidth := 0
	for i, layer := range layout.middleNodes {
		for j, cellId := range layer {
			cellName := cellMap[cellId]
			cell := top.Cells[cellName]

			cell.Attributes["pnr_layer"] = strconv.Itoa(i)
			cell.Attributes["pnr_cell"] = strconv.Itoa(j)
		}
		maxLayerWidth = max(maxLayerWidth, len(layer))
	}
	top.Attributes["has_pnr"] = "00000000000000000000000000000001"
	top.Attributes["pnr_layer_count"] = strconv.Itoa(len(layout.middleNodes))
	top.Attributes["pnr_max_layer_width"] = strconv.Itoa(maxLayerWidth)

	newNetlistRaw, err := netlist.MarshalNetlist(nt)
	if err != nil {
		panic(err)
	}

	outputFilename := strings.TrimSuffix(os.Args[1], ".json") + "_pnr.json"

	err = os.WriteFile(outputFilename, newNetlistRaw, 0644)
	if err != nil {
		panic(err)
	}
}

func findTop(netlist *netlist.Netlist) (string, *netlist.Module, bool) {
	for moduleName, module := range netlist.Modules {
		topAttr, ok := module.Attributes["top"]
		if ok && (topAttr[len(topAttr)-1] == '1') {
			return moduleName, &module, true
		}
	}

	return "", nil, false
}

func moduleToNodeGraph(module *netlist.Module) (nodeGraph *NodeGraph, layout *Layout, netMap map[NodeId]int, cellMap map[NodeId]string) {
	// This assumes the netlist is organized in nice layers where all
	//  the inputs of one layer connect to the outputs of the next
	//  layer and no other layer. This will be true for a combinational
	//  circuit that was synthesized with buffers enabled.

	// input/output nodes map to bits while middle layers map to cells
	inputNets := getPortBits(module, netlist.Direction_Input)
	outputNets := getPortBits(module, netlist.Direction_Output)
	layerCells := []mapset.Set[string]{}

	fmt.Printf("%d inputs found \n", len(inputNets))
	fmt.Printf("%d outputs found \n", len(outputNets))

	// Map cells to the bits they are connected to
	cellsByInput := map[int][]string{}
	cellOutputNet := map[string]int{}
	cellsByOutput := map[int]string{}

	for cellId, cell := range module.Cells {
		for connectionId, connections := range cell.Connections {
			// Since this netlist should only have bit-level gates, we expect only one bit per connection
			if len(connections) != 1 {
				panic("stinky stinky")
			}
			connection := connections[0]
			if connection.IsConstant() {
				panic("stinky stinky 2 electric boogaloo")
			}

			if cell.PortDirections[connectionId] == netlist.Direction_Input {
				cells, ok := cellsByInput[connection.Index]
				if !ok {
					cellsByInput[connection.Index] = []string{cellId}
				} else {
					cellsByInput[connection.Index] = append(cells, cellId)
				}
			} else {
				cellOutputNet[cellId] = connection.Index
				cellsByOutput[connection.Index] = cellId
			}
		}
	}

	cellsConnectedTo := func(nets []int) (mapset.Set[string], bool) {
		allCells := mapset.NewSet[string]()
		for _, net := range nets {
			cells, ok := cellsByInput[net]
			if !ok {
				return nil, false
			}

			allCells.Append(cells...)
		}

		return allCells, true
	}

	// Find first layer
	layer, ok := cellsConnectedTo(inputNets)
	if !ok {
		panic("input nets not connected to cell")
	}

	//Until we reach output nets, keep looking for new layers
	previousLayer := mapset.NewSet[string]()
	done := false
	for !done {
		previousLayer = layer
		layerCells = append(layerCells, previousLayer)

		nets := []int{}
		for _, cell := range previousLayer.ToSlice() {
			nets = append(nets, cellOutputNet[cell])
		}

		layer, ok = cellsConnectedTo(nets)
		if !ok {
			done = true
		}
	}

	fmt.Printf("%d layers found \n", len(layerCells))

	// Map nodes back to netlist
	netMap = map[NodeId]int{}
	cellMap = map[NodeId]string{}

	nodeGraph = NewNodeGraph()
	layout = NewLayout(len(layerCells) + 2)

	// Input layer
	inputLayerNodes := map[int]NodeId{}
	for _, net := range inputNets {
		inputLayerNodes[net] = layout.AddInputNode(nodeGraph)
		netMap[inputLayerNodes[net]] = net
	}

	// Middle layers
	layerNodes := inputLayerNodes
	for i, middleLayer := range layerCells {
		previousLayerNodes := layerNodes
		layerNodes = map[int]NodeId{}

		for _, cellName := range middleLayer.ToSlice() {
			cellInputNets := getCellBits(module.Cells[cellName], netlist.Direction_Input)
			outputNet := getCellBits(module.Cells[cellName], netlist.Direction_Output)[0]

			cellInputNodes := []NodeId{}
			for _, net := range cellInputNets {
				cellInputNodes = append(cellInputNodes, previousLayerNodes[net])
			}

			layerNodes[outputNet] = layout.AddMiddleNode(nodeGraph, module.Cells[cellName].Type, i+1, cellInputNodes)
			cellMap[layerNodes[outputNet]] = cellName
		}
	}

	// Last layer
	for _, net := range outputNets {
		inputNode := layerNodes[net]

		outputNode := layout.AddOutputNode(nodeGraph, []NodeId{inputNode})
		netMap[outputNode] = net
	}

	fmt.Println(nodeGraph.Nodes)

	return nodeGraph, layout, netMap, cellMap
}

func getPortBits(module *netlist.Module, direction netlist.Direction) []int {
	bits := []int{}

	for _, port := range module.Ports {
		if port.Direction == direction {
			for _, bit := range port.Bits {
				if bit.IsConstant() {
					panic("ew gross")
				}

				bits = append(bits, bit.Index)
			}
		}
	}

	return bits
}

func getCellBits(cell netlist.Cell, direction netlist.Direction) []int {
	bits := []int{}

	for connectionId, connections := range cell.Connections {
		// Since this netlist should only have bit-level gates, we expect only one bit per connection
		if len(connections) != 1 {
			panic("stinky stinky")
		}
		connection := connections[0]
		if connection.IsConstant() {
			panic("stinky stinky 2 electric boogaloo")
		}

		if cell.PortDirections[connectionId] == direction {
			bits = append(bits, connection.Index)
		}
	}

	return bits
}
