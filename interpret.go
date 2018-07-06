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
	booleanType
)

var (
	types = map[builtinType]string{
		numberType:  "number",
		stringType:  "string",
		booleanType: "boolean",
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
	case booleanType:
		fmt.Printf("%t\n", v.value.(bool))
	default:
		fmt.Fprintf(os.Stderr, "unexpected type %s\n", types[v.typeValue])
		os.Exit(1)
	}
	return nil
}

func (i *IfStatement) visit() *adt {
	b := i.booleanExpression.visit()
	typeCheck(booleanType, b)
	if b.value.(bool) {
		for _, s := range i.block {
			s.visit()
		}
	}
	return nil
}

func (b *BooleanExpression) visit() *adt {
	left := b.left.visit()
	if b.right == nil {
		return left
	}
	right := b.right.visit()
	expectSameType(left, right)

	switch b.operator {
	case "&&":
		typeCheck(booleanType, left, right)
		return &adt{booleanType, left.value.(bool) && right.value.(bool)}
	case "||":
		typeCheck(booleanType, left, right)
		return &adt{booleanType, left.value.(bool) || right.value.(bool)}
	}

	switch left.typeValue {
	case numberType:
		return evaluateNumberComparison(left.value.(int), b.operator, right.value.(int))
	case stringType:
		return evaluateStringComparison(left.value.(string), b.operator, right.value.(string))
	case booleanType:
		return evaluateBooleanComparison(left.value.(bool), b.operator, right.value.(bool))
	default:
		fmt.Fprintln(os.Stderr, "unrecognized type")
		os.Exit(1)
	}
	return &adt{booleanType, false}
}

func evaluateNumberComparison(left int, operator string, right int) *adt {
	var b bool
	switch operator {
	case "==":
		b = left == right
	case "!=":
		b = left != right
	case ">=":
		b = left >= right
	case ">":
		b = left > right
	case "<":
		b = left < right
	case "<=":
		b = left <= right
	default:
		fmt.Fprintln(os.Stderr, "unrecognized operator")
		os.Exit(1)
	}
	return &adt{booleanType, b}
}

func evaluateStringComparison(left string, operator string, right string) *adt {
	var b bool
	switch operator {
	case "==":
		b = left == right
	case "!=":
		b = left != right
	case ">=":
		b = left >= right
	case ">":
		b = left > right
	case "<":
		b = left < right
	case "<=":
		b = left <= right
	default:
		fmt.Fprintln(os.Stderr, "unrecognized operator")
		os.Exit(1)
	}
	return &adt{booleanType, b}
}

func evaluateBooleanComparison(left bool, operator string, right bool) *adt {
	var b bool
	switch operator {
	case "==":
		b = left == right
	case "!=":
		b = left != right
	default:
		fmt.Fprintln(os.Stderr, "unrecognized operator")
		os.Exit(1)
	}
	return &adt{booleanType, b}
}

func (e *MathExpression) visit() *adt {
	left := e.left.visit()
	if e.right != nil {
		right := e.right.visit()
		typeCheck(numberType, left, right)
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
	if t.right != nil {
		right := t.right.visit()
		typeCheck(numberType, left, right)
		switch t.operator {
		case "*":
			return &adt{numberType, left.value.(int) * right.value.(int)}
		case "/":
			return &adt{numberType, left.value.(int) / right.value.(int)}
		}
	}
	return t.left.visit()
}

func (e *LogicalNotExpression) visit() *adt {
	b := e.booleanExpression.visit()
	typeCheck(booleanType, b)
	return &adt{booleanType, !b.value.(bool)}
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

func (b *BooleanLiteral) visit() *adt {
	return &adt{booleanType, b.value}
}

func typeCheck(b builtinType, args ...*adt) {
	for _, arg := range args {
		if arg.typeValue != b {
			fmt.Fprintf(os.Stderr, "type mismatch: %s != %s\n", types[arg.typeValue], types[b])
			os.Exit(1)
		}
	}
}

func expectSameType(args ...*adt) {
	var firstType builtinType
	for _, arg := range args {
		if firstType == 0 {
			firstType = arg.typeValue
			continue
		}
		typeCheck(firstType, arg)
	}
}
