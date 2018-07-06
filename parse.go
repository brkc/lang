package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type parser struct {
	token  *tokenInfo
	lexOut <-chan string
}

type tokenInfo struct {
	symbol string
	line   int
	column int
	value  string
}

func newParser(lexOut <-chan string) *parser {
	return &parser{newTokenInfo(lexOut), lexOut}
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

func (p *parser) accept(expected string) bool {
	return p.token.symbol == expected
}

func (p *parser) expect(expected string) string {
	if p.token.symbol != expected {
		fmt.Fprintf(os.Stderr, "expected '%s', got '%s' at line %d, column %d\n", expected, p.token.symbol, p.token.line, p.token.column)
		os.Exit(1)
	}
	value := p.token.value
	p.token = newTokenInfo(p.lexOut)
	return value
}

func (p *parser) block() []visitor {
	statements := make([]visitor, 0)
	for !p.accept("eof") && !p.accept("}") {
		statements = append(statements, p.statement())
	}
	return statements
}

func (p *parser) statement() visitor {
	if p.accept("let") {
		return p.assignment()
	} else if p.accept("print") {
		return p.print()
	} else if p.accept("if") {
		return p.ifStatement()
	} else {
		p.expect("let|print")
		return nil
	}
}

func (p *parser) assignment() *AssignmentStatement {
	p.expect("let")
	id := p.expect("id")
	p.expect("=")
	n := p.booleanExpression()
	p.expect(";")
	return &AssignmentStatement{id, n}
}

func (p *parser) print() *PrintStatement {
	p.expect("print")
	expression := p.booleanExpression()
	p.expect(";")
	return &PrintStatement{expression}
}

func (p *parser) ifStatement() *IfStatement {
	p.expect("if")
	b := p.booleanExpression()
	p.expect("{")
	block := p.block()
	p.expect("}")
	return &IfStatement{b, block}
}

func (p *parser) booleanExpression() *BooleanExpression {
	b := &BooleanExpression{p.andExpression(), "", nil}
	for {
		if p.accept("||") {
			p.expect("||")
			b = &BooleanExpression{b, "||", p.andExpression()}
		} else {
			return b
		}
	}
}

func (p *parser) andExpression() *BooleanExpression {
	b := &BooleanExpression{p.condition(), "", nil}
	for {
		if p.accept("&&") {
			p.expect("&&")
			b = &BooleanExpression{b, "&&", p.condition()}
		} else {
			return b
		}
	}
}

func (p *parser) condition() *BooleanExpression {
	var operator string
	left := p.expression()
	if p.accept("==") {
		operator = p.expect("==")
	} else if p.accept("!=") {
		operator = p.expect("!=")
	} else if p.accept(">=") {
		operator = p.expect(">=")
	} else if p.accept(">") {
		operator = p.expect(">")
	} else if p.accept("<") {
		operator = p.expect("<")
	} else if p.accept("<=") {
		operator = p.expect("<=")
	} else {
		return &BooleanExpression{left, "", nil}
	}
	return &BooleanExpression{left, operator, p.expression()}
}

func (p *parser) expression() visitor {
	if p.accept("string") {
		return &StringLiteral{p.expect("string")}
	} else if p.accept("true") {
		p.expect("true")
		return &BooleanLiteral{true}
	} else if p.accept("false") {
		p.expect("false")
		return &BooleanLiteral{false}
	} else if p.accept("id") {
		return &Identifier{p.expect("id")}
	}
	return p.mathExpression()
}

func (p *parser) mathExpression() *MathExpression {
	e := &MathExpression{p.term(), "", nil}
	for {
		if p.accept("+") {
			p.expect("+")
			e = &MathExpression{e, "+", p.term()}
		} else if p.accept("-") {
			p.expect("-")
			e = &MathExpression{e, "-", p.term()}
		} else {
			return e
		}
	}
}

func (p *parser) term() *Term {
	t := &Term{p.atom(), "", nil}
	for {
		if p.accept("*") {
			p.expect("*")
			t = &Term{t, "*", p.atom()}
		} else if p.accept("/") {
			p.expect("/")
			t = &Term{t, "/", p.atom()}
		} else {
			return t
		}
	}
}

func (p *parser) atom() visitor {
	if p.accept("id") {
		return &Identifier{p.expect("id")}
	} else if p.accept("number") {
		return &NumberLiteral{p.expect("number")}
	} else if p.accept("(") {
		p.expect("(")
		n := p.booleanExpression()
		p.expect(")")
		return n
	} else if p.accept("!") {
		return p.logicalNotExpression()
	} else {
		p.expect("id|number")
		return nil
	}
}

func (p *parser) logicalNotExpression() visitor {
	p.expect("!")
	return &LogicalNotExpression{p.booleanExpression()}
}

// Parse executes the output from Lex
func Parse(lexOut <-chan string) []visitor {
	return newParser(lexOut).block()
}
