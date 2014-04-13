package goline

import (
	"errors"
	"fmt"
	"syscall"
)

const (
	MAX_LINE = 4096
)

// Error represents when a user has termianted a Goline Line operation with
// an terminating signal like CTRL-C
var UserTerminatedError = errors.New("User terminated.")

// Prompter type define an interface that supports the Prompt() function that
// returns a string
type Prompter interface {
	Prompt() string
}

// StringPrompter is a Prompter that returns a string as the prompt
type StringPrompt string

// Returns the underlying string to be used as the prompt
func (p StringPrompt) Prompt() string {
	return string(p)
}

type Handler func(*GoLine) (bool, error)

// GoLine stores the current internal state of the Line operation
type GoLine struct {
	tty            *Tty
	prompter       Prompter
	Handlers       map[rune]Handler
	DefaultHandler Handler
	LastPrompt     string
	CurLine        []rune
	Pos            int
	Len            int
}

// Creates a new goline with the input set to STDIN and the prompt set to the
// prompter p
func NewGoLine(p Prompter) *GoLine {
	tty, _ := NewTty(syscall.Stdin)
	l := &GoLine{
		tty:      tty,
		Handlers: make(map[rune]Handler),
		prompter: p}

	l.AddHandler(CHAR_ENTER, Finish)
	l.AddHandler(CHAR_CTRLC, UserTerminated)
	l.AddHandler(CHAR_BACKSPACE, Backspace)
	l.AddHandler(CHAR_CTRLH, Backspace)

	// Movement
	l.AddHandler(CHAR_CTRLB, MoveLeft)
	l.AddHandler(CHAR_CTRLF, MoveRight)
	l.AddHandler(CHAR_CTRLA, MoveStartofLine)
	l.AddHandler(CHAR_CTRLE, MoveEndofLine)

	//Edit
	l.AddHandler(CHAR_CTRLL, ClearScreen)
	l.AddHandler(CHAR_CTRLU, DeleteLine)

	//	l.DefaultHandler = DefaultHandler
	return l
}

// Refreshes the current line by first moving the cursor to the left edge, then
// writing the current prompt and contents of the buffer, erasing the remaining
// line to the right and the placing the cusor back at the original (or updated)
// position
func (l *GoLine) RefreshLine() {
	// Cursor to left edge
	l.tty.Write([]byte("\x1b[0G"))

	// Write the prompt and the current buffer content
	l.tty.WriteString(l.LastPrompt)
	l.tty.WriteString(string(l.CurLine[:l.Len]))

	// Erase to right
	l.tty.Write([]byte("\x1b[0K"))

	// Move cursor back to original position including the prompt
	pos := l.Pos + len(l.LastPrompt)
	l.tty.Write([]byte(fmt.Sprintf("\x1b[0G\x1b[%dC", pos)))
}

// Inserts the unicode character r at the current position on the line
func (l *GoLine) Insert(r rune) {
	if l.Len == l.Pos {
		l.CurLine[l.Pos] = r
		l.Pos++
		l.Len++
		//	l.tty.Write([]byte{c})
		l.RefreshLine()
	} else {
		n := len(l.CurLine)
		copy(l.CurLine[l.Pos+1:n], l.CurLine[l.Pos:n])
		l.CurLine[l.Pos] = r
		l.Pos++
		l.Len++
	}
}

// Clears the entire screen
func (l *GoLine) ClearScreen() {
	l.tty.WriteString("\x1b[H\x1b[2J")
}

// Add a custom handler function to be invoked when rune r is encontered.
//
// Function is passed a pointer to the current GoLine to be able to read or
// modify the current buffer, cursor position, or length
//
// Custom handlers are evaulated before built-in functions which allows you to
// override built-in functionality
func (l *GoLine) AddHandler(r rune, f Handler) {
	l.Handlers[r] = f
}

// Removes a handler with rune of r
func (l *GoLine) RemoveHanlder(r rune) {
	delete(l.Handlers, r)
}

// Print the current prompt and handle each character as it received from
// the underlying terminal.  Returns a string when input has been completed
// (i.e when the user hits enter/return) and an error (if one exists).
//
// Error can be UserTerminatedError which means the user terminated the current
// operation with CTRL-C or similar terminating signal.
func (l *GoLine) Line() (string, error) {
	// Write out the current prompt
	l.LastPrompt = l.prompter.Prompt()
	l.tty.WriteString(l.LastPrompt)

	// Go into RawMode and leave after we are finished
	l.tty.EnableRawMode()
	defer l.tty.DisableRawMode()

	l.Len = 0
	l.Pos = 0
	l.CurLine = make([]rune, MAX_LINE)

	for {
		r, _ := l.tty.ReadRune()

		if f, found := l.Handlers[r]; found {
			l.tty.DisableRawMode()
			stop, err := f(l)
			l.tty.EnableRawMode()

			if stop || err != nil {
				return string(l.CurLine[:l.Len]), err
			}
		} else {
			l.Insert(r)
		}
		l.RefreshLine()
	}
	// We should never get here
	panic("Unproccessed input")
}
