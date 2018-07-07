package main

type (
	AssignmentStatement struct {
		id         string
		expression visitor
	}

	BooleanExpression struct {
		left     visitor
		operator string
		right    visitor
	}

	BooleanLiteral struct {
		value bool
	}

	DeclarationStatement struct {
		id         string
		expression visitor
	}

	Expression struct {
		left     visitor
		operator string
		right    visitor
	}

	Identifier struct {
		value string
	}

	IfStatement struct {
		booleanExpression *BooleanExpression
		block             []visitor
	}

	LogicalNotExpression struct {
		booleanExpression visitor
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

	WhileStatement struct {
		booleanExpression *BooleanExpression
		block             []visitor
	}
)
