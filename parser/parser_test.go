package parser

import (
	"fmt"
	"testing"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
)

func TestSimpleTree(t *testing.T) {
	input := "1+2+4-5;" //2*7+3;3+2*7;

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	fmt.Printf("%v\n", program)
}

func TestPrecedence(t *testing.T) {
	input := "1+2+4*5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	fmt.Printf("%v\n", program)
}
func TestPrecedenceTwo(t *testing.T) {
	input := "1+2*4+5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	// fmt.Printf("%v\n", tok)
	fmt.Printf("%v\n", program)
}
