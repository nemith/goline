package main

import (
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

		fmt.Printf("\nGot: '%s' (%d)\n", data, len(data))

		if data == "exit" || data == "quit" {
			fmt.Printf("Exiting.\n")
			return
		}

	}
}
