package lex

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
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
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(itemAssign, ch, l.ch)
		} else {
			tok = newToken(COLON, l.ch)
		}
	case '=':
		tok = newToken(itemEq, l.ch)
	case '+':
		tok = newToken(itemPlus, l.ch)
	case '-':
		tok = newToken(itemMinus, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(itemNotEq, ch, l.ch)
		} else {
			tok = newToken(itemBang, l.ch)
		}
	case '/':
		tok = newToken(itemSlash, l.ch)
	case '*':
		tok = newToken(itemAsterisk, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(itemLtOrEq, ch, l.ch)
		} else {
			tok = newToken(itemLt, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(itemGtOrEq, ch, l.ch)
		} else {
			tok = newToken(itemGt, l.ch)
		}
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case '[':
		tok = newToken(LBRAKET, l.ch)
	case ']':
		tok = newToken(RBRAKET, l.ch)
	case '"':
		tok.Type = itemString
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = itemNumber
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(itemIllegal, l.ch)
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

func newToken(tokenType TokenType, chars ...byte) Token {
	literal := ""
	for _, ch := range chars {
		literal += string(ch)
	}
	return Token{
		Type:    tokenType,
		Literal: literal,
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
