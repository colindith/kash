package cmd_reader

import (
	"fmt"
	"testing"
)

func Test_VirtualTerm(t *testing.T) {
	vt.printVirtualTerm(rune(1), rune(2), rune(3))
	fmt.Println(vt)  // {[1 2 3 32 32 32 32 32 32 32 32 32 32 32 32 32] 3}
	vt.printVirtualTerm(rune(4))
	fmt.Println(vt)  // {[1 2 3 32 32 32 32 32 32 32 32 32 32 32 32 32] 3}
	vt.printVirtualTerm('\b')
	fmt.Println(vt)  // {[1 2 3 32 32 32 32 32 32 32 32 32 32 32 32 32] 3}
	vt.printVirtualTerm("123")
	fmt.Println(vt)  // {[1 2 3 32 32 32 32 32 32 32 32 32 32 32 32 32] 3}
	vt.printVirtualTerm("\b\b\b\b\b\b\b\b")
	fmt.Println(vt)  // {[1 2 3 32 32 32 32 32 32 32 32 32 32 32 32 32] 3}
}

