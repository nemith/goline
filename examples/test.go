package main

import (
	"bytes"
	"fmt"
	"github.com/nemith/go-goline/goline"
)

func main() {
	gl := goline.NewGoLine(goline.StringPrompt("prompt> "))
	for {
		data, err := gl.Line()
		if err != nil {
			panic(err)
		}

		fmt.Printf("\nGot: '%s' (%d)\n", data, gl.Len)

		if bytes.Equal(data[:gl.Len], []byte("exit")) ||
			bytes.Equal(data[:gl.Len], []byte("quit")) {
			fmt.Printf("Exiting.\n")
			return
		}

	}
}
