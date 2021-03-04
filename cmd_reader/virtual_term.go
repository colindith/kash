// The virtualTerm is used to simulate the behavior of Linux terminal.
// In cmd_reader_test.go, all of the operations sent to terminal will be sent to virtualTerm instead.
// So that in unit test we can ensure that the operations to terminal are correct.
// To use virtual term, just set debugFlag to be true.
// The len of buf of virtualTerm is limited by its initial length. It panic if out of range.

package cmd_reader

var vt = virtualTerm{}

type virtualTerm struct {
	buf []rune
	ptr int
}

func (vt *virtualTerm) printVirtualTerm(a... interface{}) {
	if len(vt.buf) == 0 {
		vt.buf = make([]rune, 16)
		for i := range vt.buf {
			vt.buf[i] = ' '
		}
	}
	for _, arg := range a {
		switch v := arg.(type) {
		case rune:
			vt.printRune(v)
		case []rune:
			for _, u := range v {
				vt.printRune(u)
			}
		case string:
			for _, u := range v {
				vt.printRune(u)
			}
		default:
			panic("unrecognizable_print_type")
		}
	}
}

func (vt *virtualTerm) printRune(r rune) {
	if r == '\b' {
		if vt.ptr > 0 {
			vt.ptr--
		}
	} else {
		vt.buf[vt.ptr] = r
		vt.ptr++
	}
}

