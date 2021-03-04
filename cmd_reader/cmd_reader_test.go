package cmd_reader

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[1;36m%s\033[0m"
)


func Test_Color(t *testing.T) {
	fmt.Printf("\033[1;47m%s\033[0m", "Info")
	fmt.Println("")
	fmt.Printf("\033[88;1m%s\033[0m", "Notice")
	fmt.Println("")
	fmt.Printf("\033[33;5m%s\033[0m", "Warning")
	fmt.Println("")
	fmt.Printf(ErrorColor, "Error")
	fmt.Println("")
	fmt.Printf(DebugColor, "Debug")
	fmt.Println("")
}

func Test_cmdLine_moveRight(t *testing.T) {
	debugFlag = true

	tests := []struct {
		name       string
		clBefore   cmdLine
		clExpected cmdLine
		vtBefore   virtualTerm
		vtExpected virtualTerm
	}{
		{
			"left_end",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 0,
			},
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 1,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 0,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 1,
			},
		},
		{
			"right_end",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 5,
			},
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 5,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 5,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 5,
			},
		},
		{
			"middle",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 3,
			},
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 4,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 3,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := tt.clBefore
			vt = tt.vtBefore
			cl.moveRight()
			if !reflect.DeepEqual(cl, tt.clExpected) {
				t.Errorf("move_right_cmd_line_not_correct | want=%v | got=%v", tt.clExpected, cl)
			}
			if !reflect.DeepEqual(vt, tt.vtExpected) {
				t.Errorf("move_right_virtual_term_not_correct | want=%v | got=%v", tt.vtExpected, vt)
			}
		})
	}
}

func Test_cmdLine_insertChar(t *testing.T) {
	debugFlag = true

	tests := []struct {
		name       string
		clBefore   cmdLine
		clExpected cmdLine
		vtBefore   virtualTerm
		vtExpected virtualTerm
	}{
		{
			"left_end",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 0,
			},
			cmdLine{
				buf: []rune{0, 1, 2, 3, 4, 5, 32},
				ptr: 1,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 0,
			},
			virtualTerm{
				buf: []rune{0, 1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 1,
			},
		},
		{
			"right_end",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 5,
			},
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 0, 32},
				ptr: 6,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 5,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 0, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 6,
			},
		},
		{
			"middle",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 3,
			},
			cmdLine{
				buf: []rune{1, 2, 3, 0, 4, 5, 32},
				ptr: 4,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 3,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 0, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := tt.clBefore
			vt = tt.vtBefore
			cl.insertChar(0)
			if !reflect.DeepEqual(cl, tt.clExpected) {
				t.Errorf("insert_char_cmd_line_not_correct | want=%v | got=%v", tt.clExpected, cl)
			}
			if !reflect.DeepEqual(vt, tt.vtExpected) {
				t.Errorf("insert_char_virtual_term_not_correct | want=%v | got=%v", tt.vtExpected, vt)
			}
		})
	}
}

func Test_cmdLine_backSpace(t *testing.T) {
	debugFlag = true

	tests := []struct {
		name       string
		clBefore   cmdLine
		clExpected cmdLine
		vtBefore   virtualTerm
		vtExpected virtualTerm
	}{
		{
			"left_end",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 0,
			},
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 0,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 0,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 0,
			},
		},
		{
			"right_end",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 5,
			},
			cmdLine{
				buf: []rune{1, 2, 3, 4, 32},
				ptr: 4,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 5,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 4,
			},
		},
		{
			"middle",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 3,
			},
			cmdLine{
				buf: []rune{1, 2, 4, 5, 32},
				ptr: 2,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 3,
			},
			virtualTerm{
				buf: []rune{1, 2, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := tt.clBefore
			vt = tt.vtBefore
			cl.backSpace()
			if !reflect.DeepEqual(cl, tt.clExpected) {
				t.Errorf("back_space_cmd_line_not_correct | want=%v | got=%v", tt.clExpected, cl)
			}
			if !reflect.DeepEqual(vt, tt.vtExpected) {
				t.Errorf("back_space_virtual_term_not_correct | want=%v | got=%v", tt.vtExpected, vt)
			}
		})
	}
}

func Test_cmdLine_searchHistoryUp(t *testing.T) {
	debugFlag = true

	node1 := &node{val: []rune{1, 2, 3}}
	node2 := &node{val: []rune{4, 5, 6}}
	node3 := &node{val: []rune{1, 2, 3, 4, 5}}
	node4 := &node{val: []rune{1, 2, 3}}

	node1.next = node2
	node2.next = node3
	node3.next = node4
	node2.prev = node1
	node3.prev = node2
	node4.prev = node3

	tests := []struct {
		name       string
		clBefore   cmdLine
		clExpected cmdLine
		vtBefore   virtualTerm
		vtExpected virtualTerm
	}{
		{
			"zero to 1st level",
			cmdLine{
				buf: []rune{1, 2, 32},
				ptr: 2,
				history: history{
					head: node1,
					len:  4,
					ptr:  nil,
					match: []rune{1, 2},
				},
			},
			cmdLine{
				buf: []rune{1, 2, 3, 32},
				ptr: 3,
				history: history{
					head: node1,
					len:  4,
					ptr:  node1,
					match: []rune{1, 2},
				},
				tmp: []rune{1, 2},
			},
			virtualTerm{
				buf: []rune{1, 2, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 2,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 3,
			},
		},
		{
			"1st to 3rd level",
			cmdLine{
				buf: []rune{1, 2, 3, 32},
				ptr: 3,
				history: history{
					head: node1,
					len:  4,
					ptr:  node1,
					match: []rune{1, 2},
				},
			},
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 5,
				history: history{
					head: node1,
					len:  4,
					ptr:  node3,
					match: []rune{1, 2},
				},
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 3,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 5,
			},
		},
		{
			"tail",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 5,
			},
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 5,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 5,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 5,
			},
		},
		{
			"3rd to 4th level",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 5,
				history: history{
					head: node1,
					len:  4,
					ptr:  node3,
					match: []rune{1, 2},
				},
			},
			cmdLine{
				buf: []rune{1, 2, 3, 32},
				ptr: 3,
				history: history{
					head: node1,
					len:  4,
					ptr:  node4,
					match: []rune{1, 2},
				},
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 5,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := tt.clBefore
			vt = tt.vtBefore
			cl.searchHistoryUp()
			if !reflect.DeepEqual(cl, tt.clExpected) {
				t.Errorf("search_history_up_cmd_line_not_correct | want=%v | got=%v", tt.clExpected, cl)
			}
			if !reflect.DeepEqual(vt, tt.vtExpected) {
				t.Errorf("search_history_up_virtual_term_not_correct | want=%v | got=%v", tt.vtExpected, vt)
			}
		})
	}
}

func Test_cmdLine_searchHistoryDown(t *testing.T) {
	debugFlag = true

	node1 := &node{val: []rune{1, 2, 3}}
	node2 := &node{val: []rune{4, 5, 6}}
	node3 := &node{val: []rune{1, 2, 3, 4, 5}}

	node1.next = node2
	node2.next = node3
	node2.prev = node1
	node3.prev = node2

	tests := []struct {
		name       string
		clBefore   cmdLine
		clExpected cmdLine
		vtBefore   virtualTerm
		vtExpected virtualTerm
	}{
		{
			"1st to zero level",
			cmdLine{
				buf: []rune{1, 2, 3, 32},
				ptr: 3,
				history: history{
					head: node1,
					len:  3,
					ptr:  node1,
					match: []rune{1, 2},
				},
				tmp: []rune{1, 2},
			},
			cmdLine{
				buf: []rune{1, 2, 32},
				ptr: 2,
				history: history{
					head: node1,
					len:  3,
					ptr:  nil,
					match: []rune{1, 2},
				},
				tmp: []rune{1, 2},
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 3,
			},
			virtualTerm{
				buf: []rune{1, 2, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 2,
			},
		},
		{
			"3rd to 1st level",
			cmdLine{
				buf: []rune{1, 2, 3, 4, 5, 32},
				ptr: 5,
				history: history{
					head: node1,
					len:  3,
					ptr:  node3,
					match: []rune{1, 2},
				},
				tmp: []rune{1, 2},
			},
			cmdLine{
				buf: []rune{1, 2, 3, 32},
				ptr: 3,
				history: history{
					head: node1,
					len:  3,
					ptr:  node1,
					match: []rune{1, 2},
				},
				tmp: []rune{1, 2},
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 4, 5, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 5,
			},
			virtualTerm{
				buf: []rune{1, 2, 3, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 3,
			},
		},
		{
			"zero to zero",
			cmdLine{
				buf: []rune{1, 2, 32},
				ptr: 2,
			},
			cmdLine{
				buf: []rune{1, 2, 32},
				ptr: 2,
			},
			virtualTerm{
				buf: []rune{1, 2, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 2,
			},
			virtualTerm{
				buf: []rune{1, 2, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
				ptr: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := tt.clBefore
			vt = tt.vtBefore
			cl.searchHistoryDown()
			if !reflect.DeepEqual(cl, tt.clExpected) {
				t.Errorf("search_history_down_cmd_line_not_correct | want=%v | got=%v", tt.clExpected, cl)
			}
			if !reflect.DeepEqual(vt, tt.vtExpected) {
				t.Errorf("search_history_down_virtual_term_not_correct | want=%v | got=%v", tt.vtExpected, vt)
			}
		})
	}
}

func Test_checkMatch(t *testing.T) {
	tests := []struct{
		name string
		ptr  *node
		s    []rune
		want bool
	}{
		{
			name: "nil_node",
			ptr:  nil,
			s:    []rune{1, 2},
			want: false,
		},
		{
			name: "short_match",
			ptr:  &node{val: []rune{1, 2, 3}},
			s:    []rune{1, 2},
			want: true,
		},
		{
			name: "short_not_match",
			ptr:  &node{val: []rune{1, 2, 3}},
			s:    []rune{1, 4},
			want: false,
		},
		{
			name: "equal_not_match",
			ptr:  &node{val: []rune{1, 2, 3}},
			s:    []rune{1, 2, 3},
			want: false,
		},
		{
			name: "long_not_match",
			ptr:  &node{val: []rune{1, 2, 3}},
			s:    []rune{1, 2, 3, 4},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if checkMatch(tt.ptr, tt.s) != tt.want {
				t.Errorf("check_match_not_correct | want=%v | got=%v", tt.want, !tt.want)
			}
		})}
}

func Test_history_searchUp(t *testing.T) {
	node1 := &node{val: []rune{1, 2, 3}}
	node2 := &node{val: []rune{4, 5, 6}}
	node3 := &node{val: []rune{1, 2, 3, 4, 5}}

	node1.next = node2
	node2.next = node3
	node2.prev = node1
	node3.prev = node2

	tests := []struct{
		name         string
		hBefore      *history
		hAfter       *history
		result       []rune
		foundResult  bool
		preserveZero bool
	}{
		{
			"empty history",
			&history{
				head: nil,
				len:  0,
				ptr:  nil,
				match: []rune{1, 2},
			},
			&history{
				head: nil,
				len:  0,
				ptr:  nil,
				match: []rune{1, 2},
			},
			nil,
			false,
			false,
		},
		{
			"zero to 1st level",
			&history{
				head: node1,
				len:  3,
				ptr:  nil,
				match: []rune{1, 2},
			},
			&history{
				head: node1,
				len:  3,
				ptr:  node1,
				match: []rune{1, 2},
			},
			[]rune{1, 2, 3},
			true,
			true,
		},
		{
			"1st to 3rd level",
			&history{
				head: node1,
				len:  3,
				ptr:  node1,
				match: []rune{1, 2},
			},
			&history{
				head: node1,
				len:  3,
				ptr:  node3,
				match: []rune{1, 2},
			},
			[]rune{1, 2, 3, 4, 5},
			true,
			false,
		},
		{
			"3rd level",
			&history{
				head: node1,
				len:  3,
				ptr:  node3,
				match: []rune{1, 2},
			},
			&history{
				head: node1,
				len:  3,
				ptr:  node3,
				match: []rune{1, 2},
			},
			nil,
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.hBefore
			result, foundResult, preserveZero := h.searchUp()
			if !reflect.DeepEqual(h, tt.hAfter) {
				t.Errorf("h_after_not_correct | want=%v | got=%v", *tt.hAfter, *h)
			}
			if !reflect.DeepEqual(result, tt.result) {
				t.Errorf("result_not_correct | want=%v | got=%v", tt.result, result)
			}
			if foundResult != tt.foundResult {
				t.Errorf("foundResult_not_correct | want=%v | got=%v", tt.foundResult, foundResult)
			}
			if preserveZero != tt.preserveZero {
				t.Errorf("preserveZero_not_correct | want=%v | got=%v", tt.preserveZero, preserveZero)
			}
		})
	}
}


func Test_history_searchDown(t *testing.T) {
	node1 := &node{val: []rune{1, 2, 3}}
	node2 := &node{val: []rune{4, 5, 6}}
	node3 := &node{val: []rune{1, 2, 3, 4, 5}}

	node1.next = node2
	node2.next = node3
	node2.prev = node1
	node3.prev = node2

	tests := []struct{
		name         string
		hBefore      *history
		hAfter       *history
		result       []rune
		foundResult  bool
		retrieveZero bool
	}{
		{
			"empty history",
			&history{
				head: nil,
				len:  0,
				ptr:  nil,
				match: []rune{1, 2},
			},
			&history{
				head: nil,
				len:  0,
				ptr:  nil,
				match: []rune{1, 2},
			},
			nil,
			false,
			false,
		},
		{
			"1st to zero level",
			&history{
				head: node1,
				len:  3,
				ptr:  node1,
				match: []rune{1, 2},
			},
			&history{
				head: node1,
				len:  3,
				ptr:  nil,
				match: []rune{1, 2},
			},
			nil,
			false,
			true,
		},
		{
			"3rd to 1st level",
			&history{
				head: node1,
				len:  3,
				ptr:  node3,
				match: []rune{1, 2},
			},
			&history{
				head: node1,
				len:  3,
				ptr:  node1,
				match: []rune{1, 2},
			},
			[]rune{1, 2, 3},
			true,
			false,
		},
		{
			"zero not empty",
			&history{
				head: node1,
				len:  3,
				ptr:  nil,
				match: []rune{1, 2},
			},
			&history{
				head: node1,
				len:  3,
				ptr:  nil,
				match: []rune{1, 2},
			},
			nil,
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.hBefore
			result, foundResult, preserveZero := h.searchDown()
			if !reflect.DeepEqual(h, tt.hAfter) {
				t.Errorf("h_after_not_correct | want=%v | got=%v", *tt.hAfter, *h)
			}
			if !reflect.DeepEqual(result, tt.result) {
				t.Errorf("result_not_correct | want=%v | got=%v", tt.result, result)
			}
			if foundResult != tt.foundResult {
				t.Errorf("foundResult_not_correct | want=%v | got=%v", tt.foundResult, foundResult)
			}
			if preserveZero != tt.retrieveZero {
				t.Errorf("retrieveZero_not_correct | want=%v | got=%v", tt.retrieveZero, preserveZero)
			}
		})
	}
}

func Test_history_pushFront(t *testing.T) {
	node1 := &node{val: []rune{1, 2, 3}}
	node2 := &node{val: []rune{4, 5, 6}}
	node3 := &node{val: []rune{1, 2, 3, 4, 5}}

	node1.next = node2
	node2.next = node3
	node2.prev = node1
	node3.prev = node2

	tests := []struct{
		name         string
		hBefore      *history
		hAfter       *history
		buf          []rune
	}{
		{
			"empty history",
			&history{
				head: nil,
				len:  0,
				ptr:  nil,
			},
			&history{
				head: &node{
					next: nil,
					prev: nil,
					val: []rune{3, 3, 3, 3},
				},
				len:  1,
				ptr:  nil,
			},
			[]rune{3, 3, 3, 3},
		},
		{
			"not empty history",
			&history{
				head: node1,
				len:  3,
				ptr:  nil,
				match: []rune{1, 2},
			},
			&history{
				head: &node{
					next: node1,
					prev: nil,
					val: []rune{3, 3, 3, 3},
				},
				len:  4,
				ptr:  nil,
				match: []rune{1, 2},
			},
			[]rune{3, 3, 3, 3},
		},
		{
			"not empty history and ptr not nil",
			&history{
				head: node1,
				len:  3,
				ptr:  node2,
				match: []rune{1, 2},
			},
			&history{
				head: &node{
					next: node1,
					prev: nil,
					val: []rune{3, 3, 3, 3},
				},
				len:  4,
				ptr:  node2,
				match: []rune{1, 2},
			},
			[]rune{3, 3, 3, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.hBefore
			h.pushFront(tt.buf)
			if !reflect.DeepEqual(h, tt.hAfter) {
				t.Errorf("h_after_not_correct | want=%v | got=%v", *tt.hAfter, *h)
			}
		})
	}
}