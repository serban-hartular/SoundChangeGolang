package main

import (
	"fmt"
	"maps"
	"slices"
	"sort"
)

type double float64

type SearchNode struct {
	vocabulary    Vocabulary
	targetVocab   Vocabulary
	abc           *Alphabet
	parent        *SearchNode
	changeApplied []ContextualChange
	children      map[string]*SearchNode
	editDistance  int
	depth         int
}

func NewSearchNode(vocab Vocabulary, target Vocabulary, abc *Alphabet, parent *SearchNode, changeApplied []ContextualChange) *SearchNode {
	cc := SearchNode{vocab, target, abc, parent, changeApplied, nil, 0, 0}
	cc.editDistance = cc.vocabulary.getDistance(cc.targetVocab)
	cc.depth = parent.depth + 1
	return &cc
}

func NewRootNode(vocab Vocabulary, target Vocabulary, abc *Alphabet) *SearchNode {
	node := SearchNode{vocab, target, abc, nil, make([]ContextualChange, 0), nil, 0, 0}
	node.editDistance = node.vocabulary.getDistance(node.targetVocab)
	return &node
}

func (node *SearchNode) ChildrenList() []*SearchNode {
	return slices.Collect(maps.Values(node.children))
}

func (node *SearchNode) expand() {
	fmt.Printf("Expanding node distance=%d, ", node.editDistance)
	changes := node.getPossibleChanges(-1, -1)
	fmt.Printf("possible changes=%d, ", len(changes))
	node.children = make(map[string]*SearchNode)
	for _, cc := range changes {
		new_vocab := node.abc.changedVocabulary(node.vocabulary, &cc)
		key := new_vocab.String()
		if child, ok := node.children[key]; ok {
			child.changeApplied = append(child.changeApplied, cc)
		} else {
			node.children[key] = NewSearchNode(new_vocab, node.targetVocab, node.abc, node, []ContextualChange{cc})
		}
	}
	fmt.Printf("expanded %d children\n", len(node.children))
	// sort changeApplied by degree of generality
	for _, child := range node.children {
		sort.SliceStable(child.changeApplied, func(i, j int) bool {
			return node.abc.generalityScoreChange(&child.changeApplied[i]) < node.abc.generalityScoreChange(&child.changeApplied[j])
		})
	}
}

func (node *SearchNode) getPossibleChanges(pre_len, post_len int) []ContextualChange {
	changesInContext := node.vocabulary.getAllChangesInContext(node.targetVocab)
	changesUnique := make(map[string]ContextualChange)
	for _, change := range changesInContext {
		combos := node.abc.getContextulChangeCombinations(&change, pre_len, post_len)
		for _, cc := range combos {
			changesUnique[cc.String()] = cc
		}
	}
	return slices.Collect(maps.Values(changesUnique))
}

func evalFn(node *SearchNode) double {
	cost := double(node.depth) * 0.99
	eval_remaining := double(node.editDistance)
	return cost + eval_remaining
}

func findSolution(root *SearchNode) *SearchNode {
	queue := []*SearchNode{root}
	best := root
	for len(queue) > 0 {
		sort.SliceStable(queue, func(i, j int) bool {
			return evalFn(queue[i]) < evalFn(queue[j])
		})
		best = queue[0]
		queue = queue[1:]
		if best.vocabulary.equals(best.targetVocab) {
			break
		} else {
			best.expand()
			queue = append(queue, slices.Collect(maps.Values(best.children))...)
		}
	}
	return best
}

func solutionToPath(solution *SearchNode) []*SearchNode {
	path := make([]*SearchNode, 0)
	for solution != nil {
		path = append(path, solution)
		solution = solution.parent
	}
	slices.Reverse(path)
	return path
}
