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
	lex.out <- line
}

func (lex *lexer) consumeString() {
	var buf bytes.Buffer
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
		buf.WriteRune(c)
	}
	lex.emit("string", buf.String())
}

func (lex *lexer) lex() {
	for lex.hasMore() {
		lex.consume("var", "var")
		lex.consume("if", "if")
		lex.consume("while", "while")
		lex.consume("break", "break")
		lex.consume("continue", "continue")
		lex.consume("fn", "fn")
		lex.consume("return", "return")
		lex.consume("true", "true")
		lex.consume("false", "false")
		lex.consume("==", "==")
		lex.consume("!=", "!=")
		lex.consume(">=", ">=")
		lex.consume("<=", "<=")
		lex.consume("not", "not")
		lex.consume("and", "and")
		lex.consume("or", "or")
		lex.consume("[a-zA-Z_][a-zA-Z_0-9]*", "id")
		lex.consume("[0-9]+", "number")
		c, _ := lex.next()
		if c == '"' {
			lex.consumeString()
		} else if c == ';' {
			lex.emit(";")
		} else if c == '\n' {
			lex.newLine()
		} else if strings.ContainsRune("=+-*/(){}<>,", c) {
			lex.emit(string(c), string(c))
		} else if !strings.ContainsRune(" \t\r\n", c) {
			fmt.Fprintf(os.Stderr, "unrecognized char '%c' at line %d, column %d\n", c, lex.line+1, lex.pos-lex.width-lex.lineIndex+1)
			os.Exit(1)
		}
	}

	lex.emit("eof")
	close(lex.out)
}

func lex(filepath string) <-chan string {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lex := newLexer(bytes)
	go lex.lex()
	return lex.out
}
