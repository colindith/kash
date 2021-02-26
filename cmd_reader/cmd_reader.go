package cmd_reader

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
)

var (
	// debugFlag if set to true, the output char will only be write to virtual term. This is used for unit test
	debugFlag = false
)

const (
	exitSignal = 4

	keyBackSpace = 127
	keyDelete = 65522
	keyUp = 65517
	keyDown = 65516
	keyLeft = 65515
	keyRight = 65514
	keyNewLine = 13
	keyEsc = 27
	keyWhiteSpace = 32
	keyCtlD = 4
	keyCtlC = 3
	keyCtlV = 22

	colorDefault   = "033[0m"
	colorCursor    = "\u001B[1;47m"
)

func myPrint(a... interface{}) {
	if debugFlag {
		vt.printVirtualTerm(a...)
	}
	for i := 0; i < len(a); i++ {
		if r, ok := a[i].(rune); ok {
			a[i] = string(r)
		} else if r, ok := a[i].([]rune); ok {
			a[i] = string(r)
		}
	}
	fmt.Fprint(os.Stdout, a...)
}

type cmdLine struct {
	history [][]rune     // TODO: the history should be saved in hard disk
	buf []rune           // The final of buf should always be ' '
	tmp []rune

	ptr int

	promptStr string
}


func newCMDLine(prompt string) *cmdLine {
	cl := &cmdLine{}
	cl.setPromptStr(prompt)
	cl.resetBufAndPrintPrompt()
	return cl
}

func (cl *cmdLine) resetBufAndPrintPrompt() {
	cl.ptr = 0
	cl.buf = []rune{' '}
	cl.printPrompt()
}

func (cl *cmdLine) addHistory() {
	cl.history = append(cl.history, cl.buf)
}

func (cl *cmdLine) printPrompt() {
	myPrint(cl.promptStr)
}

func (cl *cmdLine) setPromptStr(prompt string) {
	cl.promptStr = prompt
}

func (cl *cmdLine) newLine() {
	myPrint("\n")
	cl.resetBufAndPrintPrompt()
}

func (cl *cmdLine) back(n int) {
	if cl.ptr < n {
		n = cl.ptr
		// should beep?
	}

	for n > 0 {
		myPrint("\b")
		cl.ptr--
		n--
	}
}

func (cl *cmdLine) movePtrTo(ptr int) bool {
	if ptr <= len(cl.buf) - 1 && ptr >= 0 {
		cl.ptr = ptr
		return true
	}
	return false
}

func (cl *cmdLine) insertChar(r rune) {
	bufLen := len(cl.buf)

	i := bufLen - 1
	cl.buf = append(cl.buf, cl.buf[i])
	for i > cl.ptr {
		cl.buf[i] = cl.buf[i-1]
		i--
	}
	cl.buf[cl.ptr] = r

	myPrint(r, cl.buf[cl.ptr+1:])
	for i = 0; i < bufLen-cl.ptr; i++ {
		myPrint("\b")
	}
	cl.ptr++
}

func (cl *cmdLine) backSpace() {
	// TODO: Code is ugly. Try to sanitize it.
	if cl.ptr <= 0 {
		// beep
		return
	}

	cl.back(1)
	tmpPtr := cl.ptr
	for i := cl.ptr; i < len(cl.buf)-1; i++ {
		cl.buf[i] = cl.buf[i+1]
		myPrint(cl.buf[i])
	}
	myPrint(cl.buf[len(cl.buf)-1])

	cl.buf = cl.buf[:len(cl.buf)-1]

	for i := 0; i < (len(cl.buf) - tmpPtr)+1; i++ {
		myPrint('\b')
	}
}

func (cl *cmdLine) moveRight() {
	if cl.movePtrTo(cl.ptr+1) {
		myPrint(cl.buf[cl.ptr-1])
	}
}

func Run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorGreen, termbox.ColorDefault)

	//termbox.Flush()

	cl := newCMDLine("kash> ")   // TODO: read the prompt from outside

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:

			if ev.Key == exitSignal {
				return
			}
			//fmt.Println("ev.Key", ev.Key, " | ", string(ev.Ch))
			//continue

			switch ev.Key {
			case keyLeft:
				cl.back(1)
			case keyRight:
				cl.moveRight()
			case keyNewLine:
				cl.newLine()
			case keyBackSpace:
				cl.backSpace()
			case keyDelete:
				fmt.Print("\b\b\b\b\b\b")
			default:
				cl.insertChar(ev.Ch)
			}
		}
	}
}
