package parser_test

import (
	"fmt"
	"time"

	"math/rand"
	"strconv"
	"testing"

	"github.com/antonikliment/go-vtwin/ast"
	"github.com/antonikliment/go-vtwin/token"
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

func evaluateUnaryExpr(ex ast.Expression, scope *ast.Scope) int {
	uax := ex.(ast.UnaryExpression)
	switch uax.Operand.(type) {
	case *ast.SimpleLiteral:
		simple := uax.Operand.(*ast.SimpleLiteral)
		if simple.Value == "true" { // TODO simplify this
			val, _ := strconv.Atoi(uax.Operator.Literal + "1")
			return val
		}
		val, _ := strconv.Atoi(uax.Operator.Literal + simple.Value)
		return val
	case *ast.Identifier:
		ident := uax.Operand.(*ast.Identifier)
		ref := scope.Lookup(ident.Value)
		if ref == nil {
			panic("Identifier '" + ident.Value + "' not found")
		}
		out := evaluateExpression(ref.Value(), scope)
		return out
	default:
		// fmt.Printf("x- %v\n", uax)
		panic("x- %v")
		// return -1

	}
	/*
		uLiteral := uax.Operand.(*ast.SimpleLiteral)
		switch uLiteral.Type {
		// case token.INT:
		// 	val, _ := strconv.Atoi(uax.Operator.Literal + uLiteral.Value)
		// 	return val
		// case token.IDENT:
		// 	ref := scope.Lookup(uLiteral.Value)
		// 	if ref == nil {
		// 		panic("Identifier '" + uLiteral.Value + "' not found")
		// 	}
		// 	return evaluateExpression(ref.Value(), scope)
		case token.TRUE:
			return 1
		case token.FALSE:
			return 0
		default:
			panic("evaluateUnaryExpr: Unknown type" + uLiteral.Value)
		}*/
}

func evaluateBinaryExpr(ex ast.Expression, scope *ast.Scope) int {
	be := ex.(ast.BinaryExpression)

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

	fmt.Printf("%v\n", be)
	fmt.Printf("%v\n", be.Left)
	fmt.Printf("%v\n", be.Right)
	panic("Opsie dopsie")
	//return 0
}

func evaluateExpression(ex ast.Expression, scope *ast.Scope) int {
	switch ex.(type) {
	case ast.ParenExpression:
		pex := ex.(ast.ParenExpression)
		return evaluateExpression(pex.Expr, scope)
	case ast.BinaryExpression:
		return evaluateBinaryExpr(ex, scope)
	case ast.UnaryExpression:
		return evaluateUnaryExpr(ex, scope)
	default:
		fmt.Printf("Unknown expression %v \n", ex)
	}
	return -1
}

type result struct {
	number int
	ident  string
}

func runStatements(stmts []ast.Statement, scope *ast.Scope) result {
	res := result{}
	for _, stmt := range stmts {
		value := runStatement(stmt, scope)
		_, ok := stmt.(ast.ReturnStatement)
		if ok == true {
			res = value
		}
	}
	return res
}

func runStatement(stmt ast.Statement, scope *ast.Scope) result {
	switch stmt.(type) {
	case ast.ExpressionStatement:
		es := stmt.(ast.ExpressionStatement)
		val := evaluateExpression(es.Expr, scope)
		return result{
			number: val,
		}
	case ast.LetStatement:
		ls := stmt.(ast.LetStatement)
		val := evaluateExpression(ls.Expr, scope)
		// fmt.Printf("%v \n", ls.Name.Value)
		return result{
			ident:  ls.Name.Value,
			number: val,
		}
	case ast.ConstStatement:
		ls := stmt.(ast.ConstStatement)
		val := evaluateExpression(ls.Expr, scope)
		return result{
			ident:  ls.Name.Value,
			number: val,
		}
	case ast.BlockStatement:
		bs := stmt.(ast.BlockStatement)
		return runStatements(bs.Statements, scope)
	case ast.ReturnStatement:
		rs := stmt.(ast.ReturnStatement)
		val := evaluateExpression(rs.ReturnVal, scope)
		return result{
			number: val,
		}
	case ast.Function:
		return result{}
	default:
		fmt.Printf("%v \n", stmt)
		panic("New statement")
	}
}
