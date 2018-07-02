package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

type lexer struct {
	out       chan string
	pos       int
	width     int
	line      int
	lineIndex int
	text      string
}

func newLexer(bytes []byte) lexer {
	lex := lexer{}
	lex.out = make(chan string)
	lex.text = string(bytes)
	return lex
}

func (lex *lexer) hasMore() bool {
	return lex.pos < len(lex.text)
}

func (lex *lexer) next() (rune, error) {
	if lex.pos >= len(lex.text) {
		return 0, io.EOF
	}
	var c rune
	c, lex.width = utf8.DecodeRuneInString(lex.text[lex.pos:])
	lex.pos += lex.width
	return c, nil
}

func (lex *lexer) consume(pattern string, token string, args ...string) (string, error) {
	if lex.pos >= len(lex.text) {
		return "", io.EOF
	}
	r := regexp.MustCompile(fmt.Sprintf("^%s", pattern))
	bytes := r.Find([]byte(lex.text[lex.pos:]))
	lex.width = len(bytes)
	lex.pos += lex.width
	text := string(bytes)
	if len(text) > 0 {
		lex.emit(token, text)
	}
	return text, nil
}

func (lex *lexer) newLine() {
	lex.line++
	lex.lineIndex = lex.pos
}

func (lex *lexer) emit(s string, args ...string) {
	line := fmt.Sprintf("%s %d %d", s, lex.line+1, lex.pos-lex.width-lex.lineIndex+1)
	for _, arg := range args {
		line += fmt.Sprintf(" %s", arg)
	}
	line += fmt.Sprintf("\n")
	if os.Getenv("LEXDEBUG") != "" {
		fmt.Fprint(os.Stderr, line)
	}
	lex.out <- line
}

func (lex *lexer) consumeString() {
	var buffer bytes.Buffer
	var prev rune
	for {
		c, err := lex.next()
		if err != nil {
			break
		}
		if prev != '\\' && c == '"' {
			break
		}
		prev = c
		buffer.WriteRune(c)
	}
	lex.emit("string", buffer.String())
}

func (lex *lexer) lex() {
	for lex.hasMore() {
		lex.consume("let", "let")
		lex.consume("print", "print")
		lex.consume("[a-z]+", "id")
		lex.consume("[0-9]+", "number")
		c, _ := lex.next()
		if c == '"' {
			lex.consumeString()
		} else if c == ';' {
			lex.emit(";")
			lex.newLine()
		} else if strings.ContainsRune("=+-*/()", c) {
			lex.emit(string(c))
		} else if !strings.ContainsRune(" \t\r\n", c) {
			fmt.Fprintf(os.Stderr, "unrecognized char '%c' at line %d, column %d\n", c, lex.line+1, lex.pos-lex.width-lex.lineIndex+1)
			os.Exit(1)
		}
	}

	lex.emit("eof")
	close(lex.out)
}

// Lex returns a channel to use for Parse
func Lex(filepath string) <-chan string {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lex := newLexer(bytes)
	go lex.lex()
	return lex.out
}
