package main

import (
	"errors"
	"fmt"
	"regexp"
)

type ContextualChange struct {
	s_in     string
	s_out    string
	pre      string
	post     string
	compiled *regexp.Regexp
}

func NewContextualChange(s_in, s_out, pre, post string) ContextualChange {
	cc := ContextualChange{s_in, s_out, pre, post, nil}
	return cc
}
func (cc *ContextualChange) compile() {
	// regex_text := fmt.Sprintf("(%s)%s(%s)", cc.pre, cc.s_in, cc.post)
	regex_text := "(" + cc.pre + ")" + cc.s_in + "(" + cc.post + ")"
	cc.compiled = regexp.MustCompile(regex_text)

}

func (cc *ContextualChange) applyString(initial SymStr) string {
	// return cc.compiled.ReplaceAllString(initial.String(), fmt.Sprintf("${1}%s${2}", cc.s_out))
	return cc.compiled.ReplaceAllString(initial.String(), "${1}"+cc.s_out+"${2}")

}

func (cc *ContextualChange) apply(initial SymStr) SymStr {
	return NewSymStr(cc.applyString(initial))
}

func (cc ContextualChange) String() string {
	IO := [2]string{cc.s_in, cc.s_out}
	for i, s := range IO {
		if s == "" {
			IO[i] = "0"
		}
	}
	return fmt.Sprintf("%s > %s / %s _ %s", IO[0], IO[1], cc.pre, cc.post)
}

func ContextualChangeFromString(text string) (ContextualChange, error) {
	// form "A > B / C0 _ C1
	var in, out, pre, post string
	var count int
	var err error

	if count, err = fmt.Sscanf(text, "%s > %s / %s _ %s", &in, &out, &pre, &post); count == 4 && err == nil {

	} else if count, err = fmt.Sscanf(text, "%s > %s / _ %s", &in, &out, &post); count == 3 && err == nil {
		pre = ""
	} else if count, err = fmt.Sscanf(text, "%s > %s / %s _", &in, &out, &pre); count == 3 && err == nil {
		post = ""
	} else if count, err = fmt.Sscanf(text, "%s > %s", &in, &out); count == 2 && err == nil {
		pre, post = "", ""
	} else {
		return ContextualChange{"", "", "", "", nil}, errors.New("Bad format " + text)
	}
	if in == "0" {
		in = ""
	}
	if out == "0" {
		out = ""
	}
	return ContextualChange{in, out, pre, post, nil}, nil
}
