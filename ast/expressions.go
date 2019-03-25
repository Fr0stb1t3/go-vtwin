package ast

import (
	"fmt"

	"github.com/antonikliment/go-vtwin/token"
)

type Expression interface {
	exprNode()
	String() string
}

type SimpleLiteral struct {
	Type  token.Type
	Value string
}
type ParenExpression struct {
	Lparen token.Token
	Expr   Expression
	Rparen token.Token
}

func (e SimpleLiteral) String() string {
	return e.Value
}
func (e SimpleLiteral) exprNode()   {}
func (e ParenExpression) exprNode() {}
func (e ParenExpression) String() string {
	return ""
}

type UnaryExpression struct {
	Operator token.Token
	Operand  Expression //token.Token
	// Expr     Expression
}

func (e UnaryExpression) exprNode() {}

func (nd UnaryExpression) String() string {
	return nd.Operand.String()
}

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (e BinaryExpression) exprNode() {}
func (e BinaryExpression) String() string {
	return e.Left.String() + e.Operator.Literal + e.Right.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Expr  Expression
	Value string
}

func (e Identifier) exprNode() {}

/*
	Moves the old expression to the left BinaryExpression
*/
func (e *BinaryExpression) ShiftNode() BinaryExpression {
	expr := *e
	return BinaryExpression{
		Left: expr,
	}
}
func (e *BinaryExpression) emptyNode() bool {
	return e.Left == nil &&
		e.Operator == token.Token{} &&
		e.Right == nil
}
func (e *BinaryExpression) CompleteNode() bool {
	return e.Left != nil &&
		e.Operator != token.Token{} &&
		e.Right != nil
}
func (e *BinaryExpression) AddSubnode(subEx Expression) {
	if e.Left == nil {
		e.Left = subEx
		return
	} else if e.Right == nil {
		e.Right = subEx
		return
	}
	fmt.Printf("\n %v", e)
	panic("BinaryExpression node is full")
}
