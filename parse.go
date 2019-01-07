package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type (
	parser struct {
		token  *token
		lexOut <-chan string
	}
	token struct {
		symbol string
		line   int
		column int
		value  string
	}
)

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

func (p *parser) block(scope *scope) *block {
	var statements []statementVisitor
	for !p.accept("eof") && !p.accept("}") {
		statements = append(statements, p.statement(scope))
	}
	return &block{statements}
}

func (p *parser) statement(scope *scope) statementVisitor {
	if p.accept("var") {
		return p.declaration(scope)
	} else if p.accept("if") {
		return p.ifStatement(scope)
	} else if p.accept("while") {
		return p.whileStatement(scope)
	} else if p.accept("break") {
		return p.breakStatement(scope)
	} else if p.accept("continue") {
		return p.continueStatement(scope)
	} else if p.accept("fn") {
		return p.functionStatement(scope)
	} else if p.accept("return") {
		return p.returnStatement(scope)
	} else if p.accept("id") {
		var v statementVisitor
		id := p.expect("id")
		if p.accept("=") {
			v = p.assignment(scope, id)
		} else if p.accept("(") {
			v = p.callExpression(scope, id)
		}
		p.expect(";")
		return v
	} else {
		p.expect("var|if|while|fn|return")
		return nil
	}
}

func (p *parser) declaration(scope *scope) *declarationStatement {
	p.expect("var")
	id := p.expect("id")
	p.expect("=")
	n := p.booleanExpression(scope)
	p.expect(";")
	scope.declare(id, true)
	return &declarationStatement{id, n}
}

func (p *parser) ifStatement(scope *scope) *ifStatement {
	p.expect("if")
	b := p.booleanExpression(scope)
	p.expect("{")
	block := p.block(newScope(scope))
	p.expect("}")
	return &ifStatement{b, block}
}

func (p *parser) whileStatement(scope *scope) *whileStatement {
	p.expect("while")
	b := p.booleanExpression(scope)
	p.expect("{")
	block := p.block(newScope(scope))
	p.expect("}")
	return &whileStatement{b, block}
}

func (p *parser) breakStatement(scope *scope) *breakStatement {
	p.expect("break")
	p.expect(";")
	return &breakStatement{}
}

func (p *parser) continueStatement(scope *scope) *continueStatement {
	p.expect("continue")
	p.expect(";")
	return &continueStatement{}
}

func (p *parser) functionStatement(scope *scope) *functionStatement {
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
	newScope := newScope(scope)
	for _, p := range parameters {
		newScope.declare(p, true)
	}
	block := p.block(newScope)
	p.expect("}")
	return &functionStatement{name, parameters, block}
}

func (p *parser) returnStatement(scope *scope) *returnStatement {
	p.expect("return")
	if p.accept(";") {
		p.expect(";")
		return &returnStatement{nil}
	}
	b := p.booleanExpression(scope)
	p.expect(";")
	return &returnStatement{b}
}

func (p *parser) assignment(scope *scope, id string) *assignmentStatement {
	p.expect("=")
	return &assignmentStatement{id, p.booleanExpression(scope)}
}

func (p *parser) callExpression(scope *scope, id string) *callExpression {
	var arguments []expressionVisitor
	p.expect("(")
	for {
		if p.accept(")") {
			break
		}
		arguments = append(arguments, p.booleanExpression(scope))
		if !p.accept(")") {
			p.expect(",")
		}
	}
	p.expect(")")
	return &callExpression{id, arguments}
}

func (p *parser) booleanExpression(scope *scope) expressionVisitor {
	b := p.andExpression(scope)
	for {
		if p.accept("or") {
			p.expect("or")
			b = &booleanExpression{b, "or", p.andExpression(scope)}
		} else {
			return b
		}
	}
}

func (p *parser) andExpression(scope *scope) expressionVisitor {
	b := p.condition(scope)
	for {
		if p.accept("and") {
			p.expect("and")
			b = &booleanExpression{b, "and", p.condition(scope)}
		} else {
			return b
		}
	}
}

func (p *parser) condition(scope *scope) expressionVisitor {
	var operator string
	left := p.logicalOperand(scope)
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
		return left
	}
	return &booleanExpression{left, operator, p.logicalOperand(scope)}
}

func (p *parser) logicalOperand(scope *scope) expressionVisitor {
	e := p.term(scope)
	for {
		if p.accept("+") {
			p.expect("+")
			e = &logicalOperand{e, "+", p.term(scope)}
		} else if p.accept("-") {
			p.expect("-")
			e = &logicalOperand{e, "-", p.term(scope)}
		} else {
			return e
		}
	}
}

func (p *parser) term(scope *scope) expressionVisitor {
	t := p.logicalNotExpression(scope)
	for {
		if p.accept("*") {
			p.expect("*")
			t = &term{t, "*", p.logicalNotExpression(scope)}
		} else if p.accept("/") {
			p.expect("/")
			t = &term{t, "/", p.logicalNotExpression(scope)}
		} else {
			return t
		}
	}
}

func (p *parser) logicalNotExpression(scope *scope) expressionVisitor {
	if p.accept("not") {
		p.expect("not")
		return &logicalNotExpression{p.logicalNotExpression(scope)}
	}
	return p.atom(scope)
}

func (p *parser) atom(scope *scope) expressionVisitor {
	if p.accept("id") {
		line, column := p.token.line, p.token.column
		id := p.expect("id")
		if p.accept("(") {
			return p.callExpression(scope, id)
		}
		if scope.resolve(id) == nil {
			fmt.Fprintf(os.Stderr, "unrecognized var '%s' at line %d, column %d\n", id, line, column)
			os.Exit(1)
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
		n := p.booleanExpression(scope)
		p.expect(")")
		return n
	} else {
		p.expect("id|number|string|true|false")
		return nil
	}
}

func parse(lexOut <-chan string) *block {
	return newParser(lexOut).block(newScope(nil))
}
