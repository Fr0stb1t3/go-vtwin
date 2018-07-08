package parser

import (
	"fmt"

	"github.com/Fr0stb1t3/go-vtwin/ast"
	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Token]prefixParseFn
	infixParseFns  map[token.Token]infixParseFn
}
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

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
	p.peekToken = p.l.NextToken().Type
}

func (p *Parser) parseExpression(precedence int) {

}

func (p *Parser) processToken(t token.Token) {
	fmt.Printf("parsing %v\n", t)
	if t == token.INT {
		fmt.Printf("push it to the output queue %v\n", t)

	}
	// while there are tokens to be read:
	//     read a token.
	//     if the token is a number, then:
	//         push it to the output queue.
	//     if the token is a function then:
	//         push it onto the operator stack
	//     if the token is an operator, then:
	//         while ((there is a function at the top of the operator stack)
	//                or (there is an operator at the top of the operator stack with greater precedence)
	//                or (the operator at the top of the operator stack has equal precedence and is left associative))
	//               and (the operator at the top of the operator stack is not a left bracket):
	//             pop operators from the operator stack onto the output queue.
	//         push it onto the operator stack.
	//     if the token is a left bracket (i.e. "("), then:
	//         push it onto the operator stack.
	//     if the token is a right bracket (i.e. ")"), then:
	//         while the operator at the top of the operator stack is not a left bracket:
	//             pop the operator from the operator stack onto the output queue.
	//         pop the left bracket from the stack.
	//         /* if the stack runs out without finding a left bracket, then there are mismatched parentheses. */
	// if there are no more tokens to read:
	//     while there are still operator tokens on the stack:
	//         /* if the operator token on the top of the stack is a bracket, then there are mismatched parentheses. */
	//         pop the operator from the operator stack onto the output queue.
	// exit.
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken != token.EOF {
		p.processToken(p.curToken)
		p.nextToken()
	}
	fmt.Printf("parsing %v\n", program)
	return program
}
func (p *Parser) peekTokenIs(t token.Token) bool {
	return p.peekToken == t
}
