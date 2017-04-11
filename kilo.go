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

	_, h, err := terminal.GetSize(0)
	if err != nil {
		die(err)
	}

	for {
		editorRefreshScreen(h)
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

func editorRefreshScreen(height int) {
	buffer := make([]byte, 0)

	buffer = append(buffer, []byte("\x1b[2J")...)
	buffer = append(buffer, []byte("\x1b[H")...)
	editorDrawRows(&buffer, height)
	buffer = append(buffer, []byte("\x1b[H")...)

	os.Stdout.Write(buffer)
}

func editorDrawRows(buffer *[]byte, height int) {
	for y := 0; y < height; y++ {
		*buffer = append(*buffer, '~')

		if y < height-1 {
			*buffer = append(*buffer, []byte("\r\n")...)
		}
	}
}
