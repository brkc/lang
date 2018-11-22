package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	lexFlag   = flag.Bool("lex", false, "lex only")
	parseFlag = flag.Bool("parse", false, "parse only")
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "missing file")
		os.Exit(2)
	}
	file := flag.Arg(0)
	if *lexFlag {
		debugLex(file)
	} else if *parseFlag {
		debugParse(file)
	} else {
		interpret(file)
	}
}

func debugLex(file string) {
	lexer := lex(file)
	for {
		line := <-lexer
		fmt.Fprintln(os.Stderr, line)
		if strings.Index(line, "eof") == 0 {
			return
		}
	}
}

func debugParse(file string) {
	fmt.Fprintln(os.Stderr, parse(lex(file)))
}
