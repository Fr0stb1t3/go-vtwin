package ast

import "github.com/antonikliment/go-vtwin/token"

type Reference interface {
	Value() Expression
}

type Function struct {
	Name *Identifier
	Body *BlockStatement
}

func (ls Function) Value() Expression {
	return nil
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
