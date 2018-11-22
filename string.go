package main

import (
	"bytes"
	"fmt"
)

func (d *declarationStatement) String() string {
	return fmt.Sprintf("(declaration %s)", d.expression)
}

func (a *assignmentStatement) String() string {
	return fmt.Sprintf("(assignment %s %s)", a.id, a.expression)
}

func (i *ifStatement) String() string {
	return fmt.Sprintf("(if %s %s)", i.booleanExpression, i.block)
}

func (w *whileStatement) String() string {
	return fmt.Sprintf("(while %s %s)", w.booleanExpression, w.block)
}

func (b *breakStatement) String() string {
	return fmt.Sprintf("(break)")
}

func (c *continueStatement) String() string {
	return fmt.Sprintf("(continue)")
}

func (f *functionStatement) String() string {
	var buf bytes.Buffer
	if len(f.parameters) > 0 {
		for i, param := range f.parameters {
			if i != 0 {
				buf.WriteRune(' ')
			}
			buf.WriteString(param)
		}
	} else {
		buf.WriteString("nil")
	}
	return fmt.Sprintf("(function %s %s %s)", f.name, buf.String(), f.block)
}

func (r *returnStatement) String() string {
	return fmt.Sprintf("(return %s)", r.expression)
}

func (b *block) String() string {
	var buf bytes.Buffer
	if len(b.statements) > 0 {
		for i, statement := range b.statements {
			if i != 0 {
				buf.WriteRune(' ')
			}
			buf.WriteString(statement.String())
		}
	} else {
		buf.WriteString("nil")
	}
	return fmt.Sprintf("(block %s)", buf.String())
}

func (b *booleanExpression) String() string {
	if b.right != nil {
		return fmt.Sprintf("(booleanExpression %s %s %s)", b.left, b.operator, b.right)
	}
	return fmt.Sprintf("(booleanExpression %s)", b.left)
}

func (l *logicalOperand) String() string {
	if l.right != nil {
		return fmt.Sprintf("(logicalOperand %s %s %s)", l.left, l.operator, l.right)
	}
	return fmt.Sprintf("(logicalOperand %s)", l.left)
}

func (t *term) String() string {
	if t.right != nil {
		return fmt.Sprintf("(term %s %s %s)", t.left, t.operator, t.right)
	}
	return fmt.Sprintf("(term %s)", t.left)
}

func (l *logicalNotExpression) String() string {
	return fmt.Sprintf("(logicalNotExpression %s)", l.booleanExpression)
}

func (c *callExpression) String() string {
	var buf bytes.Buffer
	if len(c.arguments) > 0 {
		for i, arg := range c.arguments {
			if i != 0 {
				buf.WriteRune(' ')
			}
			buf.WriteString(arg.String())
		}
	} else {
		buf.WriteString("nil")
	}
	return fmt.Sprintf("(callExpression %s %s)", c.name, buf.String())
}

func (i *identifier) String() string {
	return fmt.Sprintf("(identifier %s)", i.value)
}

func (n *numberLiteral) String() string {
	return fmt.Sprintf("(numberLiteral %s)", n.value)
}

func (s *stringLiteral) String() string {
	return fmt.Sprintf("(stringLiteral \"%s\")", s.value)
}

func (b *booleanLiteral) String() string {
	return fmt.Sprintf("(booleanLiteral %t)", b.value)
}
