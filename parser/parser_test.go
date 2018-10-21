package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/token"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
func add(A int, B int) int {
	return A + B
}
func subtract(A int, B int) int {
	return A - B
}
func multiply(A int, B int) int {
	return A * B
}
func divide(A int, B int) int {
	return A / B
}

func evaluateBranch(ex BinaryExpression) int {
	if ex.Operator.Type.IsOpertor() {
		return evaluateExpression(ex)
	}
	val, _ := strconv.Atoi(ex.Operator.Literal)
	return val

}

func evaluate(op token.Token, left BinaryExpression, right BinaryExpression) int {

	a := evaluateBranch(left)
	b := evaluateBranch(right)

	switch op.Type {
	case token.ADD:
		return add(a, b)
	case token.SUBT:
		return subtract(a, b)
	case token.MULT:
		return multiply(a, b)
	case token.DIV:
		return divide(a, b)
	}
	return 0
}

func evaluateExpression(ex Expression) int {
	switch ex.(type) {
	case BinaryExpression:
		be := ex.(BinaryExpression)
		op := be.Operator
		left := *be.Left
		right := *be.Right
		res := evaluate(op, left, right)
		return res
	case UnaryExpression:
		fmt.Printf("Tree stringified %v\n", ex)
	}
	return 0
}

func TestEvaluate(t *testing.T) {
	input := "1+1+2;" //2*7+3;3+2*7;

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0].getTree()
	res := evaluateExpression(tree)
	assertEqual(t, res, 4)
}
func TestSimpleTree(t *testing.T) {
	input := "1+2+4-5;" //2*7+3;3+2*7;

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0].getTree()
	res := evaluateExpression(tree)
	assertEqual(t, res, 2)
}

func TestPrecedence(t *testing.T) {
	input := "1+2+4*5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0].getTree()
	res := evaluateExpression(tree)
	assertEqual(t, res, 23)
}

func TestPrecedenceTwo(t *testing.T) {
	input := "1+2*4+5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0].getTree()

	res := evaluateExpression(tree)
	assertEqual(t, res, 14)
}

func TestPrecedenceBraces(t *testing.T) {
	input := "(2+1)*(4+5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0].getTree()

	res := evaluateExpression(tree)
	assertEqual(t, res, 27)
}

func TestPrecedenceBracesTwo(t *testing.T) {
	input := "(2+1)*(4+5)-6/3+5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0].getTree()

	res := evaluateExpression(tree)

	assertEqual(t, res, 30)

}

func TestLetAssignment(t *testing.T) {
	input := "let test := 1;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0]
	fmt.Printf("Tree stringified %v\n", tree)
	//
	// rootToken := token.NewToken(token.SUBT, '-')
	// tokLeft := token.NewToken(token.MULT, '*')
	// tokRight := token.NewToken(token.ADD, '+')
	/**
	tokLeft := token.NewToken(token.ADD, '+')

	tokLeftRight := token.NewToken(token.INT, '1')
	tokLeftLeft := token.NewToken(token.INT, '2')
	tokRightLeft := token.NewToken(token.INT, '4')
	tokRightRight := token.NewToken(token.INT, '5')
	/**/
	// assertEqual(t, tree.Operator, rootToken)
	// assertEqual(t, tree.Left.Operator, tokLeft)
	// assertEqual(t, tree.Right.Operator, tokRight)
	/*
		assertEqual(t, tree.Left.Right.Operator, tokLeftRight)
		assertEqual(t, tree.Left.Left.Operator, tokLeftLeft)
		assertEqual(t, tree.Right.Left.Operator, tokRightLeft)
		assertEqual(t, tree.Right.Right.Operator, tokRightRight)
	*/
}
