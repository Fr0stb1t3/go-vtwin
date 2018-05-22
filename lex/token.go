package lex

// TokenType of tokens for iota
type TokenType int

const (
	// EOF char
	EOF TokenType = iota
	itemIllegal

	// Identifiers + literals
	itemIdentifier // = "IDENTIFIER" // add, foobar, x, y, ...
	itemNumber     // = "NUMBER"     // 1343456
	itemString     // = "STRING"

	// Operators

	itemAssign   // = ":="
	itemPlus     //  = "+"
	itemMinus    // = "-"
	itemBang     // = "!"
	itemAsterisk // = "*"
	itemSlash    // = "/"

	itemLt     // = "<"
	itemLtOrEq // = "<="
	itemGt     // = ">"
	itemGtOrEq // = ">="

	itemEq    // = "="
	itemNotEq // = "!="

	// Delimiters

	COMMA     // = ","
	COLON     // = ":"
	SEMICOLON // = ";"
	LPAREN    // = "("
	RPAREN    // = ")"
	LBRACE    // = "{"
	RBRACE    // = "}"
	LBRAKET   // = "["
	RBRAKET   // = "]"

	// Keywords

	FUNCTION // = "FUNCTION"
	LET      // = "LET"
	CONST    // = "CONST"
	TRUE     // = "TRUE"
	FALSE    // = "FALSE"
	NIL      // = "NIL"

	//	OR       // = "OR"
	IF     // = "IF"
	ELSE   // = "ELSE"
	RETURN // = "RETURN"
	PRINT  // = "PRINT"
)

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
	//	"or":     OR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return itemIdentifier
}
