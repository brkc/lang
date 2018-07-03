package main

type (
	AssignmentStatement struct {
		id         string
		expression visitor
	}

	Identifier struct {
		value string
	}

	IfStatement struct {
		left     visitor
		operator string
		right    visitor
		block    []visitor
	}

	MathExpression struct {
		left     visitor
		operator string
		right    visitor
	}

	NumberLiteral struct {
		value string
	}

	PrintStatement struct {
		expression visitor
	}

	StringLiteral struct {
		value string
	}

	Term struct {
		left     visitor
		operator string
		right    visitor
	}

	visitor interface {
		visit() *adt
	}
)
