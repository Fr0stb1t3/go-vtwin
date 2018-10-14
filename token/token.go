package token

import (
	"strconv"
)

// Type of tokens for iota
type Type int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Type = iota
	EOF
	COMMENT

	literal_beg
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	CHAR   // 'a'
	STRING // "abc"
	literal_end

	operator_beg
	// Operators
	AND  // &
	OR   // |
	XOR  // ^
	REM  // %
	ADD  // +
	SUBT // -
	NOT  // !
	MULT // *
	DIV  // /

	LAND    // &&
	LOR     // ||
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ASSIGN     // :=
	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	DIV_ASSIGN // /=
	LSS        // "<"
	GTR        // ">"

	EQL // "="
	LEQ // "<="
	GEQ // ">="
	NEQ // "!="

	operator_end
	// Delimiters

	COMMA     // ","
	COLON     // ":"
	SEMICOLON // ";"
	LPAREN    // "("
	RPAREN    // ")"
	LBRACE    // "{"
	RBRACE    // "}"
	LBRACK    // "["
	RBRACK    // "]"

	keyword_beg
	FUNCTION // "FUNCTION"
	LET      // "LET"
	CONST    // "CONST"
	TRUE     // "TRUE"
	FALSE    // "FALSE"
	NIL      // "NIL"
	RETURN   // "RETURN"
	IMPORT   // "IMPORT"
	FROM     // "FROM"
	IF       // "IF"
	ELSE     // "ELSE"
	keyword_end
)

var tokens = [...]string{
	// Special tokens
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	// Operators
	AND:  "&",
	OR:   "|",
	XOR:  "^",
	REM:  "%",
	ADD:  "+",
	SUBT: "-",
	NOT:  "!",
	MULT: "*",
	DIV:  "/",

	LAND:    "&&",
	LOR:     "||",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ASSIGN:     ":=",
	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	DIV_ASSIGN: "/=",
	LSS:        "<",
	GTR:        ">",

	EQL: "=",
	LEQ: "<=",
	GEQ: ">=",
	NEQ: "!=",

	COMMA:     ",",
	COLON:     ":",
	SEMICOLON: ";",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACK:    "[",
	RBRACK:    "]",

	FUNCTION: "func",
	LET:      "let",
	CONST:    "const",
	TRUE:     "true",
	FALSE:    "false",
	NIL:      "nil",
	RETURN:   "return",
	IMPORT:   "import",
	FROM:     "from",
	IF:       "if",
	ELSE:     "else",
}

func (tok Type) String() string {
	s := ""
	// fmt.Printf("tok %#v \n", tok)
	if 0 <= tok && tok < Type(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

// Precedence returns the precedence level for all operations
func (op Type) Precedence() int {
	switch op {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUBT, OR, XOR:
		return 4
	case MULT, DIV, SHL, SHR, AND, AND_NOT, REM:
		return 5
	}
	return LowestPrec
}

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Type {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

func (tok Type) IsLiteral() bool {
	return literal_beg < tok && tok < literal_end
}

func (tok Type) IsOpertor() bool {
	return operator_beg < tok && tok < operator_end
}

func (tok Type) isKeyword() bool {
	return keyword_beg < tok && tok < keyword_end
}

type Token struct {
	Type    Type
	Literal string
}

func NewToken(t Type, chars ...byte) Token {
	literal := ""
	for _, ch := range chars {
		literal += string(ch)
	}
	return Token{
		Type:    t,
		Literal: literal,
	}
}
