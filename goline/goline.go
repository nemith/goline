package goline

import (
	"bytes"
	"errors"
	"fmt"
	"syscall"
	"unicode/utf8"
)

const (
	MAX_LINE = 4096
)

var UserTerminatedError = errors.New("User terminated.")

type Prompter interface {
	Prompt() string
}

type StringPrompt string

func (p StringPrompt) Prompt() string {
	return string(p)
}

type GoLine struct {
	tty        *Tty
	prompter   Prompter
	LastPrompt string
	CurLine    []rune
	Pos        int
	Len        int
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
	l.tty.WriteString(l.LastPrompt)
	l.tty.WriteString(string(l.CurLine[:l.Len]))

	// Erase to right
	l.tty.Write([]byte("\x1b[0K"))

	// Move cursor back to original position including the prompt
	pos := l.Pos + len(l.LastPrompt)
	l.tty.Write([]byte(fmt.Sprintf("\x1b[0G\x1b[%dC", pos)))
}

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
		c, _ := l.tty.ReadChar()

		switch c {
		case CHAR_ENTER:
			return string(l.CurLine[:l.Len]), nil
		case CHAR_CTRLC:
			// TODO: Identify this as a user escape.
			return string(l.CurLine[:l.Len]), UserTerminatedError
		case CHAR_BACKSPACE, CHAR_CTRLH:
			l.Backspace()
		case CHAR_CTRLB:
			l.MoveLeft()
		case CHAR_CTRLF:
			l.MoveRight()
		case CHAR_CTRLU: // Delete whole line
			l.CurLine = make([]rune, MAX_LINE)
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

			switch {
			case bytes.Equal(chars, []byte(ESCAPE_RIGHT)): // Right arrow
				l.MoveRight()
			case bytes.Equal(chars, []byte(ESCAPE_LEFT)): // Left arrow
				l.MoveLeft()
			default:
				break
			}
		default:
			b := []byte{c}
			switch {
			case c < 0x80: // One byte utf8 character (ASCII)
			case c < 0xe0: // Two byte utf8 character
				char, _ := l.tty.ReadChar()
				b = append(b, char)
			case c < 0xf0: // Three byte utf8 character
				chars, _ := l.tty.ReadChars(2)
				b = append(b, chars...)
			case c < 0xf8: // Four byte utf8 character
				chars, _ := l.tty.ReadChars(3)
				b = append(b, chars...)

			}

			r, _ := utf8.DecodeRune(b)
			l.Insert(r)
		}
	}
	// We should never get here
	panic("Unproccessed input")
}
