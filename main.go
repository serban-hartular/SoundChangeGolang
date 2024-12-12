package main

import (
	"fmt"
	// "maps"
	// "slices"
	// "sort"
)

func main() {

	abc := NewSimpleAlphabet("ō ū ā ē ī", "")

	initial := listComprehension([]string{"pellem", "celum", "culum", "mōlam", "caballum", "clama"}, abc.SymStrFromTightString)
	final := listComprehension([]string{"piele", "cer", "cur", "moara", "cal", "cheama"}, abc.SymStrFromTightString)
	vocab := Vocabulary(initial)
	target := Vocabulary(final)
	for i, w0 := range vocab {
		fmt.Printf("%s\t%s\n", w0, target[i])
	}
	changes := wordPairChangeSequencesAll(vocab[1], target[1])
	change_seq := changes[0]
	fmt.Println(change_seq)
	cc_list := editSequence2ContextualChanges(change_seq)
	fmt.Println(cc_list[0])
	cc_options := abc.getContextualChangeCombos(cc_list[0])
	for _, cc := range cc_options {
		fmt.Println(cc, cc.compiled)
	}
}
