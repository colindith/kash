package cmd_reader

import (
	"fmt"
	"github.com/nsf/termbox-go"
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

	ptr uint8

	promptStr string
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
	fmt.Print(cl.promptStr)
}

func (cl *cmdLine) setPromptStr(prompt string) {
	cl.promptStr = prompt
}

func (cl *cmdLine) newLine() {
	fmt.Print("\n")
	cl.resetBufAndPrintPrompt()
}

func (cl *cmdLine) back(n int) {
	for n > 0 {
		fmt.Print("\b")
		n--
	}
}

func (cl *cmdLine) insertChar(r rune) {
	fmt.Print(string(r))

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
			fmt.Println("ev.Key", ev.Key, " | ", string(ev.Ch))
			if ev.Key == exitSignal {
				return
			}
			continue

			switch ev.Key {
			case keyLeft:
			case keyRight:
			case keyBackSpace:
				fmt.Print("\b \b")
			case keyDelete:
				fmt.Print("\b\b\b\b\b\b")
			default:
				cl.insertChar(ev.Ch)
			}
		}
	}
}
