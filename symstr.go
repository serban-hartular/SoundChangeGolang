package main

import (
	"slices"
	"strings"

	"golang.org/x/text/unicode/norm"
)

// type Symbol string

type SymStr []string

const SYMSEP = ' '

func SymStrFromString(s string) SymStr {
	s = norm.NFC.String(s)
	return strings.FieldsFunc(s, func(c rune) bool { return c == SYMSEP })
}

func SingleSymbol(s string) SymStr {
	s = norm.NFC.String(s)
	symstr := SymStr{s}
	return symstr
}

func EmptySymStr() SymStr {
	return SymStr{}
}

func (symstr SymStr) equals(other SymStr) bool {
	return symstr.String() == other.String()
}

func (symstr SymStr) String() string {
	return strings.Join(symstr, string(SYMSEP))
}

func (symstr SymStr) Tight() string {
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

func SymStrConcat(list ...SymStr) SymStr {
	return slices.Concat(list...)
}
