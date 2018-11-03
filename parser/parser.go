package parser

import (
	"fmt"

	"github.com/Fr0stb1t3/go-vtwin/ast"
	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/token"
)

type Parser struct {
	l         *lexer.Lexer
	errors    []string
	TopScope  ast.Scope
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.TopScope = ast.Scope{Objects: make(map[string]ast.Reference)}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseUnaryExpr() ast.Expression {
	operator := token.NewToken(token.ADD, '+')
	switch p.curToken.Type {
	case token.ADD, token.SUBT, token.NOT, token.XOR, token.AND:
		operator = p.curToken
		p.nextToken()
	}
	return ast.UnaryExpression{
		Operator: operator,
		Operand:  p.curToken,
	}
}
func (p *Parser) parseParenExpr() ast.Expression {
	var parenCounter int
	pExpr := ast.ParenExpression{}
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
 Create ast.Expression
	While there are tokens to be read
		Set the operator to the operator Value
		Set Identifiers or Ints to the branches of the ast.Expression
		When the node is full:
			If there are more tokens left shift the previous expresion as a left node and continue
			else return the expression
*/
func (p *Parser) parseBinaryExpr(endToken token.Type) ast.Expression {
	expression := ast.BinaryExpression{}
	lowPrecedenceExpr := ast.BinaryExpression{}
	for !p.tokenIs(endToken) {
		if expression.CompleteNode() {
			if p.curToken.Type.Precedence() > expression.Operator.Type.Precedence() {
				lowPrecedenceExpr.Left = expression.Left
				lowPrecedenceExpr.Operator = expression.Operator
				expression.Operator = token.Token{}
				expression.Left = expression.Right
				expression.Right = nil
			} else if lowPrecedenceExpr.Left != nil && p.curToken.Type.Precedence() <= lowPrecedenceExpr.Operator.Type.Precedence() {
				lowPrecedenceExpr.Right = expression
				expression = lowPrecedenceExpr
				lowPrecedenceExpr = ast.BinaryExpression{}
				expression = expression.ShiftNode()
			} else if !p.tokenIs(endToken) {
				expression = expression.ShiftNode()
			}
		}
		switch {
		case p.tokenIs(token.LPAREN):
			expr := p.parseParenExpr()
			expression.AddSubnode(expr)

		case p.curToken.Type.IsOpertor() && !expression.Operator.Type.IsOpertor():
			expression.Operator = p.curToken

		case expression.Left == nil, expression.Right == nil:
			expr := p.parseUnaryExpr()
			expression.AddSubnode(expr)
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

func (p *Parser) parseExpression(endToken token.Type) ast.Expression {
	switch {
	case p.tokenIs(token.LPAREN), p.peekToken.Type.IsOpertor():
		return p.parseBinaryExpr(endToken)
	default:
		return p.parseUnaryExpr()
	}
}

/*
	Can be extended to type checking
*/
func (p *Parser) checkScope(ident string) {
	if p.TopScope.Objects[ident] != nil {
		reference := p.TopScope.Objects[ident]
		switch reference.(type) {
		case *ast.ConstStatement:
			panic("Const cannot be reassigned")
		case *ast.LetStatement:
			return
		}
	}
}

func (p *Parser) parseLetAssignment(ident *ast.Identifier, expr ast.Expression) ast.LetStatement {
	return ast.LetStatement{
		Token: ident.Token,
		Name:  ident,
		Expr:  expr,
	}
}
func (p *Parser) parseConstAssignment(ident *ast.Identifier, expr ast.Expression) ast.ConstStatement {
	return ast.ConstStatement{
		Token: ident.Token,
		Name:  ident,
		Expr:  expr,
	}
}
func (p *Parser) parseAssignment() ast.Reference {
	var immutable bool
	var ident *ast.Identifier
	if p.tokenIs(token.LET) {
		p.nextToken()
		immutable = false
	}
	if p.tokenIs(token.CONST) {
		p.nextToken()
		immutable = true
	}
	p.checkScope(p.curToken.Literal)

	if p.tokenIs(token.IDENT) {
		ident = &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
	}

	if !p.peekTokenIs(token.ASSIGN) {
		panic("Invalid Let assignment")
	}
	p.nextToken() // SKIP assignment
	p.nextToken()

	expr := p.parseExpression(token.SEMICOLON)

	if !immutable {
		assignment := p.parseLetAssignment(ident, expr)
		p.TopScope.Objects[ident.Value] = &assignment
		return assignment
	} else {
		assignment := p.parseConstAssignment(ident, expr)
		p.TopScope.Objects[ident.Value] = &assignment
		return assignment
	}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET, token.CONST, token.IDENT:
		return p.parseAssignment()
	case token.RETURN:
		fmt.Printf("parse as return statement %v\n", p.curToken.Literal)
	case token.LPAREN, token.INT:
		start := p.curToken
		expression := p.parseExpression(token.SEMICOLON)
		return ast.ExpressionStatement{
			Token: start,
			Expr:  expression,
		}
	}
	return nil
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
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
