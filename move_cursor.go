package main

import "fmt"

func moveCursorUp(n int) {
	fmt.Printf("\x1b[%dA", n)
}

func moveCursorFront[T any](options []T) {
	fmt.Printf("\x1b[%d;0H", len(options)+1)
}
