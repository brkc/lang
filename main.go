package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "missing file")
		os.Exit(2)
	}
	Parse(Lex(os.Args[1]))
}
