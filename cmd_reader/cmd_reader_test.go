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