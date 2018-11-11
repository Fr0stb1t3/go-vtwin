package ast

import (
	"fmt"

	"github.com/antonikliment/go-vtwin/token"
)

type Expression interface {
	exprNode()
}

type ParenExpression struct {
	Lparen token.Token
	Expr   Expression
	Rparen token.Token
}

func (e ParenExpression) exprNode() {}

type UnaryExpression struct {
	Operator token.Token
	Operand  token.Token
}

func (e UnaryExpression) exprNode() {}

func (nd UnaryExpression) String() string {
	return nd.Operand.Literal
}

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (e BinaryExpression) exprNode() {}

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
