package parser

import (
	"fmt"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/token"
)

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

type Node struct {
	Left  *Node
	Value token.Token
	Right *Node
}

func (nd *Node) String() string {
	var out string

	if nd.Left != nil {
		out = out + nd.Left.String()
	}
	out = out + nd.Value.Literal
	if nd.Right != nil {
		out = out + nd.Right.String()
	}
	return out
}

func (p *Parser) parseExpression() Node {
	empty := token.Token{}
	expression := Node{}

	// 1+2*4+5;
	// 		 +
	// 	+		5
	// 1	*
	//  2  4
	for !p.tokenIs(token.SEMICOLON) {
		if expression.Left != nil && expression.Value != empty && expression.Right != nil {
			oldExpression := *(&expression)
			expression = Node{Left: &oldExpression}
		}
		if p.curToken.Type.IsOpertor() {
			if expression.Left.Value.Type.IsOpertor() &&
				expression.Left.Value.Type.Precedence() < p.curToken.Type.Precedence() {
				oldValue := *(&expression.Left.Right)
				newRight := Node{
					Left:  oldValue,
					Value: p.curToken,
				}
				p.nextToken()
				newRight.Right = &Node{Value: p.curToken}
				expression.Left.Right = &newRight
				expression = *expression.Left
			} else if expression.Value == empty {
				expression.Value = p.curToken
			}
		} else {
			if expression.Left == nil {
				expression.Left = &Node{Value: p.curToken}
			} else if expression.Right == nil {
				expression.Right = &Node{Value: p.curToken}
			}
		}
		// log.Info("CUR EXPRESSION", expression)
		// log.Info("oldExpression", oldExpression)
		// log.Info("-----")
		p.nextToken()
	}
	// log.Info("expression", expression)

	return expression
}

func (p *Parser) parseStatement() *Node {
	switch p.curToken.Type {
	case token.CONST:
		fmt.Printf("parse as immutable assignment %v\n", p.curToken.Literal)
	case token.LET:
		fmt.Printf("parse as mutable assignment %v\n", p.curToken.Literal)
	case token.RETURN:
		fmt.Printf("parse as return statement %v\n", p.curToken.Literal)
	case token.INT:
		if p.peekToken.Type.IsOpertor() {
			parent := p.parseExpression()
			return &parent
		}
	default:
		fmt.Printf("parse as expression statement %v \n", p.curToken.Literal)
	}
	return nil
}

func (p *Parser) parseToken(t token.Type) {
	fmt.Printf("parsing %v\n", t)
	if t == token.INT {
		fmt.Printf("push it to the output queue %v\n", t)
		return
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

// type Statement struct{}
type Program struct {
	Statements []Node
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Node{}

	for p.curToken.Type != token.EOF {
		// stmt := p.parseToken(t)(p.curToken.Type)
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, *stmt)
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
