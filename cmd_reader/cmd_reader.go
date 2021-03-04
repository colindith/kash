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
	exitSignal = keyCtlD     // Ctrl + D to exit

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

	colorDefault   = "\u001B[0m"
	colorCursor    = "\u001B[1;47;35m"
)

func myPrint(a... interface{}) {
	if debugFlag {
		vt.printVirtualTerm(a...)
	}
	for i := 0; i < len(a); i++ {
		switch r := a[i].(type) {
		case rune:
			a[i] = string(r)
		case []rune:
			a[i] = string(r)
		}
	}
	fmt.Fprint(os.Stdout, a...)
}

type Config struct {
	prompt string
	handlers []*Handler
}

func (cfg *Config) SetPromptStr(prompt string) {
	cfg.prompt = prompt
}

func (cfg *Config) RegistryHandler(h *Handler) error {
	if err := validateHandler(h); err != nil {
		return err
	}
	cfg.handlers = append(cfg.handlers, h)
	return nil
}

func validateHandler(h *Handler) error {
	if h == nil {
		return fmt.Errorf("invalid nil handler")
	}
	// TODO: implement this
	return nil
}

type cmdLine struct {
	history history

	buf []rune           // The final of buf should always be ' '
	tmp []rune
	ptr int

	promptStr string
	handlers []*Handler
}


func newCMDLine(cfg *Config) *cmdLine {
	if cfg == nil {
		cfg = &Config{
			prompt: ">",
		}
	}
	cl := &cmdLine{}
	cl.setPromptStr(cfg.prompt)
	cl.resetBufAndPrintPrompt()
	cl.initHistory()
	cl.blink(true)
	cl.handlers = cfg.handlers
	return cl
}

func (cl *cmdLine) resetBufAndPrintPrompt() {
	cl.ptr = 0
	cl.buf = []rune{' '}
	cl.printPrompt()
	cl.history.resetPtr()
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

func (cl *cmdLine) blink(on bool) {
	if on {
		myPrint(colorCursor, cl.buf[cl.ptr], colorDefault, '\b')
		return
	}
	myPrint(colorDefault, cl.buf[cl.ptr], '\b')
}

// *** History Manipulation *** //

func (cl *cmdLine) initHistory() {
	cl.history = history{}
	// TODO: Find the history from hard disk
}

func (cl *cmdLine) searchHistoryUp() {
	result, foundResult, preserveZero := cl.history.searchUp()
	if !foundResult {
		// not find any match
		// beep?
		return
	}
	// move cursor back to 0
	cl.back(cl.ptr)
	if preserveZero {
		cl.preserve()
	}
	cl.copyToBuf(result)

	// print new buf
	myPrint(cl.buf)    // buf contains ' ' in the end
	myPrint('\b')
	cl.ptr = len(cl.buf) - 1
}

func (cl *cmdLine) searchHistoryDown() {
	// move cursor back to 0
	cl.back(cl.ptr)

	result, foundResult, retrieveZero := cl.history.searchDown()
	if foundResult {
		cl.copyToBuf(result)

		// print new buf
		myPrint(cl.buf)
		cl.ptr = len(cl.buf)
	}
	if retrieveZero {
		cl.copyToBuf(cl.tmp)
	}

	// print new buf
	myPrint(cl.buf)
	cl.ptr = len(cl.buf)
}

func (cl *cmdLine) copyToBuf(s []rune) {
	if len(s)+1 > len(cl.buf) {
		cl.buf = make([]rune, len(s)+1)
	}
	copy(cl.buf, s)
	cl.buf[len(s)] = ' '
	cl.buf = cl.buf[:len(s)+1]
}

// pushHistory push the current typing back to the head of history list
func (cl *cmdLine) pushHistory() {
	cl.history.pushFront(cl.buf)
}

// preserve store the current typing into a temporary location so that cl.buf can be load with historical cmds
func (cl *cmdLine) preserve() {
	cl.tmp = make([]rune, len(cl.buf))
	copy(cl.tmp, cl.buf)
}

// *** cmdLine handle *** //
func (cl *cmdLine) handle(buf string) {
	if len(cl.handlers) == 0 {
		return
	}
	// TODO: should implement a router to the different handlers
	result, err := cl.handlers[0].Serv(buf)
	if err != nil {
		// TODO: log the error
		return
	}
	if len(result) != 0 {
		myPrint("\n")
		myPrint(result)
	}
}

func Run(cfg *Config) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorGreen, termbox.ColorDefault)

	//termbox.Flush()

	cl := newCMDLine(cfg)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			cl.blink(false)
			if ev.Key == exitSignal {
				return
			}

			switch ev.Key {
			case keyLeft:
				cl.back(1)
			case keyRight:
				cl.moveRight()
			case keyNewLine:
				buf := cl.buf
				cl.pushHistory()
				cl.handle(string(buf))
				cl.newLine()
			case keyBackSpace:
				cl.backSpace()
			case keyWhiteSpace:
				cl.insertChar(' ')
			case keyUp:
				cl.searchHistoryUp()
			case keyDown:
				cl.searchHistoryDown()
			case keyCtlC:
				// TODO: quit current task
			default:
				cl.insertChar(ev.Ch)
			}
		}
		cl.history.setMatchStr(cl.buf[:cl.ptr])
		cl.blink(true)
	}
}

// *** history struct *** //

// node is a linked list node
type node struct {
	next *node
	prev *node

	val  []rune
}

type history struct {
	head *node
	len  int
	ptr  *node     // ptr is the pointer used when navigating the history list

	match []rune   // match is used when searching next/prev history cmd
}

func (h *history) resetPtr() {
	// This method should be called when create a new line (enter or Ctl+C)
	h.ptr = h.head
}

func (h *history) setMatchStr(s []rune) {
	// This method should be called when move cursor right or left (including insert char and new line)
	h.match = s
}

func (h *history) searchUp() (result []rune, foundResult bool, preserveZero bool) {
	// The "preserveZero" means the cmdLine is moving from current tying stuff to the first level history
	// So the cmdLine should preserve the current buffer stuff in a temporary location
	var ptr *node
	if h.ptr == nil {
		if h.head == nil {
			return nil, false, false
		}
		// (ptr == nil && head != nil) means move from current typing to 1st level history
		preserveZero = true
		ptr = h.head
	} else {
		ptr = h.ptr.next
	}
	for {
		if ptr == nil {
			// reach tail
			return nil, false, false
		}
		if checkMatch(ptr, h.match) {
			h.ptr = ptr
			return ptr.val, true, preserveZero
		}
		ptr = ptr.next
	}
}

func (h *history) searchDown() (result []rune, foundResult bool, retrieveZero bool) {
	if h.ptr == nil {
		// in zero
		return nil, false, false
	}
	ptr := h.ptr.prev
	for {
		if ptr == nil {
			// reach head
			h.ptr = nil      // set ptr to nil for next searchUp call
			return nil, false, true
		}
		if checkMatch(ptr, h.match) {
			h.ptr = ptr
			return ptr.val, true, false
		}
		ptr = ptr.prev
	}
}

func (h *history) pushFront(s []rune) {
	// TODO: This function should not push the repeated rune slice into history??

	newHead := &node{
		next: h.head,
		prev: nil,
		val:  s,
	}
	h.head = newHead
	h.len++
	if h.head != nil {
		// Not empty history
		h.head.prev = newHead
	}
}

func (h *history) autoComplete(s []rune) ([]rune, bool) {
	// TODO: implement this
	return nil, false
}

func checkMatch(ptr *node, s []rune) bool {
	if ptr == nil {
		return false
	}
	if len(ptr.val) <= len(s) {
		return false
	}
	match := true
	for i := 0; i < len(s); i++ {
		if ptr.val[i] != s[i] {
			match = false
			break
		}
	}
	return match
}


// *** Handler struct *** //
type Handler struct {
	// TODO: check how to route the cmd to handler
	Serv Serv
}

type Serv func(cmd string) (result string, err error)