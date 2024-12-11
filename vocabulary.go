package main

import (
	// "maps"
	// "slices"
	"maps"
	"slices"
	"strings"
)

// type WordTransition struct {
// 	initial          SymStr
// 	final            SymStr
// 	change_score     int
// 	change_sequences []ChangeSequence
// }

// func NewWordTransition(initial, final SymStr) WordTransition {
// 	return WordTransition{initial, final, string_distance(initial, final), make([]ChangeSequence, 0)}
// }

// func (wt WordTransition) addChangeSequences() {
// 	matrix := GenerateChangeMatrix(wt.initial, wt.final)
// 	wt.change_sequences = find_change_sequences(matrix, NO_POS)
// }

type Vocabulary []SymStr

func (v Vocabulary) equals(other Vocabulary) bool {
	if len(v) != len(other) {
		return false
	}
	for i, word := range v {
		if word.String() != other[i].String() {
			return false
		}
	}
	return true
}

func (v Vocabulary) String() string {
	s_list := listComprehension(v, func(ss SymStr) string { return ss.String() })
	return strings.Join(s_list, ", ")
}

func (v Vocabulary) applyChange(cc ContextualChange) Vocabulary {
	v1 := make(Vocabulary, len(v))
	for i, word := range v {
		v1[i] = cc.applyString(word)
	}
	return v1
}

func (v Vocabulary) getDistance(other Vocabulary) int {
	if len(v) != len(other) {
		panic("Unequal vocabularies being compared")
	}
	dist := 0
	for i := range v {
		dist += string_distance(v[i], other[i])
	}
	return dist
}

func (v Vocabulary) getAllChangesInContext(other Vocabulary) []ContextualChange {
	if len(v) != len(other) {
		panic("Unequal vocabularies being compared")
	}
	all_ccs := make(map[string]ContextualChange) //make([]ContextualChange, 0)
	for i := range v {
		initial := v[i]
		final := other[i]
		change_seqs := wordPairChangeSequencesAll(initial, final) //word_pair_change_sequences(initial, final)
		for _, seq := range change_seqs {
			//all_ccs = append(all_ccs, seq.getChangesInContext()...)
			for _, cc := range seq.getChangesInContext() {
				all_ccs[cc.String()] = cc
			}
		}
	}
	return slices.Collect(maps.Values(all_ccs))
}

func (chg_seq EditSequence) getChangesInContext() []ContextualChange {
	initialStr := make(SymStr, 0, len(chg_seq))
	chg_indices := make([]int, 0, len(chg_seq))
	chg_string_positions := make([]int, 0)
	//traverse change sequence to reconstruct initial string and mark positions in initial stirng where
	//changes occur
	for i, chg := range chg_seq {
		if !chg.NoChange() {
			chg_indices = append(chg_indices, i)
			chg_string_positions = append(chg_string_positions, len(initialStr))
		}
		initialStr = SymStrConcat(initialStr, chg.s_in)
	}
	cc_list := make([]ContextualChange, len(chg_indices))
	for i, chg_index := range chg_indices {
		chg := chg_seq[chg_index]
		chg_position := chg_string_positions[i]
		pre := initialStr[:chg_position]
		post := initialStr[chg_position+len(chg.s_in):]
		cc_list[i] = NewContextualChange(chg.s_in, chg.s_out, pre, post)
	}
	return cc_list
}
