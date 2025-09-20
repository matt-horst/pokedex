package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)

	return strings.Fields(text)
}

func main() {
	fmt.Println("Hello, World!")
}
