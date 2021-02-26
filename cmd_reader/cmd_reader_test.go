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

func Test_cmdLine_backSpace(t *testing.T) {
	tests := []struct {
		name     string
		clBefore *cmdLine
		clExpected  *cmdLine
	}{
		{
			"left_end",
			&cmdLine{
				buf: []rune{1,2,3,4,5,32},
				tmp: nil,
				ptr: 0,
			},
			&cmdLine{
				buf: []rune{1,2,3,4,5,32},
				tmp: nil,
				ptr: 0,
			},
		},
		{
			"right_end",
			&cmdLine{
				buf: []rune{1,2,3,4,5,32},
				tmp: nil,
				ptr: 5,
			},
			&cmdLine{
				buf: []rune{1,2,3,4,32},
				tmp: nil,
				ptr: 4,
			},
		},
		{
			"middle",
			&cmdLine{
				buf: []rune{1,2,3,4,5,32},
				tmp: nil,
				ptr: 3,
			},
			&cmdLine{
				buf: []rune{1,2,4,5,32},
				tmp: nil,
				ptr: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := tt.clBefore
			cl.backSpace()
			if !reflect.DeepEqual(cl, tt.clExpected) {
				t.Errorf("back_space_not_correct | want=%v | got=%v", tt.clExpected, cl)
			}
		})
	}
}