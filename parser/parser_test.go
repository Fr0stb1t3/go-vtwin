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

func evaluateUnaryExpr(ex Expression) int {
	uax := ex.(UnaryExpression)
	switch uax.Operand.Type {
	case token.INT:
		val, _ := strconv.Atoi(uax.Operator.Literal + uax.Operand.Literal)
		return val
	case token.IDENT:
		fmt.Printf("\n %v:", uax.Operand)
		// val, _ := strconv.Atoi(uax.Operator.Literal + uax.Operand.Literal)
		return 0
	default:
		panic("evaluateUnaryExpr: Unknown type")
	}
}

func evaluateBinaryExpr(ex Expression) int {
	be := ex.(BinaryExpression)
	a := evaluateExpression(be.Left)
	b := evaluateExpression(be.Right)

	switch be.Operator.Type {
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
	case ParenExpression:
		pex := ex.(ParenExpression)
		return evaluateExpression(pex.Expr)
	case BinaryExpression:
		return evaluateBinaryExpr(ex)
	case UnaryExpression:
		return evaluateUnaryExpr(ex)
	}
	return 0
}

type result struct {
	number int
	ident  string
}

func runStatement(stmt Statement) result {
	switch stmt.(type) {
	case ExpressionStatement:
		es := stmt.(ExpressionStatement)
		val := evaluateExpression(es.Expr)
		return result{
			number: val,
		}
	case LetStatement:
		ls := stmt.(LetStatement)
		val := evaluateExpression(ls.Expr)
		return result{
			ident:  ls.Name.Value,
			number: val,
		}
	}
	return result{}
}

/**/

func TestBasicMathTrees(t *testing.T) {
	input := `1+2+4-5;
						1+2+4*5;
						1+2*4+5;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0])
	resPrecedence := runStatement(program.Statements[1])
	resPrecedenceTwo := runStatement(program.Statements[2])
	assertEqual(t, res.number, 2)
	assertEqual(t, resPrecedence.number, 23)
	assertEqual(t, resPrecedenceTwo.number, 14)
}

func TestBracesTwo(t *testing.T) {
	input := `
		(2+1)*(4+5);
		((2+1))*(4+5);
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0])
	assertEqual(t, res.number, 27)
	restwo := runStatement(program.Statements[1])
	assertEqual(t, restwo.number, 27)
}

func TestPrecedenceTwo(t *testing.T) {
	input := `
		27-6/3+5;
		27-6/3*5;
		27-6/3*5/2;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0])
	resTwo := runStatement(program.Statements[1])
	resThree := runStatement(program.Statements[2])
	// printExpressionStatement(program.Statements[2])
	assertEqual(t, res.number, 30)
	assertEqual(t, resTwo.number, 17)
	assertEqual(t, resThree.number, 22)

}

func TestNegativeNumbers(t *testing.T) {
	input := `3+-1;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0])
	assertEqual(t, res.number, 2)
}

/**/
func printExpressionStatement(stmt Statement) {
	fmt.Printf("\n%v \n", stmt)
}

func TestLetAssignment(t *testing.T) {
	input := `let test <- 1;
						let two <- 1+2;
						let three <- test+two;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0])
	assertEqual(t, res.number, 1)
	assertEqual(t, res.ident, "test")
	resTwo := runStatement(program.Statements[1])
	assertEqual(t, resTwo.ident, "two")
	assertEqual(t, resTwo.number, 3)

	resThree := runStatement(program.Statements[2])
	printExpressionStatement(program.Statements[2])
	assertEqual(t, resThree.ident, "three")
	assertEqual(t, resThree.number, 4)
}
