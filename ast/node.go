package ast

import (
	"github.com/Fr0stb1t3/go-vtwin/token"
)

type Node interface {
	Token() token.Token
	String() string
}
