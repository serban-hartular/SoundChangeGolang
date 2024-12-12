package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type Alphabet struct {
	symbols      []string
	groups       map[string][]string // list of symbols that belong to group 'key'
	sym2groupMap map[string][]string //maps symbol to a list of groups
	symbolsByLen map[int]Set[string] //maps symbols by length (in runes)
	maxSymbolLen int
}

func NewAlphabet(symbols []string, groups map[string][]string) Alphabet {
	symbols = listComprehension(symbols, norm.NFC.String)
	for k := range groups {
		groups[k] = listComprehension(groups[k], norm.NFC.String)
	}
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
		for i := 0; i < alph.maxSymbolLen; i++ {
			if _, ok := alph.symbolsByLen[i]; !ok {
				alph.symbolsByLen[i] = NewSet[string]()
			}
		}
	}
	return alph
}

func (abc *Alphabet) SymStrFromTightString(s string) SymStr {
	s = norm.NFC.String(s)
	symstr := EmptySymStr()
	for i := 0; i < len(s); {
		var sym_len int
		for sym_len = abc.maxSymbolLen; sym_len > 0; sym_len-- {
			if i+sym_len > len(s) {
				continue
			}
			candidate_str := s[i : i+sym_len]
			if abc.symbolsByLen[sym_len].contains(candidate_str) {
				fmt.Printf("Contained at %d\n", sym_len)
				symstr = append(symstr, candidate_str)
				i += sym_len
				break
			}
		}
		if sym_len == 0 { // not found, just add the rune
			symstr = append(symstr, string(s[i]))
			i++
		}
	}
	return symstr
}

func (abc *Alphabet) NewContextualChange(s_in, s_out, pre, post SymStr) ContextualChange {
	pre = abc.toRegexSymstr(pre)
	post = abc.toRegexSymstr(post)
	cc := NewContextualChange(s_in, s_out, pre, post)
	cc.compile()

	return cc
}

func (abc *Alphabet) toRegexSymstr(ss SymStr) SymStr {
	r_ss := make(SymStr, len(ss))
	for i := range ss {
		symbol := ss[i]
		if group, ok := abc.groups[symbol]; ok {
			group_regex := "(?:" + strings.Join(group, "|") + ")" //non-capturing group
			r_ss[i] = group_regex
		} else {
			r_ss[i] = ss[i]
		}
	}
	return r_ss
}

func NewSimpleAlphabet(vowels_spaced string, consonants_spaced string) Alphabet {
	vowels_default := NewSetFromList(SymStrFromString("a e i o u ă î â"))
	consonants_default := NewSetFromList(SymStrFromString("b c d f g h i j k l m n p q r s t v w x y z ș ț"))
	vowels_in := NewSetFromList(SymStrFromString(vowels_spaced))
	consonants_in := NewSetFromList(SymStrFromString(consonants_spaced))
	vowels := vowels_default.union(vowels_in).difference(consonants_in).toList()
	consonants := consonants_default.union(consonants_in).difference(vowels_in).toList()
	all_symbols := SymStrConcat(vowels, consonants)
	return NewAlphabet(all_symbols, map[string][]string{"V": vowels, "C": consonants})
}

func (abc *Alphabet) getSymStrGroupCombos(ss SymStr) []SymStr {
	if len(ss) < 1 {
		return []SymStr{}
	}
	options := []string{ss[0]}
	if groups, ok := abc.sym2groupMap[ss[0]]; ok {
		options = append(options, groups...)
	}
	options_symstrs := listComprehension(options, SingleSymbol) //each a symbol string
	if len(ss) == 1 {
		return options_symstrs
	}
	combo_options := make([]SymStr, 0)
	for _, this_symbol := range options_symstrs {
		for _, next_combo := range abc.getSymStrGroupCombos(ss[1:]) {
			combo_options = append(combo_options, SymStrConcat(this_symbol, next_combo))
		}
	}
	return combo_options
}
