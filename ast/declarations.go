package ast

import "github.com/antonikliment/go-vtwin/token"

type Reference interface {
	Value() Expression
}

type Function struct {
	Name *Identifier
	Body *BlockStatement
}

func (fs Function) stmtNode() {}

func (fs Function) Value() Expression {
	return nil
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Expr  Expression
}
type AssignmentStatement struct {
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

func (e AssignmentStatement) Value() Expression { return e.Expr }
func (e AssignmentStatement) stmtNode()         {}

func (cs ConstStatement) stmtNode() {}
func (ls LetStatement) stmtNode()   {}
