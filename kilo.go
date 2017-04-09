package main

import (
	"bytes"
	"fmt"
	"os"

	"unicode"

	"golang.org/x/crypto/ssh/terminal"
)

var state *terminal.State

func main() {
	err := enableRawMode()
	if err != nil {
		panic(err)
	}
	defer disableRawMode()

	b := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(b)
		// EOF or error
		if n != 1 || err != nil {
			break
		}

		// quit
		if b[0] == 'q' {
			break
		}

		// print
		c := bytes.Runes(b)[0]
		if unicode.IsControl(c) {
			fmt.Printf("%d\r\n", c)
		} else {
			fmt.Printf("%d ('%c')\r\n", c, c)
		}
	}
}

func enableRawMode() error {
	var err error
	state, err = terminal.MakeRaw(0)
	return err
}

func disableRawMode() {
	terminal.Restore(0, state)
}
