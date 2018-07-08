package parser

import (
	"fmt"
	"testing"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
)

func TestNextToken(t *testing.T) {
	input := "5+4;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	// fmt.Printf("%v\n", tok)
	fmt.Printf("%v\n", program)
	/*
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		literal, ok := stmt.Expression.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
		}

		if literal.Value != 5 {
			t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
		}

		if literal.TokenLiteral() != "5" {
			t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
		}*/
}
