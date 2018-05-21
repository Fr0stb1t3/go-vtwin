package lex

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENTIFIER = "IDENTIFIER" // add, foobar, x, y, ...
	NUMBER     = "NUMBER"     // 1343456
	STRING     = "STRING"

	// Operators
	ASSIGN   = ":="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT    = "<"
	LT_EQ = "<="
	GT    = ">"
	GT_EQ = ">="

	EQ     = "="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	COLON     = ":"
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRAKET   = "["
	RBRAKET   = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONST    = "CONST"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NIL      = "NIL"
	OR       = "OR"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	PRINT    = "PRINT"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"const":  CONST,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NIL,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"print":  PRINT,
	"or":     OR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}
