package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

const ctrlQ byte = 'q' & 0x1f

func main() {
	enableRawMode()
	defer disableRawMode()

	for {
		editorRefreshScreen()
		if quit := editorProcessKeypress(); quit {
			break
		}
	}
}

// terminal

var originalState *terminal.State

func enableRawMode() {
	var err error
	originalState, err = terminal.MakeRaw(0)
	if err != nil {
		die(err)
	}
}

func disableRawMode() {
	terminal.Restore(0, originalState)
}

func die(err error) {
	os.Stdout.WriteString("\x1b[2J")
	os.Stdout.WriteString("\x1b[H")

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func editorReadKey() byte {
	b := make([]byte, 1)
	n, err := os.Stdin.Read(b)
	// EOF or error
	if n != 1 || err != nil {
		die(err)
	}
	return b[0]
}

// input

func editorProcessKeypress() (quit bool) {
	c := editorReadKey()
	switch c {
	case ctrlQ:
		return true
	}

	return false
}

// output

func editorRefreshScreen() {
	os.Stdout.WriteString("\x1b[2J")
	os.Stdout.WriteString("\x1b[H")
	editorDrawRows()
	os.Stdout.WriteString("\x1b[H")
}

func editorDrawRows() {
	_, h, err := terminal.GetSize(0)
	if err != nil {
		die(err)
	}
	for y := 0; y < h; y++ {
		os.Stdout.WriteString("~")

		if y < h-1 {
			os.Stdout.WriteString("\r\n")
		}
	}
}
