package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type builtinType int

const (
	numberType builtinType = 1 << iota
	stringType
)

var types = map[builtinType]string{
	numberType: "number",
	stringType: "string",
}

type state struct {
	token     *tokenInfo
	lexOut    <-chan string
	variables map[string]*valueInfo
}

type valueInfo struct {
	token     *tokenInfo
	typeValue builtinType
	value     interface{}
}

type tokenInfo struct {
	symbol string
	line   int
	column int
	value  string
}

func newState(lexOut <-chan string) *state {
	return &state{newTokenInfo(lexOut), lexOut, map[string]*valueInfo{}}
}

func newValueInfo(token *tokenInfo, typeValue builtinType, value interface{}) *valueInfo {
	return &valueInfo{token, typeValue, value}
}

func newTokenInfo(lexOut <-chan string) *tokenInfo {
	var err error
	text := <-lexOut
	fields := strings.Fields(text)
	tokenInfo := &tokenInfo{}
	tokenInfo.symbol = fields[0]
	tokenInfo.line, err = strconv.Atoi(fields[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "err reading token: %s\n", err)
		os.Exit(1)
	}
	tokenInfo.column, err = strconv.Atoi(fields[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "err reading token: %s\n", err)
		os.Exit(1)
	}
	if len(fields) == 4 {
		tokenInfo.value = fields[3]
	}
	return tokenInfo
}

func (s *state) accept(expected string) bool {
	return s.token.symbol == expected
}

func (s *state) acceptType(t builtinType) bool {
	return s.variables[s.token.value].typeValue == t
}

func (s *state) expect(expected string) string {
	if s.token.symbol != expected {
		fmt.Fprintf(os.Stderr, "expected '%s', got '%s' at line %d, column %d\n", expected, s.token.symbol, s.token.line, s.token.column)
		os.Exit(1)
	}
	value := s.token.value
	s.token = newTokenInfo(s.lexOut)
	return value
}

func (s *state) expectType(t builtinType, v *valueInfo) {
	if v.typeValue != t {
		fmt.Fprintf(os.Stderr, "expected '%s', got '%s' at line %d, column %d\n", types[t], types[v.typeValue], v.token.line, v.token.column)
		os.Exit(1)
	}
}

func (s *state) expectVariable() *valueInfo {
	token := s.token
	variable := s.variables[s.expect("id")]
	return newValueInfo(token, variable.typeValue, variable.value)
}

func (s *state) expectNumber() *valueInfo {
	token := s.token
	n, err := strconv.Atoi(s.expect("number"))
	if err != nil {
		fmt.Printf("err getting number: %s\n", err)
		os.Exit(1)
	}
	return newValueInfo(token, numberType, n)
}

func (s *state) expectString() *valueInfo {
	token := s.token
	return newValueInfo(token, stringType, s.expect("string"))
}

func (s *state) root() {
	for !s.accept("eof") {
		if s.accept("let") {
			s.assignment()
		} else if s.accept("print") {
			s.print()
		} else {
			s.expect("let|number")
		}
	}
}

func (s *state) statement() {
	if s.accept("let") {
		s.assignment()
	} else if s.accept("print") {
		s.print()
	} else {
		s.expect("let|print")
	}
}

func (s *state) assignment() {
	s.expect("let")
	id := s.expect("id")
	s.expect("=")
	n := s.expression()
	s.expect(";")
	s.variables[id] = n
}

func (s *state) print() {
	s.expect("print")
	expression := s.expression()
	s.expect(";")
	switch expression.typeValue {
	case numberType:
		fmt.Printf("%d\n", expression.value.(int))
	case stringType:
		fmt.Printf("%s\n", expression.value.(string))
	}
}

func (s *state) expression() *valueInfo {
	if s.accept("string") {
		return s.expectString()
	} else if s.accept("id") && s.acceptType(stringType) {
		return s.expectVariable()
	}
	return s.mathExpression()
}

func (s *state) mathExpression() *valueInfo {
	term := s.term()
	s.expectType(numberType, term)
	for {
		if s.accept("+") {
			s.expect("+")
			right := s.term()
			s.expectType(numberType, right)
			term.value = term.value.(int) + right.value.(int)
		} else if s.accept("-") {
			s.expect("-")
			right := s.term()
			s.expectType(numberType, right)
			term.value = term.value.(int) - right.value.(int)
		} else {
			return term
		}
	}
}

func (s *state) term() *valueInfo {
	atom := s.atom()
	s.expectType(numberType, atom)
	for {
		if s.accept("*") {
			s.expect("*")
			right := s.atom()
			s.expectType(numberType, right)
			atom.value = atom.value.(int) * right.value.(int)
		} else if s.accept("/") {
			s.expect("/")
			right := s.atom()
			s.expectType(numberType, right)
			atom.value = atom.value.(int) / right.value.(int)
		} else {
			return atom
		}
	}
}

func (s *state) atom() *valueInfo {
	if s.accept("id") {
		return s.expectVariable()
	} else if s.accept("number") {
		return s.expectNumber()
	} else if s.accept("(") {
		s.expect("(")
		n := s.mathExpression()
		s.expect(")")
		return n
	} else {
		s.expect("id|number")
		return nil
	}
}

// Parse executes the output from Lex
func Parse(lexOut <-chan string) {
	newState(lexOut).root()
}
