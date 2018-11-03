package lexer

import (
	"testing"

	"github.com/Fr0stb1t3/go-vtwin/token"
)

func TestNextToken(t *testing.T) {
	input := `
				let five <- 5;
				const ten <- 10;
				const add <- func(x, y) {
					return x + y;
				};

				const result <- add(five, ten);
				!-/*2;
				5 < 6 > 4;
				if (5 < 10) {
					return true;
				} else {
					return false;
				}
				10 = 10;
				10 != 9;
				"Hello"
				"Hello world!"
				const arr <- [1, 2, 3];
				six <- 6;
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "<-"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.CONST, "const"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "<-"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.CONST, "const"},
		{token.IDENT, "add"},
		{token.ASSIGN, "<-"},
		{token.FUNCTION, "func"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.ADD, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.CONST, "const"},
		{token.IDENT, "result"},
		{token.ASSIGN, "<-"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.NOT, "!"},
		{token.SUBT, "-"},
		{token.DIV, "/"},
		{token.MULT, "*"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LSS, "<"},
		{token.INT, "6"},
		{token.GTR, ">"},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LSS, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQL, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.STRING, "Hello"},
		{token.STRING, "Hello world!"},

		{token.CONST, "const"},
		{token.IDENT, "arr"},
		{token.ASSIGN, "<-"},
		{token.LBRACK, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RBRACK, "]"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "six"},
		{token.ASSIGN, "<-"},
		{token.INT, "6"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q,  %#v", i, tt.expectedType, tok.Type, tok)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q, %#v", i, tt.expectedLiteral, tok.Literal, tok)
		}
	}
}
