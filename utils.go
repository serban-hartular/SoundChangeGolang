package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

func filter[T any](input []T, test func(T) bool) []T {
	output := make([]T, len(input))
	out_len := 0
	for _, item := range input {
		if test(item) {
			output[out_len] = item
			out_len++
		}
	}
	return output[:out_len]
}

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

var _empty struct{}

func (set Set[T]) put(item T) {
	set[item] = _empty
}

func (set Set[T]) putList(items []T) {
	for _, item := range items {
		set.put(item)
	}
}

func NewSetFromList[T comparable](items []T) Set[T] {
	set := NewSet[T]()
	set.putList(items)
	return set
}

func (set Set[T]) contains(item T) bool {
	_, ok := set[item]
	return ok
}

func (set Set[T]) pop(item T) {
	delete(set, item)
}

func (set Set[T]) toList() []T {
	return slices.Collect(maps.Keys(set))
}

func (this Set[T]) union(other Set[T]) Set[T] {
	u := NewSet[T]()
	for _, set := range []Set[T]{this, other} {
		for _, k := range set.toList() {
			u.put(k)
		}
	}
	return u
}

func (this Set[T]) difference(other Set[T]) Set[T] {
	u := NewSet[T]()
	for _, k := range this.toList() {
		if !other.contains(k) {
			u.put(k)
		}
	}
	return u
}

func (this Set[T]) intersection(other Set[T]) Set[T] {
	u := NewSet[T]()
	for _, k := range this.toList() {
		if other.contains(k) {
			u.put(k)
		}
	}
	return u
}

type ComparableStringer interface {
	comparable
	String() string
}

func Set2String[T ComparableStringer](set Set[T]) string {
	return "{" + strings.Join(listComprehension(set.toList(), func(item T) string { return item.String() }), ", ") + "}"
}

func InitMatrix[T any](num_rows, row_length int) [][]T {
	matrix := make([][]T, num_rows)
	for i := range matrix {
		matrix[i] = make([]T, row_length)
	}
	return matrix
}

func textToMatrix(text string) [][]string {
	if text[len(text)-1] == '\n' {
		text = text[:len(text)-1]
	}
	lines := strings.Split(text, "\n")
	rows := make([][]string, len(lines))
	for i, line := range lines {
		rows[i] = strings.Split(line, "\t")
	}
	return rows
}

func textToColumns(text string) [][]string {
	matrix := textToMatrix(text)
	num_columns := len(matrix[0])
	columns := InitMatrix[string](num_columns, len(matrix))
	for j := 0; j < num_columns; j++ {
		for i := range columns[j] {
			columns[j][i] = matrix[i][j]
		}
	}
	return columns
}

func ListProduct[T any](listOfLists [][]T) [][]T {
	if len(listOfLists) == 0 {
		return [][]T{}
	}
	if len(listOfLists) == 1 {
		r := make([][]T, len(listOfLists[0]))
		for i, item := range listOfLists[0] {
			r[i] = []T{item}
		}
		return r
	}
	current := listOfLists[0]
	prod_next := ListProduct(listOfLists[1:])
	// fmt.Printf("prod_next(%s) = %s\n", listOfLists[1:], prod_next)
	product := make([][]T, len(current)*len(prod_next))
	i := 0
	for _, c0 := range current {
		for _, next := range prod_next {
			product[i] = append([]T{c0}, next...)
			i++
		}
	}
	return product
}

func testTextToColumns(text string) {
	if text == "" {
		text = "abc\tdef\tghi\njklmn\top\tqrs\ntuv\twx\tyz\n123\t45\t67890"
	}
	text = text + "\n"
	fmt.Println(text)
	columns := textToColumns(text)
	for i, col := range columns {
		fmt.Printf("%d\t%s\n", i, col)
	}
}

func testListCombos() {
	combos0 := []string{"a", "V"}
	combos1 := [3]string{"b", "C", "P"}
	combos2 := []string{"$"}

	listOfLists := [][]string{combos0, combos1[:], combos2[:]}
	combos := ListProduct(listOfLists)
	for _, combo := range combos {
		fmt.Println(combo)
	}
}

func listComprehension[I any, O any](list_in []I, mapFn func(I) O) []O {
	list_out := make([]O, len(list_in))
	for i, v := range list_in {
		list_out[i] = mapFn(v)
	}
	return list_out
}

func load_file_as_string(filename string) string {
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(dat), "\r\n\t ")
}

func string2set(s string) Set[string] {
	set := NewSet[string]()
	for _, c := range s {
		set.put(string(c))
	}
	return set
}

func quick_alphabet(special_vowels string, special_consonants string) Alphabet {
	std_vowels := string2set("aeiouăîây").union(NewSetFromList([]string{"ā", "ē", "ī", "ō", "ū"}))
	//multichar_vowels := []string{"ā", "ē", "ī", "ō", "ū"}
	std_consonants := string2set("bcdfghklmnpqrstvxzwj")
	spec_vowel_set := string2set(special_vowels)
	spec_cons_set := string2set(special_consonants)
	vowel_set := std_vowels.union(spec_vowel_set).difference(spec_cons_set)
	consonant_set := std_consonants.union(spec_cons_set).intersection(spec_vowel_set)
	vowels := vowel_set.toList()
	consonants := consonant_set.toList()
	symbols := slices.Concat(vowels, consonants)
	return NewAlphabet(symbols, map[string][]string{"V": vowels, "C": consonants})
}
