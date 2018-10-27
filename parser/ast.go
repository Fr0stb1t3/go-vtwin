package parser

import "github.com/Fr0stb1t3/go-vtwin/token"

type Program struct {
	Statements []Statement
}
type Expression interface {
	exprNode()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i Identifier) String() string {
	return i.Value
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Expr  Expression
}

func (eS LetStatement) getTree() Expression {
	return eS.Expr
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
func (e *BinaryExpression) shiftNode() BinaryExpression {
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
func (e *BinaryExpression) completeNode() bool {
	return e.Left != nil &&
		e.Operator != token.Token{} &&
		e.Right != nil
}
func (e *BinaryExpression) addSubnode(subEx Expression) {
	if e.Left == nil {
		e.Left = subEx
		return
	} else if e.Right == nil {
		e.Right = subEx
		return
	}

	panic("BinaryExpression node is full")
}

type Statement interface{}

type ExpressionStatement struct {
	Token token.Token // the first token of the expression
	Expr  Expression
}

type Scope struct {
	Outer   *Scope
	Objects map[string]*LetStatement
}

func (s Scope) Lookup(ident string) *LetStatement {
	return s.Objects[ident]
}
