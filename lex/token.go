package lex

// Token of tokens for iota
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	// Identifiers and basic type literals (these tokens stand for classes of literals)
	IDENT // main
	INT   // 12345
	FLOAT // 123.45
	IMAG  // 123.45i
	// CHAR   // 'a'
	STRING // "abc"
	literal_end

	operator_beg
	// Operators
	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD        // +
	SUBT       // -
	NOT        // !
	MULT       // *
	DIV        // /
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

	operator_end

	keyword_beg
	// Keywords
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
	RETURN   // "RETURN"
	keyword_end
)

type TokenStruct struct {
	Type    Token
	Literal string
}

var keywords = map[string]Token{
	"func":   FUNCTION,
	"const":  CONST,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NIL,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	//	"or":     OR,
}

const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

func (op Token) Precedence() int {
	switch op {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB, OR, XOR:
		return 4
	case MUL, QUO, REM, SHL, SHR, AND, AND_NOT:
		return 5
	}
	return LowestPrec
}

func LookupIdent(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return Identifier
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}
