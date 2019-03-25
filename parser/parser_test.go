package parser_test

import (
	"fmt"

	"testing"

	"github.com/Knetic/govaluate"
	"github.com/antonikliment/go-vtwin/ast"
	"github.com/antonikliment/go-vtwin/lexer"
	"github.com/antonikliment/go-vtwin/parser"
)

func parseInput(input string) ([]ast.Statement, *ast.Scope) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return program.Statements, p.TopScope
}

func TestBasicMathTrees(t *testing.T) {
	input := `1+2+4-5;
						1+2+4*5;
						1+2*4+5;
	`
	statements, scope := parseInput(input)

	res := runStatement(statements[0], scope)
	resPrecedence := runStatement(statements[1], scope)
	resPrecedenceTwo := runStatement(statements[2], scope)
	assertEqual(t, res.number, 2, "")
	assertEqual(t, resPrecedence.number, 23, "")
	assertEqual(t, resPrecedenceTwo.number, 14, "")
}

func TestBraces(t *testing.T) {
	input := `
		(2+1)*(4+5);
		((2+1))*(4+5);
		1+6+(6*1)*5*5*4-6;
		8+4+2-4-(3*5)*(3-7);
	`
	statements, scope := parseInput(input)

	res := runStatement(statements[0], scope)
	assertEqual(t, res.number, 27, "")
	restwo := runStatement(statements[1], scope)
	assertEqual(t, restwo.number, 27, "")
	restwo = runStatement(statements[2], scope)
	assertEqual(t, restwo.number, 601, "1+6+(6*1)*5*5*4-7")

	restwo = runStatement(statements[3], scope)

	assertEqual(t, restwo.number, 70, "8+4+2-4-(3*5)*(3-7)")
}

func TestPrecedence(t *testing.T) {
	input := `
		27-6/3+5;
		27-6/3*5;
		27-6/3*5/2;
	`

	statements, scope := parseInput(input)

	res := runStatement(statements[0], scope)
	resTwo := runStatement(statements[1], scope)
	resThree := runStatement(statements[2], scope)
	assertEqual(t, res.number, 30, "27-6/3+5")
	assertEqual(t, resTwo.number, 17, "")
	assertEqual(t, resThree.number, 22, "")
}

func TestNegativeNumbers(t *testing.T) {
	input := `3+(-1);`

	statements, scope := parseInput(input)

	res := runStatement(statements[0], scope)
	assertEqual(t, res.number, 2, "")
}

func randomExpressionTest(t *testing.T) {
	input := GenerateRandomExpression(15)
	expression, _ := govaluate.NewEvaluableExpression(input)

	resultAny, _ := expression.Evaluate(nil)
	resultFloat := resultAny.(float64)
	result := int(resultFloat)

	statements, scope := parseInput(input + ";")

	res := runStatement(statements[0], scope)
	assertEqual(t, res.number, result, input)
}

/*
func TestRandomLoopSet(t *testing.T) {
	count := 1
	for count < 10000 {
		randomExpressionTest(t)
		count++
	}
}
*/
func printExpressionStatement(stmt ast.Statement) {
	fmt.Printf("\n%v \n", stmt)
}

func TestLetAssignment(t *testing.T) {
	input := `const test <- 1;
						let two <- 1+2;
						let three <- test+two+1;
						three <- 3;
						let four <- true;
	`

	statements, scope := parseInput(input + ";")

	res := runStatement(statements[0], scope)
	assertEqual(t, res.number, 1, "const test <- 1")
	assertEqual(t, res.ident, "test", "")

	resTwo := runStatement(statements[1], scope)
	assertEqual(t, resTwo.ident, "two", "")
	assertEqual(t, resTwo.number, 3, "")

	resThree := runStatement(statements[2], scope)
	assertEqual(t, resThree.ident, "three", "")
	assertEqual(t, resThree.number, 5, "let three <- test+two+1")

	resThreeTwo := runStatement(statements[3], scope)
	assertEqual(t, resThreeTwo.ident, "three", "")
	assertEqual(t, resThreeTwo.number, 3, "")

	resThreeFour := runStatement(statements[4], scope)
	assertEqual(t, resThreeFour.ident, "four", "")
	assertEqual(t, resThreeFour.number, 1, "")
}

/*
func TestBlockStatement(t *testing.T) {
	input := `{
		const test <- 3;
		let out <- 6 - test;
		return out + 1;
	}
	` // TODO FIX SEMICOLON
	statements, scope := parseInput(input + ";")
	res := runStatement(statements[0], scope)
	assertEqual(t, res.number, 4, "const test <- 1")
}

*/
func TestConditionStatement(t *testing.T) {

}
func TestLoopBlock(t *testing.T) {

}

func TestFuncBlock(t *testing.T) {
	input := `
	func returnSomething() {
		return 1;
	};

	` //let A <- returnSomething();
	statements, scope := parseInput(input + ";")

	runStatements(statements, scope)
	// assertEqual(t, res.number, 1, "const test <- 1")
}
