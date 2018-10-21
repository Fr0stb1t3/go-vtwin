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
		val, _ := strconv.Atoi(uax.Operand.Literal)
		return val
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

func TestEvaluate(t *testing.T) {
	input := "1+1+2;" //2*7+3;3+2*7;

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0])
	assertEqual(t, res.number, 4)
}
func TestBasicMathTrees(t *testing.T) {
	// TODO fix expressionParse
	input := `1+2+4-5;
						1+2+4*5;;
						1+2*4+5;;
						(2+1)*(4+5);
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0])
	resPrecedence := runStatement(program.Statements[1])
	resPrecedenceTwo := runStatement(program.Statements[2])
	resBraces := runStatement(program.Statements[3])

	assertEqual(t, res.number, 2)
	assertEqual(t, resPrecedence.number, 23)
	assertEqual(t, resPrecedenceTwo.number, 14)
	assertEqual(t, resBraces.number, 27)
}

func TestPrecedenceTwo(t *testing.T) {
	input := "27-6/3+5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0])

	assertEqual(t, res.number, 30)

}

/**/
func TestLetAssignment(t *testing.T) {
	input := "let test := 1;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	tree := program.Statements[0]
	fmt.Printf("Tree stringified %v\n", tree)
	res := runStatement(program.Statements[0])
	assertEqual(t, res.number, 1)
	assertEqual(t, res.ident, "test")
}
