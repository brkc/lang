package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type parser struct {
	token  *token
	lexOut <-chan string
}

type token struct {
	symbol string
	line   int
	column int
	value  string
}

func newParser(lexOut <-chan string) *parser {
	return &parser{newTokenInfo(lexOut), lexOut}
}

func newTokenInfo(lexOut <-chan string) *token {
	var err error
	text := <-lexOut
	fields := regexp.MustCompile(" ").Split(text, 4)
	token := &token{}
	token.symbol = fields[0]
	token.line, err = strconv.Atoi(fields[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "err reading token: %s\n", err)
		os.Exit(1)
	}
	token.column, err = strconv.Atoi(fields[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "err reading token: %s\n", err)
		os.Exit(1)
	}
	if len(fields) == 4 {
		token.value = fields[3]
	}
	return token
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

func (p *parser) block() *block {
	var statements []statementVisitor
	for !p.accept("eof") && !p.accept("}") {
		statements = append(statements, p.statement())
	}
	return &block{statements}
}

func (p *parser) statement() statementVisitor {
	if p.accept("var") {
		return p.declaration()
	} else if p.accept("print") {
		return p.print()
	} else if p.accept("if") {
		return p.ifStatement()
	} else if p.accept("while") {
		return p.whileStatement()
	} else if p.accept("break") {
		return p.breakStatement()
	} else if p.accept("continue") {
		return p.continueStatement()
	} else if p.accept("fn") {
		return p.functionStatement()
	} else if p.accept("return") {
		return p.returnStatement()
	} else if p.accept("id") {
		var v statementVisitor
		id := p.expect("id")
		if p.accept("=") {
			v = p.assignment(id)
		} else if p.accept("(") {
			v = p.callExpression(id)
		}
		p.expect(";")
		return v
	} else {
		p.expect("var|print|if|while|fn|return")
		return nil
	}
}

func (p *parser) declaration() *declarationStatement {
	p.expect("var")
	id := p.expect("id")
	p.expect("=")
	n := p.booleanExpression()
	p.expect(";")
	return &declarationStatement{id, n}
}

func (p *parser) print() *printStatement {
	p.expect("print")
	expression := p.booleanExpression()
	p.expect(";")
	return &printStatement{expression}
}

func (p *parser) ifStatement() *ifStatement {
	p.expect("if")
	b := p.booleanExpression()
	p.expect("{")
	block := p.block()
	p.expect("}")
	return &ifStatement{b, block}
}

func (p *parser) whileStatement() *whileStatement {
	p.expect("while")
	b := p.booleanExpression()
	p.expect("{")
	block := p.block()
	p.expect("}")
	return &whileStatement{b, block}
}

func (p *parser) breakStatement() *breakStatement {
	p.expect("break")
	p.expect(";")
	return &breakStatement{}
}

func (p *parser) continueStatement() *continueStatement {
	p.expect("continue")
	p.expect(";")
	return &continueStatement{}
}

func (p *parser) functionStatement() *functionStatement {
	var parameters []string
	p.expect("fn")
	name := p.expect("id")
	p.expect("(")
	if p.accept("id") {
		parameters = append(parameters, p.expect("id"))
		for {
			if !p.accept(",") {
				break
			}
			p.expect(",")
			parameters = append(parameters, p.expect("id"))
		}
	}
	p.expect(")")
	p.expect("{")
	block := p.block()
	p.expect("}")
	return &functionStatement{name, parameters, block}
}

func (p *parser) returnStatement() *returnStatement {
	p.expect("return")
	if p.accept(";") {
		p.expect(";")
		return &returnStatement{nil}
	}
	b := p.booleanExpression()
	p.expect(";")
	return &returnStatement{b}
}

func (p *parser) assignment(id string) *assignmentStatement {
	p.expect("=")
	return &assignmentStatement{id, p.booleanExpression()}
}

func (p *parser) callExpression(id string) *callExpression {
	var arguments []*booleanExpression
	p.expect("(")
	for {
		if p.accept(")") {
			break
		}
		arguments = append(arguments, p.booleanExpression())
		if !p.accept(")") {
			p.expect(",")
		}
	}
	p.expect(")")
	return &callExpression{id, arguments}
}

func (p *parser) booleanExpression() *booleanExpression {
	b := &booleanExpression{p.andExpression(), "", nil}
	for {
		if p.accept("or") {
			p.expect("or")
			b = &booleanExpression{b, "or", p.andExpression()}
		} else {
			return b
		}
	}
}

func (p *parser) andExpression() *booleanExpression {
	b := &booleanExpression{p.condition(), "", nil}
	for {
		if p.accept("and") {
			p.expect("and")
			b = &booleanExpression{b, "and", p.condition()}
		} else {
			return b
		}
	}
}

func (p *parser) condition() *booleanExpression {
	var operator string
	left := p.logicalOperand()
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
		return &booleanExpression{left, "", nil}
	}
	return &booleanExpression{left, operator, p.logicalOperand()}
}

func (p *parser) logicalOperand() *logicalOperand {
	e := &logicalOperand{p.term(), "", nil}
	for {
		if p.accept("+") {
			p.expect("+")
			e = &logicalOperand{e, "+", p.term()}
		} else if p.accept("-") {
			p.expect("-")
			e = &logicalOperand{e, "-", p.term()}
		} else {
			return e
		}
	}
}

func (p *parser) term() *term {
	t := &term{p.logicalNotExpression(), "", nil}
	for {
		if p.accept("*") {
			p.expect("*")
			t = &term{t, "*", p.logicalNotExpression()}
		} else if p.accept("/") {
			p.expect("/")
			t = &term{t, "/", p.logicalNotExpression()}
		} else {
			return t
		}
	}
}

func (p *parser) logicalNotExpression() expressionVisitor {
	if p.accept("not") {
		p.expect("not")
		return &logicalNotExpression{p.logicalNotExpression()}
	}
	return p.atom()
}

func (p *parser) atom() expressionVisitor {
	if p.accept("id") {
		id := p.expect("id")
		if p.accept("(") {
			return p.callExpression(id)
		}
		return &identifier{id}
	} else if p.accept("number") {
		return &numberLiteral{p.expect("number")}
	} else if p.accept("string") {
		return &stringLiteral{p.expect("string")}
	} else if p.accept("true") {
		p.expect("true")
		return &booleanLiteral{true}
	} else if p.accept("false") {
		p.expect("false")
		return &booleanLiteral{false}
	} else if p.accept("(") {
		p.expect("(")
		n := p.booleanExpression()
		p.expect(")")
		return n
	} else {
		p.expect("id|number|string|true|false")
		return nil
	}
}

func parse(lexOut <-chan string) *block {
	return newParser(lexOut).block()
}
