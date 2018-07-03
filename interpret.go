package main

import (
	"fmt"
	"os"
	"strconv"
)

type (
	adt struct {
		typeValue builtinType
		value     interface{}
	}
	builtinType int
)

const (
	numberType builtinType = 1 << iota
	stringType
)

var (
	types = map[builtinType]string{
		numberType: "number",
		stringType: "string",
	}
	variables = map[string]*adt{}
)

func Interpret(s string) {
	for _, s := range Parse(Lex(s)) {
		s.visit()
	}
}

func (a *AssignmentStatement) visit() *adt {
	variables[a.id] = a.expression.visit()
	return nil
}

func (p *PrintStatement) visit() *adt {
	v := p.expression.visit()
	switch v.typeValue {
	case stringType:
		fmt.Printf("%s\n", v.value.(string))
	case numberType:
		fmt.Printf("%d\n", v.value.(int))
	default:
		fmt.Fprintf(os.Stderr, "unexpected type %s\n", types[v.typeValue])
		os.Exit(1)
	}
	return nil
}

func (i *IfStatement) visit() *adt {
	var b bool
	left := i.left.visit()
	right := i.right.visit()
	switch left.typeValue {
	case numberType:
		typeCheck(numberType, left, right)
		b = evaluateNumberComparison(left.value.(int), i.operator, right.value.(int))
	case stringType:
		typeCheck(stringType, left, right)
		b = evaluateStringComparison(left.value.(string), i.operator, right.value.(string))
	}
	if b {
		for _, s := range i.block {
			s.visit()
		}
	}
	return nil
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

func (e *MathExpression) visit() *adt {
	left := e.left.visit()
	typeCheck(numberType, left)
	if e.right != nil {
		right := e.right.visit()
		typeCheck(numberType, right)
		switch e.operator {
		case "+":
			return &adt{numberType, left.value.(int) + right.value.(int)}
		case "-":
			return &adt{numberType, left.value.(int) - right.value.(int)}
		}
	}
	return e.left.visit()
}

func (t *Term) visit() *adt {
	left := t.left.visit()
	typeCheck(numberType, left)
	if t.right != nil {
		right := t.right.visit()
		typeCheck(numberType, right)
		switch t.operator {
		case "*":
			return &adt{numberType, left.value.(int) * right.value.(int)}
		case "/":
			return &adt{numberType, left.value.(int) / right.value.(int)}
		}
	}
	return t.left.visit()
}

func (i *Identifier) visit() *adt {
	return variables[i.value]
}

func (nl *NumberLiteral) visit() *adt {
	n, err := strconv.Atoi(nl.value)
	if err != nil {
		fmt.Fprintln(os.Stderr, "expected number")
		os.Exit(1)
	}
	return &adt{numberType, n}
}

func (s *StringLiteral) visit() *adt {
	return &adt{stringType, s.value}
}

func typeCheck(b builtinType, args ...*adt) {
	for _, arg := range args {
		if arg.typeValue != b {
			fmt.Fprintf(os.Stderr, "type mismatch: %s != %s\n", types[arg.typeValue], types[b])
			os.Exit(1)
		}
	}
}
