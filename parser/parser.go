package parser

import (
	"fmt"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/token"
)

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

type UnaryExpression struct {
	Operator token.Token
	Operand  token.Token
}

// switch l := e.X.(type) {
func (e UnaryExpression) exprNode() {

}

//
// func (nd *UnaryExpression) String() string {
// 	return nd.Operand.Literal
// }

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

// switch l := e.X.(type) {
func (e BinaryExpression) exprNode() {

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

func (nd *BinaryExpression) String() string {
	return nd.Operator.Literal
}

type Statement interface {
}
type ExpressionStatement struct {
	Token token.Token // the first token of the expression
	Expr  Expression
}

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseExpression(endToken token.Type) Expression {
	if endToken == token.RPAREN {
		if p.tokenIs(token.LPAREN) {
			p.nextToken()
		}
	}
	expression := BinaryExpression{}

	for !p.tokenIs(endToken) {

		/*
			if there is a open brace call parse expression (recursion)
		*/
		if p.tokenIs(token.LPAREN) {
			subExpression := p.parseExpression(token.RPAREN)
			expression.addSubnode(subExpression)
			p.nextToken()
		}

		/*
			If there are more tokens
			Moves the old expression to the left BinaryExpression
		*/
		if !p.tokenIs(endToken) &&
			expression.completeNode() {
			oldExpression := expression
			expression = BinaryExpression{Left: oldExpression}
		}
		/*
			If the BinaryExpression has an operator next operator precedence
		*/
		if expression.Operator.Type.IsOpertor() &&
			p.peekPrecedence() > expression.Operator.Type.Precedence() {
			subExpression := p.parseExpression(endToken)
			expression.Right = subExpression
		}
		if p.curToken.Type.IsOpertor() {
			expression.Operator = p.curToken
		} else if !p.tokenIs(endToken) {
			operand := token.NewToken(token.ADD, '+')
			leaf := UnaryExpression{
				Operator: operand,
				Operand:  p.curToken,
			}
			expression.addSubnode(leaf)
		}
		if !p.peekTokenIs(token.EOF) {
			p.nextToken()
		}
	}
	return expression
}

func (p *Parser) parseLetStatement() LetStatement {
	assignment := LetStatement{
		Token: p.curToken,
	}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		assignment.Name = &Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
	}
	if !p.peekTokenIs(token.ASSIGN) {
		panic("Invalid Let assignment")
	}
	p.nextToken() // TODO
	p.nextToken()
	operand := token.NewToken(token.ADD, '+')
	assignment.Expr = UnaryExpression{
		Operator: operand,
		Operand:  p.curToken,
	}
	return assignment
}
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case token.CONST:
		fmt.Printf("parse as immutable assignment %v\n", p.curToken.Literal)
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		fmt.Printf("parse as return statement %v\n", p.curToken.Literal)
	case token.LPAREN:
		expression := p.parseExpression(token.SEMICOLON)
		return ExpressionStatement{
			Expr: expression,
		}
	case token.INT:
		if p.peekToken.Type.IsOpertor() {
			expression := p.parseExpression(token.SEMICOLON)
			return ExpressionStatement{
				Expr: expression,
			}
			// return p.parseExpression(token.SEMICOLON)
		}
	default:
		// log.Info("Default", p.curToken)
		// fmt.Printf("parse as expression statement %v \n", p.curToken.Literal)
	}
	return nil
}

type Program struct {
	Statements []Statement
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	// fmt.Printf("parsing %v\n", program)
	return program
}
func (p *Parser) peekPrecedence() int {
	return p.peekToken.Type.Precedence()
}
func (p *Parser) tokenIs(t token.Type) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}
