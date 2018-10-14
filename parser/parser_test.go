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
	tree := program.Statements[0]
	rootToken := token.NewToken(token.SUBT, '-')
	tokLeft := token.NewToken(token.ADD, '+')
	tokRight := token.NewToken(token.INT, '5')
	tokLeftRight := token.NewToken(token.INT, '4')
	tokLeftLeft := token.NewToken(token.ADD, '+')
	tokLeftLeftLeft := token.NewToken(token.INT, '1')
	tokLeftLeftRight := token.NewToken(token.INT, '2')
	assertEqual(t, tree.Value, rootToken)
	assertEqual(t, tree.Left.Value, tokLeft)
	assertEqual(t, tree.Left.Value, tokLeft)
	assertEqual(t, tree.Right.Value, tokRight)
	assertEqual(t, tree.Left.Right.Value, tokLeftRight)
	assertEqual(t, tree.Left.Left.Value, tokLeftLeft)
	assertEqual(t, tree.Left.Left.Left.Value, tokLeftLeftLeft)
	assertEqual(t, tree.Left.Left.Right.Value, tokLeftLeftRight)
}

func TestPrecedence(t *testing.T) {
	input := "1+2+4*5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0]
	rootToken := token.NewToken(token.ADD, '+')
	tokLeft := token.NewToken(token.ADD, '+')
	tokRight := token.NewToken(token.MULT, '*')
	tokLeftRight := token.NewToken(token.INT, '2')
	tokLeftLeft := token.NewToken(token.INT, '1')
	tokRightLeft := token.NewToken(token.INT, '4')
	tokRightRight := token.NewToken(token.INT, '5')
	assertEqual(t, tree.Value, rootToken)
	assertEqual(t, tree.Left.Value, tokLeft)
	assertEqual(t, tree.Left.Value, tokLeft)
	assertEqual(t, tree.Right.Value, tokRight)
	assertEqual(t, tree.Left.Right.Value, tokLeftRight)
	assertEqual(t, tree.Left.Left.Value, tokLeftLeft)
	assertEqual(t, tree.Right.Left.Value, tokRightLeft)
	assertEqual(t, tree.Right.Right.Value, tokRightRight)
	//fmt.Printf("Tree stringified %v\n", program.Statements)
}

func TestPrecedenceTwo(t *testing.T) {
	input := "1+2*4+5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0]
	rootToken := token.NewToken(token.ADD, '+')
	tokLeft := token.NewToken(token.INT, '1')
	tokRight := token.NewToken(token.ADD, '+')

	tokRightLeft := token.NewToken(token.MULT, '*')
	tokRightRight := token.NewToken(token.INT, '5')
	tokRightLeftLeft := token.NewToken(token.INT, '2')
	tokRightLeftRight := token.NewToken(token.INT, '4')
	// fmt.Printf("Tree stringified %v\n", program.Statements)
	assertEqual(t, tree.Value, rootToken)
	assertEqual(t, tree.Left.Value, tokLeft)
	assertEqual(t, tree.Right.Value, tokRight)
	assertEqual(t, tree.Right.Left.Value, tokRightLeft)
	assertEqual(t, tree.Right.Right.Value, tokRightRight)
	assertEqual(t, tree.Right.Left.Left.Value, tokRightLeftLeft)
	assertEqual(t, tree.Right.Left.Right.Value, tokRightLeftRight)

}

func TestPrecedenceBraces(t *testing.T) {
	input := "(2+1)*(4+5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0]
	rootToken := token.NewToken(token.MULT, '*')
	tokLeft := token.NewToken(token.ADD, '+')
	tokRight := token.NewToken(token.ADD, '+')

	tokLeftRight := token.NewToken(token.INT, '1')
	tokLeftLeft := token.NewToken(token.INT, '2')
	tokRightLeft := token.NewToken(token.INT, '4')
	tokRightRight := token.NewToken(token.INT, '5')

	fmt.Printf("Tree stringified %v\n", program.Statements)
	assertEqual(t, tree.Value, rootToken)
	assertEqual(t, tree.Left.Value, tokLeft)
	assertEqual(t, tree.Right.Value, tokRight)
	assertEqual(t, tree.Left.Right.Value, tokLeftRight)
	assertEqual(t, tree.Left.Left.Value, tokLeftLeft)
	assertEqual(t, tree.Right.Left.Value, tokRightLeft)
	assertEqual(t, tree.Right.Right.Value, tokRightRight)
}

func TestPrecedenceBracesTwo(t *testing.T) {
	input := "(2+1)*(4+5)-5/3+5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0]
	/*
		rootToken := token.NewToken(token.MULT, '*')
		tokLeft := token.NewToken(token.ADD, '+')
		tokRight := token.NewToken(token.ADD, '+')

		tokLeftRight := token.NewToken(token.INT, '1')
		tokLeftLeft := token.NewToken(token.INT, '2')
		tokRightLeft := token.NewToken(token.INT, '4')
		tokRightRight := token.NewToken(token.INT, '5')
	*/

	fmt.Printf("Tree stringified %v\n", tree)
	/*
		assertEqual(t, tree.Value, rootToken)
		assertEqual(t, tree.Left.Value, tokLeft)
		assertEqual(t, tree.Right.Value, tokRight)
		assertEqual(t, tree.Left.Right.Value, tokLeftRight)
		assertEqual(t, tree.Left.Left.Value, tokLeftLeft)
		assertEqual(t, tree.Right.Left.Value, tokRightLeft)
		assertEqual(t, tree.Right.Right.Value, tokRightRight)*/
}
