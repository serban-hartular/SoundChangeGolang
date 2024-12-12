package main

import (
	"fmt"
	"strings"
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
	changes := wordPairChangeSequencesAll(vocab[0], target[0])
	for _, change := range changes {
		fmt.Println(change)
		fmt.Println(strings.Join(listComprehension(editSequence2ContextualChanges(change), func(cc ContextualChange) string { return cc.String() }), "; "))
	}
}
