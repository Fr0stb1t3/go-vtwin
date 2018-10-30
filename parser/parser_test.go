package parser_test

import (
	"fmt"
	"time"

	"math/rand"
	"strconv"
	"testing"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/parser"
	"github.com/Fr0stb1t3/go-vtwin/token"
	"github.com/Knetic/govaluate"
)

func RandomOperator() string {
	ran := rand.Intn(3)
	switch ran {
	case 0:
		return "+"
	case 1:
		return "-"
	case 2:
		return "*"
	case 3:
		return "/"
	}
	return "+"
}
func GenerateRandomExpression(operationCount int) string {
	rand.Seed(time.Now().UnixNano())
	count := 1
	expression := ""
	var longparen bool = false
	for count < operationCount {
		A := rand.Intn(100) + 1
		B := rand.Intn(100) + 1
		parens := (A % 29) == 0
		secondParen := (B % 23) == 0 // Adds parentheses around multiple subexpressions
		lparen := ""
		rparen := ""
		if longparen == true {
			rparen = ")"
			longparen = false
		} else if longparen == false && parens == true {
			lparen = "("
			rparen = ")"
		} else if secondParen {
			lparen = "("
			longparen = true
		}

		Op := RandomOperator()
		expression = expression + lparen + strconv.Itoa(A) + Op + strconv.Itoa(B) + rparen + Op
		count += count
	}

	expression = expression[:len(expression)-1] // Remove last operator
	if longparen {
		expression = expression + ")"
	}
	return expression
}
func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		t.Fatalf("%s != %s input:"+message, a, b)
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

func evaluateUnaryExpr(ex parser.Expression, scope parser.Scope) int {
	uax := ex.(parser.UnaryExpression)
	switch uax.Operand.Type {
	case token.INT:
		val, _ := strconv.Atoi(uax.Operator.Literal + uax.Operand.Literal)
		return val
	case token.IDENT:
		val := scope.Lookup(uax.Operand.Literal)
		if val == nil {
			panic("Identifier not found")
		}
		return evaluateExpression(val.Expr, scope)
	default:
		panic("evaluateUnaryExpr: Unknown type")
	}
}

func evaluateBinaryExpr(ex parser.Expression, scope parser.Scope) int {
	be := ex.(parser.BinaryExpression)
	a := evaluateExpression(be.Left, scope)
	b := evaluateExpression(be.Right, scope)

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

func evaluateExpression(ex parser.Expression, scope parser.Scope) int {
	switch ex.(type) {
	case parser.ParenExpression:
		pex := ex.(parser.ParenExpression)
		return evaluateExpression(pex.Expr, scope)
	case parser.BinaryExpression:
		return evaluateBinaryExpr(ex, scope)
	case parser.UnaryExpression:
		return evaluateUnaryExpr(ex, scope)
	}
	return 0
}

type result struct {
	number int
	ident  string
}

func runStatement(stmt parser.Statement, scope parser.Scope) result {
	switch stmt.(type) {
	case parser.ExpressionStatement:
		es := stmt.(parser.ExpressionStatement)
		val := evaluateExpression(es.Expr, scope)
		return result{
			number: val,
		}
	case parser.LetStatement:
		ls := stmt.(parser.LetStatement)
		val := evaluateExpression(ls.Expr, scope)
		return result{
			ident:  ls.Name.Value,
			number: val,
		}
	}
	return result{}
}

func TestBasicMathTrees(t *testing.T) {
	input := `1+2+4-5;
						1+2+4*5;
						1+2*4+5;
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0], p.TopScope)
	resPrecedence := runStatement(program.Statements[1], p.TopScope)
	resPrecedenceTwo := runStatement(program.Statements[2], p.TopScope)
	assertEqual(t, res.number, 2, "")
	assertEqual(t, resPrecedence.number, 23, "")
	assertEqual(t, resPrecedenceTwo.number, 14, "")
}

func TestBracesTwo(t *testing.T) {
	input := `
		(2+1)*(4+5);
		((2+1))*(4+5);
		1+6+(6*1)*5*5*4-6;
		8+4+2-4-(3*5)*(3-7);
	`
	//%!s(int=-60) != %!s(int=70) input:8+4+2-4-(3*5)*(3-7)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0], p.TopScope)
	assertEqual(t, res.number, 27, "")
	restwo := runStatement(program.Statements[1], p.TopScope)
	assertEqual(t, restwo.number, 27, "")
	restwo = runStatement(program.Statements[2], p.TopScope)
	assertEqual(t, restwo.number, 601, "1+6+(6*1)*5*5*4-7")

	restwo = runStatement(program.Statements[3], p.TopScope)

	assertEqual(t, restwo.number, 70, "8+4+2-4-(3*5)*(3-7)")
}

func TestPrecedenceTwo(t *testing.T) {
	input := `
		27-6/3+5;
		27-6/3*5;
		27-6/3*5/2;
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0], p.TopScope)
	resTwo := runStatement(program.Statements[1], p.TopScope)
	resThree := runStatement(program.Statements[2], p.TopScope)
	assertEqual(t, res.number, 30, "27-6/3+5")
	assertEqual(t, resTwo.number, 17, "")
	assertEqual(t, resThree.number, 22, "")
}

func TestNegativeNumbers(t *testing.T) {
	input := `3+-1;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0], p.TopScope)
	assertEqual(t, res.number, 2, "")
}

// 9223372036854775807
// 801128428168068000
func randomExpressionTest(t *testing.T) {
	input := GenerateRandomExpression(15)
	expression, _ := govaluate.NewEvaluableExpression(input)

	resultAny, _ := expression.Evaluate(nil)
	resultFloat := resultAny.(float64)
	result := int(resultFloat)
	// fmt.Printf("\n%v", input)
	// fmt.Printf("\n%v", result)
	l := lexer.New(input + ";")
	p := parser.New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0], p.TopScope)
	assertEqual(t, res.number, result, input)
}

func TestRandomLoopSet(t *testing.T) {
	count := 1
	for count < 10000 {
		randomExpressionTest(t)
		count++
	}
}

func printExpressionStatement(stmt parser.Statement) {
	fmt.Printf("\n%v \n", stmt)
}

func TestLetAssignment(t *testing.T) {
	input := `let test <- 1;
						let two <- 1+2;
						let three <- test+two+1;
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	res := runStatement(program.Statements[0], p.TopScope)
	assertEqual(t, res.number, 1, "")
	assertEqual(t, res.ident, "test", "")

	resTwo := runStatement(program.Statements[1], p.TopScope)
	assertEqual(t, resTwo.ident, "two", "")
	assertEqual(t, resTwo.number, 3, "")

	resThree := runStatement(program.Statements[2], p.TopScope)
	assertEqual(t, resThree.ident, "three", "")
	assertEqual(t, resThree.number, 5, "")
}
