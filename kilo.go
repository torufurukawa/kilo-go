package main

import (
	"fmt"
	"os"

	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	ctrlQ   byte   = 'q' & 0x1f
	version string = "0.0.1"
)

type screen struct {
	Width  int
	Height int
	CX     int
	CY     int
}

func main() {
	enableRawMode()
	defer disableRawMode()

	w, h, err := terminal.GetSize(0)
	if err != nil {
		die(err)
	}
	s := &screen{Width: w, Height: h, CX: 0, CY: 0}

	for {
		editorRefreshScreen(s)
		if quit := editorProcessKeypress(s); quit {
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

func editorProcessKeypress(s *screen) (quit bool) {
	c := editorReadKey()
	switch c {
	case ctrlQ:
		return true
	case 'a':
		s.CX--
	case 'd':
		s.CX++
	case 'w':
		s.CY--
	case 's':
		s.CY++
	}

	return false
}

// output

func editorRefreshScreen(s *screen) {
	buffer := make([]byte, 0)

	buffer = append(buffer, []byte("\x1b[?25l")...)
	buffer = append(buffer, []byte("\x1b[H")...)

	editorDrawRows(&buffer, s)

	cursor := fmt.Sprintf("\x1b[%d;%dH", s.CY+1, s.CX+1)
	buffer = append(buffer, []byte(cursor)...)

	buffer = append(buffer, []byte("\x1b[?25h")...)

	os.Stdout.Write(buffer)
}

func editorDrawRows(buffer *[]byte, s *screen) {
	for y := 0; y < s.Height; y++ {
		if y == s.Height/3 {
			message := fmt.Sprintf("Kilo editor -- version %s", version)
			if s.Width < len(message) {
				message = message[:s.Width]
			}
			*buffer = append(*buffer, '~')
			padding := strings.Repeat(" ", (s.Width-len(message))/2)
			*buffer = append(*buffer, padding...)
			*buffer = append(*buffer, []byte(message)...)
		} else {
			*buffer = append(*buffer, '~')
		}

		*buffer = append(*buffer, []byte("\x1b[K")...)
		if y < s.Height-1 {
			*buffer = append(*buffer, []byte("\r\n")...)
		}
	}
}
