package main

type AStar_Search struct {
	abc         *Alphabet
	queue       []*SearchNode
	costFn      func(*SearchNode) double
	heuristicFn func(*SearchNode) double
}
