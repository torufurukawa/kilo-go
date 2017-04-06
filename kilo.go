package main

import (
	"os"
)

func main() {
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
	}
}
