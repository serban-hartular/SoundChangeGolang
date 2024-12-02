package main

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type Alphabet struct {
	symbols      []string
	groups       map[string][]string // list of symbols that belong to group 'key'
	sym2groupMap map[string][]string
	symbolsByLen map[int]Set[string]
	maxSymbolLen int
}

func NewAlphabet(symbols []string, groups map[string][]string) Alphabet {
	alph := Alphabet{symbols, groups, make(map[string][]string), make(map[int]Set[string]), 0}
	//initialize groups
	for groupName, symList := range alph.groups {
		for _, symbol := range symList {
			alph.sym2groupMap[symbol] = append(alph.sym2groupMap[symbol], groupName)
		}
	}
	for _, symbol := range alph.symbols {
		s_len := len(symbol)
		if len_set, ok := alph.symbolsByLen[s_len]; ok {
			//len_set[symbol] = emptyStruct
			len_set.put(symbol)
		} else {
			new_set := make(Set[string])
			new_set.put(symbol)
			alph.symbolsByLen[s_len] = new_set
		}
		// alph.symbolsByLen[s_len][symbol] = emptyStruct
		if s_len > alph.maxSymbolLen {
			alph.maxSymbolLen = s_len
		}
	}
	return alph
}

func (abc *Alphabet) SymStr(s string) SymStr {
	if s == "" {
		return []string{""}
	}
	s = norm.NFC.String(s)
	symStr := make(SymStr, 0, len(s))
	for s != "" {
		s_len := abc.maxSymbolLen
		// process symbols by checking (by length) what length fits start of string
		for ; s_len > 1; s_len-- {
			if s_len > len(s) {
				continue
			}
			key := s[:s_len]
			if abc.symbolsByLen[s_len].contains(key) {
				symStr = append(symStr, key)
				s = s[s_len:]
				break
			}
		}
		if s_len <= 1 {
			key := s[:1]
			symStr = append(symStr, key)
			s = s[1:]
		}
	}
	return symStr
}

func (abc *Alphabet) validateSymStr(ss SymStr) bool {
	for i, s := range ss {

	}
}

func (abc *Alphabet) stringToRegex(s string) string {
	regex_string := ""
	for _, r := range s {
		c := string(r)
		if group, ok := abc.groups[c]; ok {
			regex_string += ("[" + strings.Join(group, "") + "]")
		} else {
			regex_string += c
		}
	}
	return regex_string
}

func (abc *Alphabet) CompiledContextualChange(s_in, s_out, pre, post string) ContextualChange {
	cc := ContextualChange{s_in, s_out, pre, post, nil}
	// regex_text := fmt.Sprintf("(%s)%s(%s)", abc.stringToRegex(pre), cc.s_in, abc.stringToRegex(post))
	regex_text := "(" + abc.stringToRegex(pre) + ")" + cc.s_in + "(" + abc.stringToRegex(post) + ")"
	var err error
	cc.compiled, err = regexp.Compile(regex_text)
	if err != nil {
		panic("Cannot compile regex '" + regex_text + "'")
	}
	return cc
}

func (abc *Alphabet) ContextualChangeFromString(text string) ContextualChange {
	cc, err := ContextualChangeFromString(text)
	if err == nil {
		cc = abc.CompiledContextualChange(cc.s_in, cc.s_out, cc.pre, cc.post)
	}
	return cc
}

func (abc *Alphabet) compileChange(cc *ContextualChange) {
	new_pre := abc.stringToRegex(cc.pre)
	new_post := abc.stringToRegex(cc.post)
	regex_text := fmt.Sprintf("(%s)%s(%s)", new_pre, cc.s_in, new_post)
	cc.compiled = regexp.MustCompile(regex_text)
	// cc.pre = new_pre
	// cc.post = new_post
}

func (abc *Alphabet) applyChange(ss SymStr, cc *ContextualChange) SymStr {
	return abc.SymStr(cc.applyString(ss))
}

func (abc *Alphabet) changedVocabulary(v Vocabulary, cc *ContextualChange) Vocabulary {
	new_vocab := make(Vocabulary, len(v))
	for i, word := range v {
		new_vocab[i] = abc.applyChange(word, cc)
	}
	return new_vocab
}

func (abc *Alphabet) getSymStrGroupCombinations(ss SymStr) [][]string {
	options := make([][]string, len(ss))
	// for each symbol in ss, extract options (based on groups it belongs to)
	for i, c := range ss {
		substitutions := []string{c}
		if group, ok := abc.sym2groupMap[c]; ok {
			substitutions = append(substitutions, group...)
		}
		options[i] = substitutions
	}
	return ListProduct(options)
}

func (abc *Alphabet) getContextulChangeCombinations(cc *ContextualChange, pre_len int, post_len int) []ContextualChange {
	if pre_len < 0 {
		pre_len = len(cc.pre)
	} else {
		pre_len = min(pre_len, len(cc.pre))
	}
	if post_len < 0 {
		post_len = len(cc.post)
	} else {
		post_len = min(post_len, len(cc.post))
	}
	combos := make([]ContextualChange, 0)
	for i := range pre_len + 1 {
		for j := range post_len + 1 {
			combos = append(combos, abc.getContextulChangeCombinationsLen(cc, i, j)...)
		}
	}
	return combos
}

func (abc *Alphabet) generalityScoreString(regex_string string) float64 {
	if regex_string == "" {
		return 0.0
	}
	score := 0.0
	ss := abc.SymStr(regex_string)
	for _, c := range ss {
		if _, ok := abc.groups[c]; ok { // is a group
			score += 0.5
		} else {
			score += 1
		}
	}
	return score
}

func (abc *Alphabet) generalityScoreChange(cc *ContextualChange) float64 {
	return abc.generalityScoreString(cc.pre) + abc.generalityScoreString(cc.post)
}

func (abc *Alphabet) getContextulChangeCombinationsLen(cc *ContextualChange, pre_len int, post_len int) []ContextualChange {
	if pre_len < 0 {
		pre_len = len(cc.pre)
	} else {
		pre_len = min(pre_len, len(cc.pre))
	}
	if post_len < 0 {
		post_len = len(cc.post)
	} else {
		post_len = min(post_len, len(cc.post))
	}
	post_str := cc.post[:post_len]
	pre_str := cc.pre[len(cc.pre)-pre_len:]
	pre_combos := abc.getSymStrGroupCombinations(abc.SymStr(pre_str))
	post_combos := abc.getSymStrGroupCombinations(abc.SymStr(post_str))
	combo_list := make([]ContextualChange, len(pre_combos)*len(post_combos))
	i := 0
	for _, pre := range pre_combos {
		for _, post := range post_combos {
			combo_list[i] = abc.CompiledContextualChange(cc.s_in, cc.s_out,
				strings.Join(pre, ""), strings.Join(post, ""))
			// combo_list[i] = ContextualChange{cc.s_in, cc.s_out,
			// 	SymStrConcat(pre).String(), SymStrConcat(post).String(), nil}
			// abc.compileChange(&combo_list[i])
			i++
		}
	}
	return combo_list
}

func testAlphabet() {
	abc := NewAlphabet([]string{"a", "a00", "a01", "a10", "a11"}, make(map[string][]string))
	input := "abcdefga0bcdefa01wera111a1ida10"
	symStr := abc.SymStr(input)
	for _, sym := range symStr {
		fmt.Printf("'%s', ", sym)
	}
	fmt.Println()

}

func testNormalize() {
	s := "/bat͡ʃʲ/"
	abc := NewAlphabet([]string{"ʲ", "t͡ʃ"}, make(map[string][]string))
	symStr := abc.SymStr(s)
	for _, c := range symStr {
		fmt.Println(c)
	}
}

func testSymStrGroupCombinations(s string) {
	abc := NewAlphabet([]string{"a", "b", "c", "d", "e"},
		map[string][]string{
			"V": {"a", "e", "i", "o", "u"},
			"C": {"b", "c", "d", "f", "g"},
			"B": {"b"},
		},
	)
	if s == "" {
		s = "ab1ed"
	}
	ss := NewSymStr(s)
	fmt.Println(ss)
	combos := abc.getSymStrGroupCombinations(ss)
	for i, combo := range combos {
		fmt.Printf("%d\t%s\n", i, combo)
	}
}
