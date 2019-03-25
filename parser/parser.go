package parser

import (
	"fmt"

	"github.com/antonikliment/go-vtwin/ast"
	"github.com/antonikliment/go-vtwin/lexer"
	"github.com/antonikliment/go-vtwin/token"
)

type Parser struct {
	l         *lexer.Lexer
	errors    []string
	TopScope  *ast.Scope
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.TopScope = &ast.Scope{
		Objects: make(map[string]ast.Reference),
	}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}
func (p *Parser) openScope() {
	p.TopScope = ast.NewScope(p.TopScope)
}

func (p *Parser) closeScope() {
	p.TopScope = p.TopScope.Outer
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
		p.expectTokenToBe(token.LPAREN)
	}
	pExpr.Expr = p.parseExpression(token.RPAREN)
	for p.peekTokenIs(token.RPAREN) {
		parenCounter--
		p.nextToken()
	}
	pExpr.Rparen = p.curToken
	// p.expectTokenToBe(pExpr.Rparen)
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

func (p *Parser) parseFunction() ast.Function {
	p.nextToken()
	ident := p.parseIdentifier()
	// p.TopScope.Objects[ident.Value] = ident
	p.expectTokenToBe(token.LPAREN)

	p.expectTokenToBe(token.RPAREN)
	p.openScope()

	body := p.parseBlockStatement()
	p.closeScope()

	fun := ast.Function{
		Name: ident,
		Body: &body,
	}
	p.TopScope.Objects[ident.Value] = &fun
	// Delcare func in
	return fun
}
func (p *Parser) parseIdentifier() *ast.Identifier {
	//pos := p.pos
	name := "_"
	tok := p.curToken
	if p.curToken.Type == token.IDENT {
		name = p.curToken.Literal
		p.nextToken()
	} /*  else {
		p.expect(token.IDENT) // use expect() error handling
	}*/
	return &ast.Identifier{
		Token: tok,
		Value: name,
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
		ident = p.parseIdentifier()
	}
	if !p.tokenIs(token.ASSIGN) {
		panic("Invalid Let assignment" + p.peekToken.Literal)
	}
	p.nextToken() // SKIP assignment

	expr := p.parseExpression(token.SEMICOLON)

	if !immutable {
		assignment := p.parseLetAssignment(ident, expr)
		p.TopScope.Objects[ident.Value] = &assignment
		return assignment
	}

	assignment := p.parseConstAssignment(ident, expr)
	p.TopScope.Objects[ident.Value] = &assignment
	return assignment

}

/*
func declare(decl, scope *ast.Scope, declType, identifiers ...*ast.Identifier) {
	for _, ident := range identifiers {
		obj = ast.New
	}
}
*/
func (p *Parser) parseBlockStatement() ast.BlockStatement {
	block := ast.BlockStatement{
		Lbrace: p.curToken,
	}
	p.nextToken()

	statements := p.parseStatementList()
	block.Statements = statements

	block.Rbrace = p.curToken

	p.expectTokenToBe(token.RBRACE)
	return block
}

func (p *Parser) expectTokenToBe(tok token.Type) {
	// pos := p.pos
	if p.curToken.Type != tok {
		p.errorExpected("expected '" + tok.String() + "' got " + p.curToken.Type.String())
		panic("Fatal")
	}
	p.nextToken() // make progress
	return        //pos
}
func (p *Parser) errorExpected(str string) {
	fmt.Printf("%v\n", str)
}
func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.ReturnVal = p.parseExpression(token.SEMICOLON)
	return stmt
}
func (p *Parser) parseStatementList() (statements []ast.Statement) {
	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}
	return
}
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET, token.CONST: //, token.IDENT:
		return p.parseAssignment()
	case token.FUNCTION:
		return p.parseFunction()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	case token.IDENT:
		return p.parseAssignment()
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
