package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type state struct {
	token  *tokenInfo
	lexOut <-chan string
}

type tokenInfo struct {
	symbol string
	line   int
	column int
	value  string
}

func newState(lexOut <-chan string) *state {
	return &state{newTokenInfo(lexOut), lexOut}
}

func newTokenInfo(lexOut <-chan string) *tokenInfo {
	var err error
	text := <-lexOut
	fields := regexp.MustCompile(" ").Split(text, 4)
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

func (s *state) expect(expected string) string {
	if s.token.symbol != expected {
		fmt.Fprintf(os.Stderr, "expected '%s', got '%s' at line %d, column %d\n", expected, s.token.symbol, s.token.line, s.token.column)
		os.Exit(1)
	}
	value := s.token.value
	s.token = newTokenInfo(s.lexOut)
	return value
}

func (s *state) block() []interface{} {
	statements := make([]interface{}, 0)
	for !s.accept("eof") && !s.accept("}") {
		statements = append(statements, s.statement())
	}
	return statements
}

func (s *state) statement() interface{} {
	if s.accept("let") {
		return s.assignment()
	} else if s.accept("print") {
		return s.print()
	} else if s.accept("if") {
		return s.ifStatement()
	} else {
		s.expect("let|print")
		return nil
	}
}

func (s *state) assignment() *AssignmentStatement {
	s.expect("let")
	id := s.expect("id")
	s.expect("=")
	n := s.expression()
	s.expect(";")
	return &AssignmentStatement{id, n}
}

func (s *state) print() *PrintStatement {
	s.expect("print")
	expression := s.expression()
	s.expect(";")
	return &PrintStatement{expression}
}

func (s *state) ifStatement() *IfStatement {
	var operator string
	s.expect("if")
	left := s.expression()
	if s.accept("==") {
		operator = s.expect("==")
	} else if s.accept("!=") {
		operator = s.expect("!=")
	} else if s.accept(">=") {
		operator = s.expect(">=")
	} else if s.accept(">") {
		operator = s.expect(">")
	} else if s.accept("<") {
		operator = s.expect("<")
	} else if s.accept("<=") {
		operator = s.expect("<=")
	}
	right := s.expression()
	s.expect("{")
	block := s.block()
	s.expect("}")
	return &IfStatement{left, operator, right, block}
}

func (s *state) expression() interface{} {
	if s.accept("string") {
		return &StringLiteral{s.expect("string")}
	} else if s.accept("id") {
		return &Identifier{s.expect("id")}
	}
	return s.mathExpression()
}

func (s *state) mathExpression() *MathExpression {
	e := &MathExpression{s.term(), "", nil}
	for {
		if s.accept("+") {
			s.expect("+")
			e = &MathExpression{e, "+", s.term()}
		} else if s.accept("-") {
			s.expect("-")
			e = &MathExpression{e, "-", s.term()}
		} else {
			return e
		}
	}
}

func (s *state) term() *Term {
	t := &Term{s.atom(), "", nil}
	for {
		if s.accept("*") {
			s.expect("*")
			t = &Term{t, "*", s.atom()}
		} else if s.accept("/") {
			s.expect("/")
			t = &Term{t, "/", s.atom()}
		} else {
			return t
		}
	}
}

func (s *state) atom() interface{} {
	if s.accept("id") {
		return &Identifier{s.expect("id")}
	} else if s.accept("number") {
		return &NumberLiteral{s.expect("number")}
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
func Parse(lexOut <-chan string) []interface{} {
	return newState(lexOut).block()
}
