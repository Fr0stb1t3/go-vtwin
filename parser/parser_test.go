package parser

import (
	"fmt"
	"testing"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/token"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
func TestSimpleTree(t *testing.T) {
	input := "1+2+4-5;" //2*7+3;3+2*7;

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	rootToken := token.NewToken(token.SUBT, '-')
	tokLeft := token.NewToken(token.ADD, '+')
	tokRight := token.NewToken(token.INT, '5')
	tokLeftRight := token.NewToken(token.INT, '4')
	tokLeftLeft := token.NewToken(token.ADD, '+')
	tokLeftLeftLeft := token.NewToken(token.INT, '1')
	tokLeftLeftRight := token.NewToken(token.INT, '2')
	assertEqual(t, program.Statements[0].Value, rootToken)
	assertEqual(t, program.Statements[0].Left.Value, tokLeft)
	assertEqual(t, program.Statements[0].Left.Value, tokLeft)
	assertEqual(t, program.Statements[0].Right.Value, tokRight)
	assertEqual(t, program.Statements[0].Left.Right.Value, tokLeftRight)
	assertEqual(t, program.Statements[0].Left.Left.Value, tokLeftLeft)
	assertEqual(t, program.Statements[0].Left.Left.Left.Value, tokLeftLeftLeft)
	assertEqual(t, program.Statements[0].Left.Left.Right.Value, tokLeftLeftRight)
}

func TestPrecedence(t *testing.T) {
	input := "1+2+4*5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	fmt.Printf("Tree stringified %v\n", program.Statements)
}
func TestPrecedenceTwo(t *testing.T) {
	input := "1+2*4+5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	// fmt.Printf("%v\n", tok)
	fmt.Printf("Tree stringified %v\n", program.Statements)
}
