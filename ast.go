package main

type (
	AssignmentStatement struct {
		id         string
		expression interface{}
	}

	Identifier struct {
		value string
	}

	MathExpression struct {
		left     interface{}
		operator string
		right    interface{}
	}

	NumberLiteral struct {
		value string
	}

	PrintStatement struct {
		expression interface{}
	}

	StringLiteral struct {
		value string
	}

	Term struct {
		left     interface{}
		operator string
		right    interface{}
	}
)
