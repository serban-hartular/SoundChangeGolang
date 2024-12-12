package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

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

type StringerSet[T fmt.Stringer] map[string]T

func NewStringerSet[T fmt.Stringer]() StringerSet[T] {
	return make(StringerSet[T])
}

func (sset StringerSet[T]) put(item T) {
	sset[item.String()] = item
}

func (sset StringerSet[T]) putList(items []T) {
	for _, item := range items {
		sset.put(item)
	}
}

func NewStringerSetFromList[T fmt.Stringer](items []T) StringerSet[T] {
	sset := NewStringerSet[T]()
	sset.putList(items)
	return sset
}

func (sset StringerSet[T]) contains(item T) bool {
	_, ok := sset[item.String()]
	return ok
}

func (sset StringerSet[T]) getByString(stringRepresentation string) (T, bool) {
	v, ok := sset[stringRepresentation]
	return v, ok
}

func (sset StringerSet[T]) pop(item T) {
	delete(sset, item.String())
}

func (sset StringerSet[T]) popByString(stringRepresentation string) (T, bool) {
	v, ok := sset.getByString(stringRepresentation)
	if ok {
		delete(sset, stringRepresentation)
	}
	return v, ok
}

func (sset StringerSet[T]) toList() []T {
	return slices.Collect(maps.Values(sset))
}

func (sset StringerSet[T]) keys() Set[string] {
	return NewSetFromList[string](slices.Collect(maps.Keys(sset)))
}

func (sset StringerSet[T]) union(other StringerSet[T]) StringerSet[T] {
	u := NewStringerSet[T]()
	for _, set := range []StringerSet[T]{sset, other} {
		for _, k := range set.toList() {
			u.put(k)
		}
	}
	return u
}
