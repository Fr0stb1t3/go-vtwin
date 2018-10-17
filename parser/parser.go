package parser

import (
	"fmt"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/token"
)

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
	Value Expression
}

func (eS LetStatement) getTree() Expression {
	return eS.Value
}

type Expression struct {
	Left  *Expression
	Value token.Token
	Right *Expression
}

type Statement interface {
	getTree() Expression
}
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (eS ExpressionStatement) getTree() Expression {
	return eS.Expression
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

func (nd *Expression) String() string {
	var out string

	if nd.Left != nil {
		out = " <" + out + nd.Left.String()
	}
	out = out + nd.Value.Literal
	if nd.Right != nil {
		out = out + nd.Right.String() + "> "
	}
	return out
}

func (p *Parser) parseExpression(endToken token.Type) *Expression {
	if endToken == token.RPAREN {
		if p.tokenIs(token.LPAREN) {
			p.nextToken()
		}
	}
	empty := token.Token{}
	expression := &Expression{}

	for !p.tokenIs(endToken) {

		/*
			if there is a open brace run call parse expression (recursion)
		*/
		if p.tokenIs(token.LPAREN) {
			subExpression := p.parseExpression(token.RPAREN)
			if expression.Left == nil {
				expression.Left = subExpression
			} else if expression.Right == nil {
				expression.Right = subExpression
			}
			p.nextToken()
		}

		/*
			If there are more tokens
			Moves the old expression to the left Expression
		*/
		if !p.tokenIs(endToken) &&
			expression.Left != nil &&
			expression.Value != empty &&
			expression.Right != nil {
			oldExpression := *(&expression)
			expression = &Expression{Left: oldExpression}
		}
		/*
			If the left Expression has an operator next operator precedence
		*/
		if expression.Left != nil &&
			expression.Value.Type.IsOpertor() &&
			p.peekPrecedence() > expression.Value.Type.Precedence() {
			subExpression := p.parseExpression(endToken)
			expression.Right = subExpression
		}
		if p.curToken.Type.IsOpertor() {
			expression.Value = p.curToken
		} else {
			if expression.Left == nil {
				expression.Left = &Expression{Value: p.curToken}
			} else if expression.Right == nil {
				expression.Right = &Expression{Value: p.curToken}
			}
		}
		if p.tokenIs(endToken) {
			return expression
		}
		if !p.peekTokenIs(token.EOF) {
			p.nextToken()
		}
	}
	return expression
}

func (p *Parser) parseLetStatement() *LetStatement {
	assignment := &LetStatement{
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
		return nil
	}
	p.nextToken() // TODO
	p.nextToken()
	assignment.Value = Expression{
		Value: p.curToken,
	} // *p.parseExpression(token.SEMICOLON)
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
			Expression: *expression,
		}
	case token.INT:
		if p.peekToken.Type.IsOpertor() {
			expression := p.parseExpression(token.SEMICOLON)
			return ExpressionStatement{
				Expression: *expression,
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
