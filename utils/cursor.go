package utils

import "fmt"

func CursorUp(line int) {
	fmt.Printf("\033[%dA", line)
}

func CursorDown(line int) {
	fmt.Printf("\033[%dB", line)
}

func ClearScreen() {
	fmt.Print("\033c")
}
