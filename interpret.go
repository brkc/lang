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
	variableScope struct {
		variables map[string]*expression
		parent    *variableScope
	}
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
	rootScope = newVariableScope()
	types     = map[expressionType]string{
		numberType:  "number",
		stringType:  "string",
		booleanType: "boolean",
	}
)

func interpret(file string) {
	parse(lex(file)).visitStatement(rootScope)
}

func newVariableScope() *variableScope {
	return &variableScope{map[string]*expression{}, nil}
}

func (a *declarationStatement) visitStatement(scope *variableScope) *statement {
	scope.variables[a.id] = a.expression.visitExpression(scope)
	return &statement{declarationType, nil}
}

func (a *assignmentStatement) visitStatement(scope *variableScope) *statement {
	for scope != nil {
		if _, ok := scope.variables[a.id]; ok {
			scope.variables[a.id] = a.expression.visitExpression(scope)
			return &statement{assignmentType, nil}
		}
		scope = scope.parent
	}
	fmt.Fprintf(os.Stderr, "unrecognized var: '%s'\n", a.id)
	os.Exit(1)
	return nil
}

func (i *ifStatement) visitStatement(scope *variableScope) *statement {
	b := i.booleanExpression.visitExpression(scope)
	typeCheck(booleanType, b)
	if b.value.(bool) {
		newScope := newVariableScope()
		newScope.parent = scope
		return i.block.visitStatement(newScope)
	}
	return &statement{ifType, nil}
}

func (i *whileStatement) visitStatement(scope *variableScope) *statement {
	for {
		b := i.booleanExpression.visitExpression(scope)
		typeCheck(booleanType, b)
		if !b.value.(bool) {
			break
		}
		newScope := newVariableScope()
		newScope.parent = scope
		v := i.block.visitStatement(newScope)
		switch v.typeValue {
		case breakType, returnType:
			return v
		}
	}
	return &statement{whileType, nil}
}

func (i *breakStatement) visitStatement(scope *variableScope) *statement {
	return &statement{breakType, nil}
}

func (i *continueStatement) visitStatement(scope *variableScope) *statement {
	return &statement{continueType, nil}
}

func (f *functionStatement) visitStatement(scope *variableScope) *statement {
	functions[f.name] = f
	return &statement{functionType, nil}
}

func (r *returnStatement) visitStatement(scope *variableScope) *statement {
	return &statement{returnType, r.expression.visitExpression(scope)}
}

func (b *block) visitStatement(scope *variableScope) *statement {
	for _, s := range b.statements {
		v := s.visitStatement(scope)
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

func (b *booleanExpression) visitExpression(scope *variableScope) *expression {
	left := b.left.visitExpression(scope)
	if b.right == nil {
		return left
	}
	right := b.right.visitExpression(scope)
	expectSameType(left, right)

	switch b.operator {
	case "and":
		typeCheck(booleanType, left, right)
		return &expression{booleanType, left.value.(bool) && right.value.(bool)}
	case "or":
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

func (e *logicalOperand) visitExpression(scope *variableScope) *expression {
	left := e.left.visitExpression(scope)
	if e.right != nil {
		right := e.right.visitExpression(scope)
		typeCheck(numberType, left, right)
		switch e.operator {
		case "+":
			return &expression{numberType, left.value.(int) + right.value.(int)}
		case "-":
			return &expression{numberType, left.value.(int) - right.value.(int)}
		}
	}
	return e.left.visitExpression(scope)
}

func (t *term) visitExpression(scope *variableScope) *expression {
	left := t.left.visitExpression(scope)
	if t.right != nil {
		right := t.right.visitExpression(scope)
		typeCheck(numberType, left, right)
		switch t.operator {
		case "*":
			return &expression{numberType, left.value.(int) * right.value.(int)}
		case "/":
			return &expression{numberType, left.value.(int) / right.value.(int)}
		}
	}
	return t.left.visitExpression(scope)
}

func (e *logicalNotExpression) visitExpression(scope *variableScope) *expression {
	b := e.booleanExpression.visitExpression(scope)
	typeCheck(booleanType, b)
	return &expression{booleanType, !b.value.(bool)}
}

func (c *callExpression) visitExpression(scope *variableScope) *expression {
	f, ok := functions[c.name]
	if !ok {
		expr, err := visitBuiltin(c, scope)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return expr
	}
	newScope := newVariableScope()
	newScope.parent = scope
	for i, p := range f.parameters {
		newScope.variables[p] = c.arguments[i].visitExpression(scope)
	}
	v := f.block.visitStatement(newScope)
	switch v.typeValue {
	case returnType:
		return v.expression
	}
	return nil
}

func (c *callExpression) visitStatement(scope *variableScope) *statement {
	return &statement{callType, c.visitExpression(scope)}
}

func (i *identifier) visitExpression(scope *variableScope) *expression {
	for scope != nil {
		if value, ok := scope.variables[i.value]; ok {
			return value
		}
		scope = scope.parent
	}
	return nil
}

func (nl *numberLiteral) visitExpression(scope *variableScope) *expression {
	n, err := strconv.Atoi(nl.value)
	if err != nil {
		fmt.Fprintln(os.Stderr, "expected number")
		os.Exit(1)
	}
	return &expression{numberType, n}
}

func (s *stringLiteral) visitExpression(scope *variableScope) *expression {
	return &expression{stringType, s.value}
}

func (b *booleanLiteral) visitExpression(scope *variableScope) *expression {
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

func visitBuiltin(c *callExpression, scope *variableScope) (*expression, error) {
	var args []*expression
	for _, arg := range c.arguments {
		args = append(args, arg.visitExpression(scope))
	}
	return builtin(c.name, args)
}
