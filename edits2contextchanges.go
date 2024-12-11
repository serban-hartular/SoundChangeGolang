package main

func editSequence2ContextualChanges(e_seq EditSequence) []ContextualChange {
	original_str := listComprehension(e_seq, func(edit Edit) SymStr { return edit.s_in })
	cc_list := make([]ContextualChange, 0)
	for i, edit := range e_seq {
		if !edit.NoChange() { //it's a change
			pre := SymStrConcat(original_str[:i]...)
			pre = filter(pre, func(s string) bool { return len(s) > 0 })
			post := SymStrConcat(original_str[i+1:]...)
			post = filter(post, func(s string) bool { return len(s) > 0 })
			cc := NewContextualChange(edit.s_in, edit.s_out, pre, post)
			cc_list = append(cc_list, cc)
		}
	}
	return cc_list
}
