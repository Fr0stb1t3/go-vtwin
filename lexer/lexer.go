package lexer

import (
	"github.com/Fr0stb1t3/go-vtwin/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

type Lexeme struct {
	Type    token.Token
	Literal string
	Name    string
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Lexeme {
	var tok Lexeme

	l.skipWhitespace()

	switch l.ch {
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.ASSIGN, ch, l.ch)
		} else {
			tok = newToken(token.COLON, l.ch)
		}
	case '=':
		tok = newToken(token.EQL, l.ch)
	case '+':
		tok = newToken(token.ADD, l.ch)
	case '-':
		tok = newToken(token.SUBT, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.NEQ, ch, l.ch)
		} else {
			tok = newToken(token.NOT, l.ch)
		}
	case '/':
		tok = newToken(token.DIV, l.ch)
	case '*':
		tok = newToken(token.MULT, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.LEQ, ch, l.ch)
		} else {
			tok = newToken(token.LSS, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.GEQ, ch, l.ch)
		} else {
			tok = newToken(token.GTR, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACK, l.ch)
	case ']':
		tok = newToken(token.RBRACK, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.Lookup(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
	}

	return l.input[position:l.position]
}

func newToken(token token.Token, chars ...byte) Lexeme {
	literal := ""
	for _, ch := range chars {
		literal += string(ch)
	}
	return Lexeme{
		Type:    token,
		Literal: literal,
		Name:    token.String(),
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
