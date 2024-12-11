package main

import (
	"fmt"
	"regexp"
)

type ContextualChange struct {
	s_in     SymStr
	s_out    SymStr
	pre      SymStr
	post     SymStr
	compiled *regexp.Regexp
}

func NewContextualChange(s_in, s_out, pre, post SymStr) ContextualChange {
	cc := ContextualChange{s_in, s_out, pre, post, nil}
	return cc
}
func (cc *ContextualChange) compile() {
	// regex_text := fmt.Sprintf("(%s)%s(%s)", cc.pre, cc.s_in, cc.post)
	regex_text := "(" + cc.pre.String() + ") " + cc.s_in.String() + " (" + cc.post.String() + ")"
	cc.compiled = regexp.MustCompile(regex_text)

}

func (cc *ContextualChange) applyString(initial SymStr) SymStr {
	// return cc.compiled.ReplaceAllString(initial.String(), fmt.Sprintf("${1}%s${2}", cc.s_out))
	return SymStrFromString(cc.compiled.ReplaceAllString(initial.String(), "${1} "+cc.s_out.String()+" ${2}"))

}

func (cc ContextualChange) String() string {
	IO := [2]string{cc.s_in.String(), cc.s_out.String()}
	for i, s := range IO {
		if s == "" {
			IO[i] = "0"
		}
	}
	return fmt.Sprintf("%s > %s / %s _ %s", IO[0], IO[1], cc.pre, cc.post)
}

// func ContextualChangeFromString(text string) (ContextualChange, error) {
// 	// form "A > B / C0 _ C1
// 	var in, out, pre, post string
// 	var count int
// 	var err error

// 	if count, err = fmt.Sscanf(text, "%s > %s / %s _ %s", &in, &out, &pre, &post); count == 4 && err == nil {

// 	} else if count, err = fmt.Sscanf(text, "%s > %s / _ %s", &in, &out, &post); count == 3 && err == nil {
// 		pre = ""
// 	} else if count, err = fmt.Sscanf(text, "%s > %s / %s _", &in, &out, &pre); count == 3 && err == nil {
// 		post = ""
// 	} else if count, err = fmt.Sscanf(text, "%s > %s", &in, &out); count == 2 && err == nil {
// 		pre, post = "", ""
// 	} else {
// 		return ContextualChange{EmptySymStr(), EmptySymStr(), EmptySymStr(), EmptySymStr(), nil}, errors.New("Bad format " + text)
// 	}
// 	if in == "0" {
// 		in = ""
// 	}
// 	if out == "0" {
// 		out = ""
// 	}
// 	return ContextualChange{NewSymStr(in), NewSymStr(out), NewSymStr(pre), NewSymStr(post), nil}, nil
// }
