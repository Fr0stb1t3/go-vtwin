package ast

import (
	"github.com/Fr0stb1t3/go-vtwin/token"
)

type Boolean struct {
	Token token.Token
	Value bool
}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

type IntegerLiteral struct {
	Token token.Token
	Value string
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.String()
}
func (il *IntegerLiteral) String() string {
	return il.Token.String()
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.String()
}
func (b *Boolean) String() string {
	return b.Token.String()
}
