package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "missing file")
		os.Exit(2)
	}
	file := os.Args[1]
	if os.Getenv("LEXDEBUG") != "" {
		lexdebug(file)
	} else if os.Getenv("PARSEDEBUG") != "" {
		parsedebug(file)
	} else {
		interpret(file)
	}
}

func lexdebug(file string) {
	lexer := lex(file)
	for {
		line := <-lexer
		fmt.Fprintln(os.Stderr, line)
		if strings.Index(line, "eof") == 0 {
			return
		}
	}
}

func parsedebug(file string) {
	fmt.Fprintln(os.Stderr, parse(lex(file)))
}
