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

	cc := NewContextualChange(SingleSymbolStr("l"), SingleSymbolStr("r"), SingleSymbolStr("u"), SingleSymbolStr("u"))
	cc.compile()
	fmt.Println(cc.compiled)

	fmt.Println(vocab.applyChange(cc))
}
