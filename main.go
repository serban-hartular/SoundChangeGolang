package main

import (
	"fmt"
	"strings"
	// "maps"
	// "slices"
	// "sort"
)

func main() {

	// data_str := load_file_as_string("nouns0.txt")
	// columns := textToColumns(data_str)
	// initial_col := columns[1][1:101] // header
	// final_col := columns[0][1:101]
	// initial := listComprehension(initial_col, NewSymStr)
	// final := listComprehension(final_col, NewSymStr)

	initial := listComprehension([]string{"pellem", "celum", "culum", "mōlam", "caballum", "clama"}, NewSymStr)
	final := listComprehension([]string{"piele", "cer", "cur", "moara", "cal", "cheama"}, NewSymStr)
	vocab := Vocabulary(initial)
	target := Vocabulary(final)
	for i, w0 := range vocab {
		fmt.Printf("%s\t%s\n", w0, target[i])
	}
	abc := quick_alphabet("yw", "")

	root := NewRootNode(vocab, target, &abc)

	// root.expand()
	// children := root.ChildrenList()
	// sort.SliceStable(children, func(i, j int) bool {
	// 	return children[i].editDistance < children[j].editDistance
	// })
	// fmt.Println(children[0].editDistance)
	// fmt.Println(children[0].changeApplied[0])
	// fmt.Println(children[0].vocabulary)

	// return
	solution := findSolution(root)
	path := solutionToPath(solution)
	for i, node := range path {
		if i == 0 {
			fmt.Printf("%d\t\t%s\n", node.editDistance, node.vocabulary)
		} else {
			fmt.Printf("%d\t%s\t%s\n", node.editDistance, node.changeApplied[0], node.vocabulary)
		}
	}
	// fmt.Println(len(path))
	// testAlphabet()
	// testTextToColumns("")
	// testListCombos()
	// testNormalize()
	// testSymStrGroupCombinations("")
	// testContextualChangeCombos()
	// testApplyContextualChange("")

	// v11 := Vocabulary([]SymStr{NewSymStr("a"), NewSymStr("b"), NewSymStr("c")})
	// v12 := Vocabulary([]SymStr{NewSymStr("a"), NewSymStr("b"), NewSymStr("c")})
	// fmt.Println(v11.equals(v12))
	// fmt.Println(v11)
	// fmt.Println(v12)

	// s2 := "d͡ʒe̯"
	// for i, w := 0, 0; i < len(s2); i += w {
	// 	runeValue, width := utf8.DecodeRuneInString(s2[i:])
	// 	fmt.Printf("%#U starts at byte position %d\n", runeValue, i)
	// 	w = width
	// }
}

func testContextualChangeCombos() {
	abc := NewAlphabet([]string{}, map[string][]string{
		"V": []string{"a", "e", "i", "o", "u"},
		"C": []string{"b", "c", "d", "f", "g"},
	})
	s := "aCrVq"
	fmt.Printf("%s -> %s\n", s, abc.stringToRegex(s))
	cc := NewContextualChange("l", "r", "^ce", "um$")
	combos := abc.getContextulChangeCombinations(&cc, 3, 3)
	for _, combo := range combos {
		fmt.Println(combo)
	}
}

func testApplyContextualChange(word string) {
	if word == "" {
		word = "celum"
	}
	abc := NewAlphabet([]string{}, map[string][]string{
		"V": []string{"a", "e", "i", "o", "u"},
		"C": []string{"b", "c", "d", "f", "g"},
	})
	cc := abc.CompiledContextualChange("l", "r", "V", "V")
	new_str := abc.applyChange(NewSymStr(word), &cc)
	fmt.Printf("%s -> %s\n", word, new_str)
}

func testChangeSequences(s1 string, s2 string) {
	if s1 == "" && s2 == "" {
		s1, s2 = "cajlum", "cer"
	}
	w1, w2 := NewSymStr(s1), NewSymStr(s2)
	c_seqs0 := word_pair_change_sequences(w1, w2)
	for _, seq := range c_seqs0 {
		fmt.Println(seq)
		c_context := seq.getChangesInContext()
		fmt.Printf("\t%s\n", strings.Join(listComprehension(c_context, func(cc ContextualChange) string { return cc.String() }), "; "))
		versions := changeSequenceGetVersions(seq)
		for _, version := range versions {
			fmt.Println(version)
			c_context := version.getChangesInContext()
			fmt.Printf("\t%s\n", strings.Join(listComprehension(c_context, func(cc ContextualChange) string { return cc.String() }), "; "))
		}
		fmt.Println()
	}
}
