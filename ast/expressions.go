package ast

import (
	"bytes"

	"github.com/Fr0stb1t3/go-vtwin/token"
)

type Node interface {
	Token() token.Token
	String() string
}

type Expression interface {
	Node
	expressionNode()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Node
	Operator string
	Right    Node
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string {
	return oe.Token.String()
}
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
