package goline

import (
	"errors"
	"fmt"
	"syscall"
)

const (
	MAX_LINE = 4096
)

type Prompter interface {
	Prompt() string
}

type StringPrompt string

func (p StringPrompt) Prompt() string {
	return string(p)
}

//type CustomFunction func()

type GoLine struct {
	tty        *Tty
	prompter   Prompter
	LastPrompt string
	CurLine    []byte
	Pos        int
	Len        int
	//	map[byte]
}

func NewGoLine(p Prompter) *GoLine {
	tty, _ := NewTty(syscall.Stdin)
	return &GoLine{
		tty:      tty,
		prompter: p}
}

func (l *GoLine) RefreshLine() {
	// Cursor to left edge
	l.tty.Write([]byte("\x1b[0G"))

	// Write the prompt and the current buffer content
	l.tty.Write([]byte(l.LastPrompt))
	l.tty.Write(l.CurLine[:l.Len])

	// Erase to right
	l.tty.Write([]byte("\x1b[0K"))

	// Move cursor back to original position including the prompt
	pos := l.Pos + len(l.LastPrompt)
	l.tty.Write([]byte(fmt.Sprintf("\x1b[0G\x1b[%dC", pos)))
}

func (l *GoLine) Insert(c byte) {
	if l.Len == l.Pos {
		l.CurLine[l.Pos] = c
		l.Pos++
		l.Len++
		//	l.tty.Write([]byte{c})
		l.RefreshLine()
	} else {
		n := len(l.CurLine)
		copy(l.CurLine[l.Pos+1:n], l.CurLine[l.Pos:n])
		l.CurLine[l.Pos] = c
		l.Pos++
		l.Len++
		l.RefreshLine()
	}
}

func (l *GoLine) Backspace() {
	if l.Len > 0 && l.Pos > 0 {
		l.CurLine = append(l.CurLine[:l.Pos-1], l.CurLine[l.Pos:]...)
		l.Len--
		l.Pos--
		l.CurLine[l.Len] = 0
		l.RefreshLine()
	}
}

func (l *GoLine) DeleteLastWord() {
	//TODO: Implement
}

func (l *GoLine) MoveLeft() {
	if l.Pos > 0 {
		l.Pos--
		l.RefreshLine()
	}
}

func (l *GoLine) MoveRight() {
	if l.Pos != l.Len {
		l.Pos++
		l.RefreshLine()
	}
}

func (l *GoLine) ClearScreen() {
	l.tty.WriteString("\x1b[H\x1b[2J")
}

func (l *GoLine) Line() ([]byte, error) {
	// Write out the current prompt
	l.LastPrompt = l.prompter.Prompt()
	l.tty.Write([]byte(l.LastPrompt))

	// Go into RawMode and leave after we are finished
	l.tty.EnableRawMode()
	defer l.tty.DisableRawMode()

	l.Len = 0
	l.Pos = 0
	l.CurLine = make([]byte, MAX_LINE)

	for {
		c, _ := l.tty.ReadChar()

		switch c {
		case CHAR_ENTER:
			return l.CurLine[:l.Len], nil
		case CHAR_CTRLC:
			// TODO: Identify this as a user escape.
			return l.CurLine[:l.Len], errors.New("CTRL-C!")
		case CHAR_Backspace, CHAR_CTRLH:
			l.Backspace()
		case CHAR_CTRLB:
			l.MoveLeft()
		case CHAR_CTRLF:
			l.MoveRight()
		case CHAR_CTRLU: // Delete whole line
			l.CurLine = make([]byte, MAX_LINE)
			l.Pos = 0
			l.Len = 0
			l.RefreshLine()
		case CHAR_CTRLK: // Delete from current position to the end of the line
			copy(l.CurLine, l.CurLine[:l.Pos])
			l.Len = l.Pos
			l.RefreshLine()
		case CHAR_CTRLA: // Go to the start of line
			l.Pos = 0
			l.RefreshLine()
		case CHAR_CTRLE: // Go to the end of line
			l.Pos = l.Len
			l.RefreshLine()
		case CHAR_CTRLL: // Clear the screen
			l.ClearScreen()
			l.RefreshLine()
		case CHAR_CTRLW: // Delete Previous Word
			l.DeleteLastWord()
			l.RefreshLine()
		case CHAR_ESCAPE:
			// Recevied an escape sequence.  Read the next two characters
			chars, err := l.tty.ReadChars(2)
			if err != nil {
				break
			}

			switch string(chars) {
			case "[C": // Right arrow
				l.MoveRight()
			case "[D": // Left arrow
				l.MoveLeft()
			default:
				break
			}
		default:
			l.Insert(c)
		}
	}
}
