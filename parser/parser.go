package parser

import (
	"fmt"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/token"
)

type Parser struct {
	l         *lexer.Lexer
	errors    []string
	TopScope  Scope
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.TopScope = Scope{Objects: make(map[string]*LetStatement)}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseUnaryExpr() Expression {
	operator := token.NewToken(token.ADD, '+')
	switch p.curToken.Type {
	case token.ADD, token.SUBT, token.NOT, token.XOR, token.AND:
		operator = p.curToken
		p.nextToken()
	}
	return UnaryExpression{
		Operator: operator,
		Operand:  p.curToken,
	}
}
func (p *Parser) parseParenExpr() Expression {
	var parenCounter int
	pExpr := ParenExpression{}
	for p.tokenIs(token.LPAREN) {
		parenCounter++
		pExpr.Lparen = p.curToken
		p.nextToken()
	}
	pExpr.Expr = p.parseExpression(token.RPAREN)
	for p.peekTokenIs(token.RPAREN) {
		parenCounter--
		p.nextToken()
	}
	pExpr.Rparen = p.curToken
	parenCounter--
	if parenCounter > 0 {
		panic("Paren mismatch")
	}

	return pExpr
}

/*
 parseBinaryExpr starts constructing the AOL tree from the bottom left.
 Create Expression
	While there are tokens to be read
		Set the operator to the operator Value
		Set Identifiers or Ints to the branches of the Expression
		When the node is full:
			If there are more tokens left shift the previous expresion as a left node and continue
			else return the expression
*/
func (p *Parser) parseBinaryExpr(endToken token.Type) Expression {
	expression := BinaryExpression{}
	lowPrecedenceExpr := BinaryExpression{}
	for !p.tokenIs(endToken) {
		if expression.completeNode() {
			if p.curToken.Type.Precedence() > expression.Operator.Type.Precedence() {
				lowPrecedenceExpr.Left = expression.Left
				lowPrecedenceExpr.Operator = expression.Operator
				expression.Operator = token.Token{}
				expression.Left = expression.Right
				expression.Right = nil
			} else if lowPrecedenceExpr.Left != nil && p.curToken.Type.Precedence() <= lowPrecedenceExpr.Operator.Type.Precedence() {
				lowPrecedenceExpr.Right = expression
				expression = lowPrecedenceExpr
				lowPrecedenceExpr = BinaryExpression{}
				expression = expression.shiftNode()
			} else if !p.tokenIs(endToken) {
				expression = expression.shiftNode()
			}

		}
		switch {
		case p.tokenIs(token.LPAREN):
			expr := p.parseParenExpr()
			expression.addSubnode(expr)

		case p.curToken.Type.IsOpertor() && !expression.Operator.Type.IsOpertor():
			expression.Operator = p.curToken

		case expression.Left == nil, expression.Right == nil:
			expr := p.parseUnaryExpr()
			expression.addSubnode(expr)
		}

		if !p.peekTokenIs(token.EOF) {
			p.nextToken()
		}
	}
	// Resolve any dangling expressions
	if lowPrecedenceExpr.Left != nil {
		lowPrecedenceExpr.Right = expression
		expression = lowPrecedenceExpr
	}
	return expression
}

func (p *Parser) parseExpression(endToken token.Type) Expression {
	switch {
	case p.tokenIs(token.LPAREN), p.peekToken.Type.IsOpertor():
		return p.parseBinaryExpr(endToken)
	default:
		return p.parseUnaryExpr()
	}
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

	assignment.Expr = p.parseExpression(token.SEMICOLON)
	return assignment
}
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case token.CONST:
		fmt.Printf("parse as immutable assignment %v\n", p.curToken.Literal)
	case token.LET:
		stmt := p.parseLetStatement()
		p.TopScope.Objects[stmt.Name.Value] = &stmt
		return stmt
	case token.RETURN:
		fmt.Printf("parse as return statement %v\n", p.curToken.Literal)
	case token.LPAREN, token.INT:
		start := p.curToken
		expression := p.parseExpression(token.SEMICOLON)
		return ExpressionStatement{
			Token: start,
			Expr:  expression,
		}
	}
	return nil
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
