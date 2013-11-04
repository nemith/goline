package goline

import (
	"syscall"
	"unsafe"
)

// ASCII codes for comonly used control characters
const (
	CHAR_CTRLA     = 1
	CHAR_CTRLB     = 2
	CHAR_CTRLC     = 3
	CHAR_CTRLE     = 5
	CHAR_CTRLF     = 6
	CHAR_CTRLH     = 8
	CHAR_CTRLK     = 11
	CHAR_CTRLL     = 12
	CHAR_ENTER     = 13
	CHAR_CTRLU     = 21
	CHAR_CTRLW     = 23
	CHAR_ESCAPE    = 27
	CHAR_BACKSPACE = 127
)

// Commonly used escape codes with out the escape character
const (
	ESCAPE_UP    = "[A"
	ESCAPE_RIGHT = "[C"
	ESCAPE_LEFT  = "[D"
)

// Get the current Termios via syscall for the given terminal at the file
// descriptor fd
func GetTermios(fd int) (*syscall.Termios, error) {
	var termios syscall.Termios
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), ioctlReadTermios,
		uintptr(unsafe.Pointer(&termios)))

	if err != 0 {
		return nil, err
	}

	return &termios, nil
}

// Sets the Termios via syscall for the given terminal at the file descriptor
// fd
func SetTermios(fd int, termios *syscall.Termios) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), ioctlWriteTermios,
		uintptr(unsafe.Pointer(termios)))

	if err != 0 {
		return err
	}

	return nil
}

// Test to see if the file descriptor is a terminal or not
func IsTerminal(fd int) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), ioctlReadTermios,
		uintptr(unsafe.Pointer(&termios)))
	return err == 0
}

// Wrapper function to keep state and perform low-level functions for a
// given terminal
type Tty struct {
	fd          int
	origTermios syscall.Termios
	rawMode     bool
}

// Create a new TTY at the given terminal
func NewTty(fd int) (*Tty, error) {
	origTermios, err := GetTermios(fd)
	if err != nil {
		return nil, err
	}

	return &Tty{fd, *origTermios, false}, nil
}

// Enable raw mode on the terminal to allow programs to perform and modify the
// terminal itself.  Saves off the old termios first to restor with
// DisableRawMode
func (t *Tty) EnableRawMode() error {
	raw := t.origTermios

	/* input modes: no break, no CR to NL, no parity check, no strip char,
	 * no start/stop output control. */
	raw.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	/* output modes - disable post processing */
	raw.Oflag &^= syscall.OPOST
	/* control modes - set 8 bit chars */
	raw.Cflag |= (syscall.CS8)
	/* local modes - choing off, canonical off, no extended functions,
	 * no signal chars (^Z,^C) */
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	/* control chars - set return condition: min number of bytes and timer.
	 * We want read to return every single byte, without timeout. */
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	if err := SetTermios(t.fd, &raw); err != nil {
		return err
	}

	t.rawMode = true
	return nil
}

// Disables Raw mode by restoring the saved termios.  NOOP if not currently in
// raw mode
func (t *Tty) DisableRawMode() {
	if t.rawMode {
		SetTermios(t.fd, &t.origTermios)
		t.rawMode = false
	}
}

// Implments a Reader interface for the tty by wrapping the syscall.Read
// function
func (t *Tty) Read(p []byte) (int, error) {
	return syscall.Read(t.fd, p)
}

// Implments a Writer interface for the tty by wrapping the syscall.Write
// function
func (t *Tty) Write(p []byte) error {
	_, err := syscall.Write(t.fd, p)
	return err
}

// Write a string instead of a byte to the terminal
func (t *Tty) WriteString(s string) error {
	_, err := syscall.Write(t.fd, []byte(s))
	return err
}

// Read a single character and return it
func (t *Tty) ReadChar() (byte, error) {
	var char [1]byte
	_, err := t.Read(char[0:])

	if err != nil {
		return 0, err
	}

	return char[0], nil
}

// Read a number of chacters from terminal and return them as a byte slice
func (t *Tty) ReadChars(numChars int) ([]byte, error) {
	chars := make([]byte, numChars)
	_, err := t.Read(chars[0:])

	if err != nil {
		return []byte(""), err
	}

	return chars, nil
}
