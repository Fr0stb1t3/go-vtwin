package ast

import (
	"github.com/antonikliment/go-vtwin/token"
)

type Program struct {
	Statements []Statement
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i Identifier) String() string {
	return i.Value
}

type Reference interface {
	Value() Expression
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Expr  Expression
}

func (ls LetStatement) Value() Expression {
	return ls.Expr
}

type ConstStatement struct {
	Token token.Token
	Name  *Identifier
	Expr  Expression
}

func (ls ConstStatement) Value() Expression {
	return ls.Expr
}

type Statement interface{}

type ExpressionStatement struct {
	Token token.Token // the first token of the expression
	Expr  Expression
}
type BlockStatement struct {
	Statements []Statement
	Lbrace     token.Token
	Rbrace     token.Token
}
type ReturnStatement struct {
	Token     token.Token // RETURN token
	ReturnVal Expression
}
type Function struct {
	Name *Identifier
	Body *BlockStatement
}
type Scope struct {
	Outer   *Scope
	Objects map[string]Reference
}

func NewScope(outer *Scope) *Scope {
	return &Scope{
		Outer:   outer,
		Objects: make(map[string]Reference),
	}
}

func (s Scope) Lookup(ident string) Reference {
	return s.Objects[ident]
}
