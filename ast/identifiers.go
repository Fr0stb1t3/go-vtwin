package ast

import (
	"github.com/Fr0stb1t3/go-vtwin/token"
)

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

type ConstStatement struct {
	Token token.Token // the token.CONST token
	Name  *Identifier
	Value Expression
}

func (i *Identifier) expressionNode() {}
