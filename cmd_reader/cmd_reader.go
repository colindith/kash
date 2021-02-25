package cmd_reader

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
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
)

type cmdLine struct {
	history [][]rune     // TODO: the history should be saved in hard disk
	buf []rune
	tmp []rune

	ptr int

	promptStr string
}

func myPrint(a... interface{}) {
	for i := 0; i < len(a); i++ {
		if r, ok := a[i].(rune); ok {
			a[i] = string(r)
		}
	}
	fmt.Fprint(os.Stdout, a...)
}

func newCMDLine(prompt string) *cmdLine {
	cl := &cmdLine{}
	cl.resetBufAndPrintPrompt()
	cl.setPromptStr(prompt)
	return cl
}

func (cl *cmdLine) resetBufAndPrintPrompt() {
	cl.ptr = 0
	cl.buf = make([]rune, 0)
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
	cl.ptr -= n
	if cl.ptr < 0 {
		cl.ptr = 0
	}
	for n > 0 {
		myPrint("\b")
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
	if cl.ptr == bufLen {
		myPrint(r)
		cl.ptr = bufLen + 1
		cl.buf = append(cl.buf, r)
		return
	}

	i := bufLen - 1
	cl.buf = append(cl.buf, cl.buf[i])
	for i > cl.ptr {
		cl.buf[i] = cl.buf[i-1]
		i--
	}
	cl.ptr++

	myPrint(r)

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
				if cl.movePtrTo(cl.ptr+1) {
					myPrint(cl.buf[cl.ptr+1])
				}
			case keyBackSpace:
				myPrint("\b \b")
			case keyDelete:
				fmt.Print("\b\b\b\b\b\b")
			default:
				cl.insertChar(ev.Ch)
			}
		}
	}
}
