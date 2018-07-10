package main

import (
	"fmt"
	"os"
	"strconv"
)

type (
	expression struct {
		typeValue expressionType
		value     interface{}
	}
	expressionType int
	statement      struct {
		typeValue  statementType
		expression *expression
	}
	statementType int
)

const (
	assignmentType statementType = 1 << iota
	blockType
	breakType
	callType
	continueType
	declarationType
	functionType
	ifType
	printType
	returnType
	whileType

	booleanType expressionType = 1 << iota
	numberType
	stringType
)

var (
	functions = map[string]*functionStatement{}
	types     = map[expressionType]string{
		numberType:  "number",
		stringType:  "string",
		booleanType: "boolean",
	}
	variables = map[string]*expression{}
)

func interpret(s string) {
	parse(lex(s)).visitStatement()
}

func (a *declarationStatement) visitStatement() *statement {
	variables[a.id] = a.expression.visitExpression()
	return &statement{declarationType, nil}
}

func (a *assignmentStatement) visitStatement() *statement {
	variables[a.id] = a.expression.visitExpression()
	return &statement{assignmentType, nil}
}

func (p *printStatement) visitStatement() *statement {
	v := p.expression.visitExpression()
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
	return &statement{printType, nil}
}

func (i *ifStatement) visitStatement() *statement {
	b := i.booleanExpression.visitExpression()
	typeCheck(booleanType, b)
	if b.value.(bool) {
		return i.block.visitStatement()
	}
	return &statement{ifType, nil}
}

func (i *whileStatement) visitStatement() *statement {
	for {
		b := i.booleanExpression.visitExpression()
		typeCheck(booleanType, b)
		if !b.value.(bool) {
			break
		}
		v := i.block.visitStatement()
		switch v.typeValue {
		case breakType, returnType:
			return v
		}
	}
	return &statement{whileType, nil}
}

func (i *breakStatement) visitStatement() *statement {
	return &statement{breakType, nil}
}

func (i *continueStatement) visitStatement() *statement {
	return &statement{continueType, nil}
}

func (f *functionStatement) visitStatement() *statement {
	functions[f.name] = f
	return &statement{functionType, nil}
}

func (r *returnStatement) visitStatement() *statement {
	return &statement{returnType, r.expression.visitExpression()}
}

func (b *block) visitStatement() *statement {
	for _, s := range b.statements {
		v := s.visitStatement()
		if v == nil {
			continue
		}
		switch v.typeValue {
		case returnType:
			return v
		}
	}
	return &statement{blockType, nil}
}

func (b *booleanExpression) visitExpression() *expression {
	left := b.left.visitExpression()
	if b.right == nil {
		return left
	}
	right := b.right.visitExpression()
	expectSameType(left, right)

	switch b.operator {
	case "&&":
		typeCheck(booleanType, left, right)
		return &expression{booleanType, left.value.(bool) && right.value.(bool)}
	case "||":
		typeCheck(booleanType, left, right)
		return &expression{booleanType, left.value.(bool) || right.value.(bool)}
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
	return &expression{booleanType, false}
}

func evaluateNumberComparison(left int, operator string, right int) *expression {
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
	return &expression{booleanType, b}
}

func evaluateStringComparison(left string, operator string, right string) *expression {
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
	return &expression{booleanType, b}
}

func evaluateBooleanComparison(left bool, operator string, right bool) *expression {
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
	return &expression{booleanType, b}
}

func (e *logicalOperand) visitExpression() *expression {
	left := e.left.visitExpression()
	if e.right != nil {
		right := e.right.visitExpression()
		typeCheck(numberType, left, right)
		switch e.operator {
		case "+":
			return &expression{numberType, left.value.(int) + right.value.(int)}
		case "-":
			return &expression{numberType, left.value.(int) - right.value.(int)}
		}
	}
	return e.left.visitExpression()
}

func (t *term) visitExpression() *expression {
	left := t.left.visitExpression()
	if t.right != nil {
		right := t.right.visitExpression()
		typeCheck(numberType, left, right)
		switch t.operator {
		case "*":
			return &expression{numberType, left.value.(int) * right.value.(int)}
		case "/":
			return &expression{numberType, left.value.(int) / right.value.(int)}
		}
	}
	return t.left.visitExpression()
}

func (e *logicalNotExpression) visitExpression() *expression {
	b := e.booleanExpression.visitExpression()
	typeCheck(booleanType, b)
	return &expression{booleanType, !b.value.(bool)}
}

func (c *callExpression) visitExpression() *expression {
	f := functions[c.name]
	for i, p := range f.parameters {
		variables[p] = c.arguments[i].visitExpression()
	}
	v := f.block.visitStatement()
	for _, p := range f.parameters {
		delete(variables, p)
	}
	switch v.typeValue {
	case returnType:
		return v.expression
	}
	return nil
}

func (c *callExpression) visitStatement() *statement {
	return &statement{callType, c.visitExpression()}
}

func (i *identifier) visitExpression() *expression {
	return variables[i.value]
}

func (nl *numberLiteral) visitExpression() *expression {
	n, err := strconv.Atoi(nl.value)
	if err != nil {
		fmt.Fprintln(os.Stderr, "expected number")
		os.Exit(1)
	}
	return &expression{numberType, n}
}

func (s *stringLiteral) visitExpression() *expression {
	return &expression{stringType, s.value}
}

func (b *booleanLiteral) visitExpression() *expression {
	return &expression{booleanType, b.value}
}

func typeCheck(b expressionType, args ...*expression) {
	for _, arg := range args {
		if arg.typeValue != b {
			fmt.Fprintf(os.Stderr, "type mismatch: %s != %s\n", types[arg.typeValue], types[b])
			os.Exit(1)
		}
	}
}

func expectSameType(args ...*expression) {
	var firstType expressionType
	for _, arg := range args {
		if firstType == 0 {
			firstType = arg.typeValue
			continue
		}
		typeCheck(firstType, arg)
	}
}
