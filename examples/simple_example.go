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
			if err == goline.UserTerminatedError {
				fmt.Println("\nUser terminated.")
				return
			} else {
				panic(err)
			}
		}

		fmt.Printf("\nGot: '%s' (%d)\n", data, len(data))

		if data == "exit" || data == "quit" {
			fmt.Println("Exiting.")
			return
		}

	}
}
