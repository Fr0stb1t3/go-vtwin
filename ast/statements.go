package ast

import "github.com/antonikliment/go-vtwin/token"

type Statement interface {
	// stmtNode()
}

type AssignmentStatement struct {
	LeftSide  []Expression
	Token     token.Token
	RightSide []Expression
}
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

func (e AssignmentStatement) stmtNode() {}
func (e ExpressionStatement) stmtNode() {}
func (e BlockStatement) stmtNode()      {}
func (e ReturnStatement) stmtNode()     {}
