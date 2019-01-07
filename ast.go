package main

type (
	assignmentStatement struct {
		id         string
		expression expressionVisitor
	}

	block struct {
		statements []statementVisitor
	}

	booleanExpression struct {
		left     expressionVisitor
		operator string
		right    expressionVisitor
	}

	booleanLiteral struct {
		value bool
	}

	breakStatement struct {
	}

	callExpression struct {
		name      string
		arguments []expressionVisitor
	}

	continueStatement struct {
	}

	declarationStatement struct {
		id         string
		expression expressionVisitor
	}

	expressionVisitor interface {
		String() string
		visitExpression(scope *scope) *expression
	}

	functionStatement struct {
		name       string
		parameters []string
		block      *block
	}

	identifier struct {
		value string
	}

	ifStatement struct {
		booleanExpression expressionVisitor
		block             *block
	}

	logicalNotExpression struct {
		booleanExpression expressionVisitor
	}

	logicalOperand struct {
		left     expressionVisitor
		operator string
		right    expressionVisitor
	}

	numberLiteral struct {
		value string
	}

	returnStatement struct {
		expression expressionVisitor
	}

	statementVisitor interface {
		String() string
		visitStatement(scope *scope) *statement
	}

	stringLiteral struct {
		value string
	}

	term struct {
		left     expressionVisitor
		operator string
		right    expressionVisitor
	}

	whileStatement struct {
		booleanExpression expressionVisitor
		block             *block
	}
)
