package main

import (
	"strings"
)

// type Symbol string

type SymStr []string

func NewSymStr(s string) SymStr {
	if s == "" {
		return []string{""}
	}
	symstr := make(SymStr, len(s))
	for i, c := range []rune(s) {
		symstr[i] = string(c)
	}
	return symstr
}

func (symstr SymStr) String() string {
	return strings.Join(symstr, "")
}

func (symstr SymStr) Empty() bool {
	return len(symstr) == 0
}

func (symstr SymStr) CopySlice(start int, end int) SymStr {
	var dest SymStr = make(SymStr, end-start)
	copy(dest, symstr[start:end])
	return dest
}
