package main

import (
	"fmt"
	"github.com/nemith/go-goline/goline"
	"syscall"
)

func main() {
	tty, _ := goline.NewTty(syscall.Stdin)

	tty.EnableRawMode()
	defer tty.DisableRawMode()

	for {
		c, _ := tty.ReadChar()
		switch c {
		case goline.CHAR_CTRLC:
			return
		default:
			tty.Write([]byte(fmt.Sprintf("Char: %c (%d)\r\n", c, c)))
		}
	}
}
