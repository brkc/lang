package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type state struct {
	token    string
	lexOut   <-chan string
	integers map[string]int
}

func newState(lexOut <-chan string) *state {
	return &state{<-lexOut, lexOut, map[string]int{}}
}

func (s *state) accept(expected string) bool {
	return strings.Fields(s.token)[0] == expected
}

func (s *state) expect(expected string) string {
	fields := strings.Fields(s.token)
	if fields[0] != expected {
		fmt.Fprintf(os.Stderr, "expected '%s', got '%s' at line %s, column %s\n", expected, fields[0], fields[1], fields[2])
		os.Exit(1)
	}
	s.token = <-s.lexOut
	if len(fields) > 3 {
		return fields[3]
	}
	return ""
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
	s.expect("\\n")
	s.integers[id] = n
}

func (s *state) print() {
	s.expect("print")
	n := s.expression()
	s.expect("\\n")
	fmt.Printf("%d\n", n)
}

func (s *state) expression() int {
	n := s.term()
	for {
		if s.accept("+") {
			s.expect("+")
			n += s.term()
		} else if s.accept("-") {
			s.expect("-")
			n -= s.term()
		} else {
			return n
		}
	}
}

func (s *state) term() int {
	n := s.atom()
	for {
		if s.accept("*") {
			s.expect("*")
			n *= s.atom()
		} else if s.accept("/") {
			s.expect("/")
			n /= s.atom()
		} else {
			return n
		}
	}
}

func (s *state) atom() int {
	if s.accept("id") {
		id := s.expect("id")
		return s.integers[id]
	} else if s.accept("number") {
		n, _ := strconv.Atoi(s.expect("number"))
		return n
	} else if s.accept("(") {
		s.expect("(")
		n := s.expression()
		s.expect(")")
		return n
	} else {
		s.expect("id|number")
		return 0
	}
}

// Parse executes the output from Lex
func Parse(lexOut <-chan string) {
	newState(lexOut).root()
}
