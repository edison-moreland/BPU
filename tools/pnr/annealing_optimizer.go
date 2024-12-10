package main

import "context"

type annealingOptimizer struct{}

func NewAnnealingOptimizer() Optimizer {
	return annealingOptimizer{}
}

func (a annealingOptimizer) Optimize(nodes *NodeGraph, layout *Layout) {
	Visualize(context.Background(), nodes, layout)
}
