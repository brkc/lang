package main

import (
	"fmt"
	"os"
	"strconv"
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

var variables = map[string]interface{}{}

func Interpret(s string) {
	for _, s := range Parse(Lex(s)) {
		visit(s)
	}
}

func visit(i interface{}) {
	switch n := i.(type) {
	case *PrintStatement:
		n.visit()
	case *AssignmentStatement:
		n.visit()
	case *IfStatement:
		n.visit()
	}
}

func visitBlock(block []interface{}) {
	for _, b := range block {
		visit(b)
	}
}

func visitExpression(i interface{}) interface{} {
	switch n := i.(type) {
	case *StringLiteral:
		return n.visit()
	case *NumberLiteral:
		return n.visit()
	case *Identifier:
		return n.visit()
	case *MathExpression:
		return n.visit()
	default:
		fmt.Fprintf(os.Stderr, "unexpected type %T\n", n)
		os.Exit(1)
	}
	return nil
}

func visitMathExpression(i interface{}) int {
	switch n := i.(type) {
	case *MathExpression:
		return n.visit()
	case *Term:
		return n.visit()
	case *NumberLiteral:
		return n.visit()
	case *Identifier:
		return n.visit().(int)
	}
	return 0
}

func (a *AssignmentStatement) visit() {
	switch e := a.expression.(type) {
	case *StringLiteral:
		variables[a.id] = e.visit()
	case *NumberLiteral:
		variables[a.id] = e.visit()
	case *MathExpression:
		variables[a.id] = e.visit()
	case *Term:
		variables[a.id] = e.visit()
	case *Identifier:
		v := variables[a.id]
		variables[a.id] = v
	default:
		fmt.Fprintf(os.Stderr, "unrecognized type: %T\n", a.expression)
		os.Exit(1)
	}
}

func (p *PrintStatement) visit() {
	v := visitExpression(p.expression)
	switch d := v.(type) {
	case string:
		fmt.Printf("%s\n", d)
	case int:
		fmt.Printf("%d\n", d)
	default:
		fmt.Fprintf(os.Stderr, "unexpected type %T\n", v)
		os.Exit(1)
	}
}

func (i *IfStatement) visit() {
	var b bool
	left := visitExpression(i.left)
	right := visitExpression(i.right)
	switch left.(type) {
	case int:
		if _, ok := right.(int); !ok {
			fmt.Fprintf(os.Stderr, "type mismatch, int != %T\n", right)
			os.Exit(1)
		}
		b = evaluateNumberComparison(left.(int), i.operator, right.(int))
	case string:
		if _, ok := right.(string); !ok {
			fmt.Fprintf(os.Stderr, "type mismatch, string != %T\n", right)
			os.Exit(1)
		}
		b = evaluateStringComparison(left.(string), i.operator, right.(string))
	}
	if b {
		visitBlock(i.block)
	}
}

func evaluateNumberComparison(left int, operator string, right int) bool {
	switch operator {
	case "==":
		return left == right
	case "!=":
		return left != right
	case ">=":
		return left >= right
	case ">":
		return left > right
	case "<":
		return left < right
	case "<=":
		return left <= right
	default:
		fmt.Fprintln(os.Stderr, "unrecognized operator")
		os.Exit(1)
		return false
	}
}

func evaluateStringComparison(left string, operator string, right string) bool {
	switch operator {
	case "==":
		return left == right
	case "!=":
		return left != right
	case ">=":
		return left >= right
	case ">":
		return left > right
	case "<":
		return left < right
	case "<=":
		return left <= right
	default:
		fmt.Fprintln(os.Stderr, "unrecognized operator")
		os.Exit(1)
		return false
	}
}

func (e *MathExpression) visit() int {
	if e.right != nil {
		switch e.operator {
		case "+":
			return visitMathExpression(e.left) + visitMathExpression(e.right)
		case "-":
			return visitMathExpression(e.left) - visitMathExpression(e.right)
		}
	}
	return visitMathExpression(e.left)
}

func (t *Term) visit() int {
	if t.right != nil {
		switch t.operator {
		case "*":
			return visitMathExpression(t.left) * visitMathExpression(t.right)
		case "/":
			return visitMathExpression(t.left) / visitMathExpression(t.right)
		}
	}
	return visitMathExpression(t.left)
}

func (i *Identifier) visit() interface{} {
	return variables[i.value]
}

func (nl *NumberLiteral) visit() int {
	n, err := strconv.Atoi(nl.value)
	if err != nil {
		fmt.Fprintln(os.Stderr, "expected number")
		os.Exit(1)
	}
	return n
}

func (s *StringLiteral) visit() string {
	return s.value
}
